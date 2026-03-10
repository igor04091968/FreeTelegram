package router

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type Options struct {
	Domains            []string
	DNSServers         []string
	AggregateThreshold int
	ApplyRoutes        bool
	Interface          string
}

type Updater struct {
	opts Options
}

func New(opts Options) *Updater {
	if opts.AggregateThreshold <= 0 {
		opts.AggregateThreshold = 3
	}
	return &Updater{opts: opts}
}

func (u *Updater) RunOnce(ctx context.Context) ([]string, error) {
	ips, err := u.resolveAll(ctx)
	if err != nil {
		return nil, err
	}
	routes := aggregateRoutes(ips, u.opts.AggregateThreshold)
	if u.opts.ApplyRoutes {
		if err := applyRoutes(routes, u.opts.Interface); err != nil {
			return routes, err
		}
	}
	return routes, nil
}

func (u *Updater) resolveAll(ctx context.Context) ([]net.IP, error) {
	if len(u.opts.Domains) == 0 {
		return nil, errors.New("no domains configured")
	}
	var out []net.IP
	for _, d := range u.opts.Domains {
		if strings.HasPrefix(d, "*.") {
			continue
		}
		ips, err := resolveDomain(ctx, d, u.opts.DNSServers)
		if err == nil {
			out = append(out, ips...)
		}
	}
	if len(out) == 0 {
		return nil, errors.New("no IPs resolved")
	}
	return uniqueIPs(out), nil
}

func resolveDomain(ctx context.Context, domain string, servers []string) ([]net.IP, error) {
	if len(servers) == 0 {
		return net.DefaultResolver.LookupIP(ctx, "ip4", domain)
	}
	var lastErr error
	for _, s := range servers {
		resolver := &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{Timeout: 3 * time.Second}
				return d.DialContext(ctx, "udp", net.JoinHostPort(s, "53"))
			},
		}
		ips, err := resolver.LookupIP(ctx, "ip4", domain)
		if err == nil {
			return ips, nil
		}
		lastErr = err
	}
	return nil, lastErr
}

func uniqueIPs(ips []net.IP) []net.IP {
	seen := map[string]struct{}{}
	out := make([]net.IP, 0, len(ips))
	for _, ip := range ips {
		v4 := ip.To4()
		if v4 == nil {
			continue
		}
		key := v4.String()
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, v4)
	}
	return out
}

func aggregateRoutes(ips []net.IP, threshold int) []string {
	count := map[string]int{}
	for _, ip := range ips {
		p := prefix24(ip)
		count[p]++
	}
	routes := map[string]struct{}{}
	for _, ip := range ips {
		p := prefix24(ip)
		if count[p] >= threshold {
			routes[p+".0/24"] = struct{}{}
			continue
		}
		routes[ip.String()+"/32"] = struct{}{}
	}
	out := make([]string, 0, len(routes))
	for r := range routes {
		out = append(out, r)
	}
	sort.Strings(out)
	return out
}

func prefix24(ip net.IP) string {
	v4 := ip.To4()
	return fmt.Sprintf("%d.%d.%d", v4[0], v4[1], v4[2])
}

func applyRoutes(routes []string, iface string) error {
	if iface == "" {
		return errors.New("interface is required")
	}
	for _, r := range routes {
		cmd := exec.Command("ip", "route", "replace", r, "dev", iface)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("apply route %s: %w", r, err)
		}
	}
	return nil
}

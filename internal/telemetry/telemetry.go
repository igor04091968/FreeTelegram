package telemetry

import (
	"sync/atomic"
	"time"
)

type Stats struct {
	lastUpdateUnix int64
	lastErrorUnix  int64
	lastErrorMsg   atomic.Value
	routesApplied  int64
}

func New() *Stats {
	return &Stats{}
}

func (s *Stats) MarkUpdate(routes int) {
	atomic.StoreInt64(&s.lastUpdateUnix, time.Now().Unix())
	atomic.StoreInt64(&s.routesApplied, int64(routes))
}

func (s *Stats) MarkError(msg string) {
	atomic.StoreInt64(&s.lastErrorUnix, time.Now().Unix())
	s.lastErrorMsg.Store(msg)
}

func (s *Stats) Snapshot() map[string]any {
	msg, _ := s.lastErrorMsg.Load().(string)
	return map[string]any{
		"last_update_unix": atomic.LoadInt64(&s.lastUpdateUnix),
		"last_error_unix":  atomic.LoadInt64(&s.lastErrorUnix),
		"last_error_msg":   msg,
		"routes_applied":   atomic.LoadInt64(&s.routesApplied),
	}
}

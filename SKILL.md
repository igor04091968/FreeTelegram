---
name: freetelegram
context: FreeTelegram OSS + startup grant prep
---

# FreeTelegram Skill

## Purpose
Maintain the FreeTelegram repository and keep grant/application materials aligned with the project narrative.

## Scope
- Repo root: `/mnt/usb_hdd1/Projects/FreeTelegram`.
- Keep README, LICENSE, and contribution docs consistent.
- Maintain a short, reusable grant/application text.
- Keep Go scaffolding, routing updater, and telemetry in sync with config.

## Key Files
- `README.md`
- `LICENSE`
- `CONTRIBUTING.md`
- `CODE_OF_CONDUCT.md`
- `SECURITY.md`
- `ROADMAP.md`
- `config.json`
- `cmd/freetelegram/main.go`
- `internal/config/config.go`
- `internal/router/updater.go`
- `internal/router/worker.go`
- `internal/telemetry/telemetry.go`
- `internal/modules/module.go`
- `internal/queue/queue.go`

## Current Architecture (Go)
- Config versioning in `internal/config/config.go`.
- Queue base in `internal/queue/queue.go`.
- Module interface in `internal/modules/module.go`.
- Route auto-updater in `internal/router/updater.go`.
- Periodic worker in `internal/router/worker.go`.
- Telemetry stats in `internal/telemetry/telemetry.go`.
- Entrypoint wires routing worker in `cmd/freetelegram/main.go`.

## Config fields (routing)
- `routing.domains` (list of domains to resolve)
- `routing.dns_server_list`
- `routing.aggregate_threshold` (default 3)
- `routing.apply_routes` (bool)
- `routing.interface` (e.g. `anet-client`)
- `routing.update_interval_sec` (e.g. 300)

## Application Pack (draft)
Use this text when asked to prepare grant submissions:

**Project name:** FreeTelegram

**Summary:**
Open-source toolkit for resilient access to communication services under network restrictions. Focus on client routing, auto-updating routes, and operational diagnostics. The project is designed to adapt to changing blocking conditions using AI-assisted analysis of network signals and logs.

**Problem:**
Access disruptions reduce the ability to communicate during network blocking.

**Solution:**
Client + routing tooling to keep connectivity stable with minimal user action. Includes automated route updates and transport tuning.

**Stage:** Idea / pre-MVP

**Maintainer:** Igor Rachkov (Developer)
GitHub: https://github.com/igor04091968/

**License:** MIT

**API usage (if requested):**
- accelerate development (code generation/refactoring, tests)
- analyze logs and config drift
- generate docs/runbooks
Expected usage: ~1–3M tokens/month during active development.

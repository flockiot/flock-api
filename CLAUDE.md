# flock-api

Go backend for the Flock device fleet management platform. One binary, multiple runtime targets selected via `--target` flag or `FLOCK_TARGET` env var.

## Build & Run

```bash
go build -o flock-api .
./flock-api --target=all
```

## Test

```bash
go test ./...
```

## Architecture

Single Go binary with these targets (each runs as a separate Kubernetes Deployment in production, or all together with `--target=all`):

- **api** — User-facing REST API (auth, CRUD, SSE, WebSocket)
- **ingester** — Subscribes to device MQTT topics, writes state/telemetry, publishes internal events
- **scheduler** — Rollout execution engine, subscribes to internal MQTT topics, drives deployments
- **builder** — Multi-arch container image builds
- **delta** — Binary delta patches between OCI image layers
- **registry-proxy** — Proxy in front of Harbor for delta serving
- **tunnel** — WireGuard overlay network for device SSH and reverse tunnels
- **proxy** — Public device URL routing (`<uuid>.devices.flock.io`)
- **events-gateway** — Customer-facing event streaming (webhooks, WebSocket)

## Data Stores (added as needed per build phase)

- **PostgreSQL** — Relational data (orgs, users, fleets, devices, releases, deployments, rollouts)
- **Valkey** — Live device state (pluggable KV store interface, added later)
- **MQTT broker (EMQX)** — Central message bus for device and internal async communication (added later)

## Auth

Delegated to Zitadel (self-hosted). API validates JWTs but never issues them for human users. Devices use provisioning keys → short-lived device-scoped JWTs. RBAC: owner/admin/developer/viewer per org.

## Conventions

- Tests are part of every task, not a follow-up
- No stubs or scaffolding unless explicitly asked
- Prefer explicit and readable over clever
- Work in small, verifiable increments

## Code Style

- Do NOT write comments that explain what the code does — the code should speak for itself
- Only write comments when explaining **why** something non-obvious is done, or linking to an issue/spec
- No redundant/repetitive comments

## Git Workflow

- Each task gets its own branch off `main`
- Use **conventional commits**: `feat:`, `fix:`, `chore:`, `docs:`, `test:`, `ci:`
- Commit message must reference the GitHub issue: `Closes #N`
- One commit per PR, squash if needed
- Open PR with `gh pr create` and add `--reviewer SyntaxSmith106`

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
- **ingester** — MQTT connections from devices, writes state/telemetry, publishes to Kafka
- **scheduler** — Rollout execution engine, consumes Kafka events, drives deployments
- **builder** — Multi-arch container image builds
- **delta** — Binary delta patches between OCI image layers
- **registry-proxy** — Proxy in front of Harbor for delta serving
- **tunnel** — WireGuard overlay network for device SSH and reverse tunnels
- **proxy** — Public device URL routing (`<uuid>.devices.flock.io`)
- **events-gateway** — Customer-facing event streaming (webhooks, WebSocket, Kafka forwarding)

## Data Stores

- **PostgreSQL** — Relational data (orgs, users, fleets, devices, releases, deployments, rollouts)
- **KV store** — Live device state (interface-based: Redis/Valkey first, ScyllaDB/DynamoDB later)
- **ClickHouse** — Time-series telemetry, container logs, connectivity events
- **Loki** — Audit logs
- **Harbor** — OCI image registry
- **S3-compatible** — Delta patches, build logs, OS image artifacts
- **Kafka** — Internal message bus + customer event system
- **Apicurio Schema Registry** — Kafka schema enforcement

## Auth

Delegated to Zitadel (self-hosted). API validates JWTs but never issues them for human users. Devices use provisioning keys → short-lived device-scoped JWTs. RBAC: owner/admin/developer/viewer per org.

## Conventions

- Tests are part of every task, not a follow-up
- No stubs or scaffolding unless explicitly asked
- Prefer explicit and readable over clever
- Work in small, verifiable increments

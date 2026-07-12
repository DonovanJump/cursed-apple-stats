# cursed-apple-stats

Deadlock player tracking and analytics platform built with React (Next.js), Go, Rust, and Postgres.

## Stack

- Web: Next.js + TypeScript
- API/Auth + ingestion orchestration: Go
- Analytics worker: Rust
- Database: Postgres

## Monorepo Layout

- `apps/web`: frontend dashboard
- `apps/api-go`: API server and Steam auth/session entry point
- `apps/worker-rust`: async analytics worker
- `db/migrations`: SQL schema and future migrations
- `docs`: architecture and route notes

## Why This Split

- Go owns request/response APIs, Steam auth, Deadlock sync orchestration, and rate-limit-safe caching boundaries.
- Rust owns scheduled aggregate computation and generated insights.
- Postgres stores canonical history and precomputed summary tables.

## Local Quick Start

1. Copy `.env.example` to `.env` and set values.
2. Start Postgres (Docker):
   - `docker compose up -d db`
3. Apply migration:
   - Migration is auto-applied on first DB initialization via `docker-entrypoint-initdb.d`.
   - If the `pg_data` volume already exists, init scripts will not rerun.
   - To re-run initialization from scratch: `docker compose down -v` then `docker compose up -d db`.
   - Schema validation notes: `docs/deadlock-schema-check.md`
4. Start API:
   - `cd apps/api-go`
   - `go run ./cmd/api`
5. Start Worker:
   - `cd apps/worker-rust`
   - `cargo run`
6. Start Web:
   - `cd apps/web`
   - `npm install`
   - `npm run dev`

## MVP Feature Checklist

- Steam login entry flow in Go API
- Player sync endpoint triggered manually
- Match history + basic hero and item aggregates
- Rust computed insights written to `generated_insights`
- Dashboard reads only from your Go API

## Next Build Steps

1. Implement Steam OpenID callback in Go (`/auth/steam/login`, `/auth/steam/callback`).
2. Add Deadlock API client integration in Go sync service.
3. Replace placeholder dashboard API calls with real summary endpoints.
4. Add a cron/scheduler for incremental sync and worker recomputation.
5. Add friend-group comparison endpoints and UI.

# Architecture Notes

## Services

- `web`: Next.js frontend for auth flow, dashboard, and friend comparison pages.
- `api-go`: Canonical API and sync orchestrator. Owns auth/session and Deadlock API integration.
- `worker-rust`: Scheduled analytics and insight generation.

## Data Flow

1. User logs in through Steam OpenID via Go API.
2. API resolves both `steam_id64` and `account_id` and stores user.
3. API sync job fetches match history and metadata from Deadlock API.
4. Normalized rows are upserted into Postgres.
5. Rust worker recomputes summary tables and generated insights.
6. Dashboard endpoints read precomputed rows for low latency and fewer upstream calls.

## Rate-Limit Strategy

- Never call Deadlock API directly from browser.
- Cache synced player summaries in Postgres.
- Incremental sync only for unseen match IDs.
- Add per-user sync cooldown and global concurrency cap in Go service.

## Resume Framing

- Multi-service telemetry analytics platform
- Steam-authenticated ingestion pipeline
- Go API and orchestration + Rust compute worker
- Postgres-backed historical aggregation and derived metrics

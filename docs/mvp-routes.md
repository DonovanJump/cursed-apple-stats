# MVP API Routes

Base path: `/api/v1`

## Health

- `GET /healthz`
  - Liveness check

## Auth (Steam)

- `GET /auth/steam/login`
  - Redirect user to Steam OpenID
- `GET /auth/steam/callback`
  - Validate OpenID response, create or update user session

## User Dashboard

- `GET /me`
  - Current user profile and last sync metadata
- `GET /me/matches`
  - Recent synced match history for the current user
- `POST /me/sync`
  - Triggers incremental sync for current user
- `GET /me/summary`
  - Returns hero stats, item stats, and latest generated insights

## Player Views

- `GET /players/:accountId/matches?limit=50&offset=0`
  - Returns normalized match rows for one player
- `GET /players/:accountId/hero-stats`
  - Returns precomputed player hero aggregates
- `GET /players/:accountId/item-stats`
  - Returns precomputed player item aggregates
- `GET /players/:accountId/insights`
  - Returns generated/funny insights from worker

## Friends/Group

- `POST /groups`
  - Creates a friend group
- `POST /groups/:groupId/members`
  - Adds tracked players to a group
- `GET /groups/:groupId/leaderboard`
  - Returns comparison metrics and funny rankings

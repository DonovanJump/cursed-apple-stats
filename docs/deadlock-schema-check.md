# Deadlock Schema Check (OpenAPI vs DB)

Date checked: 2026-07-12
OpenAPI source: local `api-1.json`

## Result

The migration now matches the critical `PlayerMatchHistoryEntry` fields used for player timelines and summary analytics.

## Endpoint Coverage

### `/v1/players/{account_id}/match-history`

Required OpenAPI fields:

- `match_id`, `account_id`, `hero_id`, `hero_level`, `start_time`
- `game_mode`, `match_mode`, `match_result`
- `player_team`, `player_kills`, `player_deaths`, `player_assists`
- `denies`, `last_hits`, `net_worth`, `match_duration_s`
- `objectives_mask_team0`, `objectives_mask_team1`

Mapped in DB:

- `matches`: `match_id`, `start_time`, `duration_seconds`, `game_mode`, `match_mode`, `match_result`, `objectives_mask_team0`, `objectives_mask_team1`, `raw_json`
- `player_matches`: `account_id`, `hero_id`, `hero_level`, `team`, `kills`, `deaths`, `assists`, `denies`, `last_hits`, `net_worth`, `abandoned_time_seconds`, `team_abandoned`

### `/v1/matches/metadata` and `/v1/matches/{match_id}/metadata`

- Metadata payload shape can evolve and is large.
- Strategy is correct: keep canonical normalized fields and persist full payload in `matches.raw_json` for reprocessing.

## Rate-Limit Notes from OpenAPI

1. Match history endpoint has strict bot-friend limits when using `force_refetch=true`.
2. Bulk metadata endpoint is limited (`IP 10/min`, `Global 100/min`).
3. Single metadata endpoint has very strict Steam fallback limits.

## Ingestion Rules (recommended)

1. Do incremental sync by unseen `match_id`, never full replay each page load.
2. Do not default to `force_refetch=true`.
3. Set `disable_steam=true` for non-critical metadata retries.
4. Backoff failed players and track cooldowns in `player_sync_state`.

## Conclusion

You can safely run migration step 3 with the updated schema.

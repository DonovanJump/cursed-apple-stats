CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    steam_id64 TEXT NOT NULL UNIQUE,
    account_id BIGINT NOT NULL UNIQUE,
    display_name TEXT NOT NULL,
    avatar_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_synced_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS player_sync_state (
    account_id BIGINT PRIMARY KEY,
    last_match_history_sync_at TIMESTAMPTZ,
    last_metadata_sync_at TIMESTAMPTZ,
    next_allowed_refetch_at TIMESTAMPTZ,
    consecutive_failures INTEGER NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS tracked_players (
    id BIGSERIAL PRIMARY KEY,
    owner_user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    account_id BIGINT NOT NULL,
    nickname TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (owner_user_id, account_id)
);

CREATE TABLE IF NOT EXISTS matches (
    match_id BIGINT PRIMARY KEY,
    start_time TIMESTAMPTZ NOT NULL,
    duration_seconds INTEGER NOT NULL,
    game_mode INTEGER,
    match_mode INTEGER,
    match_result INTEGER,
    objectives_mask_team0 INTEGER,
    objectives_mask_team1 INTEGER,
    raw_json JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS player_matches (
    id BIGSERIAL PRIMARY KEY,
    match_id BIGINT NOT NULL REFERENCES matches(match_id) ON DELETE CASCADE,
    account_id BIGINT NOT NULL,
    hero_id INTEGER NOT NULL,
    hero_level INTEGER,
    team INTEGER,
    kills INTEGER,
    deaths INTEGER,
    assists INTEGER,
    denies INTEGER,
    last_hits INTEGER,
    net_worth INTEGER,
    won BOOLEAN,
    abandoned_time_seconds INTEGER,
    team_abandoned BOOLEAN,
    first_death_at_seconds INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (match_id, account_id)
);

CREATE TABLE IF NOT EXISTS match_items (
    id BIGSERIAL PRIMARY KEY,
    match_id BIGINT NOT NULL REFERENCES matches(match_id) ON DELETE CASCADE,
    account_id BIGINT NOT NULL,
    item_id INTEGER NOT NULL,
    purchase_time_seconds INTEGER,
    slot_index INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS player_hero_stats (
    id BIGSERIAL PRIMARY KEY,
    account_id BIGINT NOT NULL,
    hero_id INTEGER NOT NULL,
    games INTEGER NOT NULL DEFAULT 0,
    wins INTEGER NOT NULL DEFAULT 0,
    avg_kda NUMERIC(8, 3) NOT NULL DEFAULT 0,
    last_played_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (account_id, hero_id)
);

CREATE TABLE IF NOT EXISTS player_item_stats (
    id BIGSERIAL PRIMARY KEY,
    account_id BIGINT NOT NULL,
    item_id INTEGER NOT NULL,
    purchases INTEGER NOT NULL DEFAULT 0,
    wins_when_bought INTEGER NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (account_id, item_id)
);

CREATE TABLE IF NOT EXISTS generated_insights (
    id BIGSERIAL PRIMARY KEY,
    account_id BIGINT NOT NULL,
    insight_type TEXT NOT NULL,
    insight_text TEXT NOT NULL,
    score NUMERIC(8, 3),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_player_matches_account_id ON player_matches(account_id);
CREATE INDEX IF NOT EXISTS idx_player_matches_match_id ON player_matches(match_id);
CREATE INDEX IF NOT EXISTS idx_matches_start_time ON matches(start_time DESC);
CREATE INDEX IF NOT EXISTS idx_match_items_account_id ON match_items(account_id);
CREATE INDEX IF NOT EXISTS idx_generated_insights_account_id_created_at ON generated_insights(account_id, created_at DESC);

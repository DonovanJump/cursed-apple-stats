package store

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNoRows = pgx.ErrNoRows

type User struct {
    SteamID64    string
    AccountID    int64
    DisplayName  string
    AvatarURL    *string
    LastSyncedAt *time.Time
}

type RecentMatch struct {
    MatchID            int64      `json:"match_id"`
    StartTime          time.Time  `json:"start_time"`
    DurationSeconds     int        `json:"duration_seconds"`
    GameMode           *int       `json:"game_mode,omitempty"`
    MatchMode          *int       `json:"match_mode,omitempty"`
    MatchResult        *int       `json:"match_result,omitempty"`
    HeroID             int        `json:"hero_id"`
    HeroLevel          *int       `json:"hero_level,omitempty"`
    Team               *int       `json:"team,omitempty"`
    Kills              *int       `json:"kills,omitempty"`
    Deaths             *int       `json:"deaths,omitempty"`
    Assists            *int       `json:"assists,omitempty"`
    Denies             *int       `json:"denies,omitempty"`
    LastHits           *int       `json:"last_hits,omitempty"`
    NetWorth           *int       `json:"net_worth,omitempty"`
    Won                *bool      `json:"won,omitempty"`
    FirstDeathAtSeconds *int       `json:"first_death_at_seconds,omitempty"`
}

type PlayerMatchHistoryEntry struct {
    AbandonedTimeS      *int   `json:"abandoned_time_s,omitempty"`
    AccountID           int64  `json:"account_id"`
    Denies              int    `json:"denies"`
    GameMode            int    `json:"game_mode"`
    HeroID              int    `json:"hero_id"`
    HeroLevel           int    `json:"hero_level"`
    LastHits            int    `json:"last_hits"`
    MatchDurationS      int    `json:"match_duration_s"`
    MatchID             int64  `json:"match_id"`
    MatchMode           int    `json:"match_mode"`
    MatchResult         int    `json:"match_result"`
    NetWorth            int    `json:"net_worth"`
    ObjectivesMaskTeam0 int    `json:"objectives_mask_team0"`
    ObjectivesMaskTeam1 int    `json:"objectives_mask_team1"`
    PlayerAssists       int    `json:"player_assists"`
    PlayerDeaths        int    `json:"player_deaths"`
    PlayerKills         int    `json:"player_kills"`
    PlayerTeam          int    `json:"player_team"`
    StartTime           int64  `json:"start_time"`
    TeamAbandoned       *bool  `json:"team_abandoned,omitempty"`
}

type Store struct {
    pool *pgxpool.Pool
}

func New(ctx context.Context, databaseURL string) (*Store, error) {
    pool, err := pgxpool.New(ctx, databaseURL)
    if err != nil {
        return nil, err
    }

    if err := pool.Ping(ctx); err != nil {
        pool.Close()
        return nil, err
    }

    return &Store{pool: pool}, nil
}

func (s *Store) Close() {
    s.pool.Close()
}

func Steam64ToAccountID(steamID64 string) (int64, error) {
    const steamID64Base = int64(76561197960265728)

    parsed, err := strconv.ParseInt(steamID64, 10, 64)
    if err != nil {
        return 0, fmt.Errorf("parse steam id64: %w", err)
    }

    if parsed < steamID64Base {
        return 0, errors.New("steam id64 is too small to convert")
    }

    return parsed - steamID64Base, nil
}

func AccountIDToSteam64(accountID int64) string {
    const steamID64Base = int64(76561197960265728)
    return strconv.FormatInt(steamID64Base+accountID, 10)
}

func (s *Store) UpsertUser(ctx context.Context, steamID64 string, displayName string, avatarURL *string) (User, error) {
    accountID, err := Steam64ToAccountID(steamID64)
    if err != nil {
        return User{}, err
    }

    if displayName == "" {
        displayName = "Tracked Deadlock Player"
    }

    row := s.pool.QueryRow(ctx, `
        INSERT INTO users (steam_id64, account_id, display_name, avatar_url)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (steam_id64)
        DO UPDATE SET
            account_id = EXCLUDED.account_id,
            display_name = EXCLUDED.display_name,
            avatar_url = EXCLUDED.avatar_url,
            updated_at = NOW()
        RETURNING steam_id64, account_id, display_name, avatar_url, last_synced_at
    `, steamID64, accountID, displayName, avatarURL)

    var user User
    if err := row.Scan(&user.SteamID64, &user.AccountID, &user.DisplayName, &user.AvatarURL, &user.LastSyncedAt); err != nil {
        return User{}, err
    }

    return user, nil
}

func (s *Store) GetUserBySteamID64(ctx context.Context, steamID64 string) (User, error) {
    row := s.pool.QueryRow(ctx, `
        SELECT steam_id64, account_id, display_name, avatar_url, last_synced_at
        FROM users
        WHERE steam_id64 = $1
    `, steamID64)

    var user User
    if err := row.Scan(&user.SteamID64, &user.AccountID, &user.DisplayName, &user.AvatarURL, &user.LastSyncedAt); err != nil {
        return User{}, err
    }

    return user, nil
}

func (s *Store) SetUserLastSyncedAt(ctx context.Context, steamID64 string, syncedAt time.Time) error {
    _, err := s.pool.Exec(ctx, `
        UPDATE users
        SET last_synced_at = $2,
            updated_at = NOW()
        WHERE steam_id64 = $1
    `, steamID64, syncedAt)
    return err
}

func (s *Store) UpsertMatchHistoryEntry(ctx context.Context, entry PlayerMatchHistoryEntry, rawJSON []byte) error {
    startTime := time.Unix(int64(entry.StartTime), 0).UTC()
    matchResult := entry.MatchResult

    _, err := s.pool.Exec(ctx, `
        INSERT INTO matches (
            match_id, start_time, duration_seconds, game_mode, match_mode,
            match_result, objectives_mask_team0, objectives_mask_team1, raw_json
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        ON CONFLICT (match_id)
        DO UPDATE SET
            start_time = EXCLUDED.start_time,
            duration_seconds = EXCLUDED.duration_seconds,
            game_mode = EXCLUDED.game_mode,
            match_mode = EXCLUDED.match_mode,
            match_result = EXCLUDED.match_result,
            objectives_mask_team0 = EXCLUDED.objectives_mask_team0,
            objectives_mask_team1 = EXCLUDED.objectives_mask_team1,
            raw_json = EXCLUDED.raw_json
    `, entry.MatchID, startTime, entry.MatchDurationS, entry.GameMode, entry.MatchMode, matchResult, entry.ObjectivesMaskTeam0, entry.ObjectivesMaskTeam1, rawJSON)
    if err != nil {
        return err
    }

    var won *bool
    if entry.MatchResult == 1 {
        win := true
        won = &win
    } else if entry.MatchResult == 0 {
        loss := false
        won = &loss
    }

    var abandonedTimeSeconds *int
    if entry.AbandonedTimeS != nil {
        abandonedTimeSeconds = entry.AbandonedTimeS
    }

    _, err = s.pool.Exec(ctx, `
        INSERT INTO player_matches (
            match_id, account_id, hero_id, hero_level, team, kills, deaths, assists,
            denies, last_hits, net_worth, won, abandoned_time_seconds, team_abandoned
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
        ON CONFLICT (match_id, account_id)
        DO UPDATE SET
            hero_id = EXCLUDED.hero_id,
            hero_level = EXCLUDED.hero_level,
            team = EXCLUDED.team,
            kills = EXCLUDED.kills,
            deaths = EXCLUDED.deaths,
            assists = EXCLUDED.assists,
            denies = EXCLUDED.denies,
            last_hits = EXCLUDED.last_hits,
            net_worth = EXCLUDED.net_worth,
            won = EXCLUDED.won,
            abandoned_time_seconds = EXCLUDED.abandoned_time_seconds,
            team_abandoned = EXCLUDED.team_abandoned
    `, entry.MatchID, entry.AccountID, entry.HeroID, entry.HeroLevel, entry.PlayerTeam, entry.PlayerKills, entry.PlayerDeaths, entry.PlayerAssists, entry.Denies, entry.LastHits, entry.NetWorth, won, abandonedTimeSeconds, entry.TeamAbandoned)
    return err
}

func (s *Store) GetRecentMatchesBySteamID64(ctx context.Context, steamID64 string, limit int) ([]RecentMatch, error) {
    if limit <= 0 {
        limit = 1000
    }
    if limit > 5000 {
        limit = 5000
    }

    rows, err := s.pool.Query(ctx, `
        SELECT
            pm.match_id, m.start_time, m.duration_seconds, m.game_mode, m.match_mode, m.match_result,
            pm.hero_id, pm.hero_level, pm.team, pm.kills, pm.deaths, pm.assists, pm.denies,
            pm.last_hits, pm.net_worth, pm.won, pm.first_death_at_seconds
        FROM player_matches pm
        JOIN matches m ON m.match_id = pm.match_id
        JOIN users u ON u.account_id = pm.account_id
        WHERE u.steam_id64 = $1
        ORDER BY m.start_time DESC
        LIMIT $2
    `, steamID64, limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    matches := make([]RecentMatch, 0, limit)
    for rows.Next() {
        var match RecentMatch
        if err := rows.Scan(
            &match.MatchID,
            &match.StartTime,
            &match.DurationSeconds,
            &match.GameMode,
            &match.MatchMode,
            &match.MatchResult,
            &match.HeroID,
            &match.HeroLevel,
            &match.Team,
            &match.Kills,
            &match.Deaths,
            &match.Assists,
            &match.Denies,
            &match.LastHits,
            &match.NetWorth,
            &match.Won,
            &match.FirstDeathAtSeconds,
        ); err != nil {
            return nil, err
        }
        matches = append(matches, match)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return matches, nil
}

package syncer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cursed-apple-stats/apps/api-go/internal/deadlock"
	"cursed-apple-stats/apps/api-go/internal/store"
)

type Result struct {
	SteamID64      string    `json:"steam_id64"`
	AccountID      int64     `json:"account_id"`
	MatchesSynced  int       `json:"matches_synced"`
	SyncedAt       time.Time `json:"synced_at"`
}

type Service struct {
	store  *store.Store
	client *deadlock.Client
}

func New(store *store.Store, client *deadlock.Client) *Service {
	return &Service{store: store, client: client}
}

func (s *Service) SyncSteamID(ctx context.Context, steamID64 string) (Result, error) {
	user, err := s.store.GetUserBySteamID64(ctx, steamID64)
	if err != nil {
		user, err = s.store.UpsertUser(ctx, steamID64, "", nil)
		if err != nil {
			return Result{}, err
		}
	}

	entries, err := s.client.GetMatchHistory(ctx, user.AccountID)
	if err != nil {
		return Result{}, err
	}

	for _, entry := range entries {
		rawJSON, _ := json.Marshal(entry)
		if err := s.store.UpsertMatchHistoryEntry(ctx, entry, rawJSON); err != nil {
			return Result{}, fmt.Errorf("store match %d: %w", entry.MatchID, err)
		}
	}

	now := time.Now().UTC()
	if err := s.store.SetUserLastSyncedAt(ctx, steamID64, now); err != nil {
		return Result{}, err
	}

	return Result{
		SteamID64:     steamID64,
		AccountID:     user.AccountID,
		MatchesSynced: len(entries),
		SyncedAt:      now,
	}, nil
}

func (s *Service) RecentMatches(ctx context.Context, steamID64 string, limit int) ([]store.RecentMatch, error) {
	return s.store.GetRecentMatchesBySteamID64(ctx, steamID64, limit)
}
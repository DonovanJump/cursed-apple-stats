package server

import (
	"encoding/json"
	"net/http"

	"cursed-apple-stats/apps/api-go/internal/store"
	"cursed-apple-stats/apps/api-go/internal/syncer"
)

type Server struct {
    store       *store.Store
    steamID64   string
    sync        *syncer.Service
}

func New(appStore *store.Store, steamID64 string, syncService *syncer.Service) *Server {
    return &Server{store: appStore, steamID64: steamID64, sync: syncService}
}

func (s *Server) Routes() http.Handler {
    mux := http.NewServeMux()

    mux.HandleFunc("GET /healthz", s.handleHealth)
    mux.HandleFunc("GET /api/v1/me", s.handleMe)
    mux.HandleFunc("GET /api/v1/me/matches", s.handleMeMatches)
    mux.HandleFunc("POST /api/v1/me/sync", s.handleSyncMe)

    return withJSON(mux)
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
    respondJSON(w, http.StatusOK, map[string]any{
        "ok": true,
    })
}

func (s *Server) handleMe(w http.ResponseWriter, r *http.Request) {
    if s.steamID64 == "" {
        respondJSON(w, http.StatusNotFound, map[string]any{"error": "MY_STEAM_ID is not set"})
        return
    }

    user, err := s.store.GetUserBySteamID64(r.Context(), s.steamID64)
    if err != nil {
        if err == store.ErrNoRows {
            respondJSON(w, http.StatusNotFound, map[string]any{"error": "user not seeded"})
            return
        }
        respondJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
        return
    }

    respondJSON(w, http.StatusOK, map[string]any{
        "steam_id64":     user.SteamID64,
        "account_id":     user.AccountID,
        "display_name":   user.DisplayName,
        "last_synced_at": user.LastSyncedAt,
    })
}

func (s *Server) handleMeMatches(w http.ResponseWriter, r *http.Request) {
    if s.steamID64 == "" {
        respondJSON(w, http.StatusNotFound, map[string]any{"error": "MY_STEAM_ID is not set"})
        return
    }

    matches, err := s.sync.RecentMatches(r.Context(), s.steamID64, 12)
    if err != nil {
        respondJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
        return
    }

    respondJSON(w, http.StatusOK, map[string]any{
        "matches": matches,
    })
}

func (s *Server) handleSyncMe(w http.ResponseWriter, r *http.Request) {
    if s.steamID64 == "" {
        respondJSON(w, http.StatusNotFound, map[string]any{"error": "MY_STEAM_ID is not set"})
        return
    }

    result, err := s.sync.SyncSteamID(r.Context(), s.steamID64)
    if err != nil {
        respondJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
        return
    }

    respondJSON(w, http.StatusOK, map[string]any{
        "steam_id64":      result.SteamID64,
        "account_id":      result.AccountID,
        "matches_synced":  result.MatchesSynced,
        "synced_at":       result.SyncedAt,
    })
}

func withJSON(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        next.ServeHTTP(w, r)
    })
}

func respondJSON(w http.ResponseWriter, status int, payload map[string]any) {
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(payload)
}

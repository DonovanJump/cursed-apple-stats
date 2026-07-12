package server

import (
    "encoding/json"
    "net/http"
)

type Server struct{}

func New() *Server {
    return &Server{}
}

func (s *Server) Routes() http.Handler {
    mux := http.NewServeMux()

    mux.HandleFunc("GET /healthz", s.handleHealth)
    mux.HandleFunc("GET /api/v1/me", s.handleMe)
    mux.HandleFunc("POST /api/v1/me/sync", s.handleSyncMe)

    return withJSON(mux)
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
    respondJSON(w, http.StatusOK, map[string]any{
        "ok": true,
    })
}

func (s *Server) handleMe(w http.ResponseWriter, _ *http.Request) {
    respondJSON(w, http.StatusOK, map[string]any{
        "steam_id64":     "76561198000000000",
        "account_id":     39734272,
        "display_name":   "Demo Deadlock Player",
        "last_synced_at": nil,
    })
}

func (s *Server) handleSyncMe(w http.ResponseWriter, _ *http.Request) {
    // Placeholder: enqueue user sync and return accepted status.
    respondJSON(w, http.StatusAccepted, map[string]any{
        "status": "sync_queued",
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

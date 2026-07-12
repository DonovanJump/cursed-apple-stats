package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"cursed-apple-stats/apps/api-go/internal/config"
	"cursed-apple-stats/apps/api-go/internal/deadlock"
	"cursed-apple-stats/apps/api-go/internal/server"
	"cursed-apple-stats/apps/api-go/internal/store"
	"cursed-apple-stats/apps/api-go/internal/syncer"
)

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    cfg := config.Load()

    appStore, err := store.New(ctx, cfg.DatabaseURL)
    if err != nil {
        log.Fatal(err)
    }
    defer appStore.Close()

    if cfg.MySteamID != "" {
        if _, err := appStore.UpsertUser(ctx, cfg.MySteamID, "", nil); err != nil {
            log.Fatal(err)
        }
    }

    deadlockClient := deadlock.New(cfg.DeadlockBaseURL, cfg.DeadlockAPIKey)
    syncService := syncer.New(appStore, deadlockClient)
    srv := server.New(appStore, cfg.MySteamID, syncService)

    log.Printf("api listening on :%s", cfg.Port)
    if err := http.ListenAndServe(":"+cfg.Port, srv.Routes()); err != nil {
        log.Fatal(err)
    }
}

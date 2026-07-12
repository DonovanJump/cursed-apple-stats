package main

import (
    "log"
    "net/http"

    "cursed-apple-stats/apps/api-go/internal/config"
    "cursed-apple-stats/apps/api-go/internal/server"
)

func main() {
    cfg := config.Load()

    srv := server.New()

    log.Printf("api listening on :%s", cfg.Port)
    if err := http.ListenAndServe(":"+cfg.Port, srv.Routes()); err != nil {
        log.Fatal(err)
    }
}

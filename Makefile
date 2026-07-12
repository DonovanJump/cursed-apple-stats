.PHONY: up down run-api run-worker

up:
	docker compose up -d db

down:
	docker compose down

run-api:
	cd apps/api-go && go run ./cmd/api

run-worker:
	cd apps/worker-rust && cargo run

.PHONY: run build test tidy docker-up docker-down sqlc-gen migrate

# ── Local dev ─────────────────────────────────────────────────────────────────
run:
	go run ./cmd/server

build:
	go build -ldflags="-s -w" -o server ./cmd/server

test:
	go test ./... -v -count=1

tidy:
	go mod tidy

# ── SQLC ──────────────────────────────────────────────────────────────────────
sqlc-gen:
	sqlc generate

# ── Docker ───────────────────────────────────────────────────────────────────
docker-up:
	docker compose up --build -d

docker-down:
	docker compose down -v

# ── DB Migrations (uses psql directly) ───────────────────────────────────────
migrate:
	psql "$$DATABASE_URL" -f db/migrations/001_create_users.up.sql

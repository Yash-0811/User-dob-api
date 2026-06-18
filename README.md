# user-dob-api

A RESTful API built with **Go + GoFiber** that manages users with their date of birth and calculates age dynamically.

## Tech Stack

| Layer | Library |
|-------|---------|
| HTTP framework | [GoFiber v2](https://gofiber.io) |
| Database | PostgreSQL |
| DB access layer | [SQLC](https://sqlc.dev) (generated) |
| Logging | [Uber Zap](https://github.com/uber-go/zap) |
| Validation | [go-playground/validator](https://github.com/go-playground/validator) |
| UUID | [google/uuid](https://github.com/google/uuid) |

---

## Project Structure

```
.
├── cmd/server/main.go          # Entry point
├── config/config.go            # Environment config
├── db/
│   ├── migrations/             # SQL migration files
│   ├── queries/users.sql       # SQLC query definitions
│   └── sqlc/                   # SQLC generated code
├── internal/
│   ├── handler/                # HTTP handlers
│   ├── middleware/             # RequestID + Logger + Recover
│   ├── models/                 # Request / Response DTOs
│   ├── repository/             # Data-access layer
│   ├── routes/                 # Route registration
│   ├── service/                # Business logic + age calculation
│   └── logger/                 # Zap wrapper
├── Dockerfile
├── docker-compose.yml
└── Makefile
```

---

## Quickstart

### Option A — Docker Compose (recommended)

```bash
git clone https://github.com/<you>/user-dob-api
cd user-dob-api

docker compose up --build
# API is available at http://localhost:8080
```

The compose file automatically:
1. Starts a PostgreSQL 16 container
2. Runs the migration (`001_create_users.up.sql`) on first boot
3. Builds and starts the API container

### Option B — Run locally

**Prerequisites:** Go 1.22+, PostgreSQL running

```bash
# 1. Clone and enter
git clone https://github.com/<you>/user-dob-api
cd user-dob-api

# 2. Copy and fill environment variables
cp .env.example .env
# Edit .env with your DB credentials

# 3. Apply migration
psql -U postgres -d userdb -f db/migrations/001_create_users.up.sql

# 4. Install dependencies
go mod tidy

# 5. Run
go run ./cmd/server
```

---

## API Reference

All endpoints use `Content-Type: application/json`.
Every response includes an `X-Request-Id` header for tracing.

### POST /users — Create user

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","dob":"1990-05-10"}'
```

Response `201 Created`:
```json
{"id":1,"name":"Alice","dob":"1990-05-10"}
```

---

### GET /users/:id — Get user (with age)

```bash
curl http://localhost:8080/users/1
```

Response `200 OK`:
```json
{"id":1,"name":"Alice","dob":"1990-05-10","age":35}
```

---

### PUT /users/:id — Update user

```bash
curl -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice Updated","dob":"1991-03-15"}'
```

Response `200 OK`:
```json
{"id":1,"name":"Alice Updated","dob":"1991-03-15"}
```

---

### DELETE /users/:id — Delete user

```bash
curl -X DELETE http://localhost:8080/users/1
```

Response `204 No Content`

---

### GET /users — List users (paginated)

```bash
curl "http://localhost:8080/users?page=1&page_size=10"
```

Response `200 OK`:
```json
{
  "data": [
    {"id":1,"name":"Alice","dob":"1990-05-10","age":35}
  ],
  "total": 1,
  "page": 1,
  "page_size": 10,
  "total_pages": 1
}
```

| Query Param | Default | Description |
|-------------|---------|-------------|
| `page`      | `1`     | Page number |
| `page_size` | `10`    | Items per page (max 100) |

---

### GET /health — Health check

```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

---

## Running Tests

```bash
make test
# or
go test ./... -v
```

The `CalculateAge` function in `internal/service` is covered by unit tests that assert correct age across boundary conditions (birthday today, tomorrow, newborn, etc.).

---

## Regenerating SQLC Code

If you modify `db/queries/users.sql` or `db/migrations/`:

```bash
# Install sqlc (if not already)
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

make sqlc-gen
```

---

## Environment Variables

| Variable      | Default      | Description                    |
|---------------|--------------|--------------------------------|
| `DB_HOST`     | `localhost`  | PostgreSQL host                |
| `DB_PORT`     | `5432`       | PostgreSQL port                |
| `DB_USER`     | `postgres`   | Database user                  |
| `DB_PASSWORD` | `postgres`   | Database password              |
| `DB_NAME`     | `userdb`     | Database name                  |
| `DB_SSLMODE`  | `disable`    | SSL mode                       |
| `APP_PORT`    | `8080`       | HTTP listen port               |
| `APP_ENV`     | `development`| `production` enables JSON logs |

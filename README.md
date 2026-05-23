# tuto-api

Backend API for the Tuto (Lumio) Flutter app. Built with Go + chi.

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.21+ |
| HTTP Router | chi |
| Database | PostgreSQL (pgx) |
| Cache / Sessions | Redis |
| Object Storage | MinIO (S3-compatible) |
| Migrations | golang-migrate |
| Config | Viper (YAML) |

## Project Structure

```
cmd/main.go           # entrypoint
internal/
  auth/               # parent + kid authentication
  kid/                # kid profile, map, lessons
  content/            # subjects, lessons, library
  progress/           # stars, streaks, hearts, badges
  badge/              # badge catalog + engine
  dashboard/          # parent dashboard + weekly report
  pairing/            # device pairing (6-digit code)
pkg/
  config/             # config loader (config.yml)
  db/                 # postgres connection pool
  redis/              # redis client
  storage/            # minio client
  response/           # standard JSON envelope
  middleware/         # JWT, scope check, rate limit
migrations/           # SQL migration files
```

## Prerequisites

Make sure the following tools are installed before running the project:

| Tool | Install | Purpose |
|---|---|---|
| Go 1.21+ | `brew install go` | Language runtime |
| Docker | [docker.com](https://www.docker.com/products/docker-desktop/) | Run infra (Postgres, Redis, MinIO) |
| golang-migrate | `brew install golang-migrate` | Run database migrations |
| sqlc | `brew install sqlc` | Generate type-safe Go from SQL |
| golangci-lint | `brew install golangci-lint` | Linter |
| yq | `brew install yq` | YAML parser (used by migrate script) |

Install all at once:
```bash
brew install go golang-migrate sqlc golangci-lint yq
```

## Getting Started

### 1. Start infrastructure

```bash
docker compose -p tuto up -d
```

### 2. Copy config

```bash
cp config.example.yml config.yml
```

### 3. Run migrations

```bash
./migrate.sh up
```

### 4. Run the server

```bash
go run ./cmd/main.go
```

Server starts on `http://localhost:8080`.

## Health Check

```bash
curl http://localhost:8080/health
```

## Development

### Lint

```bash
golangci-lint run
```

### Test

```bash
go test -race ./...
```

### Build

```bash
go build ./...
```

## Infrastructure UIs

| Service | URL | Credentials |
|---|---|---|
| Adminer (PostgreSQL) | http://localhost:8888 | user / password |
| RedisInsight | http://localhost:5540 | — |
| MinIO Console | http://localhost:9001 | minioadmin / minioadmin |
| Mailpit | http://localhost:8025 | — |

## API Design

See [System Design](https://www.notion.so/01-System-Design-369258e2a09780f08a45fef37ffabbeb) for the full API reference, domain entities, and sequence diagrams.

## Database Schema

See `schema.dbml` — paste into [dbdiagram.io](https://dbdiagram.io) to visualize.

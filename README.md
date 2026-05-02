# Quote Schedule System

REST API built with Go + Gin + GORM + MySQL, with a Next.js frontend.

## Requirements

- [Go 1.21+](https://golang.org/dl/)
- [Docker](https://www.docker.com/)

## Backend setup

**1. Start MySQL**

```bash
docker compose up -d
```

**2. Configure environment**

```bash
cp br-api/.env.example br-api/.env.local
# Edit br-api/.env.local with your credentials
```

**3. Run the API**

```bash
cd br-api
go run ./cmd/api
```

The API starts on `http://localhost:8081` by default (configurable via `PORT` in `.env.local`).  
On first start, tables are created and seed data (managers, technicians, quotes) is inserted automatically.

## API

Swagger UI: `http://localhost:8081/swagger/index.html`

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/health` | Health check |
| GET | `/api/v1/quotes` | List unscheduled quotes |
| POST | `/api/v1/quotes` | Create a quote |
| GET | `/api/v1/technicians` | List technicians |
| GET | `/api/v1/technicians/:id/jobs` | Technician schedule (accepts `?date=YYYY-MM-DD&timezone=Australia/Sydney`) |
| GET | `/api/v1/jobs/:id` | Get job by ID |
| POST | `/api/v1/jobs` | Assign quote to technician |
| PATCH | `/api/v1/jobs/:id/complete` | Mark job as complete |
| GET | `/api/v1/notifications/stream` | SSE stream of notifications |
| PATCH | `/api/v1/notifications/:id/read` | Mark notification as read |

Pagination is supported on list endpoints via `?page=1&page_size=20` (max page_size: 100).

## Testing

**Unit tests** (no DB required):

```bash
cd br-api
go test ./internal/service/...
```

**Integration tests** (requires Docker):

```bash
cd br-api
go test ./internal/integration/... -timeout 300s
```

Integration tests spin up an isolated MySQL container per test via `testcontainers-go` — no local DB setup needed. Each container is torn down automatically after the test completes.

## Environment variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DATABASE_URL` | `root:root@tcp(localhost:3307)/brix?parseTime=true` | MySQL DSN |
| `PORT` | `8080` | Server port |

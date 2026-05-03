# Quote Schedule System

REST API built with Go + Gin + GORM + MySQL, with a Next.js frontend.

## Requirements

- [Go 1.21+](https://golang.org/dl/)
- [Node.js 18+](https://nodejs.org/)
- [Docker](https://www.docker.com/)

## Backend setup

**1. Start MySQL**

```bash
docker compose up -d
```

**2. Configure environment**

```bash
cp br-api/.env.example br-api/.env.local
```

Edit `br-api/.env.local`:

| Variable | Description |
|----------|-------------|
| `DATABASE_URL` | MySQL DSN — format: `user:password@tcp(host:port)/dbname?parseTime=true` |
| `PORT` | Port the API listens on (default: `8080`) |
| `CORS_ORIGIN` | Allowed frontend origin (default: `http://localhost:3000`) |
| `MYSQL_ROOT_PASSWORD` | Root password used by the Docker MySQL container |
| `MYSQL_DATABASE` | Database name created in the Docker container (default: `brix`) |

**3. Run the API**

```bash
cd br-api
go run ./cmd/api
```

The API starts on `http://localhost:8080` by default.  
On first start, tables are created and seed data (managers, technicians, quotes) is inserted automatically.

## Frontend setup

**1. Configure environment**

```bash
cp br-app/.env.example br-app/.env.local
```

Edit `br-app/.env.local`:

| Variable | Description |
|----------|-------------|
| `NEXT_PUBLIC_API_URL` | Base URL of the backend API (default: `http://localhost:8080`) |

**2. Install dependencies**

```bash
cd br-app
npm install
```

**3. Run the dev server**

```bash
npm run dev
```

The app starts on `http://localhost:3000`.

## How it works

### Data model

A `job` is the central record. It holds a foreign key to `quote`, `technician`, and `manager`. `quote_id` is unique on the jobs table, so a quote can only ever be assigned once. A `notification` belongs to a job and targets a recipient via `recipient_type` + `recipient_id` rather than separate FK columns. Using a polymorphic reference keeps the table flat and avoids adding a new nullable FK every time a new recipient type is introduced.

### Assignment flow

The manager sees a list of unscheduled quotes and a list of technicians with their existing slots for the day. They pick a quote, click a technician, and choose a 2-hour window. That hits `POST /api/v1/jobs`, which creates the job, updates the quote status, and sends the technician a notification.

### Conflict prevention

The API rejects any assignment that would give a technician overlapping jobs. If the quote was already scheduled by the time the request arrives, that's caught too and returns a 409.

### Notifications

When a job is assigned the technician gets a notification. When it's completed the manager gets one. Notifications are delivered in real time to the sidebar without needing a page refresh.

## API

Swagger UI: `http://localhost:8080/swagger/index.html`

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/health` | Health check |
| GET | `/api/v1/quotes` | List unscheduled quotes |
| POST | `/api/v1/quotes` | Create a quote |
| GET | `/api/v1/managers` | List managers |
| GET | `/api/v1/technicians` | List technicians |
| GET | `/api/v1/technicians/:id/jobs` | Technician schedule (accepts `?date=YYYY-MM-DD`) |
| GET | `/api/v1/jobs/:id` | Get job by ID |
| POST | `/api/v1/jobs` | Assign quote to technician |
| PATCH | `/api/v1/jobs/:id/complete` | Mark job as complete |
| GET | `/api/v1/notifications` | List notifications |
| GET | `/api/v1/notifications/stream` | SSE stream of notifications |
| PATCH | `/api/v1/notifications/:id/read` | Mark notification as read |

Pagination is supported on list endpoints via `?page=1&page_size=20` (max page_size: 100).

## Testing

**Backend unit tests** (no DB required):

```bash
cd br-api
go test ./internal/service/...
```

**Backend integration tests** (requires Docker):

```bash
cd br-api
go test ./internal/integration/... -timeout 300s
```

Integration tests spin up an isolated MySQL container per test via `testcontainers-go`.

**Frontend unit tests**:

```bash
cd br-app
npx jest
```

With coverage:

```bash
npx jest --coverage
```

## Architecture decisions & trade-offs

### Husky
Pre-commit hooks run Prettier, ESLint, and the frontend unit tests before every commit. Keeps formatting consistent and catches regressions before they land. In this project it is used to run local checks before the code enters github.

### Custom hooks and Api layer
All API calls live in `src/api/`, one file per resource. When an API call involves stateful logic, like subscribing to the SSE stream or managing pagination, I wrap it in a custom hook: `useNotification` for live notifications, `usePaginated` for list pages. Components stay thin and just call the hook.

### URL state and routing

The manager and technician dashboards are separate routes (`/dashboard/manager/[id]` and `/dashboard/technician/[id]`). This means refreshing the page keeps you in the right view and doesn't drop you back to a default screen. Identity is part of the URL, so the sidebar can derive it from the path without any extra context, including which recipient to subscribe to for the SSE notification stream.

Pagination is currently in local component state rather than the URL. It works fine for this use case but means the page number resets on refresh. Given more time I'd push it into the URL as a query param so the user can refresh or share a link to a specific page without losing their place.

### Offset pagination over cursor pagination

List endpoints accept `page` and `page_size` as API query params for offset pagination. For this small dataset of quotes and technicians that a manager browses manually, the dataset is too small for the offset scan cost to matter. You can jump to any page, the UI is straightforward, and the data doesn't change fast enough for the classic cursor pagination problems to matter.

**Trade-off:** Offset pagination degrades at scale. If rows are inserted or deleted between page requests, items can shift and cause duplicates or gaps across pages. Cursor (keyset) pagination avoids this by anchoring on a stable column like `created_at` + `id`, but it means you can only go forward, and the UI can't show "page 3 of 12".

### Swagger
Swagger UI is auto-generated from Go annotations on each handler. It gives a live, interactive view of every endpoint. It helps me to inspect request/response shapes and fire test requests directly from the browser.

### No authentication

Auth was out of scope, so identity is passed as query parameters (`recipient_type` + `recipient_id`) to keep things simple and focused on the scheduling logic. This is one of the first things I would add with more time.

If I were to add it, both `Manager` and `Technician` would embed a shared `User` struct. Login would issue a JWT in an HTTP-only cookie, a `/me` endpoint would identify the caller from the token, and middleware at the router level would handle the rest, removing identity in query params.

**Trade-off:** Any client can impersonate anyone, so this is not production-safe.

### SSE over WebSocket

SSE is simpler because it is just HTTP with chunked streaming, no library needed, and browsers reconnect automatically. For a notification feed that only flows one way, WebSockets would be overkill.

**Trade-off:** Some reverse proxies buffer the stream. The `X-Accel-Buffering: no` header is already set to handle that.

### In-process pub/sub hub (SSE delivery)

Rather than polling the DB every 5 seconds, the SSE handler subscribes to an in-process Go channel hub (`br-api/internal/hub/hub.go`). When a job is assigned or completed, the service publishes directly to the hub and connected clients get the notification instantly.

**Trade-off:** The hub lives in memory, so it only works on a single instance. Scaling horizontally would require an external broker like Redis pub/sub.

### Pessimistic locking for conflict detection

Both the quote and any conflicting jobs are locked at the start of the transaction with `SELECT ... FOR UPDATE`. The quote is locked first. Without that, two concurrent requests for the same quote could both read it as `unscheduled`, pass the check, and attempt to update the same quote. The jobs lock handles the technician double-booking case.

**Trade-off:** Simple and correct, but concurrent bookings for the same technician are serialised, which could become a performance bottleneck under heavy load. Pessimistic locking also makes deadlocks a real risk that needs to be debugged via SQL. An alternative would be optimistic locking with a retry loop, which may scale better but adds more complexity.

### GORM over raw SQL

GORM handles migrations, seeding, and standard queries without writing repetitive SQL. For the scope of this project it's a good fit, the schema is simple and the query patterns are straightforward.

**Trade-off:** GORM hides what SQL is actually being generated, which makes debugging slow queries harder. For the conflict check and pessimistic locking, I used GORM's lower-level APIs (`clause.Locking`, explicit WHERE strings) to keep the generated SQL predictable and easy to reason about.

### testcontainers-go for integration tests

Integration tests set up a real MySQL container per test using `testcontainers-go`. Tests run against the actual DB engine rather than a mock, so things like transactions, locking behaviour, and constraint violations are tested as they would be in production.

**Trade-off:** Requires Docker and is slower than unit tests. For that reason, integration tests are not run in the Husky pre-commit script — unit tests run on every commit, integration tests are run manually. In the future they could be wired into a CI pipeline.

## Use of AI tools

Claude was used throughout this project as a pair programming assistant. I configured a claude.md and a custom skill upfront so Claude could stay aligned with my engineering preferences, including code review style, naming conventions, commit format, and how I structure technical decisions. This reduced the need to repeatedly explain project context and made the feedback feel closer to working with a teammate.

I also used Claude Design early in the project to explore the typography scale and colour/theme direction before writing the Tailwind config. Having a visual reference helped me settle on a simple four-size type scale and the brass accent palette more quickly.

The final architecture and implementation decisions were still my own. Claude was mainly used to compare trade-offs, review code, sanity-check the locking strategy, debug issues such as the deadlock around quote scheduling, and help refine the README.
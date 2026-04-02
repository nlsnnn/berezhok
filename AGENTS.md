# AGENTS.md — Berezhok

Guide for AI coding agents working in this repository.

## Project Overview

Berezhok is a full-stack web application (Go backend + React frontend) — a platform for surprise box pickup from local partners (restaurants, cafes).

**Go 1.25** backend with chi router, PostgreSQL (pgx), Redis, S3/Yandex Object Storage, JWT auth.
**React 19** frontend with Vite, Tailwind CSS, MobX 6, Axios.

---

## Build / Lint / Test Commands

### Backend (Go)

```bash
# Run the server
make run                      # go run ./cmd/api

# Lint (uses golangci-lint)
task lint                      # runs ./bin/golangci-lint run ./... --config=.golangci.yml

# Format (gofumpt + gci import sorting)
task format                   # gofumpt -extra + gci (skips mocks/)

# Run all tests
task tests                    # go test ./internal/... -v

# Run a single test
go test ./internal/modules/catalog/service/ -run TestCreateBox -v

# Run E2E tests (require DB)
task e2e_tests                # go test ./internal/tests/... -v -run TestAPI

# Pre-commit (runs all linters + formatters)
make pre-commit               # pre-commit run --all-files

# SQL code generation (sqlc)
make sql-gen                  # sqlc generate

# Migrations
make migrate-create name=foo  # create migration pair
make migrate-up               # apply migrations
make migrate-down             # rollback 1 migration
```

### Frontend (React)

```bash
cd frontend

npm run dev                   # vite dev server
npm run build                 # production build
npm run lint                  # eslint .
```

### Infrastructure

```bash
docker compose up -d          # start Postgres (PostGIS) + Redis
docker compose down            # stop services
```

---

## Project Structure (Go)

```
cmd/api/                      # main entrypoint (main.go, api.go)
internal/
  adapters/                   # external integrations
    postgresql/               # pgx connection + sqlc generated code
    redis/                    # Redis client
    s3/                       # Yandex Object Storage
    sms/                      # SMS sender
    yookassa/                 # payment
  lib/                        # shared internal libs
    logger/                   # slog helpers
    pgconverter/              # pgx type converters
    validator/                # request validation wrapper
  modules/                    # feature modules (DDD-style)
    <module>/
      domain/                 # domain types + business entities
      repository/             # data access (implements sqlc queries)
      service/                # business logic
      handlers/               # HTTP handlers (chi)
        dto/                  # request/response DTOs + converters
      errors/                 # module-specific sentinel errors
  shared/                     # cross-cutting concerns
    auth/                     # bcrypt helpers
    config/                   # config loading (cleanenv)
    domain/                   # shared value objects (phone, geo, pickup_time)
    errors/                   # shared sentinel errors
    generator/                # ID generators
    jwt/                      # JWT token service
    middleware/                # auth middleware
    response/                 # HTTP JSON response helpers
migrations/                   # SQL migration files (golang-migrate)
```

---

## Project Structure (React)

```
frontend/src/
  api/                        # axios API clients (client.js, partner.js, etc.)
  components/
    ui/                       # reusable UI primitives (Button, Input, Modal, Spinner...)
  context/                    # React context providers (AuthContext)
  hooks/                      # custom hooks (useAddressSearch)
  lib/                        # utils (cn, formatDate, getErrorMessage), constants
  pages/
    landing/                  # public landing page
    partner/                  # partner dashboard pages
    admin/                    # admin pages
```

---

## Code Style — Go

### Imports
- Sorted with `gci`: standard → third-party → project (`github.com/nlsnnn/berezhok/...`)
- Use import aliases for disambiguation: `catalogErrors`, `partnerRepos`, `redisAdapter`, `sharedDomain`
- No dot imports

### Naming
- Exported: `PascalCase` — unexported: `camelCase`
- Constructor: `NewXxx(...)` returns `*xxx` (pointer to unexported struct)
- Interfaces defined at consumer side, one-method preferred (`boxService` uses `BoxRepository`, `locationFinder`)
- Constants use typed enums: `type BoxStatus string` with `BoxStatusActive BoxStatus = "active"`

### Error Handling
- Sentinel errors per module in `errors/errors.go`: `ErrBoxNotFound`, `ErrInvalidBoxID`
- Shared errors in `internal/shared/errors/errors.go`
- Wrap with `fmt.Errorf("context: %w", err)` for internal errors
- Map pgx errors: `errors.Is(err, pgx.ErrNoRows) → module.ErrNotFound`
- HTTP handlers use `switch/errors.Is` to map to appropriate response codes
- Never expose internal errors to client — use `response.InternalError(w, nil)`

### HTTP Handlers
- Pattern: `const op = "module.handler.Method"` → `log := h.log.With(slog.String("op", op))`
- Validate with `h.validator.DecodeAndValidate(r, &req)` → returns `map[string]any` or nil
- Use `response.Success(w, data)`, `response.Created(w, data)`, `response.NotFound(w, msg)` helpers
- Extract path params with `chi.URLParam(r, "id")`
- Extract auth context with `contextx.PartnerID(r)` / `contextx.CustomerID(r)` / `contextx.EmployeeID(r)` (returns `uuid.UUID, error`)

### Logging
- Use `slog` (stdlib) exclusively
- Helper: `sl.Err(err)` for error attrs, `sl.Errs(errs)` for validation maps
- Levels: `Info` for normal flow, `Warn` for client errors, `Error` for server errors

### DB / SQL
- SQL queries managed by sqlc — never write raw SQL in Go code
- Queries in `internal/adapters/postgresql/queries/*.sql`
- Generated code in `internal/adapters/postgresql/sqlc/`
- Use `pgconverter` package for pgx type conversions

---

## Code Style — React (JSX)

- **No TypeScript** — plain JS with JSX
- Functional components, `export default function ComponentName`
- State management: **MobX 6** stores (NOT React Query)
  - Stores are classes with `makeAutoObservable`, async actions use `runInAction`
  - Exported as singletons: `export const boxesStore = new BoxesStore()`
  - Access via `useStores()` hook or `StoresContext`
  - Pages wrapped with `observer()` from `mobx-react-lite`
  - Pattern: `function PageBase() { ... }` then `export default observer(PageBase)`
- Styling: Tailwind CSS + `cn()` utility (clsx + tailwind-merge)
  - Reusable CSS classes: `btn-primary`, `btn-secondary`, `btn-danger`, `btn-ghost`, `input-base`, `badge`, `card`
  - Custom colors: `brand` (green), `cream` (warm neutral)
- Path alias: `@/` → `./src/`
- API calls through centralized axios client (`api/client.js`) with JWT interceptor
  - Response unwrapping: `.then((r) => r.data.data)` — strips `{success, data}` envelope
- Error display via `sonner` toast: `toast.success('...')`, `toast.error('...')`
- Component props destructured inline, spread `...props` for passthrough
- Unused vars error disabled for uppercase/underscore patterns
- UI text is in Russian

---

## Key Conventions

- Domain types are plain Go structs — no ORM tags, no JSON tags on domain layer
- DTOs in `handlers/dto/` carry JSON/validation tags and converter methods (`ToInput()`, `ToResponse()`)
  - Three-file pattern per module: `request.go`, `response.go`, `converter.go`
- Services accept input structs, return domain types
- Repositories accept/return domain types, internally map to/from sqlc generated types
- JWT tokens: partner (email+password), customer (phone+SMS code), admin — role in middleware
- Config via `.env` + `cleanenv` → `config/local.yaml`
- Response envelope: `{ "success": bool, "data": any }` or `{ "success": false, "error": { code, message, details } }`
- Prefer `contextx.PartnerID(r)` / `contextx.CustomerID(r)` / `contextx.EmployeeID(r)` over raw `r.Context().Value(...)`

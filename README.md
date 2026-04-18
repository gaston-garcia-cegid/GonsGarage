# GonsGarage

Auto repair shop management system: **Go** API (Gin, GORM, PostgreSQL, Redis) and **Next.js** web app (App Router, React, Zustand).

[![CI](https://img.shields.io/badge/CI-GitHub_Actions-blue)](.github/workflows/ci.yml)
[![Go](https://img.shields.io/badge/Go-1.25-blue)](backend/go.mod)
[![Next.js](https://img.shields.io/badge/Next.js-15.5-black)](frontend/package.json)
[![React](https://img.shields.io/badge/React-19-61dafb)](frontend/package.json)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

## Documentation

Technical deep-dives (architecture, env vars, TDD, roadmap) live under **[`docs/`](docs/)** — start at [`docs/README.md`](docs/README.md).

| Topic | Link |
|--------|------|
| Development setup | [`docs/development-guide.md`](docs/development-guide.md) |
| TDD / testing | [`docs/testing-tdd.md`](docs/testing-tdd.md) |
| Contributing | [`CONTRIBUTING.md`](CONTRIBUTING.md) |
| Agent / code standards | [`Agent.md`](Agent.md) |

## Verified stack (source of truth)

Versions below are taken from the repo manifests as of the last README refresh. Upgrade the files first, then update this table.

| Layer | Technology | Where it is defined |
|--------|------------|---------------------|
| Backend runtime | **Go 1.25** (`go 1.25.3` directive) | [`backend/go.mod`](backend/go.mod) |
| HTTP | Gin, JWT, GORM, sqlx, Redis client | `backend/go.mod` |
| Database (local) | **PostgreSQL 16** (`postgres:16-alpine`) | [`docker-compose.yml`](docker-compose.yml) |
| Cache (local) | **Redis 7** (`redis:7-alpine`) | [`docker-compose.yml`](docker-compose.yml) |
| Frontend | **Next.js 15.5.5**, **React 19.1.0**, TypeScript **^5** | [`frontend/package.json`](frontend/package.json) |
| Package manager | **pnpm 9.15.4** (`packageManager` field) | [`frontend/package.json`](frontend/package.json) |
| Unit / component tests (default) | **Vitest** + Testing Library | `frontend/package.json` → `pnpm test` |
| Lint / types | ESLint 9, `eslint-config-next` aligned with Next | `frontend/package.json` |

CI runs **Node 22**, **pnpm 9**, `pnpm lint` → `pnpm typecheck` → `pnpm test` → `pnpm build` for the frontend, and **Go from `go.mod`** for `vet` / `test -race` / `build` on the backend (see [`.github/workflows/ci.yml`](.github/workflows/ci.yml)).

## Prerequisites

- **Node.js** 22+ and **pnpm** 9+ ([`corepack enable`](https://nodejs.org/api/corepack.html) recommended)
- **Go** 1.25+ (match `backend/go.mod`)
- **Docker** with Compose v2 (`docker compose`) for local PostgreSQL and Redis
- **Git**

## Quick start

### 1. Clone

```bash
git clone https://github.com/gaston-garcia-cegid/GonsGarage.git
cd GonsGarage
```

### 2. Environment files

```bash
cp backend/.env.example backend/.env
cp frontend/.env.local.example frontend/.env.local
```

### 3. Databases (recommended)

From the **repository root**:

```bash
docker compose up -d
```

Services: **postgres** (port `5432`) and **redis** (`6379`). Defaults match `backend/.env.example` and `backend/cmd/api/main.go` when `DATABASE_URL` / `REDIS_URL` are unset.

### 4. Backend API

```bash
cd backend
go mod download
go run ./cmd/api
```

- API: <http://localhost:8080> (or `SERVER_PORT`)
- Swagger UI: <http://localhost:8080/swagger/index.html>
- Health: <http://localhost:8080/health> · Readiness: <http://localhost:8080/ready>

Regenerate Swagger after changing `// @Summary` / `// @Router` annotations:

```bash
cd backend
go run github.com/swaggo/swag/cmd/swag@v1.8.12 init -g main.go -o docs -d ./cmd/api,./internal/handler,./internal/core/ports --parseInternal
```

### 5. Frontend

```bash
cd frontend
pnpm install
pnpm dev
```

App: <http://localhost:3000>. Set `NEXT_PUBLIC_API_URL` in `frontend/.env.local` if the API is not on `http://localhost:8080`.

## Testing (mirrors CI)

**Backend**

```bash
cd backend
go vet ./...
go test ./... -count=1 -race
```

**Frontend**

```bash
cd frontend
pnpm lint
pnpm typecheck
pnpm test
pnpm build
```

The default **`pnpm test`** script runs **Vitest**. Jest remains available for legacy scripts (`pnpm test:jest`) but is not the primary runner documented in CI.

## Demo users

| Role | Email | Password | Notes |
|------|-------|----------|--------|
| Admin | `admin@gonsgarage.com` | `admin123` | Seeded / default for local dev (see codebase and `docs/`) |

**Demo client (optional seed)** from `backend/` with PostgreSQL up and schema migrated (e.g. after the API has run once):

```bash
cd backend
go run ./cmd/seed-test-client
```

Default client: `cliente.demo@gonsgarage.local` / `ClienteDemo123`. Override with `SEED_CLIENT_EMAIL` and `SEED_CLIENT_PASSWORD`. Command is idempotent if the email already exists.

## Project layout

```text
backend/          Go API (cmd/api, internal/, docs/swagger)
frontend/         Next.js App Router (src/app, components, stores)
docs/             Technical documentation index
openspec/         Spec-driven development (OpenSpec)
docker-compose.yml   Local Postgres + Redis
```

## Contributing

See [`CONTRIBUTING.md`](CONTRIBUTING.md) and [`docs/testing-tdd.md`](docs/testing-tdd.md) for branch workflow and TDD expectations.

## License

[MIT](LICENSE).

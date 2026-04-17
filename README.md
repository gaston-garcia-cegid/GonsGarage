# GonsGarage

**Dominio (`{{BUSINESS_DOMAIN}}`):** auto repair shop management system — *valor fijo en [`template_project.md`](template_project.md); en forks nuevos sustituir el placeholder.*

**Locales UI/docs (`{{LOCALE}}`):** `pt_PT`, `es_ES`, `en_GB` (prioridad según producto; ver tabla en `template_project.md`).

Monorepo alineado con la plantilla **§1–§4** en [`template_project.md`](template_project.md) y reglas compactas en [`gonsgarage-rules/`](gonsgarage-rules/). Plan de adopción por fases: [`docs/template-adoption-plan.md`](docs/template-adoption-plan.md).

## Arquitectura (resumen)

| Capa | Tecnología |
|------|------------|
| **Frontend** | Next.js 16 (App Router), React 19, TypeScript 5.9+, Tailwind CSS v4, Zustand |
| **Backend** | Go 1.25, Gin, Clean Architecture (`internal/...`), OpenAPI/Swagger |
| **Datos** | PostgreSQL 16, Redis 7 |
| **API** | REST JSON bajo `/api/v1` (camelCase) |

Índice de documentación: [`docs/DOCUMENTATION_INDEX.md`](docs/DOCUMENTATION_INDEX.md).

## Requisitos

- **Go** 1.25+ ([`backend/go.mod`](backend/go.mod))
- **Node** 22+ y **pnpm** 9+
- **Docker** (opcional, recomendado para Postgres + Redis)

## Desarrollo local

### Base de datos y Redis

Desde la raíz del repositorio:

```bash
docker compose up -d
```

Servicios por defecto: Postgres (`localhost:5432`), Redis (`localhost:6379`). Ver [`docker-compose.yml`](docker-compose.yml).

### Backend

```bash
cd backend
go mod download
go vet ./...
go test ./... -race -count=1
go run ./cmd/api
```

Variables: copiar ejemplos de entorno según [`backend/`](backend/) y [`docs/development-guide.md`](docs/development-guide.md).

### Frontend

```bash
cd frontend
pnpm install --frozen-lockfile
pnpm lint
pnpm typecheck
pnpm test
pnpm dev
```

Tests con **Vitest** (plantilla §3). Variables: [`frontend/.env.local.example`](frontend/.env.local.example).

## CI

[`.github/workflows/ci.yml`](.github/workflows/ci.yml):

- **backend:** `go vet` → `go test -race` → `go build` (`cmd/api`, `cmd/worker`)
- **frontend:** `pnpm lint` → `tsc --noEmit` → `vitest run` → `pnpm build`

## Gobernanza

- Ramas, commits, PRs, API y seguridad: [`gonsgarage-rules/04-governance.md`](gonsgarage-rules/04-governance.md)
- Contribución: [`CONTRIBUTING.md`](CONTRIBUTING.md)
- Changelog: [`CHANGELOG.md`](CHANGELOG.md)

## Desviaciones documentadas respecto a la plantilla

| Tema | Nota |
|------|------|
| **Redis Go** | Plantilla cita `go-redis/v8`; el repo usa **v9** (ver `gonsgarage-rules/02-stack.md`). |
| **ORM** | Plantilla cita **sqlx**; migración desde GORM: **Fase 2** del plan de adopción. |
| **Logging** | Plantilla cita **Zerolog**; migración desde slog: **Fase 4**. |

No se añadieron stacks opcionales (p. ej. Sentry, xlsx) sin decisión explícita en la plantilla §5.

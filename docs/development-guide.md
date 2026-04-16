# Guía de desarrollo GonsGarage

## Requisitos

- Go 1.21+ (según README del proyecto)
- Node.js 22+ con **pnpm** 9+ (`corepack enable` recomendado)
- PostgreSQL 16+ (o Docker)
- Redis opcional (recomendado para probar cache real)

## Variables de entorno (backend)

Definidas o usadas en `backend/cmd/api/main.go`:

| Variable | Uso | Valor por defecto (si vacío) |
|----------|-----|-------------------------------|
| `DATABASE_URL` | DSN PostgreSQL | `postgres://admindb:gonsgarage123@localhost:5432/gonsgarage?sslmode=disable` |
| `REDIS_URL` | Host Redis | `localhost:6379` |
| `JWT_SECRET` | Firma JWT | Clave de desarrollo (log de advertencia) |
| `SERVER_PORT` | Puerto HTTP | `8080` |
| `GIN_MODE` | `release` desactiva modo debug Gin | — |
| `RESET_DATABASE` | `true` elimina tablas antes de migrar (solo desarrollo) | — |

## Levantar PostgreSQL y Redis (Docker)

Desde la **raíz del repositorio** (archivo `docker-compose.yml`):

```powershell
docker compose up -d
```

Credenciales y nombres coinciden con `backend/.env.example` y con los valores por defecto en `backend/cmd/api/main.go` si no defines `DATABASE_URL`.

## Backend

```powershell
Set-Location backend
Copy-Item .env.example .env   # primera vez
go mod download
go run ./cmd/api
```

- API: `http://localhost:8080` (o `SERVER_PORT`)
- Swagger: `http://localhost:8080/swagger/index.html`
- Health: `http://localhost:8080/health`
- Perfil autenticado: `GET http://localhost:8080/api/v1/auth/me` con cabecera `Authorization: Bearer <token>`
- Empleados: `GET/POST /api/v1/employees/...` solo para roles **admin** o **manager**

## Frontend

```powershell
Set-Location frontend
Copy-Item .env.local.example .env.local   # primera vez
pnpm install
pnpm dev
```

`NEXT_PUBLIC_API_URL` por defecto apunta a `http://localhost:8080` (coincide con el backend en el mismo host). El cliente en `src/lib/api` añade `/api/v1` donde corresponde.

## Tests (backend)

```powershell
Set-Location backend
go test ./...
```

## Documentación de código y OpenAPI

Si el proyecto regenera Swagger con `swag`, mantener anotaciones en handlers y el `main` según `Agent.md` y el README raíz.

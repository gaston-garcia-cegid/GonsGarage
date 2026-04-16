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
| `CORS_ALLOWED_ORIGINS` | Con `GIN_MODE=release`, orígenes permitidos del navegador (coma-separada, URL exacta) | Vacío: ver advertencia al arrancar; navegador cross-origin fallará salvo mismo sitio |
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

### Coches y citas (roles)

- **GET /api/v1/cars** — si el JWT es **client**, solo devuelve sus coches. Si el rol es **employee**, **manager** o **admin**, admite `?ownerId=<uuid>` (coches de ese cliente) o, sin `ownerId`, lista paginada del taller (`limit` por defecto 50, `offset` 0).
- **POST /api/v1/cars** — el cliente no envía `ownerID` (el dueño es siempre él). El taller puede enviar **`ownerID`** (UUID del usuario cliente) para registrar un coche a su nombre.
- **GET /api/v1/appointments** — el cliente solo ve sus citas (el backend fuerza su `customer_id`). El taller puede usar `customerId`, `carId`, `status`, `limit`, `offset`, `sortBy`, `sortOrder`.
- **POST /api/v1/appointments** — el cliente reserva para sí; el taller puede enviar **`customerID`** (camelCase) para crear la cita en nombre del cliente. El coche debe pertenecer a ese cliente.

## Frontend

```powershell
Set-Location frontend
Copy-Item .env.local.example .env.local   # primera vez
pnpm install
pnpm dev
```

`NEXT_PUBLIC_API_URL` por defecto apunta a `http://localhost:8080` (coincide con el backend en el mismo host). El cliente en `src/lib/api` añade `/api/v1` donde corresponde.

## Staging / Docker (Fase D)

- Guía: **`docs/deployment-staging.md`** (secretos, HTTPS delante del proxy, backup `pg_dump`).
- **Producción mínima:** `docker-compose.prod.yml` + variables desde `deploy/.env.production.example`.
- **Smoke (CI o local):** `docker compose -f docker-compose.smoke.yml up --build` y comprobar `http://localhost:8080/health`.

## Tests (backend)

```powershell
Set-Location backend
go test ./...
```

## Documentación de código y OpenAPI (swag)

Los artefactos generados viven en **`backend/docs/`** (`docs.go`, `swagger.json`, `swagger.yaml`). El servidor expone la UI en **`http://localhost:8080/swagger/index.html`**.

Regenerar tras cambiar anotaciones `// @Summary`, `// @Router`, etc.:

```powershell
Set-Location backend
go run github.com/swaggo/swag/cmd/swag@v1.8.12 init -g main.go -o docs -d ./cmd/api,./internal/adapters/http/handlers,./internal/core/ports --parseInternal
```

- El **general API** (título, `BasePath`, seguridad `BearerAuth`) está en `cmd/api/main.go`.
- Las rutas documentadas están en `internal/adapters/http/handlers` y el ancla de **`/health`** en `cmd/api/swagger_health.go`.
- Tipos de petición compartidos (`ports.RegisterRequest`, empleados, etc.) requieren incluir **`internal/core/ports`** en `-d`.

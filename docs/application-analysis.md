# Análisis de la aplicación GonsGarage

## Propósito

Sistema de gestión para taller mecánico: usuarios con roles, vehículos, citas, empleados y dominio de reparaciones (modelo presente; API de reparaciones pendiente de cablear).

## Stack real (código actual)

| Capa | Tecnología |
|------|------------|
| Backend | Go, Gin, GORM (PostgreSQL), JWT, Redis opcional (cache con fallback si no hay conexión) |
| Frontend | Next.js (App Router), TypeScript, Zustand (según `Agent.md` y estructura `src/`) |
| Datos | PostgreSQL; migraciones vía `AutoMigrate` en `backend/cmd/api/main.go` y scripts SQL en `backend/scripts/` |

## Punto de entrada del backend

- **Ejecutable principal**: `backend/cmd/api/main.go` (no `cmd/server` en el árbol actual).
- **Swagger**: ruta `GET /swagger/*any` (swaggo).
- **Salud**: `GET /health`.
- **Prefijo API**: `/api/v1`.

## Rutas API expuestas (resumen)

**Públicas**

- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`

**Protegidas** (middleware JWT en contexto Gin)

- Sesión: `GET /api/v1/auth/me` (usuario actual, camelCase)
- Empleados (**só admin/manager**): `POST|GET|GET/:id|PUT|DELETE /api/v1/employees/...`
- Coches: `POST|GET|GET/:id|PUT|DELETE /api/v1/cars/...`
- Citas: `POST|GET|GET/:id|PUT|DELETE /api/v1/appointments/...`

**Reparaciones**: el dominio y migraciones incluyen `Repair`; las rutas REST en `main.go` están comentadas (pendiente `repairHandler`).

## Frontend (App Router)

Rutas de página localizadas bajo `frontend/src/app/`:

- `/`, `/dashboard`, `/client`, `/employees`
- `/auth/login`, `/auth/register`
- `/cars`, `/cars/[id]`
- `/appointments`, `/appointments/new`

La integración con la API debe alinearse con el prefijo `/api/v1` y JSON en camelCase (`Agent.md`).

## Infraestructura y documentación

1. **Docker Compose** en la raíz del repo (`docker-compose.yml`): PostgreSQL (`gonsgarage` / `admindb`) y Redis, alineado con los defaults del backend y `backend/.env.example`.
2. **Entrada del servidor**: `backend/cmd/api/main.go` (AutoMigrate al arranque; no hay `cmd/migrate` obligatorio).
3. **README del frontend** sigue siendo en parte la plantilla de `create-next-app`; ver `.env.local.example` y [development-guide.md](./development-guide.md).

Guía operativa: [development-guide.md](./development-guide.md). Planificación: [roadmap.md](./roadmap.md).

## Estructura de código (alto nivel)

```
backend/
  cmd/api/main.go          # Arranque, migraciones GORM, rutas Gin
  internal/core/domain     # Entidades
  internal/core/services   # Casos de uso / servicios
  internal/handler         # Handlers HTTP (Gin)
  internal/middleware      # Auth, CORS, rate limit, etc.
  internal/adapters/repository/postgres|redis

frontend/
  src/app/                 # Rutas UI
  src/stores/              # Estado (Zustand)
  src/lib/api/             # Cliente HTTP
```

## Referencias cruzadas

- Detalle de convenciones: [../Agent.md](../Agent.md).
- Cliente API frontend: [../frontend/docs/api-client.md](../frontend/docs/api-client.md).

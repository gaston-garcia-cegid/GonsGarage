# Análisis de la aplicación GonsGarage

## Propósito

Sistema de gestión para taller mecánico: usuarios con roles, vehículos, citas, empleados y dominio de reparaciones (lectura por coche + **escritura staff** en API y UI mínima en ficha de coche).

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
- **CORS**: con `GIN_MODE=release`, orígenes permitidos = localhost:3000/3001 más los de la variable **`CORS_ORIGINS`** (coma-separado). Ver `corsMiddleware` en `backend/cmd/api/main.go`. Despliegue Docker/LAN: [`deploy/README.md`](./deploy/README.md).

## Rutas API expuestas (resumen)

**Públicas**

- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`

**Protegidas** (middleware JWT en contexto Gin)

- Sesión: `GET /api/v1/auth/me` (usuario actual, camelCase)
- Empleados (**só admin/manager**): `POST|GET|GET/:id|PUT|DELETE /api/v1/employees/...`
- Coches: `POST|GET|GET/:id|PUT|DELETE /api/v1/cars/...` — listado de flota de un cliente (staff): `GET /api/v1/cars?ownerId=<uuid>&limit=&offset=` (no existe `GET /cars/owner/:id`).
- Citas: `POST|GET|GET/:id|PUT|DELETE /api/v1/appointments/...` — cambiar estado (cancelar / confirmar / completar): **`PUT /api/v1/appointments/:id`** con JSON parcial (`status`, etc.); **no** hay `PATCH …/cancel|/confirm|/complete`.

**Reparaciones** (JWT): `GET /api/v1/repairs/car/:carId` (cliente: solo su coche; staff: cualquier coche); **`POST /api/v1/repairs`**, `GET|PUT|DELETE /api/v1/repairs/:id` — escritura solo personal (`RepairService` exige `IsEmployee()` para crear/actualizar/borrar). UI staff: formulario y acciones en `/cars/[id]` cuando el usuario no es `client`.

**Regresión por rol (tests):** `GET /api/v1/employees` solo `admin`/`manager` (middleware `RequireStaffManagers`); `POST /api/v1/repairs` devuelve **403** para JWT `client` y **201** para `employee` en condiciones válidas — cubierto en `backend/internal/handler/mvp_role_access_test.go`; spec [`openspec/specs/mvp-role-access/spec.md`](../openspec/specs/mvp-role-access/spec.md).

## Frontend (App Router)

Rutas de página localizadas bajo `frontend/src/app/`:

- `/`, `/dashboard`, `/client`, `/employees`
- `/auth/login`, `/auth/register`
- `/cars`, `/cars/[id]`
- `/appointments`, `/appointments/new`

La integración con la API debe alinearse con el prefijo `/api/v1` y JSON en camelCase (`Agent.md`). El cliente HTTP principal está en `frontend/src/lib/api-client.ts`; citas en `lib/api/appointment.api.ts`; coches en `lib/services/car.service.ts`; el módulo legacy `lib/api.ts` (usado por dashboard/coches) debe limitarse a rutas realmente expuestas (p. ej. repairs solo `GET /repairs/car/:carId`).

## Infraestructura y documentación

1. **Docker Compose** en la raíz (`docker-compose.yml`): PostgreSQL y Redis para desarrollo local. **`docker-compose.prod.yml`**: stack prod (API + Redis + Next + nginx en 8102) con Postgres en el host vía `DATABASE_URL`; ver [`deploy/README.md`](./deploy/README.md).
2. **Entrada del servidor**: `backend/cmd/api/main.go` (AutoMigrate al arranque; no hay `cmd/migrate` obligatorio).
3. **README del frontend** sigue siendo en parte la plantilla de `create-next-app`; ver `.env.local.example` y [development-guide.md](./development-guide.md).

Guía operativa: [development-guide.md](./development-guide.md). Planificación: [roadmap.md](./roadmap.md).

## Estructura de código (alto nivel)

```
backend/
  cmd/api/main.go          # Arranque, migraciones GORM, rutas Gin
  internal/domain          # Entidades
  internal/service         # Casos de uso (subpaquetes auth, car, …)
  internal/handler         # Handlers HTTP (Gin)
  internal/middleware      # Auth, CORS, etc.
  internal/repository/postgres|redis|mock
  internal/core/ports      # Contratos (interfaces)

frontend/
  src/app/                 # Rutas UI
  src/stores/              # Estado (Zustand)
  src/lib/api/             # Cliente HTTP
```

## Referencias cruzadas

- Detalle de convenciones: [../Agent.md](../Agent.md).
- Cliente API frontend: [../frontend/docs/api-client.md](../frontend/docs/api-client.md).

# Fases mínimas para un primer MVP (GonsGarage)

Objetivo del MVP: **un taller puede registrar usuarios, clientes con coches, pedir citas y el personal puede ver/gestionar lo básico** con despliegue reproducible y calidad mínima aceptable.

## Fase A — Base técnica (cerrada)

- [x] Compose local (Postgres + Redis) y variables documentadas (`docker-compose.yml`, `backend/.env.example`, `frontend/.env.local.example`, `docs/development-guide.md`).
- [x] CI en GitHub (backend + frontend) en `.github/workflows/ci.yml`.
- [x] Política **TDD** en `docs/testing-tdd.md` y tests base: backend `auth`, `car`, `appointment` (`*_service_test.go`); frontend stores/API (`__tests__/auth.*`, `car.store`, `appointment.*`).

**Criterio de salida:** `main` con CI verde; README y `docs/development-guide.md` reflejan **pnpm** y **`go run ./cmd/api`**.

## Fase B — Identidad y datos coherentes (cerrada)

- Roles claros (`client` / `employee` / `admin` / `manager`) y comprobaciones en API alineadas con la UI.
- Flujo **register + login + JWT** estable end-to-end.
- **Coches y citas** CRUD coherentes con el contrato JSON camelCase; **OpenAPI generado con swag** (`backend/docs/`, UI `/swagger/index.html`).

**Avances en código (iteración actual):**

- [x] `GET /api/v1/auth/me` (JWT) — perfil del usuario autenticado; JSON **camelCase** en registro y en `me`.
- [x] Registro público: rol por defecto **`client`** si no se envía `role`; conflicto con `409` usando `ErrUserAlreadyExists`.
- [x] Rutas **`/employees/*`** restringidas a **admin** y **manager** (`RequireStaffManagers`).
- [x] **Coches:** cliente solo su flota; personal (`employee` / `manager` / `admin`) puede **listar inventario** (`GET /cars?limit=&offset=&ownerId=`) y **crear coche para un cliente** (`ownerID` en el body si no es `client`). Lectura de cualquier coche para staff en servicio de dominio.
- [x] **Citas:** listado acotado a **customer_id** del cliente; staff ve/filtra con query params; creación valida que el **coche pertenezca al cliente** de la cita; actualización conserva `customer_id` y admite `scheduledAt` opcional (RFC3339).
- [x] **Frontend:** `UserRole.manager`, `canManageUsers` alineado con backend; **login** y **checkAuthStatus** validan sesión con **`/auth/me`**.

**Criterio de salida:** demo en local: cliente crea coche y cita; empleado/admin lista y actualiza estados básicos.

## Fase C — Reparaciones (diferenciador de taller) (cerrada)

- Exponer **repairs** en REST (hoy el dominio existe pero las rutas en `cmd/api` están comentadas) o recortar el alcance del MVP y documentar “sin repairs en API hasta vX”.
- UI mínima: listado / detalle de reparación por coche (aunque sea solo lectura en MVP).

**Criterio de salida:** historia de reparación visible para el cliente o decisión explícita de aplazar con issue enlazado.

**Avances en código (iteración actual):**

- [x] **API Gin:** `POST /api/v1/repairs` (personal del taller), `GET /api/v1/repairs/:id`, `PUT /api/v1/repairs/:id`, `GET /api/v1/cars/:id/repairs` (ruta anidada registrada antes de `GET /cars/:id`); JSON **camelCase**; OpenAPI regenerado con **swag**.
- [x] **Persistencia:** `PostgresRepairRepository` alineado con `domain.Repair` y soft delete; permisos de creación para coches de clientes (técnico no tiene que ser dueño).
- [x] **Dominio / tests:** `repair_service_test.go` (stubs); `go test ./...` en backend.
- [x] **Frontend:** detalle de coche consume `GET .../cars/:id/repairs` y tipos camelCase en `api.ts` / vista de historial.

## Fase D — MVP en producción (mínimo viable “de verdad”) (cerrada en repo)

- **Secrets** fuera del repo; `JWT_SECRET` y DB en variables del entorno.
- Imagen Docker o plataforma elegida (un solo `docker compose` de prod o PaaS).
- **Deploy workflow** sustituyendo el placeholder `.github/workflows/deploy.yml` (migraciones, healthcheck, rollback básico).

**Criterio de salida:** URL pública o entorno staging con HTTPS y backup de BD documentado.

**Avances en código (iteración actual):**

- [x] **Docker:** `backend/Dockerfile`, `frontend/Dockerfile` (Next **standalone**), `docker-compose.prod.yml`, `docker-compose.smoke.yml`, `deploy/.env.production.example`.
- [x] **Secretos:** `JWT_SECRET` obligatorio distinto del default cuando `GIN_MODE=release` en `cmd/api/main.go`.
- [x] **Deploy:** `.github/workflows/deploy.yml` — smoke con Compose + `/health`; job opcional **push a GHCR** (`push_ghcr`).
- [x] **Documentación:** `docs/deployment-staging.md` (HTTPS vía proxy, backup, rollback básico) y enlace desde `docs/development-guide.md`.
- [x] **Deuda:** dependencia `gorilla/mux` eliminada (`user_handler` migrado a Gin). Plantilla de issue en **`docs/github/issue-remove-gorilla-mux.md`** (crear en GitHub si se desea ticket formal).

**Nota sobre el criterio “URL pública HTTPS”:** el repositorio entrega compose + reverse proxy documentado; la URL y certificados dependen del entorno del operador (DNS + Caddy/Traefik/nginx).

## Orden recomendado

`A → B → (C o decisión explícita) → D`. No hace falta “Arnela-paridad” (Tailwind v4, colas pesadas, etc.) para cerrar el MVP; eso entra como deuda o fase posterior.

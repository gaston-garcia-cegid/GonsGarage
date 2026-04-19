# Fases mínimas para un primer MVP (GonsGarage)

Objetivo del MVP: **un taller puede registrar usuarios, clientes con coches, pedir citas y el personal puede ver/gestionar lo básico** con despliegue reproducible y calidad mínima aceptable.

**Plan consolidado hasta “MVP funcionando” (staging/prod, contrato, demo):** [openspec/changes/mvp-funcionando-plan/proposal.md](../openspec/changes/mvp-funcionando-plan/proposal.md) (aprobado; usar como épica para issues).

**Brecha MVP vs Entra 1.1 + roadmap P0–P3 (tareas SDD):** [openspec/changes/mvp-gap-roadmap-2026/proposal.md](../openspec/changes/mvp-gap-roadmap-2026/proposal.md) · [`tasks.md`](../openspec/changes/mvp-gap-roadmap-2026/tasks.md).

**Checklist operativo (equipo de una persona, servidor de pruebas = staging):** [mvp-solo-checklist.md](./mvp-solo-checklist.md).

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

## Fase C — Reparaciones (diferenciador de taller)

- [x] **API:** `GET /api/v1/repairs/car/:carId` activo en `cmd/api` (lista por coche; permisos en servicio).
- [x] **UI cliente:** detalle de coche carga historial; **dashboard** agrega reparaciones recientes de la flota del usuario (solo lectura).
- [ ] **Opcional MVP+:** `POST`/`PATCH`/`DELETE` de repairs en Gin + UI mínima para staff; o documentar aplazamiento con issue.

**Criterio de salida:** historia de reparación visible para el cliente (**cumplido** en lectura); ampliaciones staff = criterio opcional arriba.

## Fase D — MVP en producción (mínimo viable “de verdad”)

- **Secrets** fuera del repo; `JWT_SECRET` y DB en variables del entorno.
- Imagen Docker o plataforma elegida (un solo `docker compose` de prod o PaaS).
- **Deploy workflow** sustituyendo el placeholder `.github/workflows/deploy.yml` (migraciones, healthcheck, rollback básico).

**Criterio de salida:** URL pública o entorno staging con HTTPS y backup de BD documentado.

## Orden recomendado

`A → B → (C o decisión explícita) → D`. No hace falta “Arnela-paridad” (Tailwind v4, colas pesadas, etc.) para cerrar el MVP; eso entra como deuda o fase posterior.

# Fases mínimas para un primer MVP (GonsGarage)

Objetivo del MVP: **un taller puede registrar usuarios, clientes con coches, pedir citas y el personal puede ver/gestionar lo básico** con despliegue reproducible y calidad mínima aceptable.

## Fase A — Base técnica (cerrada)

- [x] Compose local (Postgres + Redis) y variables documentadas (`docker-compose.yml`, `backend/.env.example`, `frontend/.env.local.example`, `docs/development-guide.md`).
- [x] CI en GitHub (backend + frontend) en `.github/workflows/ci.yml`.
- [x] Política **TDD** en `docs/testing-tdd.md` y tests base: backend `auth`, `car`, `appointment` (`*_service_test.go`); frontend stores/API (`__tests__/auth.*`, `car.store`, `appointment.*`).

**Criterio de salida:** `main` con CI verde; README y `docs/development-guide.md` reflejan **pnpm** y **`go run ./cmd/api`**.

## Fase B — Identidad y datos coherentes

- Roles claros (`client` / `employee` / `admin` / `manager`) y comprobaciones en API alineadas con la UI.
- Flujo **register + login + JWT** estable end-to-end.
- **Coches y citas** CRUD coherentes con el contrato JSON camelCase y Swagger actualizado.

**Avances en código (iteración actual):**

- [x] `GET /api/v1/auth/me` (JWT) — perfil del usuario autenticado; JSON **camelCase** en registro y en `me`.
- [x] Registro público: rol por defecto **`client`** si no se envía `role`; conflicto con `409` usando `ErrUserAlreadyExists`.
- [x] Rutas **`/employees/*`** restringidas a **admin** y **manager** (`RequireStaffManagers`).
- [ ] Revisar el mismo criterio de roles en citas / coches si el negocio lo exige (hoy la lógica vive sobre todo en servicios de dominio).

**Criterio de salida:** demo en local: cliente crea coche y cita; empleado/admin lista y actualiza estados básicos.

## Fase C — Reparaciones (diferenciador de taller)

- Exponer **repairs** en REST (hoy el dominio existe pero las rutas en `cmd/api` están comentadas) o recortar el alcance del MVP y documentar “sin repairs en API hasta vX”.
- UI mínima: listado / detalle de reparación por coche (aunque sea solo lectura en MVP).

**Criterio de salida:** historia de reparación visible para el cliente o decisión explícita de aplazar con issue enlazado.

## Fase D — MVP en producción (mínimo viable “de verdad”)

- **Secrets** fuera del repo; `JWT_SECRET` y DB en variables del entorno.
- Imagen Docker o plataforma elegida (un solo `docker compose` de prod o PaaS).
- **Deploy workflow** sustituyendo el placeholder `.github/workflows/deploy.yml` (migraciones, healthcheck, rollback básico).

**Criterio de salida:** URL pública o entorno staging con HTTPS y backup de BD documentado.

## Orden recomendado

`A → B → (C o decisión explícita) → D`. No hace falta “Arnela-paridad” (Tailwind v4, colas pesadas, etc.) para cerrar el MVP; eso entra como deuda o fase posterior.

# Plan de adopción — `template_project.md` (secciones 1–4)

Este documento descompone la alineación del repositorio **GonsGarage** con la plantilla tecnológica (estilo Arnela). Los campos `{{LOCALE}}` y `{{BUSINESS_DOMAIN}}` del prompt genérico están resueltos en `template_project.md` (tabla); para nuevos forks, conservá esos placeholders hasta rellenar la tabla.

| Placeholder plantilla | Valor en este repo (tabla) |
|------------------------|----------------------------|
| `{{BUSINESS_DOMAIN}}` | auto repair shop management system |
| `{{LOCALE}}` | `pt_PT`, `es_ES`, `en_GB` (UI/docs principal) |

---

## Fase 0 — Hecho (baseline)

- `docker-compose.yml` en raíz: **PostgreSQL 16** + **Redis 7** (dev).
- `gonsgarage-rules/` y `docs/DOCUMENTATION_INDEX.md` como gobernanza y mapa de docs.
- **CI**: jobs `backend` / `frontend` con `working-directory`; frontend **Vitest** en lugar de Jest.
- **README** raíz alineado con arquitectura y stack de la plantilla (resumen + enlaces).

---

## Fase 1 — Estructura de carpetas backend (plantilla §1, §4.7) — **completada (fachadas + tests de contrato)**

**Objetivo:** `internal/domain`, `internal/service`, `internal/repository`, `internal/handler`, `internal/middleware`, `pkg/` como contrato principal.

**Implementado:**

- Paquetes de fachada con **aliases** a `internal/adapters/http/handlers`, `.../middleware`, `internal/adapters/repository/postgres` y `internal/core/services/*`.
- `internal/domain` re-exporta modelos usados por `cmd/api` (AutoMigrate) y constantes de rol/estado alineadas con `internal/core/domain`.
- `pkg/` reservado con `doc.go` y test de módulo.
- **TDD:** tests `contract_test.go` por paquete (`reflect` sobre punteros de funciones / asignabilidad de tipos) para garantizar que **no se duplica** lógica en la fase de layout.
- `cmd/api` y `tests/integration` importan ya los nuevos paths públicos.

**Pendiente (fases posteriores):** mover implementación física de archivos desde `adapters/` / `core/` hacia los árboles finales (sin aliases).

---

## Fase 2 — Persistencia: **sqlx** + `lib/pq` (plantilla §2)

**Objetivo:** acceso SQL con **github.com/jmoiron/sqlx** y driver **pq**, migraciones con **golang-migrate** (ya presente).

**Estado actual:** **GORM** + driver postgres de GORM.

**Propuesta explícita (stack no listado hoy):** mantener GORM **temporalmente** solo si documentamos la deuda; la plantilla pide sqlx — la migración es el trabajo más grande del backend.

**Acciones:** introducir `sqlx` en un bounded context (p. ej. `auth` o `health`), tabla de migraciones compartida, retirar GORM por vertical.

---

## Fase 3 — Redis cliente (plantilla §2 vs repo)

**Plantilla:** `github.com/go-redis/redis/v8` + **miniredis** en tests.

**Repo:** `github.com/redis/go-redis/v9`.

**Decisión recomendada:** **permanecer en v9** (mantenimiento activo, API cercana). Documentado en `gonsgarage-rules/02-stack.md` como desviación justificada; tests con **miniredis** o contenedor en CI cuando toque.

---

## Fase 4 — Logging **Zerolog** (plantilla §2)

**Objetivo:** logging estructurado con **rs/zerolog**.

**Estado actual:** mezcla **slog** / helpers legacy.

**Acciones:** `pkg/logging` (zerolog), sustituir puntos de entrada (`cmd/api`, `cmd/worker`) y bajar por capas.

---

## Fase 5 — OpenAPI / CORS / validación (plantilla §3)

- Asegurar **swag** + **gin-swagger** publicados bajo `/api/v1` o ruta acordada.
- **gin-contrib/cors** y validación **validator** donde falte.
- Contrato JSON **camelCase** (§4.4) auditado en handlers.

---

## Fase 6 — Frontend stack §3 (librerías)

**Objetivo:** dependencias listadas: **next-themes**, **react-hook-form** + **zod** + **@hookform/resolvers**, **date-fns** + **react-day-picker**, **sonner**, **tailwindcss-animate**; Radix/Shadcn ampliados según pantallas.

**Acciones:** añadir paquetes (Fase 6a), adoptar en formularios nuevos y refactor progresivo.

---

## Fase 7 — `{{VALIDATION_RULES}}` y MVP (plantilla §4.8, §5)

- Definir `docs/validation-rules.md` (o similar) cuando exista mercado fiscal/teléfono.
- Rellenar `{{MVP_SCOPE}}`, `{{INTEGRATIONS}}`, `{{TENANCY}}` en `template_project.md` §5 para el siguiente ciclo de producto.

---

## Orden sugerido de ejecución

0 → 1 (estructura) → 6a (deps UI) en paralelo con 4 (zerolog en cmd) → 2 (sqlx) → 5 → 7.

# Siguientes pasos (post‑MVP v1)

Priorización después del cierre operativo del checklist MVP (fases 1–5) y de [`mvp-role-access`](../openspec/specs/mvp-role-access/spec.md). **Accounting** HTTP/UI sigue fuera: [`p1-accounting-defer`](../openspec/specs/p1-accounting-defer/spec.md).

## Orden sugerido

1. **P0** — filas P0 de la [tabla](#tabla-p0--p1--p2) (seeds, CI `-race`, secretos).  
2. **Fiabilidad deploy (Arnela)** — si `DATABASE_URL` usa `arnela-postgres`, cerrar red y `COMPOSE_OVERRIDE` **antes** de depender del API en servidor: [Incidente DNS arnela-postgres](#incidente-dns-arnela-postgres) y [`deploy/README.md`](../deploy/README.md) (script del servidor).  
3. **P1** — misma tabla (admin `/employees`, issue GitHub).  
4. **P2** — misma tabla y [`roadmap.md`](./roadmap.md).

## Incidente DNS arnela-postgres

Con **`.env.prod`** apuntando a **`arnela-postgres`** pero un `docker compose` **solo** con `docker-compose.prod.yml`, el API queda en la red `gonsgarage-network` y **no** resuelve ese hostname (Postgres vive en la red de Arnela): *lookup arnela-postgres … no such host*, contenedor en **reinicio**, **502** en `/health`. **Fix:** `export COMPOSE_OVERRIDE=docker-compose.prod.arnela-network.yml` (y verificar `name:` de la red externa en ese YAML) antes de `up`, o el equivalente manual con dos `-f`. Detalle: [`deploy/README.md`](../deploy/README.md).

## Tabla P0 / P1 / P2

| P | Ítem | Criterio “hecho” observable |
|---|------|----------------------------|
| **P0** | Seeds en Postgres dev | `go run ./cmd/seed-mvp-users` dos veces seguidas: segunda ejecución sin error y sin duplicar usuarios (log “ya existe”). |
| **P0** | CI backend con carrera de datos | Job Linux: `go test ./... -race` en verde (CGO disponible en runner). |
| **P0** | Secretos en servidor real | Ningún DSN ni `JWT_SECRET` en repo; solo env en host (ya documentado en `deploy/README.md` — revisión periódica). |
| **P1** | Cobertura middleware `admin` en `/employees` | Test `GET /api/v1/employees` con JWT `role=admin` → **200** (no 403 por `RequireStaffManagers`) — ver `backend/internal/handler/mvp_role_access_test.go`. |
| **P1** | Issue en GitHub (tracker) | Issue abierto con título sugerido: **`test: GET /api/v1/employees accepts admin JWT`** (cuerpo: enlace a este doc + spec `mvp-role-access`). Marcar checkbox abajo cuando exista. |
| **P2** | Matriz Arnela Fase 0 | De [`roadmap.md`](./roadmap.md): filas de `arnela-specs.md` convertidas en issues priorizados. |
| **P2** | Changelog / versionado API | Documento o sección en `docs/` + regla en PR (roadmap Fase 4). |

## P2 — Paridad Arnela (issues sugeridos)

Change SDD (archivado): [`openspec/changes/archive/2026-04-21-arnela-parity/`](../openspec/changes/archive/2026-04-21-arnela-parity/) — matriz en [`arnela-specs.md`](./arnela-specs.md). Crear issues (o tareas) con estos títulos; cuerpo: enlace a este doc + fila de la matriz.

1. **`docs: evaluar golang-migrate vs GORM AutoMigrate`** — Trazabilidad SQL y revisión de esquema al estilo Arnela sin romper el arranque actual en `cmd/api`.
2. **`docs: índice DOCUMENTATION_INDEX.md`** — Espejo opcional de Arnela para descubribilidad (`docs/DOCUMENTATION_INDEX.md`).
3. **`api: rate limiting (paridad Arnela)`** — Alinear con políticas de Arnela si el producto lo exige.
4. **`spike: Next / Tailwind / Shadcn (alcance mínimo)`** — Upgrade solo si el beneficio justifica el churn (matriz: Arnela Next 16 + Tailwind v4).
5. **`docs: matriz roles GonsGarage vs arnela-rules`** — Comparación explícita de permisos por rol frente a `arnela-rules/`.

## Checklist issue GitHub (P1)

- [ ] Issue creado en el repositorio con enlace a `docs/mvp-next-steps.md` y a `openspec/specs/mvp-role-access/spec.md`.

## Referencias

- Checklist MVP: [`mvp-solo-checklist.md`](./mvp-solo-checklist.md)
- Plan épico archivado: `openspec/changes/archive/2026-04-20-mvp-funcionando-plan/proposal.md`
- Change SDD archivado: [`mvp-post-v1-followup` proposal](../openspec/changes/archive/2026-04-20-mvp-post-v1-followup/proposal.md)

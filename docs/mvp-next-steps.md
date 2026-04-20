# Siguientes pasos (post‑MVP v1)

Priorización después del cierre operativo del checklist MVP (fases 1–5) y de [`mvp-role-access`](../openspec/specs/mvp-role-access/spec.md). **Accounting** HTTP/UI sigue fuera: [`p1-accounting-defer`](../openspec/specs/p1-accounting-defer/spec.md).

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

## Checklist issue GitHub (P1)

- [ ] Issue creado en el repositorio con enlace a `docs/mvp-next-steps.md` y a `openspec/specs/mvp-role-access/spec.md`.

## Referencias

- Checklist MVP: [`mvp-solo-checklist.md`](./mvp-solo-checklist.md)
- Plan épico archivado: `openspec/changes/archive/2026-04-20-mvp-funcionando-plan/proposal.md`
- Change SDD archivado: [`mvp-post-v1-followup` proposal](../openspec/changes/archive/2026-04-20-mvp-post-v1-followup/proposal.md)

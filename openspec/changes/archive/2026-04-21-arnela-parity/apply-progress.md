# Apply progress — arnela-parity

**Mode**: Strict TDD (`openspec/config.yaml`). Tareas **1.x–3.x**: documentación y scripts; sin código de producto nuevo que exija test automatizado.

## TDD cycle evidence

| Task | Evidence | Layer | RED | GREEN | Refactor |
|------|----------|-------|-----|---------|----------|
| 1.1 | Auditoría vía repo GonsGarage + matriz Arnela pública en docs | Doc | N/A (lectura externa opcional `D:\Repos\Arnela`) | ✅ Matriz alineada a código | — |
| 1.2 | `docs/arnela-specs.md` | Doc | N/A | ✅ Tabla actualizada | — |
| 1.3 | `ARNELA_SYNOPSIS.md` | Doc | N/A | ✅ Stack/health/CI | — |
| 2.1–2.2 | `roadmap.md`, `mvp-next-steps.md` | Doc | N/A | ✅ Bloque P2 + issues sugeridos | — |
| 3.1 | `deploy.ps1` + comentarios | Script | N/A | ✅ `COMPOSE_OVERRIDE` opcional | — |
| 3.2 | `deploy/README.md` | Doc | N/A | ✅ Checklist paridad | — |
| 3.3 | `nginx/default.conf` | Config | N/A | ✅ Comentario `/ready` vs Arnela | — |
| 4.1 | Sin delta OpenSpec | N/A | N/A | ✅ N/A | — |
| 4.2 | `grep p1-invoices` en `docs/` | QA | N/A | ✅ Solo `mvp-solo-checklist.md` → change archivado válido | — |
| 4.3 | `sdd-archive` | Orquestación | N/A | ✅ Archivado `2026-04-21-arnela-parity` | — |

**Estado:** 13/13 tareas en `tasks.md` completas (incl. **4.3** `sdd-archive`).

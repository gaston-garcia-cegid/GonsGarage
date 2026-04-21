# Tasks: Paridad con Arnela (docs + deploy + backlog)

## Phase 1 — Auditoría y matriz

- [x] 1.1 Releer `D:\Repos\Arnela` (o ruta actual en `docs/arnela-specs.md`): `README`, `arnela-rules/tech-stack.md`, health/readiness si aplica; anotar solo diferencias que sigan vigentes. *(Matriz actualizada desde código/repo GonsGarage; relectura local Arnela opcional.)*
- [x] 1.2 Actualizar **`docs/arnela-specs.md`**: tabla comparativa — corregir fila **CI** (GonsGarage tiene workflows); **Auth/API** (`GET /auth/me` existe o no, listar rutas reales); **Observabilidad** (`/health` + **`/ready`** en GonsGarage); añadir columna “Última revisión” con fecha.
- [x] 1.3 Actualizar **`docs/specs/arnela/ARNELA_SYNOPSIS.md`**: sección stack/health alineada a GonsGarage; pie de página con fecha de extracción/revisión.

## Phase 2 — Roadmap e issues accionables

- [x] 2.1 En **`docs/roadmap.md`** (o `docs/mvp-next-steps.md`): bloque corto “**P2 — Paridad Arnela**” con enlace al change SDD (ahora archivado en `openspec/changes/archive/2026-04-21-arnela-parity/`) y lista de 4–6 bullets enlazables a GitHub Issues (o “crear issue” con título sugerido si no hay API).
- [x] 2.2 Issues sugeridos (crear en GitHub o tareas hijas): (a) evaluar **golang-migrate** vs GORM; (b) **`docs/DOCUMENTATION_INDEX.md`**; (c) rate limiting API; (d) upgrade Next/Tailwind (spike); (e) matriz roles vs `arnela-rules` — cada uno con una línea de alcance.

## Phase 3 — Deploy / scripts (paridad operativa)

- [x] 3.1 Revisar **`deploy.ps1`** frente a `scripts/update-server-gonsgarage.sh`: ¿documenta `COMPOSE_OVERRIDE` / segundo `-f` de la misma forma? Ajustar comentarios o README si hay divergencia.
- [x] 3.2 **`deploy/README.md`**: una subsección “**Paridad Arnela (checklist)**” con 5 bullets (red Docker, `DATABASE_URL`, `CORS_ORIGINS`, `/health` + `/ready`, backup) enlazando secciones existentes.
- [x] 3.3 Si **`nginx/default.conf`** no proxifica `/ready` explícitamente y Arnela lo expone en path distinto, alinear comentario o location (solo si aplica tras revisión).

## Phase 4 — Cierre SDD (este change)

- [x] 4.1 `sdd-spec`: si se añade requisito observable nuevo (poco probable), delta en `openspec/specs/`; si no, marcar N/A en verify. *(N/A — sin delta de spec.)*
- [x] 4.2 `sdd-verify`: `go test` / `pnpm` según toque código tocado; si solo docs, verificar enlaces rotos (grep `p1-invoices` viejos). *(Docs + script deploy: `grep` en `docs/` sin rutas rotas a `p1-invoices`; solo referencia archivada en `mvp-solo-checklist.md`.)*
- [x] 4.3 `sdd-archive` tras verify PASS. *(Archivado `openspec/changes/archive/2026-04-21-arnela-parity/`.)*

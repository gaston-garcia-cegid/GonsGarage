# Tasks: MVP status — next iteration

## Phase 1: Narrativa y orden en docs

- [x] 1.1 En `docs/mvp-next-steps.md`, añadir sección breve **«Orden sugerido»** (P0 → fiabilidad deploy Arnela/`COMPOSE_OVERRIDE` → P1 → P2) con enlaces a subsecciones existentes; **no** duplicar la tabla P0/P1/P2 entera.
- [x] 1.2 En el mismo archivo, subsección **«Incidente: DNS `arnela-postgres`»** (3–6 líneas): causa (red distinta), fix (`export COMPOSE_OVERRIDE=docker-compose.prod.arnela-network.yml` o equivalente), enlace a `deploy/README.md`.
- [x] 1.3 En `docs/roadmap.md`, una línea o párrafo corto que remita a `mvp-next-steps.md` para orden post‑MVP y al apartado deploy/Arnela.

## Phase 2: Deploy README

- [x] 2.1 En `deploy/README.md`, bloque **«Script del servidor»** (o ampliar Postgres Arnela): comando con `COMPOSE_OVERRIDE` / `update-server-gonsgarage.sh` y advertencia de **502/restart** si falta el segundo `-f`.
- [x] 2.2 Enlazar desde ese bloque a `scripts/update-server-gonsgarage.sh` (ruta relativa repo).

## Phase 3: Script `update-server-gonsgarage.sh`

- [x] 3.1 Sustituir el **WARN** actual (cuando `.env.prod` referencia `arnela-postgres` y `COMPOSE_OVERRIDE` vacío) por **`exit 1`** con mensaje claro **antes** de `docker compose up`, si existe `docker-compose.prod.arnela-network.yml` en `GONSGARAGE_DIR` (fail‑fast alineado al riesgo “wrong network” mitigado sin auto‑attach).
- [x] 3.2 Documentar en comentario al inicio del script (debajo del bloque Postgres Arnela) que el operador puede `export COMPOSE_OVERRIDE=docker-compose.prod.arnela-network.yml` antes de invocar.

## Phase 4: Compose y checklist proposal

- [x] 4.1 En cabecera de `docker-compose.prod.yml` (comentarios), una línea: si `DATABASE_URL` usa `arnela-postgres`, compose **debe** incluir `-f docker-compose.prod.arnela-network.yml` (o variable `COMPOSE_OVERRIDE` del script).
- [x] 4.2 Revisar checkbox **issue GitHub** en `docs/mvp-next-steps.md`: dejar `[ ]` si no hay issue; si ya existe en GitHub, marcar `[x]` y poner URL en la misma línea o debajo.

## Phase 5: Cierre del change (docs)

- [x] 5.1 En `openspec/changes/mvp-status-next-iteration/proposal.md`, marcar **Success Criteria** con `[x]` cuando 1.1–4.2 estén verificados.
- [x] 5.2 Confirmar que ningún cambio introduce rutas HTTP ni UI de accounting (criterio `p1-accounting-defer`).

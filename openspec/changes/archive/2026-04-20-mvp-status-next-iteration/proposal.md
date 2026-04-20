# Proposal: MVP status — next iteration

## Intent

Sequence **post–MVP v1** work after checklist closure (**1–5** done; **6** partial in `docs/mvp-solo-checklist.md`). Align with **`docs/mvp-next-steps.md`** (refine; do not paste the P0/P1/P2 table), **`docs/roadmap.md`**, and **`openspec/specs/mvp-role-access/spec.md`**. Fold in **prod deploy lessons**: Postgres on Arnela (`arnela-postgres`) needs **`docker-compose.prod.arnela-network.yml`**; **empty `COMPOSE_OVERRIDE`** can yield API **restart loops** (DNS “no such host”). **`openspec/specs/p1-accounting-defer/spec.md`** remains in force: **no accounting HTTP/UI** here.

## Scope

### In Scope

- **State:** MVP v1 frozen; priorities in `mvp-next-steps.md` + roadmap; role matrix in `mvp-role-access`; invoices row **Deferred** via `p1-accounting-defer`.
- **Order:** (1) **P0** from `mvp-next-steps.md` (Postgres seeds idempotency, CI `-race`, secrets); (2) **reliability** — `deploy/README.md` + Compose guidance for Arnela; optionally **harden `scripts/update-server-gonsgarage.sh`**: **fail-fast** when `.env.prod` names `arnela-postgres` and override is absent, and/or **auto `-f docker-compose.prod.arnela-network.yml`** only if file exists and `name:` is documented (avoid wrong external network); (3) **P1** — `/employees` admin path + **GitHub issue** checkbox if still open; (4) **P2** — Arnela matrix → issues; API changelog (roadmap Fase 4).

### Out of Scope

Accounting routes/UI; multi-tenant; i18n; payment integrations.

## Capabilities

### New Capabilities

None

### Modified Capabilities

None

## Approach

**Docs-first:** update `mvp-next-steps.md`, `roadmap.md`, `deploy/README.md` for Arnela + `COMPOSE_OVERRIDE`. Script: today **WARN** only (lines 54–58); optional **exit 1** before `compose up`, or documented auto-override. No OpenSpec requirement edits unless a later change touches **`mvp-role-access`**.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `docs/mvp-next-steps.md`, `docs/roadmap.md` | Modified | Ordering; deploy caveat; link table, don’t duplicate. |
| `deploy/README.md` | Modified | Arnela DNS / `COMPOSE_OVERRIDE` troubleshooting. |
| `scripts/update-server-gonsgarage.sh` | Modified (optional) | Fail-fast or safe auto-override vs warn-only. |
| `docker-compose.prod*.yml` | Reference | Comments / `name:` clarity if needed. |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Auto-override attaches wrong network | Med | Prefer fail-fast; verify `docker inspect` on host. |
| Accounting creep | Low | Keep `p1-accounting-defer` in every summary. |

## Rollback Plan

Revert doc/script **commit**; no DB migrations. Restore prior script and redeploy if behavior regresses.

## Dependencies

Host **external** network name matches `docker-compose.prod.arnela-network.yml` (`name:`).

## Success Criteria

- [x] Narrative orders **P0 → deploy reliability → P1 → P2** with pointers (not a verbatim P0/P1/P2 copy).
- [x] **`arnela-postgres` + empty `COMPOSE_OVERRIDE`** incident and fix documented.
- [x] Accounting stays **out**; `p1-accounting-defer` intent unchanged.
- [x] GitHub issue checkbox in `mvp-next-steps.md` matches repo state.

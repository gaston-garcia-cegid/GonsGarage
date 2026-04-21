# Proposal: Paridad operativa y documentación con Arnela

## Intent

Cerrar la **brecha declarada** entre GonsGarage y el stack/plantilla Arnela (`docs/arnela-specs.md`, `ARNELA_SYNOPSIS.md`): matriz desactualizada (p. ej. CI ya existe), observabilidad mal descrita (`/ready` ya existe), y **P2** del roadmap sin issues accionables. Objetivo: **paridad documentada y reproducible** en deploy LAN + lista priorizada de divergencias técnicas (sin reescribir el producto).

## Scope

### In Scope

- Actualizar **matriz Arnela ↔ GonsGarage** y sinopsis con estado real (2026-04).
- **Issues o tareas** concretas desde filas P2 (migraciones GORM vs migrate, docs index, rate limit, versiones front).
- Ajustes menores **deploy**: alinear nombres/comentarios (`/ready` vs “readiness”), `deploy.ps1` / ejemplo env si falta paridad con `update-server-gonsgarage.sh`.
- Opcional acotado: **un** endpoint o contrato documentado en OpenSpec solo si se redefine comportamiento observable.

### Out of Scope

- Upgrade Next 16 / Tailwind v4 / Shadcn (coste alto; decisión explícita aparte).
- Sustituir GORM AutoMigrate por golang-migrate **en este change** (solo decisión o spike doc).
- Copiar `arnela-rules/` al repo (se sigue leyendo `D:\Repos\Arnela`).

## Capabilities

### New Capabilities

- None *(cambios principalmente docs + checklist + issues; sin nuevo dominio de producto).*

### Modified Capabilities

- None *(no se alteran requisitos de `mvp-role-access`, `invoices`, `billing`, etc.).*

## Approach

1. Auditoría rápida repo vs `docs/specs/arnela/ARNELA_SYNOPSIS.md` y tabla en `arnela-specs.md`.  
2. PR o commits pequeños: docs + `mvp-next-steps` / `roadmap` enlaces a issues.  
3. Deploy: revisar `deploy.ps1`, `.env.prod.example`, health en nginx si Arnela expone patrón distinto.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `docs/arnela-specs.md` | Modified | Matriz al día; CI, `/ready`, compose Arnela. |
| `docs/specs/arnela/ARNELA_SYNOPSIS.md` | Modified | Fecha + filas que cambiaron. |
| `docs/roadmap.md`, `docs/mvp-next-steps.md` | Modified | Enlaces a issues P2 / este change. |
| `deploy/README.md`, `deploy.ps1` | Modified | Solo si se detecta gap vs script bash. |
| `.github/` | Optional | Issues creados manualmente en GitHub (documentados en tasks). |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Matriz sigue obsoleta al cabo de semanas | Med | Fecha “última revisión” y enlace a este change archivado. |
| Confundir `/ready` con naming Arnela | Baja | Una línea en synopsis unificando vocabulario. |

## Rollback Plan

Revertir commits de docs/deploy; no hay migración de datos ni flag de producto.

## Dependencies

- Acceso de lectura al repo **Arnela** local (ruta en docs) para contrastar.
- Opcional: permisos para crear **issues** en GitHub GonsGarage.

## Success Criteria

- [x] Matriz `arnela-specs.md` sin filas falsas (CI, health).
- [x] `ARNELA_SYNOPSIS.md` refleja `/ready` y stack GonsGarage actual.
- [x] `tasks.md` de este change con **≥5** ítems cerrables con criterio de hecho.
- [x] Roadmap o `mvp-next-steps` enlaza este change o issues derivados.

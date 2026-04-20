# Proposal: MVP docs alignment — siguientes pasos

## Intent

El **MVP checklist** (`docs/mvp-solo-checklist.md`) tiene fases **1–5** y endurecimiento **hechos**; **Fase 6** solo **parcial** (repairs staff + doc deploy manual; CI→servidor sigue opcional). **`mvp-minimum-phases.md`** aún marca Fase C opcional staff sin tachar pese a 6.1 cerrada — **deuda de documentación**. El **roadmap** deja abiertas Fase 0 (matriz→issues), Fase 4 (changelog, plantillas, observabilidad) y Fase 5 (secretos/CORS/imagen). **`mvp-next-steps.md`** ya ordena P0→deploy→P1→P2; falta **ejecutar** filas (issue GitHub, seeds smoke, etc.). Este change **prioriza y acota** el siguiente tramo sin nuevo alcance de producto.

## Scope

### In Scope

- **Alinear** `mvp-minimum-phases.md` Fase C con el checklist (staff repairs **hecho**; opcional restante = solo CI→servidor u otra mejora explícita).
- **Refrescar** “estado rápido” en checklist Fase 6 si hace falta una línea explícita “CI→servidor: opcional / no priorizado”.
- **Ejecutar** filas pendientes de `mvp-next-steps.md` en el orden ya definido: P0 (smoke seeds, confiar CI `-race` Linux), P1 (issue GitHub si checkbox sigue vacío), P2 (issues desde `arnela-specs.md` + borrador changelog en `docs/`).
- **Roadmap Fase 0**: crear issues o tabla de seguimiento enlazada desde `roadmap.md` (sin paridad Arnela completa).

### Out of Scope

Accounting HTTP/UI; nuevas rutas invoice/billing; multi-tenant; i18n; pasarelas de pago.

## Capabilities

### New Capabilities

None

### Modified Capabilities

None

## Approach

**Docs primero** (1–2 h): PR único que actualiza `mvp-minimum-phases.md` + checklist Fase 6 si aplica. **Luego** tareas operativas fuera del repo o en tracker: issue GitHub P1, smoke seeds, issues matriz Arnela. **Changelog**: sección nueva en `docs/` o `CHANGELOG.md` + una línea en `roadmap.md` Fase 4 cuando exista el doc.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `docs/mvp-minimum-phases.md` | Modified | Fase C / criterios alineados con MVP actual. |
| `docs/mvp-solo-checklist.md` | Modified (opc.) | Aclaración Fase 6 / CI→servidor. |
| `docs/roadmap.md`, `docs/mvp-next-steps.md` | Modified (opc.) | Enlaces a issues o changelog cuando existan. |
| GitHub | External | Issues desde matriz + issue P1 si falta. |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Docs desactualizan de nuevo | Med | Checklist “última revisión” fecha en cabecera de fases. |
| Scope creep a accounting | Low | Enlazar `p1-accounting-defer` en cada PR de este change. |

## Rollback Plan

Revertir commits solo de docs; issues en GitHub se cierran o archivan manualmente.

## Dependencies

Ninguna bloqueante; acceso al repo GitHub para issues.

## Success Criteria

- [ ] `mvp-minimum-phases.md` Fase C coherente con repairs staff en producción del checklist.
- [ ] Al menos **una** fila P0 o P1 de `mvp-next-steps.md` cerrada con evidencia (log, CI verde, o issue creado).
- [ ] Roadmap Fase 0: **≥1** issue o enlace de tracking desde `roadmap.md`, o justificación breve en doc si se difiere.

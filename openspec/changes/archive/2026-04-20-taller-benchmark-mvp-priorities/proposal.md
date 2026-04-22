# Proposal: Backlog MVP (referencia TALLER ALPHA)

## Intent

Priorizar **features** de GonsGarage frente a material público de [TALLER ALPHA](https://www.youtube.com/@DesignSoftSA) (vídeos: [organizar taller](https://www.youtube.com/watch?v=ZMKI_Tna8UI), [recepción vehicular app](https://www.youtube.com/watch?v=rv7wD9rtheI)). **No hay análisis de vídeo aquí** — solo títulos oEmbed. La tabla mezcla esas pistas con el **típico** de software de tallerest. Validar con visionado o capturas antes de norma SHALL nueva.

## Scope

### In Scope

- Backlog P0–P2 y mapa a specs; deltas en otro change (`sdd-spec`).

### Out of Scope

- Código; paridad 1:1 con TALLER ALPHA sin criterio acordado.

### Lista priorizada

| Pri | Feature | Nota |
|-----|---------|------|
| P0 | Recepción (km, fluidos, neumáticos, notas) + **flujo “app”** (rápido; móvil/ pantalla dedicada) | 2.º título; `workshop-repair-execution` ya cubre recepción/cierre — reforzar UX/API de error. |
| P0 | Visita **open/closed** + trazas | Especificado; cruzar `mvp-role-access`. |
| P1 | **Cola / hoy en taller** | 1.er título “organizar”; lista o agenda día sin inflar a ERP. |
| P1 | **Repairs ↔ visita** | Cerrar contrato si falta en UI. |
| P2 | OBD, presupuesto, **fotos** | Stubs o post-MVP1 según spec actual. |
| — | Facturación | Diferido (`p1-accounting-defer`). |

## Capabilities

### New Capabilities

- None.

### Modified Capabilities

- `workshop-repair-execution`: requisitos de **experiencia** de recepción (incl. flujo estilo app) y, si se valida, **vista “hoy”** sin romper visita como unidad.
- `mvp-role-access`: solo si nuevo rol o ruta; si no, sin cambio.

## Approach

Validar P0 con el equipo → deltas en `changes/.../specs/workshop-repair-execution/` → UI staff existente + `service-jobs`.

## Affected Areas

| Area | Impact |
|------|--------|
| `openspec/specs/workshop-repair-execution` | Delta futuro |
| `frontend` (taller) | UI recepción / opcional cola |
| `backend` `service-jobs` | Si filtros o consultas “hoy” nuevas |

## Risks

| Risk | L/M | Mitigation |
|------|-----|------------|
| Features inventadas vs vídeo | M | Lista firmada por quien vio el material |
| “Organizar” → scope ERP | M | Acotar a cola/día, no stock completo |

## Rollback Plan

Borrar o sustituir este `proposal.md`; no hay entrega de código.

## Dependencies

- Opcional: capturas o bullet list de pantallas de los vídeos.

## Success Criteria

- [ ] P0 consensuado y trazable a tareas o delta.
- [ ] Sin contradicción con visita como trazabilidad principal.

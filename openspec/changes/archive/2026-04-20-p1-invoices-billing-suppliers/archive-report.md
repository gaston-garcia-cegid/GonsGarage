# Archive report — p1-invoices-billing-suppliers

**Archived at**: 2026-04-20  
**Destination**: `openspec/changes/archive/2026-04-20-p1-invoices-billing-suppliers/`  
**Mode**: OpenSpec (filesystem merge + move)

## Preconditions

- `tasks.md`: **23/23** `[x]`
- `verify-report.md`: **PASS** (sen issues CRITICAL)

## Specs synced to main (`openspec/specs/`)

| Domain | Action | Details |
|--------|--------|---------|
| `invoices/` | **Created** | `spec.md` — facturas recibidas (P1). |
| `billing/` | **Created** | `spec.md` — billing emitido (P1). |
| `suppliers/` | **Created** | `spec.md` — proveedores (P1). |
| `p1-accounting-defer/` | **Updated** | Fusionado delta do change: requisito *Programme P1 tracking*; R-1/R-2 ampliados con glosario e referencia ao change archivado. |

## Archive contents (audit trail)

Tras o `move`, esta carpeta conserva:

- `proposal.md`
- `design.md`
- `tasks.md`
- `apply-progress.md`
- `verify-report.md`
- `archive-report.md` (este ficheiro)
- `specs/` (delta histórico)

## Source of truth

Os requisitos P1 viven agora en:

- `openspec/specs/invoices/spec.md`
- `openspec/specs/billing/spec.md`
- `openspec/specs/suppliers/spec.md`
- `openspec/specs/p1-accounting-defer/spec.md` (contexto MVP + ponte P1)

## SDD cycle

Plan → apply → verify → **archive** completado. Lista para o seguinte change.

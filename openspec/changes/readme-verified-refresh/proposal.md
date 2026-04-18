# Proposal: README raíz alineado al stack real

## Intent

El `README.md` de la raíz es largo, mezcla plantilla genérica con detalle útil, y **no refleja versiones ni herramientas actuales** del repo. Eso confunde a nuevos colaboradores y desincroniza badges, prerequisitos y comandos respecto a `go.mod`, `frontend/package.json` y `openspec/config.yaml`.

## Verificación de versiones (estado actual)

| Afirmación en README raíz (o implícita) | Fuente de verdad | ¿Alineado? |
|----------------------------------------|------------------|------------|
| Go 1.21+ (badges y sección Backend) | `backend/go.mod` → `go 1.25.3` | **No** |
| Next.js 15.0+ (badge) | `frontend/package.json` → `next` **15.5.5** | Parcial (15.x sí; badge “15.0+” engañoso frente a pin) |
| Next.js 15 (App Router) en stack | Coincide con dependencia | Sí |
| React | `react` / `react-dom` **19.1.0** | README no lo destaca; conviene documentarlo |
| TypeScript 5+ | `typescript` **^5** en frontend | Sí (rango) |
| Tests frontend “Jest + RTL” | `pnpm test` → **Vitest**; Jest sigue en `devDependencies` | **No** (prioridad mal descrita) |
| `pnpm run generate-types` en Quick Start | **No existe** script en `frontend/package.json` | **No** |
| Contexto monorepo en OpenSpec | `openspec/config.yaml` ya dice Go 1.25 + Next 15 + React 19 | README debería alinear con esto |

El `frontend/README.md` es corto y no fija versión de Next; puede quedar como enlace al README raíz o alinearse en una tarea opcional.

## Scope

### In Scope

- **Sustituir por completo** el archivo raíz `README.md` por uno nuevo: estructura clara (qué es el proyecto, stack verificado con tablas/citas de versiones, prerequisitos, arranque rápido con comandos que existan, enlaces a `docs/`, `CONTRIBUTING.md`, CI).
- **Auditar y corregir** badges y bullets para que coincidan con `go.mod`, `frontend/package.json` y flujo real de tests (Vitest como principal).
- Eliminar o mover a `docs/` solo si encaja en el diseño del nuevo README (evitar duplicar novelas largas: el README debe ser **entrada**, no manual completo).

### Out of Scope

- Reescribir todo `docs/` ni specs OpenSpec.
- Cambiar versiones en código (solo documentación).
- Docker/compose de producción del README salvo corregir errores factuales mínimos si se citan.

## Capabilities

### New Capabilities

- None

### Modified Capabilities

- None

## Approach

1. Leer `backend/go.mod`, `frontend/package.json`, `.github/workflows/ci.yml` (comandos CI) y `docs/README.md` para no contradecer la documentación existente.
2. Redactar README nuevo en Markdown: español o neutro según convención del repo (el código/UI usa pt en producto; el README histórico está en inglés con trozos en español — **unificar** en un idioma coherente, p. ej. inglés para README raíz salvo decisión contraria en spec/design).
3. Opcional (tasks): una pasada breve en `frontend/README.md` (versión Next o enlace “ver raíz”).

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `README.md` | Removed / Replaced | Nuevo README verificado frente a manifests. |
| `frontend/README.md` | Optional | Ajuste mínimo de enlaces o versión. |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Enlaces rotos al reorganizar secciones | Med | Buscar referencias internas (`grep`) tras el cambio. |
| Contenido útil solo en el README antiguo | Med | Conservar en el nuevo README enlaces a `docs/` donde ya exista profundidad. |

## Rollback Plan

`git revert` del commit que reemplace `README.md` (y `frontend/README.md` si se tocó).

## Dependencies

- Ninguna externa.

## Success Criteria

- [x] Badges y tabla de stack citan **Go 1.25.x** (o la línea del `go` directive) y **Next 15.5.x** / **React 19.x** según `package.json`.
- [x] Comandos de “Quick start” existen en `package.json` / Makefile / scripts documentados (sin `generate-types` salvo que se añada el script en otro change).
- [x] Tests frontend descritos como **Vitest** (primario), Jest opcional si aplica.
- [x] Enlace claro al hub `docs/README.md` y a contribución.

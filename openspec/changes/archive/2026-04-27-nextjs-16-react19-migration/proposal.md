# Proposal: Next.js 16 + React 19 patterns (RSC, `use`, `useOptimistic`)

## Intent

Subir **Next.js 15.5 → 16.x** (o superior estable) con **React 19**, **Tailwind v4** y **refactor amplio de Zustand** en el **mismo** cambio, sin regresiones de producto. Reducir `useEffect` usado solo para fetch o sincronización que encaje en **Server Components**. Adoptar **`use`** / **`useOptimistic`** donde aporten claridad o UX, documentando cada cambio estructural.

## Scope

### In Scope

- Bump `next`, `eslint-config-next`, `@types/*`; **Tailwind v4** (config, tokens, migración de clases); `pnpm build`, `pnpm test`, `pnpm lint`, `pnpm typecheck` y **CI** verdes.
- **Mega-refactor Zustand**: reestructura de stores, límites client/server, hidratación y consumo alineados con RSC/islands; mismo comportamiento observable en UI y llamadas API.
- Auditoría por rutas: **Server Components** + datos en servidor donde proceda; **islands** `"use client"` para estado local, modales y eventos.
- **React 19**: `use()` / `useOptimistic` en flujos acordados en diseño.
- **ADR** + mapa “antes/después” por archivo o dominio importante.

### Out of Scope

- Playwright/E2E nuevos (salvo decisión explícita aparte).
- Cambios de contrato HTTP backend fuera de lo estrictamente necesario para el cliente.

## Capabilities

### New Capabilities

- None

### Modified Capabilities

- **ui-brand-shell**: Migración **Next 16 + Tailwind v4** en `main`, tokens y shell coherentes; ADR **GO** (sin DEFER de TW v4).
- **ui-component-system**: Límites **RSC/client**, uso de **React 19** (`use`, `useOptimistic`) y directrices del **refactor Zustand** (stores, primitivas UI en client).

## Approach

Trabajo en **`main`** (no se exige rama dedicada); commits atómicos razonables aunque el destino sea un solo gran entregable. Orden sugerido: (1) bump Next + fix build; (2) Tailwind v4 + ajuste de estilos globales; (3) refactor Zustand por dominio con tests en verde entre pasos; (4) RSC / menos `useEffect`; (5) `use` / `useOptimistic` donde aplique. Cada paso **documentado** en ADR o design con el “por qué”.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `frontend/package.json` + lockfile | Modified | Next 16, React ecosystem, Tailwind v4 |
| `frontend/tailwind.config.*`, `postcss.config.*`, `src/app/globals.css` | Modified | v4 pipeline y tokens |
| `frontend/src/stores/**` | Modified | Mega-refactor Zustand |
| `frontend/src/app/**` | Modified | RSC, islands, menos `useEffect` data-only |
| `frontend/next.config.*` | Modified | Next 16 / bundler |
| `docs/` o `design.md` del change | Modified | ADR GO, tabla de cambios importantes |
| `.github/workflows/*` | Maybe | Node / comandos CI |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Big bang en `main` bloquea a otros | Alto | Ventana acordada; commits reversibles; comunicar en equipo |
| TW v4 + Radix + clases legacy | Alto | Guía de migración; búsqueda de patrones rotos; revisión visual |
| Zustand + RSC (hydration, duplicación) | Alto | Diseño explícito por store; tests RTL y smoke manual |
| Terceros (Radix, auth) | Med | Release notes; aislar fallos con tests |

## Rollback Plan

Revertir la(s) serie(s) de commits en `main` (o `git revert` del rango) restaurando `package.json`, lockfile, configs Tailwind/PostCSS, stores y rutas; ADR **ROLLBACK** con SHA.

## Dependencies

- CI Node ≥ requisito Next 16 / toolchain Tailwind v4.
- Documentación Next 16, Tailwind v4 upgrade, React 19.1.

## Success Criteria

- [ ] `pnpm build` + `pnpm test` + `pnpm typecheck` + `pnpm lint` y **CI** verdes.
- [ ] **Tailwind v4** activo y build sin deprecaciones críticas sin resolver.
- [ ] **Refactor Zustand** completado según `design.md` / checklist de stores.
- [ ] Regresión funcional nula en checklist crítico (auth, taller, peças, accounting).
- [ ] ADR con **cambios importantes** (archivo + motivo) y **GO**.
- [ ] `useEffect` data-only reducidos o documentados como no candidatos con razón.

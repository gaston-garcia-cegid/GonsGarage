# Proposal: Frontend ESLint warnings cleanup (React Hooks)

## Intent

Tras a migración Next 16 / React 19 ([ADR 0003](docs/adr/0003-nextjs16-react19-migration-main.md)), `react-hooks/set-state-in-effect` e regras relacionadas quedaron en **`warn`** para non bloquear o merge. Hoxe hai **~36 avisos** en `pnpm lint`, maiormente **setState síncrono dentro de `useEffect`** e un caso **`react-hooks/exhaustive-deps`**. O obxectivo é **eliminar eses avisos** refactorando ou patrón-alternativos (sen silenciar masivamente con `eslint-disable`), de xeito que o árbore pase **lint con 0 warnings** (ou elevar a regra a `error` cando estea limpo).

## Scope

### In scope

- Inventario por ficheiro/regra; agrupar por dominio (accounting `[id]`/listas, `ThemeSwitcher`, `useAuthHydrationReady`, `dashboard`, `workshop`, etc.).
- Corrixir `set-state-in-effect`: patróns recomendados por React (evitar setState directo no corpo do effect cando dispara re-render en cadea), p.ex. **subscribirse** a hidratación/async con callback, **`queueMicrotask`/`startTransition`** só cando documentado, ou **mover estado derivado** fóra do effect.
- Corrixir `exhaustive-deps` onde a dependencia sexa segura (ou extraer lóxica a `useCallback` estable).
- `pnpm lint`, `pnpm typecheck`, `pnpm test` verdes ao pechar cada fase.

### Out of scope

- Novas features de produto, cambio de contratos API, E2E Playwright.
- **Backend** `go vet`/golangci (separado).
- Refactor masivo de UX accounting (só o mínimo para silenciar hooks de forma lexítima).

## Capabilities

### New capabilities

- None (hixiene de código e calidade de lint).

### Modified capabilities

- None (non se cambia comportamento observable contractual máis alá de cumprir a barra de calidade xa descrita no ADR; se no futuro se quere **SHALL 0 warnings** en spec, facer delta `ui-brand-shell` nunha fase `sdd-spec` opcional).

## Approach

1. Exportar lista: `pnpm lint --format stylish` (ou JSON) e mapa `ruta → regra`.  
2. Fases por cartafol: **accounting** (maior volume), logo **hooks compartidos** (`useAuthHydrationReady`, `ThemeSwitcher`), logo **dashboard/workshop/client**.  
3. Por cada warning: ler o fluxo (load en effect, redirect, Zustand persist); aplicar unha das estratexias da documentación React 19 / guías do equipo (evitar effect só para setear estado que podería ser inicialización ou evento externo).  
4. Cando `frontend/eslint.config.mjs` estea limpo en local/CI, valorar pasar `set-state-in-effect` de `warn` a `error` nun commit final pequeno.

## Affected areas

| Area | Impact | Description |
|------|--------|-------------|
| `frontend/src/app/accounting/**/[id]/page.tsx` e listas | Modified | Varios `void load()` en `useEffect` con setState |
| `frontend/src/components/theme/ThemeSwitcher.tsx` | Modified | `setMounted` / `setPref` no mount effect |
| `frontend/src/hooks/useAuthHydrationReady.ts` | Modified | `setReady` tras hidratación |
| `frontend/eslint.config.mjs` | Modified | Posible `warn` → `error` ao final |
| Tests RTL existentes | Touched só se comportamento cambia | Manter verdes |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Dobre fetch ou flash de UI tras refactor de effects | Med | Probar rutas tocadas; manter `pnpm test`; smoke manual checklist ADR §6.2 onde aplique |
| `eslint-disable` abusivo | Med | Revisión PR: prohibir salvo 1 liña + comentario coa razón RFC |

## Rollback plan

`git revert` do rango de commits da fase; se só se subiu a regra a `error`, revertir ese commit e volver a `warn` mantendo fixes útiles.

## Dependencies

- Ningunha externa; ramo `main` actualizado.

## Success criteria

- [ ] `pnpm lint` en `frontend/` reporta **0 warnings** (ou só excepcións documentadas en ≤3 ficheiros).
- [ ] `pnpm typecheck` e `pnpm test -- --passWithNoTests` verdes.
- [ ] Sen novos `eslint-disable` masivos; cada excepción xustificada.

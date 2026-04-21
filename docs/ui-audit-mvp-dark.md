# Auditoria UI — tema escuro (MVP)

Change: `mvp-ui-visual-parity` (Fase 1).  
Objetivo: comprovar contraste e superfícies nas rotas MVP com `html[data-theme='dark']` (ver delta `openspec/changes/mvp-ui-visual-parity/specs/ui-brand-shell/spec.md`).

## Ajustes de tokens (2026-04-21)

| Token / área | Antes (resumo) | Depois | Notas |
|--------------|----------------|--------|--------|
| `--surface-panel` / `--surface-muted` / `--surface-elevated` | Muito próximos da página | Valores em `#0e1522`–`#141d2d` | Mais hierarquia entre página e painéis |
| `--text-muted` (só dark) | `gray-500` | `gray-400` | Texto secundário mais legível em fundos muito escuros |
| Spinners | `3rem` / `2.5rem` soltos | `--app-loading-size-*` + `AppLoading` | Ver `frontend/src/styles/tokens.css`, `AppLoading.tsx` |
| `employees/page.tsx` (Fase 2) | Hex Tailwind-style inline | `var(--chip-*)`, `var(--color-error)`, superfícies `var(--surface-*)` | Sem `#rrggbb` na primeira vista |
| `DashboardLayout.module.css` (área cliente) | Hex cinzentos / vermelho | Tokens `--surface-*`, `--text-*`, `--brand-signal` | Alinha `/client` ao tema |

## Checklist manual (marcar após smoke)

| Rota | Dark OK | Light OK | Notas |
|------|---------|----------|--------|
| `/dashboard` | [ ] | [ ] | Painel, cartões, reparações recentes |
| `/accounting` (entrada + 1 sub-rota) | [ ] | [ ] | |
| `/client` | [ ] | [ ] | |
| `/employees` | [ ] | [ ] | |
| `/cars` | [ ] | [ ] | Fase 2 — loaders + tema |
| `/appointments` | [ ] | [ ] | Idem |

Instruções: alternar tema no UI; verificar texto legível, bordas visíveis, sem “blobs” claros sem intenção.

### Smoke Fase 2 (dev)

Após `pnpm dev`, percorrer **`/cars`** e **`/appointments`** em **light** e **dark**: listas carregam, `AppLoading` visível só durante fetch, sem erros de consola.

## Loaders unificados

- Componente: `frontend/src/components/ui/AppLoading.tsx` (`sm` \| `md` \| `lg`).
- Pantalla completa (pre-auth): contenedor `aria-busy="true"` + `AppLoading` `lg` + `label` sr-only onde aplique.
- Dentro de `AppShell`: `md` + texto visível junto ao spinner quando existir.

Próximas fases: completar a matriz para rotas migradas (tarefa 5.2 em `tasks.md` do change).

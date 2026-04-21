# ADR 0002 — Stack Shadcn/ui (Fase 4)

**Estado:** Aceptado para el MVP GonsGarage.  
**Fecha:** 2026-04-21  
**Relación Fase 3:** El [ADR 0001](./0001-next16-tailwind4-spike.md) documenta **NO-GO** a Next 16 + Tailwind v4 en `main`. La fundación Shadcn adopta **Tailwind CSS v3** estable con Next 15.

## Decisión

- **Componentes:** [shadcn/ui](https://ui.shadcn.com/) estilo **New York**, copiados en `frontend/src/components/ui/` (`button.tsx`, `input.tsx`, `label.tsx`, `dialog.tsx`, …).
- **Utilidades:** `class-variance-authority`, `clsx`, `tailwind-merge`; helper `cn()` en `frontend/src/lib/utils.ts`.
- **Radix:** dependencias que exigen los primitives (p. ej. `@radix-ui/react-slot`, `react-label`, `react-dialog`).
- **Tema:** variables HSL en `frontend/src/styles/shadcn-theme.css` (light en `:root`, dark en `html[data-theme='dark']`), consumidas por `tailwind.config.ts` (`hsl(var(--primary))`, etc.).
- **Coexistencia:** `tokens.css` sigue siendo la fuente de verdad para layout legacy, marketing y `AppLoading`; el mapa marca ↔ variables Shadcn está en `docs/ui-shadcn-theme.md`.
- **Carpeta canónica:** nuevas primitives en **kebab-case** al estilo shadcn (`button.tsx`); componentes legacy por carpetas (`Loading/`, `Modal/`) se sustituyen progresivamente según `tasks.md`.

## Consecuencias

- Aumento de tamaño de bundle en rutas que importan Radix + Lucide (aceptado para MVP).
- Los módulos CSS globales (`*.module.css`) pueden convivir con utilidades Tailwind en el mismo árbol; priorizar tokens para color semántico fuera de componentes Shadcn.

## Alternativas descartadas

- **Tailwind v4-only (Fase 3 GO):** aplazada por decisión ADR 0001.
- **Solo CSS modules sin design system:** rechazada en el change `mvp-ui-visual-parity` (spec `ui-component-system`).

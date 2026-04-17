# Design: UI minimalista alineada al logotipo

## Technical Approach

Refactor **presentational CSS only**: keep Next.js / React structure unchanged. Align implementation with `ui-brand-shell` by (1) documenting brand roles in `tokens.css`, (2) replacing colour literals in the authenticated shell and priority routes with **existing CSS variables** (`--color-*`, `--brand-*`, `--surface-*`, `--text-*`) or with **shared utilities** in `utilities.css` following the same pattern as `.btn-primary` and `html[data-theme='dark'] .alertError`.

Theme is already driven by `html[data-theme]` set from `frontend/src/lib/theme.ts` and the inline boot script; no new runtime API.

## Architecture Decisions

| Decision | Choice | Alternatives | Rationale |
|----------|--------|--------------|-----------|
| Colour source of truth | Extend/clarify `:root` and `html[data-theme='dark']` in `tokens.css` | CSS-in-JS theme object | Matches current stack (global CSS + CSS modules). |
| Route styling | Map hex → `var(--…)` inside existing `*.module.css` | New Tailwind / component library | Out of scope; avoids dependency churn. |
| Repeated status chips | Optional utilities (e.g. `.badgeInfo`) in `utilities.css` | Duplicate per module | Reduces drift; mirrors existing utility pattern. |
| Logo sampling | Manual compare JPG ↔ tokens; adjust `--brand-*` if off | Automated extraction script | MVP-appropriate; script not required for spec compliance. |

### Decision: Logout / destructive actions

**Choice**: Use `--brand-signal`, `--color-error`, or a dedicated token (e.g. hover via `color-mix` on `--brand-signal`) for logout hover instead of `#b91c1c`.

**Alternatives**: Keep hex; use `opacity` only.

**Rationale**: Satisfies “no orphan brand hex in shell” and keeps dark-mode overrides in one file.

### Decision: Cars / appointments migration order

**Choice**: `AppShell.module.css` first, then `appointments` (smaller surface of status colours), then `cars` (largest hex footprint).

**Alternatives**: Cars first.

**Rationale**: Validates shell spec early; appointments exercises semantic status mapping before the heaviest file.

## Data flow

No new server or client data. Styling cascade only:

```text
localStorage ('gonsgarage-theme')
       → theme.ts / inline script sets html[data-theme]
       → tokens.css switches variable values
       → globals.css + utilities.css + *.module.css consume var(--*)
```

## File changes

| File | Action | Description |
|------|--------|-------------|
| `frontend/src/styles/tokens.css` | Modify | Brand comment block; tweak `--brand-*` if logo audit differs. |
| `frontend/src/styles/utilities.css` | Modify | Optional `.badge*` / status helpers if duplication remains after token swap. |
| `frontend/src/components/layout/AppShell.module.css` | Modify | Logout hover/nav: tokens only. |
| `frontend/src/app/landing.module.css` | Modify | Only if audit finds literals conflicting with spec (likely minimal). |
| `frontend/src/app/cars/*.module.css` | Modify | Replace hardcoded palette with tokens/utilities. |
| `frontend/src/app/appointments/*.module.css` | Modify | Same for status surfaces. |

## Interfaces / contracts

None. No changes to TypeScript public APIs, Zustand stores, or REST contracts.

## Testing strategy

| Layer | What | Approach |
|-------|------|----------|
| Manual | Light/dark toggle on shell, cars, appointments | Smoke after CSS edits. |
| Regression | Lint/build | `pnpm lint`, `pnpm build` in `frontend` (spec gate). |
| Unit | N/A for pure CSS | No new Vitest unless a component gains logic. |

## Migration / rollout

No migration, flags, or backend rollout. Revert = git revert of CSS commits.

## Open questions

- [ ] Confirm target browsers for `color-mix` (already used on landing) if expanding use.
- [ ] Whether to add `.badgeSuccess` / `.badgeWarning` utilities or only tokenize in-place (decide during apply based on duplication count).

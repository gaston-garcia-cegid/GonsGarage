# Apply progress — auth-ui-unified

**Mode**: Strict TDD (per `openspec/config.yaml`).

## TDD Cycle Evidence

| Task | Test file | Layer | Safety Net | RED | GREEN | TRIANGULATE | REFACTOR |
|------|-----------|-------|------------|-----|-------|---------------|----------|
| 1.1–1.3 | `AuthShell.test.tsx` | Integration | Vitest baseline 1/1 | Tests import missing module → fail | `AuthShell.tsx` + module CSS | Multiple cases: title/subtitle, success `status`, error `alert`, logo `src` | Readonly props, banner `if` clarity |
| 2.1–2.4 | `LoginForm.test.tsx` | Integration | Same | Expect new shell copy vs old UI → fail | `LoginForm` + `login.module.css` | Query param success + submit flow | Negated `if` / banner block cleanup |
| 3.1–3.3 | `page.test.tsx` (register) | Integration | Same | N/A (added after GREEN for regression) | Register uses `@/stores` + `AuthShell` | Heading + footer button | `UserRole` cast for `register()` |
| 4.1 | — | — | — | — | `pnpm test`, `pnpm typecheck` | — | — |
| 4.2 | (covered by 1.x) | — | — | — | — | — | — |
| 4.3 | Manual | — | — | — | Documented for verify | — | — |

## Test summary

- **Vitest**: 11 tests passing (`auth.client-role`, `AuthShell`, `LoginForm`, `RegisterPage`).
- **Lint**: `pnpm lint` — 0 errors (pre-existing warnings elsewhere).
- **Typecheck**: `pnpm typecheck` — pass after `UserRole` cast on register payload.

## Deviations from design

- **Button label on register**: children fixed to `Criar conta`; loading state handled by `Button` spinner (no duplicate text in children). Matches `Button` API.

## Manual checks (4.3)

Perform during `sdd-verify`: login → dashboard `replace`; register → `/auth/login?message=…`; invalid fields on both routes; visual parity at one viewport width.

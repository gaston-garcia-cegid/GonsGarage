# Tasks: Cópia em português na página de registo

## Phase 1: Tests (Strict TDD)

- [x] 1.1 Run baseline: `pnpm vitest run src/app/auth/register/page.test.tsx` (safety net).
- [x] 1.2 RED: add tests for pt e-mail label/placeholder, confirm label, and banner text without English `Error:` prefix when `register` rejects with `Error`.
- [x] 1.3 GREEN: update `frontend/src/app/auth/register/page.tsx` copy and catch handling until tests pass.
- [x] 1.4 REFACTOR: only if needed; keep tests green.

## Phase 2: Verification

- [x] 2.1 `pnpm lint` and `pnpm typecheck` in `frontend`.

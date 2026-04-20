# Tasks: Unificar layout e UX das superfícies de login e registo

## Phase 1: AuthShell (foundation)

- [x] 1.1 Create `frontend/src/components/auth/AuthShell.module.css`: outer page (`--surface-page`), centred column, card (`--surface-panel`, `--border-default`, `--shadow-lg`, `--radius-lg`), header spacing, `.bannerSuccess` / `.bannerError` using tokens (no raw hex for roles covered by tokens).
- [x] 1.2 Create `frontend/src/components/auth/AuthShell.tsx` (`'use client'`): props `title`, `subtitle?`, `banner?`, `children`; render `next/image` logo `/images/LogoGonsGarage.jpg` (32×32 or design token size) in header; map `banner` to success/error classes.
- [x] 1.3 Add optional shared class for secondary text link row (cross-auth navigation) in the shell module if used inside `children`, or export a small `AuthShellFooter` fragment in the same file — keep one module for shell chrome.

## Phase 2: Login

- [x] 2.1 Create `frontend/src/app/auth/login/login.module.css`: form-only layout (field stack, demo account block, spacing); reuse shell tokens for any local surfaces.
- [x] 2.2 Refactor `frontend/src/app/auth/login/LoginForm.tsx`: wrap content in `<AuthShell title="…" subtitle="…" banner={…} />` where `banner` derives from `successMessage` or `errors.general`; remove large inline `style={{}}` blocks.
- [x] 2.3 Replace native inputs/button with `Input` and `Button` from `frontend/src/components/ui/` where straightforward; keep `showPassword` toggle accessible (aria-labels).
- [x] 2.4 Style cross-link “Criar conta” using the same secondary/link pattern as register’s “Iniciar sessão” (shared class from shell module or `login.module.css` importing composition).

## Phase 3: Register

- [x] 3.1 In `frontend/src/app/auth/register/page.tsx`: change `useAuth` import from `@/contexts/AuthContext` to `@/stores`; verify `register()` payload still matches `RegisterRequest`.
- [x] 3.2 Wrap register UI in `<AuthShell title="Criar conta" subtitle="…" banner={errors.general ? { variant: 'error', message } : null}>`; move form body only inside `children` (drop outer `.pageContainer`/`.formContainer` wrappers from JSX).
- [x] 3.3 Edit `frontend/src/app/auth/register/register.module.css`: remove rules duplicated by `AuthShell.module.css` (page outer, card shell, header block if fully in shell); keep `.nameFieldsRow`, password rows, `Input`/`.submitButton` overrides.

## Phase 4: Verification

- [x] 4.1 Run `cd frontend && pnpm lint` and `pnpm typecheck`; fix any regressions.
- [x] 4.2 (Optional) Add `frontend/src/components/auth/AuthShell.test.tsx`: render with title + children; assert heading text (covers spec “Shared page shell” smoke).
- [x] 4.3 Manual: login success → `router.replace('/dashboard')` unchanged; register success → redirect to `/auth/login?message=…` and banner shows with success styling; invalid fields show errors on both routes; compare login vs register layout at same viewport (spec structural parity).

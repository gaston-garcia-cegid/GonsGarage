# Verify — auth-register-pt-copy

## Spec scenarios

| Requirement / Scenario | Evidence |
|--------------------------|----------|
| Email field matches login wording | `page.test.tsx` — label `E-mail`, placeholder `O seu e-mail` |
| Confirm password label in Portuguese | `page.test.tsx` — `getByLabelText('Confirmar palavra-passe')` |
| No English `Error:` prefix | `page.test.tsx` — `register` rejeita com `Error`; alert sem prefixo `Error:` |

## Commands

- `pnpm vitest run src/app/auth/register/page.test.tsx` — passou
- `pnpm lint` — 0 erros (avisos pré-existentes noutros ficheiros)
- `pnpm typecheck` — passou

## Verdict

Implementação alinhada ao delta `specs/client-auth-shell/spec.md` desta change.

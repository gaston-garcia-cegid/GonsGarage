# Design: auth-register-pt-copy

## Decisions

1. **Copy source of truth**: Reuse the same strings as `LoginForm.tsx` for the e-mail field (`E-mail`, `O seu e-mail`).
2. **Confirm password label**: `Confirmar palavra-passe`, consistent with validation copy already using “palavra-passe”.
3. **Catch block**: Stop prefixing with English `Error:`; surface `Error.message` alone or fall back to the existing generic Portuguese sentence so the shell banner stays pt-first.

## Non-goals

- No new i18n layer, no shared string module unless a later change consolidates auth copy.

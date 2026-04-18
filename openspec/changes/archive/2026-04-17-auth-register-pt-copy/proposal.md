# Proposal: Cópia em português na página de registo

## Intent

A rota `/auth/register` mistura **inglês** (rótulo de e-mail, placeholder, “Confirm Password”) com o resto da UI já em **português europeu**. Isso quebra a coerência com `/auth/login` e com mensagens de validação já em pt. Corrigir só o que o utilizador vê nesta página, sem introduzir i18n.

## Scope

### In Scope

- Substituir rótulos e placeholders em inglês no formulário de registo (`page.tsx`) por pt alinhado ao login (ex.: “E-mail”, “O seu e-mail”, “Confirmar palavra-passe”).
- Rever **texto visível** no mesmo ficheiro para inglês residual (ex.: prefixo `Error:` no ramo `catch` → mensagem só em pt ou sem prefixo em inglês).

### Out of Scope

- Biblioteca de traduções ou ficheiros de locale; outras rotas; mensagens devolvidas pela API em inglês (tratamento genérico mantém-se salvo o prefixo local).

## Capabilities

`openspec/specs/` na raiz está vazio (só `.gitkeep`). Existe delta **`client-auth-shell`** em `openspec/changes/auth-ui-unified/specs/client-auth-shell/spec.md` que já cobre shell e consistência login/registo.

### New Capabilities

- None

### Modified Capabilities

- `client-auth-shell`: acrescentar requisito explícito de que **cópia visível** do formulário de registo (rótulos, placeholders, prefixos de erro construídos no cliente nesta página) está em **português europeu** e alinhada terminologicamente ao login (e-mail, palavra-passe).

## Approach

1. Editar `frontend/src/app/auth/register/page.tsx`: strings hardcoded nos componentes `Input`/labels e no `catch`.
2. Opcional: alinhar capitalização (“E-mail” vs “Email”) ao que já existe no login — conferir `LoginForm` ou `AuthShell` antes de fixar texto final.
3. Se existir teste de registo que asserta texto em inglês, atualizar expectativas.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `frontend/src/app/auth/register/page.tsx` | Modified | Cópia pt; possível ajuste de mensagem de erro genérica. |
| `frontend/src/app/auth/register/*.test.*` | Modified (se existir) | Expectativas de texto. |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Testes E2E ou unitários a procurar strings antigas em inglês | Low | Correr `pnpm test` / lint no `frontend` após implementação. |

## Rollback Plan

Reverter o commit que altere só `page.tsx` (e testes); sem migrações nem alteração de API.

## Dependencies

- Nenhuma.

## Success Criteria

- [ ] Não há inglês visível nos rótulos/placeholders do formulário de registo (campos cobertos por esta alteração).
- [ ] Mensagem de exceção apresentada ao utilizador nesta página não começa por “Error:” em inglês.
- [ ] `pnpm lint` e `pnpm typecheck` no frontend passam após a alteração.

# Proposal: Unificar layout e UX das superfícies de login e registo

## Intent

As rotas `/auth/login` e `/auth/register` parecem produtos diferentes: login em `LoginForm.tsx` com estilos inline e ícone genérico; registo com `register.module.css` e componentes `Input`/`Button`. O utilizador espera a mesma hierarquia visual, tipografia e padrões de formulário que no resto da app (tokens, shell minimalista). Alinhar também a fonte de estado de auth: login usa `@/stores`, registo usa `@/contexts/AuthContext`, o que aumenta risco de comportamentos divergentes.

## Scope

### In Scope

- Shell partilhado (layout centrado, fundo, cartão, cabeçalho com marca) para login e registo, reutilizando `tokens.css` / utilitários existentes.
- Refatorar `LoginForm` para CSS modules (ou componentes partilhados) em par com registo; mensagens de sucesso/erro com o mesmo padrão visual.
- Unificar `useAuth` para o mesmo import que o resto das páginas autenticadas (`@/stores`), mantendo `AuthProvider` em `layout.tsx` se ainda for necessário para hidratação.
- Links cruzados “Criar conta” / “Já tem conta?” com estilo coerente com botões secundários da landing.

### Out of Scope

- Transformar login/registo em modais sobre a landing (nova UX de navegação); i18n além do pt actual.
- Alterar contratos da API de auth ou regras de validação de negócio (só apresentação e wiring de store).

## Capabilities

`openspec/specs/` só contém `.gitkeep`; não há capabilities principais a modificar.

### New Capabilities

- `client-auth-shell`: Layout e padrões visuais partilhados para fluxos de autenticação no cliente (login, registo, mensagens inline).

### Modified Capabilities

- None

## Approach

1. Extrair `AuthShell` (ou `AuthPageLayout`) com wrapper de página, título opcional, área de alertas e `children` para o formulário.
2. Migrar estilos inline de `LoginForm.tsx` para `login.module.css` (ou módulo partilhado `auth-shell.module.css`) alinhado a `register.module.css`.
3. Opcional: usar `Input`/`Button` no login para paridade com registo.
4. Trocar import de `useAuth` em `register/page.tsx` para `@/stores` e validar redirect pós-registo/login com o mesmo fluxo que `LoginForm` (`router.replace` onde fizer sentido).

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `frontend/src/app/auth/login/LoginForm.tsx` | Modified | Remover inline styles; encaixar no shell partilhado. |
| `frontend/src/app/auth/login/` | New | CSS module (ou partilhado com registo). |
| `frontend/src/app/auth/register/page.tsx` | Modified | Shell comum; alinhar `useAuth`. |
| `frontend/src/app/auth/register/register.module.css` | Modified | Fundir ou modularizar tokens comuns com login. |
| `frontend/src/components/...` (novo) | New | Layout/superfície auth reutilizável. |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Regressão no redirect ou hidratação ao mudar `useAuth` | Med | Testar manualmente login, registo e primeira carga do dashboard; `pnpm test` onde houver cobertura. |
| CSS duplicado entre módulos | Low | Um único `auth-shell.module.css` importado por ambas as páginas. |

## Rollback Plan

Revert do commit (ou PR) que introduza o shell e os imports; restaurar `LoginForm.tsx` e `register/page.tsx` da revisão anterior. Sem migrações de dados.

## Dependencies

- Nenhuma externa; depende dos design tokens já em `frontend/src/styles/tokens.css`.

## Success Criteria

- [ ] Login e registo partilham o mesmo layout de página (marca, cartão, espaçamentos, sombras/radius alinhados à app).
- [ ] Não há estilos inline “grandes” no login salvo exceções justificadas (ex. animação pontual).
- [ ] Ambos os fluxos usam `useAuth` de `@/stores` (ou documentado equivalente único).
- [ ] `pnpm lint` e `pnpm typecheck` no `frontend` passam.

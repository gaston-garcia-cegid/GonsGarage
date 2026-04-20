# Design: Unificar layout e UX das superfícies de login e registo

## Technical Approach

Introduzir um **`AuthShell`** cliente (`'use client'`) que concentra o layout de página já usado semanticamente em `register.module.css` (`pageContainer`, `formContainer`, `header`, tipografia). **Login** deixa de usar estilos inline e passa a renderizar o mesmo shell + classes partilhadas. **Registo** passa a compor o shell partilhado (em vez de duplicar `.pageContainer`/`.formContainer` no seu módulo) e troca `useAuth` de `@/contexts/AuthContext` para `@/stores`, alinhado ao `LoginForm` e ao resto das rotas autenticadas.

Alertas (sucesso no login vindo da query, erro geral, erros de campo) usam **classes partilhadas** no módulo do shell (variantes success/error) mapeadas a tokens (`--color-error`, superfícies semânticas), eliminando hex soltos no login.

## Architecture Decisions

| Decision | Alternatives | Choice & rationale |
|----------|--------------|-------------------|
| Onde vive o shell | Colocar só em `app/auth/` | **`frontend/src/components/auth/`** — reutilizável, testável, padrão do repo (`components/`). |
| CSS | CSS-in-JS / Tailwind | **CSS Modules** (`AuthShell.module.css`) — igual ao registo actual; login migra para o mesmo padrão. |
| Fonte de auth nas páginas auth | Manter Context no registo | **`useAuth` de `@/stores`** — contrato do spec; store já expõe `login`/`register` com a mesma forma mental que o resto da app. |
| `AuthProvider` no `layout.tsx` | Remover já neste PR | **Manter neste PR** — `employees` e outros ainda importam Context; remoção é **follow-up** para evitar regressões fora do âmbito auth-shell. |
| Controlo primário no login | Manter `<button>` nativo | **Preferir `Button` + `Input`** no login — paridade com registo e menos CSS manual (opcional se o esforço for alto: manter um único botão estilizado via classes partilhadas). |
| Marca no cabeçalho | SVG genérico | **`next/image` + `/images/LogoGonsGarage.jpg`** no shell — coerência com landing e spec de marca. |

## Data Flow

```
┌─────────────────────────────────────────────────────────┐
│  RootLayout: AuthProvider (Context) — inalterado        │
└──────────────────────────┬──────────────────────────────┘
                           │
     ┌─────────────────────┴─────────────────────┐
     ▼                                           ▼
 /auth/login                              /auth/register
 LoginForm                                RegisterPage
     │                                           │
     └──────────► useAuth() @/stores ◄──────────┘
                           │
                    login / register
                           │
                    apiClient (store)
```

Ambas as páginas leem/escrevem o **mesmo** estado persistido do Zustand após acções; o Context continua a existir para código legado até migração futura.

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `frontend/src/components/auth/AuthShell.tsx` | Create | Props: `title`, `subtitle`, `children`; opcional slot para alertas acima do cartão. |
| `frontend/src/components/auth/AuthShell.module.css` | Create | Shell + variantes de alerta + link secundário; tokens como `--surface-page`. |
| `frontend/src/app/auth/login/LoginForm.tsx` | Modify | Compor `AuthShell`; remover inline styles; opcional `Input`/`Button`. |
| `frontend/src/app/auth/login/login.module.css` | Create | Estilos específicos do formulário de login (campos, demo copy). |
| `frontend/src/app/auth/register/page.tsx` | Modify | `useAuth` → `@/stores`; envolver conteúdo em `AuthShell`; remover duplicação de layout do CSS local. |
| `frontend/src/app/auth/register/register.module.css` | Modify | Retirar classes movidas para `AuthShell.module.css`; manter grelhas/campos específicos. |

## Interfaces / Contracts

```tsx
// AuthShell — superfície mínima (exact props a afinar na implementação)
type AuthShellProps = {
  title: string;
  subtitle?: string;
  /** Alerta global opcional (ex.: sucesso vindo da query) */
  banner?: { variant: 'success' | 'error'; message: string } | null;
  children: React.ReactNode;
};
```

`register()` no store aceita `RegisterRequest` (`@/types`); o payload actual do registo já é compatível — **não** alterar o contrato da API.

## Testing Strategy

| Layer | What | Approach |
|-------|------|----------|
| Manual | Login, registo, redirect com `?message=` | Smoke após refactor. |
| Unit | `AuthShell` renderiza título + children | Vitest + Testing Library se existir teste semelhante; caso contrário, adicionar um teste mínimo opcional. |
| E2E | — | Não requerido neste change. |

## Migration / Rollout

No migration required. Deploy contínuo; rollback = revert do commit.

## Open Questions

- [ ] Migrar `employees` (e restantes consumidores de `AuthContext`) para `@/stores` e remover `AuthProvider` — **fora deste change**, documentar como debt técnico.

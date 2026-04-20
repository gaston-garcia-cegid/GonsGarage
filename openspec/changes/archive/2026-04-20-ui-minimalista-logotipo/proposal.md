# Proposal: UI minimalista alineada ao logótipo

## Intent

Unificar a UI para um aspeto **profissional e minimalista**, com **cores derivadas do logótipo** GonsGarage, reduzindo ruído visual e **cores hardcoded** que quebram o tema claro/escuro.

## Scope

### In Scope

- Rever `tokens.css` face ao asset do logo (ajustes documentados de `--brand-*` se necessário).
- Harmonizar **AppShell** (incl. hover do logout) e **landing** com tokens.
- Migrar os **CSS modules com mais hex soltos** (prioridade: `cars`, `appointments`) para variáveis ou utilitários partilhados.
- Critérios de sucesso mensuráveis: menos hex fora de tokens nos módulos alvo; `pnpm build` e lint sem regressão.

### Out of Scope

- Redesign completo de fluxos ou novos componentes de biblioteca (ex. shadcn).
- i18n (`es_ES` / `en_GB`).
- Testes E2E visuais / Percy.

## Capabilities

### New Capabilities

- `ui-brand-shell`: Comportamento e tokens da shell autenticada, navegação e uso de cor de marca/accento/sinal.

### Modified Capabilities

- None *(ainda não existem `openspec/specs/*/spec.md` no repositório; deltas futuros podem referir `ui-brand-shell`.)*

## Approach

1. Amostragem visual do JPG do logo → atualizar comentários e valores em `tokens.css` se preciso.
2. Substituir hex em `AppShell.module.css` e módulos prioritários por `var(--*)` ou classes em `utilities.css`.
3. Passo opcional: extrair padrões repetidos (badges de estado) para utilitários reutilizáveis.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `frontend/src/styles/tokens.css` | Modified | Paleta e superfícies alinhadas à marca |
| `frontend/src/styles/utilities.css` | Modified | Utilitários de estado se fizer sentido |
| `frontend/src/components/layout/AppShell.module.css` | Modified | Logout/nav sem hex solto |
| `frontend/src/app/cars/*.module.css` | Modified | Tokens em vez de paleta Tailwind-like solta |
| `frontend/src/app/appointments/*.module.css` | Modified | Idem |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Regressão de contraste WCAG | Med | Rever focus rings e estados hover com tema escuro |
| Ficheiros grandes difíceis de rever | Baixa | Mudanças por commits pequenos por rota |

## Rollback Plan

Revert do commit (ou série) que toque em tokens/CSS; não há migração de dados.

## Dependencies

- Nenhuma externa; apenas build frontend.

## Success Criteria

- [ ] Paleta documentada em `tokens.css` coerente com o logo (comentário ou secção “Brand”).
- [ ] Zero ou mínimo de hex **não token** nos módulos `cars` e `appointments` acordados.
- [ ] `pnpm lint` e `pnpm build` no `frontend` passam.

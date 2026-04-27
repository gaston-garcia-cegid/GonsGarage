# Delta for ui-brand-shell

## ADDED Requirements

### Requirement: Production stack on main (Next.js 16 + Tailwind v4)

Tras este change, a linha **oficial** em `main` **SHALL** incluir **Next.js 16.x** (ou superior acordado) e **Tailwind CSS v4** como ferramentas de build de estilos do frontend; **MUST NOT** ficar o frontend preso a Next 15.5 + Tailwind v3 só por omissão pós-merge. O produto **SHALL** manter **ADR GO** (ou documento equivalente no change) com hipótese de benefício, passos de migração e resultados de `pnpm build`, `pnpm test`, `pnpm lint` e **CI** alinhados ao `openspec/config.yaml`.

#### Scenario: CI and local quality gate

- **GIVEN** o merge do change concluído em `main`
- **WHEN** corren os comandos de qualidade frontend acordados no CI (lint, typecheck, test, build)
- **THEN** **SHALL** completar sem erros introduzidos por esta migración
- **AND** o ADR **GO** **SHALL** estar ligado ou referenciado no repositório

#### Scenario: Theme and shell without functional regression

- **GIVEN** tema claro e escuro e rotas MVP en escopo do checklist do change
- **WHEN** un usuario percorre shell, navegación primaria e vistas tocadas
- **THEN** contraste e legibilidade **SHALL** permanecer alinhados a tokens documentados
- **AND** **MUST NOT** introducirse regresión funcional documentada como bloqueante no `verify-report` do change

### Requirement: Documentación de cambios importantes na UI

Cada alteración **estructural** na capa de presentación (config Tailwind v4, tokens globais, `next.config`, pipeline CSS) que afecte o comportamento observábel **SHALL** aparecer na tabla “antes/después” do ADR ou `design.md` do change, con unha liña de motivación testábel.

#### Scenario: Maintainer finds rationale

- **GIVEN** un mantenedor revisa un diff non obvio (p. ex. renomeo de utilidades Tailwind)
- **WHEN** abre o ADR ou design do change
- **THEN** **SHALL** atopar entrada que ligue o cambio a requisito ou risco mitigado

## MODIFIED Requirements

### Requirement: Upgrade spike decision record

Before adopting **Next.js 16** and **Tailwind CSS v4** on `main`, the project **SHALL** keep a document (ADR in `docs/` or change section) with: benefit hypothesis, branches tried, `pnpm build` + `pnpm test` outcomes, and **GO** / **NO-GO** / **DEFER** decision. When the ADR records **GO** for this change, **`package.json` on `main` SHALL list Next.js 16.x and Tailwind CSS v4** as the adopted dependency versions (not only an optional spike without merging those dependencies).

(Previously: the document required a decision and steps but did not explicitly tie GO to mandatory versions on `main`.)

##### Scenario: NO-GO preserves current stack

- **GIVEN** a NO-GO outcome
- **WHEN** the phase is archived
- **THEN** `package.json` on `main` stays on Next 15.x without mandatory Tailwind v4
- **AND** the ADR explains why

##### Scenario: GO documents migration steps

- **GIVEN** a GO outcome
- **WHEN** the spike merges
- **THEN** the ADR lists contributor migration steps (commands, known breaking changes)
- **AND** `package.json` on `main` **SHALL** list Next.js 16.x and Tailwind CSS v4 as adopted dependency versions

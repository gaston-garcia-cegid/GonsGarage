# ui-component-system Specification

## Purpose

Define the **greenfield** UI layer for the MVP using a **Shadcn-style** system (copy-paste components, Radix primitives, theme tokens) as the **canonical** source for new UI and migrations, replacing ad hoc patterns in a phased way.

## Requirements

### Requirement: Greenfield Shadcn foundation

The frontend **SHALL** incorporate the stack-recommended base (e.g. **shadcn/ui** + Tailwind, with version per ADR) in a **dedicated directory** (`frontend/src/components/ui/` or agreed path), without mixing new primitives with legacy styles in the same source file.

#### Scenario: Foundation PR is self-contained

- **GIVEN** the first Phase-4 foundation PR
- **WHEN** the diff is reviewed
- **THEN** it introduces configuration, dependencies, and base primitives **without** rewriting all MVP routes yet
- **AND** `pnpm build` and `pnpm test` remain green

### Requirement: MVP route migration to system components

Each authenticated MVP view **SHALL** migrate to the new system’s components (buttons, inputs, dialogs, tables) per `tasks.md` in the owning change; “greenfield complete” **SHALL NOT** be claimed while critical forms remain only in legacy CSS modules without a tracked replacement plan.

#### Scenario: Auth shell on new system

- **GIVEN** an agreed Phase-4 milestone in tasks
- **WHEN** a user opens login or register
- **THEN** the UI uses system primitives exclusively for fields and primary actions **OR** verify documents a time-boxed exception

### Requirement: Theme parity with brand

The Shadcn theme **SHALL** map GonsGarage brand colours (navy, accent, signal) to stack theme variables so light/dark stay coherent with `ui-brand-shell` after merge.

#### Scenario: Token mapping documented

- **GIVEN** Phase 4 is merged to `main`
- **WHEN** a maintainer opens the theme guide (`docs/ui-shadcn-theme.md` or successor)
- **THEN** a table or commentary links brand tokens ↔ component-system variables

---

## Merged additions (change `ui-homogeneity-modal-workshop-parts`, archived 2026-04-27)

### Requirement: Primitivas do sistema para criação em peças e taller

Os fluxos de **criação** na UI de inventário de peças e na **lista taller** (incluindo confirmação antes ou depois da mutação, conforme desenho) **SHALL** usar primitivas do **sistema canónico** (`Dialog`, `Button`, `Input` ou equivalentes documentados em `frontend/src/components/ui/`) para campos e acções primárias **novas ou migradas** por esta mudança; **MUST NOT** acrescentar formulários de criação apenas com marcação legacy na mesma alteração sem plano de substituição referenciado no relatório de verificação.

#### Scenario: Superfície de nova peça

- **GIVEN** o fluxo de criação de peça a partir da lista
- **WHEN** o diálogo ou painel de criação está aberto
- **THEN** campos editáveis e botões de submeter/cancelar **SHALL** ser componentes do sistema canónico

#### Scenario: Fluxo de nova visita na lista taller

- **GIVEN** staff com permissão de taller na lista de visitas
- **WHEN** inicia **Nova visita**
- **THEN** qualquer confirmação ou captura mínima obrigatória antes de navegar **SHALL** usar o mesmo conjunto de primitivas do sistema para acções e campos expostos

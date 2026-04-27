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

---

## Merged additions (change `nextjs-16-react19-migration`, archived 2026-04-27)

### Requirement: Server and client boundaries for App Router

Novas ou migradas vistas en `frontend/src/app/` **SHALL** colocar **fetch inicial** e datos derivábeis sen estado de cliente en **Server Components** cando non dependan de APIs só dispoñíbeis no navegador (p. ex. `localStorage`). Estado interactivo, Zustand, eventos DOM, modais Radix e integracións que requiran `window` **SHALL** vivir en **ilhas** `"use client"` explícitas. **MUST NOT** usar `useEffect` só para cargar datos que poidan resolverse no servidor sen perder paridade funcional.

#### Scenario: Data loads without spurious client effect

- **GIVEN** unha vista cuxo listado inicial poida obterse no servidor co mesmo contrato de datos ca hoxe
- **WHEN** o usuario abre a ruta
- **THEN** o HTML inicial **SHALL** incluir os datos ou loading acordado sen depender dun `useEffect` exclusivamente para ese fetch
- **AND** as accións interactivas posteriores **MAY** seguir en cliente

#### Scenario: Client island for store-driven UI

- **GIVEN** unha vista que consume Zustand ou modais
- **WHEN** inspecciona o árbore de compoñentes
- **THEN** o módulo que subscribe o store **SHALL** estar baixo `"use client"` ou importado só dende client boundaries

### Requirement: Zustand store refactor contract

O **mega-refactor** de stores en `frontend/src/stores/**` **SHALL** manter **semántica observable** equivalente: mesmos fluxos de login, datos de usuario, dominios expostos aos compoñentes e chamadas HTTP desde o cliente onde aplique, salvo **documentación explícita** de cambio intencional no `design.md`. Stores **SHALL** ser compatíbeis con hidratación e límites RSC (sen asumir estado global no servidor); **MUST NOT** introducir duplicación de fonte de verdade sen reconciliación documentada.

#### Scenario: Auth flow still works

- **GIVEN** fluxos de login e sesión cubertos por tests ou checklist
- **WHEN** executan despois do refactor
- **THEN** o usuario **SHALL** autenticarse e acceder ás rutas permitidas como antes

#### Scenario: Store boundary documented

- **GIVEN** un store que deixa de ser importábel dende certos módulos
- **WHEN** un desenvolvedor le o design do change
- **THEN** **SHALL** atopar razón e patrón de substitución (p. ex. props dende server, context client)

### Requirement: React 19 APIs for async and optimistic UI

O frontend **MAY** usar **`use`** para ler **Promises** ou **Context** según React 19 dentro de client boundaries onde simplifique o fluxo respecto a `useEffect`+estado manual. **MAY** usar **`useOptimistic`** para actualizacións optimistas que **SHALL** reconciliarse coa resposta real da API ou revertir de xeito explícito; **MUST NOT** presentar como persistido o que só é optimista sen feedback de erro.

#### Scenario: Optimistic update reconciles

- **GIVEN** unha acción de mutación cuberta por `useOptimistic`
- **WHEN** a API responde éxito ou erro
- **THEN** a UI **SHALL** reflectir o estado final coherente coa resposta

#### Scenario: Documented use of use()

- **GIVEN** un fluxo que adopta `use()`
- **WHEN** o change arquívase
- **THEN** o design ou ADR **SHALL** mencionar o motivo fronte a alternativas

# Delta for ui-component-system

## ADDED Requirements

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

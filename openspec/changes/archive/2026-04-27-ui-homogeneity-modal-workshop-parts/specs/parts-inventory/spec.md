# Delta for parts-inventory

## ADDED Requirements

### Requirement: Ponto de entrada principal para criação de ítem

A criação de um novo ítem de inventário **SHALL** ser iniciada a partir do **contexto da lista** (mesma rota de listagem com modal, query, ou padrão equivalente documentado no produto). Uma rota dedicada só de criação **MAY** existir como **redireccionamento** ou compatibilidade de marcadores; **MUST NOT** ser o único caminho suportado após esta mudança.

#### Scenario: Alta desde a lista

- **GIVEN** `manager` ou `admin` com sessão válida
- **WHEN** escolhe criar novo ítem a partir da UI de inventário
- **THEN** o ponto de entrada primário **SHALL** ser a lista com fluxo de criação em sobreposição ou query na mesma vista

#### Scenario: Marcador legado para `/new`

- **GIVEN** um marcador ou hiperligação antiga para a rota dedicada de criação
- **WHEN** o utilizador a abre
- **THEN** o sistema **SHALL** encaminhar para o fluxo suportado na lista **sem** perda de capacidade de criar ítem

# Delta for ui-component-system

## ADDED Requirements

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

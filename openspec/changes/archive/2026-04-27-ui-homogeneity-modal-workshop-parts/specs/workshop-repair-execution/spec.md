# Delta for workshop-repair-execution

## ADDED Requirements

### Requirement: Leitura de visita perceptível na UI

Depois de navegar para o **detalhe de uma visita** (identificador válido na URL da app), o staff autorizado **SHALL** ver **ou** os dados mínimos da visita (estado, recepção/cierre conforme modelo) **ou** uma **mensagem de erro accionável** (p. ex. permissão, não encontrado, falha de rede). A UI **MUST NOT** ficar indefinidamente num estado vazio sem texto ou indicador de carregamento quando a solicitação terminou.

#### Scenario: Carregamento com sucesso

- **GIVEN** visita existente e `GET` de visita devolve 2xx com corpo válido
- **WHEN** o staff abre o detalhe dessa visita na app
- **THEN** o estado da visita **SHALL** tornar-se legível (texto ou componente equivalente) após o fim do carregamento

#### Scenario: Falha ou negação

- **GIVEN** resposta não-2xx ou corpo inválido para o `GET` de visita
- **WHEN** o staff permanece no detalhe
- **THEN** a UI **SHALL** mostrar erro interpretável ou orientação para voltar; **MUST NOT** mostrar apenas área em branco como se não houvesse visita

### Requirement: Nova visita alinhada ao contexto da lista

A acção **Nova visita** desde a **lista taller** **SHALL** seguir um padrão **lista-primário**: confirmação ou passo breve na lista (modal ou painel na mesma rota) **SHOULD** preceder ou acompanhar a criação antes de depender só de navegação para detalhe; **MUST NOT** ser apenas um salto sem feedback coerente com outras áreas staff que criam recursos a partir de lista.

#### Scenario: Criação com viatura seleccionada

- **GIVEN** viatura seleccionada na lista e staff com permissão
- **WHEN** confirma **Nova visita**
- **THEN** o utilizador **SHALL** receber feedback de progresso ou resultado antes de ficar só num ecrã de detalhe sem contexto

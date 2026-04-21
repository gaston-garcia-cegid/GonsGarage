# ui-accounting-staff Specification

> **Source:** promovido desde `openspec/changes/archive/2026-04-21-ui-accounting-modal-create-flows/` (2026-04-21).

## Purpose

Definir a **UX staff** da área **Contabilidade** (`/accounting` e as quatro listas operativas): criação de registos **no contexto da listagem** mediante **modal**, alinhada ao comportamento de outras áreas staff autenticadas que já usam criação in-context (sem alterar regras de domínio P1 em `suppliers`, `invoices` ou `billing`).

## Requirements

### Requirement: Criação em modal nas quatro listas

Nas vistas de listagem staff de **Fornecedores**, **Faturas recebidas**, **Documentos de faturação** e **Faturas emitidas**, o fluxo primário de **criar** um novo registo **SHALL** abrir um **modal** sobre a lista; o utilizador **MUST NOT** depender de uma página dedicada `*/new` como único caminho feliz.

#### Scenario: Criar a partir da toolbar

- GIVEN um utilizador staff na listagem de um dos quatro domínios
- WHEN escolhe a acção explícita de criação (ex.: «Novo …») na toolbar
- THEN abre-se um modal com o formulário de criação
- AND a lista permanece por baixo do overlay (contexto preservado)

#### Scenario: Criar a partir do estado vazio

- GIVEN a lista sem linhas e sem erro de carga
- WHEN o utilizador usa o CTA de criação no estado vazio
- THEN o mesmo modal de criação **SHALL** abrir (equivalente ao da toolbar)

#### Scenario: Sucesso fecha e actualiza a lista

- GIVEN o modal de criação aberto com dados válidos
- WHEN a criação completa com sucesso no backend
- THEN o modal **SHALL** fechar
- AND a listagem **SHALL** reflectir o novo registo sem navegação para `*/new`

#### Scenario: Cancelar sem persistir

- GIVEN o modal de criação aberto
- WHEN o utilizador cancela ou fecha sem confirmar criação
- THEN nenhum registo novo **SHALL** aparecer na lista
- AND o utilizador permanece na mesma rota de listagem

### Requirement: Coerência com padrões de modal staff

Comportamento de **abrir**, **fechar**, **acção primária** vs **cancelar** e tratamento de **erro de submissão** nestes modais **SHALL** ser observavelmente coerente com outras vistas staff do mesmo shell que já usam criação em modal (mesma expectativa de dismiss e foco, sem especificar biblioteca).

#### Scenario: Fecho por teclado ou controlo explícito

- GIVEN o modal de criação aberto e foco no diálogo
- WHEN o utilizador activa fecho por teclado (ex.: ESC) ou botão de fechar explícito
- THEN o modal **SHALL** fechar sem criar registo (salvo fluxo que exija confirmação adicional já existente noutras áreas — nesse caso **SHALL** seguir o mesmo padrão global)

### Requirement: Rotas legacy `*/new`

Qualquer rota `*/new` sob `accounting` que exista por compatibilidade **SHALL** redireccionar para a listagem correspondente **ou** apresentar equivalente sem UX morta (navegação clara de volta à lista); **MUST NOT** deixar o utilizador preso num formulário full-page como único destino sem alternativa.

#### Scenario: Atalho ou bookmark para `*/new`

- GIVEN um pedido directo a uma URL `…/accounting/…/new`
- WHEN a página resolve
- THEN o utilizador **SHALL** poder concluir ou abandonar a criação sem estado inconsistente
- AND **SHALL** poder regressar ao contexto da lista (redirect ou fluxo equivalente documentado em deseño)

### Requirement: Língua da interface

Rótulos e mensagens destes fluxos (títulos do modal, botões primários/secundários, erros inline) **SHALL** estar em **português (pt_PT)**, coerente com o resto da área Contabilidade.

#### Scenario: Acções visíveis em pt_PT

- GIVEN o modal de criação aberto
- WHEN o utilizador vê título e botões de acção
- THEN o texto **SHALL** estar em pt_PT (sem mistura desnecessária com outras línguas)

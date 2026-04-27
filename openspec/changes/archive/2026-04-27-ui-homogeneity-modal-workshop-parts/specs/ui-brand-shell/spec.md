# Delta for ui-brand-shell

## ADDED Requirements

### Requirement: Staff inventário e taller alinhados ao shell homogéneo

As vistas autenticadas de **inventário de peças** (`admin` / `manager`) e **taller** (staff de oficina) **SHALL** seguir o mesmo padrão de **AppShell + toolbar** que áreas staff já homogéneas: **criação** de recurso primário **SHALL NOT** ser a única acção disponível como página isolada sem continuidade com a lista, quando a lista é o contexto principal de trabalho.

#### Scenario: Criar peça desde o contexto da lista

- **GIVEN** um utilizador `admin` ou `manager` na lista de peças
- **WHEN** inicia **Nova peça**
- **THEN** o fluxo de criação **SHALL** manter a coerência com rotas que usam **sobreposição** (modal ou query na mesma rota) alinhada a accounting/appointments
- **AND** o utilizador **SHALL** conseguir cancelar e regressar à lista sem perda de contexto de navegação principal

#### Scenario: Tema e tokens nas superfícies tocadas

- **GIVEN** tema claro ou escuro activo
- **WHEN** o utilizador percorre formulários ou modais de inventário ou taller alterados por esta mudança
- **THEN** fundos, bordas e texto de interacção primária **SHALL NOT** introduzir literais de cor que dupliquem o papel de tokens já definidos para o mesmo semântica

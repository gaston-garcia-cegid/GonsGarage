# parts-inventory Specification

> **Promoción:** catálogo principal desde el change archivado `openspec/changes/archive/2026-04-24-spare-parts-inventory/` (2026-04-24).

## Purpose

Stock repuestos: cantidad+UoM, barcode HID; roles `mvp-role-access`.

## Requirements

### Requirement: Manager and admin only

Solo `manager` y `admin` **SHALL** usar la UI y obtener 2xx en rutas HTTP de inventario (CRUD + lecturas). `employee` y `client` **MUST NOT** obtener 2xx ni ver enlace de navegación principal a la sección.

#### Scenario: Admin passes role gate

- GIVEN JWT `admin` y payload válido
- WHEN mutador inventario
- THEN **MUST NOT** `403` por solo rol

#### Scenario: Client mutator denied

- GIVEN JWT `client`
- WHEN mutador HTTP de inventario
- THEN **MUST NOT** ser 2xx

### Requirement: CRUD ítem

El sistema **SHALL** CRUD ítems con: id estable, referencia o código interno, marca, nombre, cantidad ≥0, UoM de conjunto cerrado (p. ej. `unit`, `liter`). **SHALL** permitir borrar o desactivar equivalente.

#### Scenario: Persist quantity UoM

- GIVEN `manager`
- WHEN crea con cantidad `3` y UoM litro
- THEN lecturas devuelven `3` y misma UoM

#### Scenario: Negative quantity

- GIVEN `manager`
- WHEN guarda cantidad negativa
- THEN **MUST NOT** 2xx

### Requirement: Barcode y búsqueda

**SHALL** campo barcode opcional. **SHALL** localizar ítem por código en flujo de búsqueda. Código desconocido **SHALL** poder pre-rellenar alta; **MUST NOT** persistir sin confirmación explícita.

#### Scenario: Find by code

- GIVEN ítem con barcode conocido o código nuevo sin persistir
- WHEN usuario introduce código en búsqueda o alta
- THEN UI **SHALL** resolver ítem existente o pre-rellenar sin persistir hasta confirmación

### Requirement: Barcode único

Dos ítems **MUST NOT** compartir barcode no vacío; duplicado **MUST** rechazarse con error claro.

#### Scenario: Duplicate assign

- GIVEN A con barcode `X`, B sin barcode
- WHEN asignar `X` a B
- THEN **MUST NOT** 2xx

### Requirement: Mínimo y tiempo (SHOULD)

**SHOULD** umbral mínimo con aviso en listado si cantidad por debajo; **SHOULD** avanzar `updated_at` en mutación válida.

#### Scenario: Aviso y timestamp

- GIVEN ítem mín `5` qty `2` y otro ítem
- WHEN listado y `manager` guarda válido en el otro
- THEN aviso bajo y `updated_at` avanza

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

<!-- Promoted from change `ui-homogeneity-modal-workshop-parts`, archived 2026-04-27 -->

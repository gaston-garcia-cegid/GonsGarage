# Delta for mvp-role-access

## MODIFIED Requirements

### Requirement: Published role–surface matrix

MVP verification docs **SHALL** expose the matrix above (or equivalent), including **staff user provisioning**, **workshop (service job)** (taller visit), and **parts inventory** (HTTP + UI `parts-inventory`) como filas distintas, alineadas con `workshop-repair-execution` y `parts-inventory`.

(Previously: sin fila explícita parts inventory.)

#### Scenario: Four roles listed

- GIVEN the linked checklist or spec excerpt
- WHEN a reviewer scans role coverage
- THEN all four roles appear with MUST/MUST NOT per row

#### Scenario: Invoices marked deferred

- GIVEN the matrix
- WHEN the invoices row is read
- THEN it defers HTTP/UI per `openspec/specs/p1-accounting-defer/spec.md`

#### Scenario: Provisioning row present

- GIVEN the matrix excerpt
- WHEN a reviewer scans API rows
- THEN the staff user provisioning row appears with correct MUST/MUST NOT per column

#### Scenario: Matrix row for workshop (service job)

- GIVEN the normative matrix summary
- WHEN the matrix is read
- THEN a row documents `Workshop / service job` and distinguishes `client` **MUST NOT** on mutating those routes from staff (`employee` / `manager` / `admin`) **MUST\*** for the HTTP surface in `workshop-repair-execution` scope

#### Scenario: Matrix row for parts inventory

- GIVEN matriz resumida
- WHEN se lee
- THEN fila **Parts inventory**: `manager`/`admin` **MUST**; `client`/`employee` **MUST NOT**

## ADDED Requirements

### Requirement: CI regression for parts inventory

Tests **SHALL** fallar CI si `client`/`employee` tienen 2xx en inventario o nav inventario sin `manager`/`admin`.

#### Scenario: Client API denied

- GIVEN CI `go test` y JWT `client`
- WHEN mutador inventario
- THEN **MUST NOT** 2xx

#### Scenario: Employee API denied

- GIVEN CI y JWT `employee`
- WHEN mutador inventario
- THEN **MUST NOT** 2xx

#### Scenario: Nav hidden

- GIVEN CI shell y usuario `client` o `employee`
- WHEN se renderiza navegación principal
- THEN sin enlace UI inventario

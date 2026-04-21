# Delta for mvp-role-access

## ADDED Requirements

### Requirement: Staff user provisioning HTTP surface

The normative matrix **SHALL** include one row for **staff user provisioning** (authenticated `POST` that creates a `User` with a discrete `role`), aligned with `openspec/specs/staff-user-provisioning/spec.md` (or its change-folder delta until archive).

#### Scenario: Matrix row exists after change

- GIVEN MVP role verification materials
- WHEN the matrix is read
- THEN a row documents `POST` staff user provisioning for `client`, `employee`, `manager`, and `admin` columns per MUST/MUST NOT

### Requirement: CI coverage for provisioning authorization

Automated tests **SHALL** fail CI if `client` or `employee` gains 2xx from the staff user provisioning `POST`, or if `admin` gains 2xx when the body requests `role=admin`.

#### Scenario: Client denied provisioning

- GIVEN `go test` in CI
- WHEN provisioning authorization tests run for JWT `client`
- THEN no 2xx success on the provisioning route

#### Scenario: Admin escalation denied

- GIVEN `go test` in CI
- WHEN a test sends `role=admin` on the provisioning route with JWT `admin`
- THEN the response **MUST NOT** be 2xx success

## MODIFIED Requirements

### Requirement: Published role–surface matrix

MVP verification docs **SHALL** expose the matrix below (or equivalent), including **staff user provisioning** as a distinct API row.

| Area | client | employee | manager | admin |
|------|--------|----------|---------|-------|
| Auth + shell (dashboard, cars, appointments nav) | MUST | MUST | MUST | MUST |
| Cars / appointments (tenant rules) | MUST | MUST | MUST | MUST |
| Repairs read by car | MUST (own car) | MUST | MUST | MUST |
| Repairs POST/PUT/DELETE | MUST NOT | MUST* | MUST* | MUST* |
| `/api/v1/employees` | MUST NOT | MUST NOT | MUST | MUST |
| Staff user provisioning `POST` | MUST NOT | MUST NOT | MUST** | MUST |
| Invoices HTTP/UI (MVP v1) | **Deferred** (`p1-accounting-defer`) | — | — | — |

\*Subject to existing service rules (car ownership, etc.).  
\*\*`manager` **MAY** create only `employee` and `client` per `staff-user-provisioning` matrix.

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

(Previously: matrix had no row for authenticated user creation / staff provisioning.)

### Requirement: CI authorization regression tests

Tests **SHALL** fail CI if employees list opens to `client`/`employee`, if `client` gains 2xx repair mutation, or if staff user provisioning rules in `staff-user-provisioning` are violated (client/employee success, or `role=admin` body success).

#### Scenario: Employees gate tested

- GIVEN `go test` in CI
- WHEN employees authorization tests run
- THEN non-manager non-admin receives forbidden on list route

#### Scenario: Repair mutation denial for client

- GIVEN `go test` in CI
- WHEN client repair mutation tests run
- THEN at least one assertion denies 2xx for `client`

#### Scenario: Provisioning gate tested

- GIVEN `go test` in CI
- WHEN staff user provisioning authorization tests run
- THEN forbidden callers do not receive 2xx and admin-in-body is rejected

(Previously: CI mandate did not include staff user provisioning route.)

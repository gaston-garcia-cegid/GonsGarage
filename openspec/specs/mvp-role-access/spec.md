# mvp-role-access Specification

> **Promoción:** catálogo principal desde el change archivado `openspec/changes/archive/2026-04-20-mvp-role-verification/` (2026-04-20). Actualizado desde `openspec/changes/archive/2026-04-20-admin-provision-user-roles/` — fila de aprovisionamiento staff y CI (2026-04-20). Actualizado desde `openspec/changes/archive/2026-04-21-ui-admin-users-discoverability/` — matriz UI `/admin/users`, CI UI y `openspec/specs/staff-user-management-ui/spec.md` (2026-04-21).

## Purpose

Reproducible **MVP verification** by **role** (`client`, `employee`, `manager`, `admin`): expected **API/UI** access, **idempotent dev seeds** (one user per role via env), and **automated tests** that MUST catch authorization regressions in CI.

## Normative matrix (summary)

| Area | client | employee | manager | admin |
|------|--------|----------|---------|-------|
| Auth + shell (dashboard, cars, appointments nav) | MUST | MUST | MUST | MUST |
| Cars / appointments (tenant rules) | MUST | MUST | MUST | MUST |
| Repairs read by car | MUST (own car) | MUST | MUST | MUST |
| Repairs POST/PUT/DELETE | MUST NOT | MUST* | MUST* | MUST* |
| `/api/v1/employees` | MUST NOT | MUST NOT | MUST | MUST |
| Staff user provisioning `POST /api/v1/admin/users` | MUST NOT | MUST NOT | MUST** | MUST |
| Staff user management UI (`/admin/users` via shell) | MUST NOT | MUST NOT | MUST | MUST |
| Invoices HTTP/UI (MVP v1) | **Deferred** (`p1-accounting-defer`) | — | — | — |

\*Subject to existing service rules (car ownership, etc.).  
\*\*`manager` **MAY** create only `employee` and `client` per `openspec/specs/staff-user-provisioning/spec.md`.

## Requirements

### Requirement: Published role–surface matrix

MVP verification docs **SHALL** expose the matrix above (or equivalent), including **staff user provisioning** as a distinct API row.

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

### Requirement: Staff user provisioning HTTP surface

The normative matrix **SHALL** include one row for **staff user provisioning** (authenticated `POST` that creates a `User` with a discrete `role`), aligned with `openspec/specs/staff-user-provisioning/spec.md`.

#### Scenario: Matrix row exists after change

- GIVEN MVP role verification materials
- WHEN the matrix is read
- THEN a row documents `POST /api/v1/admin/users` for `client`, `employee`, `manager`, and `admin` columns per MUST/MUST NOT

### Requirement: Staff user management UI in published matrix

The normative matrix (summary) **SHALL** include one row for **staff user management UI**: in-app navigation to `/admin/users` for provisioning/management, aligned with `openspec/specs/staff-user-provisioning/spec.md` and `openspec/specs/staff-user-management-ui/spec.md` (see matrix row above).

#### Scenario: Matrix documents UI access

- GIVEN MVP role verification materials
- WHEN the matrix is read
- THEN a row documents shell access to `/admin/users` with MUST for `manager` and `admin`, MUST NOT for `client` and `employee`

#### Scenario: Discoverability not URL-only

- GIVEN `admin` or `manager` using only the shell (no manual URL)
- WHEN they use primary navigation
- THEN they **SHALL** be able to open `/admin/users`

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

### Requirement: Idempotent dev seeds per role

The repo **SHALL** ship idempotent seed command(s) creating ≥1 user per role when missing; email/password **SHALL** come from env (no secrets in git).

#### Scenario: First run creates users

- GIVEN dev DB without those emails
- WHEN seed runs with valid `DATABASE_URL` and `SEED_*`
- THEN rows exist for all four roles

#### Scenario: Re-run is safe

- GIVEN users already exist for those emails
- WHEN seed runs again
- THEN exit success without duplicate errors or silent password overwrite

### Requirement: Employees API for manager and admin only

`/api/v1/employees` **SHALL** return forbidden for `client` and `employee` before business logic.

#### Scenario: Client forbidden

- GIVEN JWT role `client`
- WHEN `GET /api/v1/employees`
- THEN response is forbidden (e.g. 403)

#### Scenario: Manager allowed past role gate

- GIVEN JWT role `manager`
- WHEN `GET /api/v1/employees`
- THEN response is not forbidden **solely** by staff-manager middleware

### Requirement: Client cannot mutate repairs via HTTP

`client` **MUST NOT** get 2xx from repair **POST/PUT/DELETE**; staff **SHALL** mutate per service rules.

#### Scenario: Client POST repair fails

- GIVEN JWT role `client`
- WHEN `POST /api/v1/repairs` with otherwise valid payload
- THEN response is not 2xx success

#### Scenario: Employee may mutate

- GIVEN JWT role `employee` and payload allowed by `RepairService`
- WHEN `POST /api/v1/repairs`
- THEN response is not rejected **only** because role is `employee`

### Requirement: Invoices not in MVP acceptance

MVP completion **SHALL NOT** depend on invoice Gin routes nor Next invoice pages until a superseding accounting change.

#### Scenario: Checklist omits invoice HTTP

- GIVEN MVP role verification sign-off
- WHEN criteria are checked
- THEN success does not require invoice handlers in `cmd/api`

### Requirement: CI authorization regression tests

Tests **SHALL** fail CI if employees list opens to `client`/`employee`, if `client` gains 2xx repair mutation, if staff user provisioning rules in `staff-user-provisioning` are violated, **or** if `AppShell` omits the `/admin/users` nav entry when `canManageUsers` is true (regression on staff user management UI discovery).

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

#### Scenario: Staff users nav present for manager in UI tests

- GIVEN `pnpm test` (or equivalent) in CI
- WHEN a test renders `AppShell` with a `manager` user
- THEN a navigation control linking to `/admin/users` is present

#### Scenario: Staff users nav absent for client in UI tests

- GIVEN `pnpm test` in CI
- WHEN a test renders `AppShell` with a `client` user
- THEN no navigation control links to `/admin/users`

# mvp-role-access Specification

> **Promoción:** catálogo principal desde el change archivado `openspec/changes/archive/2026-04-20-mvp-role-verification/` (2026-04-20).

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
| Invoices HTTP/UI (MVP v1) | **Deferred** (`p1-accounting-defer`) | — | — | — |

\*Subject to existing service rules (car ownership, etc.).

## Requirements

### Requirement: Published role–surface matrix

MVP verification docs **SHALL** expose the matrix above (or equivalent).

#### Scenario: Four roles listed

- GIVEN the linked checklist or spec excerpt
- WHEN a reviewer scans role coverage
- THEN all four roles appear with MUST/MUST NOT per row

#### Scenario: Invoices marked deferred

- GIVEN the matrix
- WHEN the invoices row is read
- THEN it defers HTTP/UI per `openspec/specs/p1-accounting-defer/spec.md`

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

Tests **SHALL** fail CI if employees list opens to `client`/`employee`, or if `client` gains 2xx repair mutation.

#### Scenario: Employees gate tested

- GIVEN `go test` in CI
- WHEN employees authorization tests run
- THEN non-manager non-admin receives forbidden on list route

#### Scenario: Repair mutation denial for client

- GIVEN `go test` in CI
- WHEN client repair mutation tests run
- THEN at least one assertion denies 2xx for `client`

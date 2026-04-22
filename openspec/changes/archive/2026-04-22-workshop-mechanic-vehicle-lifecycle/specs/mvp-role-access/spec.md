# Delta for mvp-role-access

> Change: `workshop-mechanic-vehicle-lifecycle`. See `workshop-repair-execution` for the service-job scope; this delta extends the published matrix and CI.

## MODIFIED Requirements

### Requirement: Client cannot mutate repairs via HTTP

`client` **MUST NOT** get 2xx from repair **POST/PUT/DELETE** **nor** from any **workshop service job** (taller visit) or subresource (checklists) **POST/PUT/PATCH** that mutates; staff **SHALL** mutate per service rules.

(Previously: `repairs` only; no normative *service job*.)

#### Scenario: Client POST repair fails

- GIVEN JWT role `client`
- WHEN `POST /api/v1/repairs` with otherwise valid payload
- THEN response is not 2xx success

#### Scenario: Client mutates service job is denied

- GIVEN JWT role `client`
- WHEN `POST`/`PUT`/`PATCH` to a documented service-job or checklist path for this change
- THEN response is not 2xx success

#### Scenario: Employee may mutate

- GIVEN JWT role `employee` and payload allowed by `RepairService` (or workshop rules when they apply to the request)
- WHEN `POST /api/v1/repairs`
- THEN response is not rejected **solely** because role is `employee`

## ADDED Requirements

### Requirement: Published matrix includes workshop / service job

MVP verification materials **SHALL** include one **Workshop (service job)** row: read for staff (subject to car/tenant rules); `client` **MUST NOT** on mutations; align with the implementing change and `workshop-repair-execution`.

#### Scenario: Matrix row for workshop

- GIVEN the normative matrix summary
- WHEN the change is merged
- THEN the `Workshop / service job` row is present and distinguishes `client` vs. staff (employee, manager, admin) for the HTTP surface in scope

### Requirement: CI must fail if client mutates service job

Automated tests **SHALL** fail CI if a `client` gets 2xx from mutating a documented service-job or checklist `POST/PUT/PATCH` route in scope of this feature.

#### Scenario: CI client denied on create visit

- GIVEN `go test` in CI
- WHEN a test sends a visit *create* with JWT `client`
- THEN the response is not 2xx success

## REMOVED Requirements

(None)

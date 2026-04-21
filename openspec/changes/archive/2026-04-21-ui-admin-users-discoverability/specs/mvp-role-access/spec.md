# Delta for mvp-role-access

## ADDED Requirements

### Requirement: Staff user management UI in published matrix

The normative matrix (summary) **SHALL** include one row for **staff user management UI**: in-app navigation to `/admin/users` for provisioning/management, aligned with `openspec/specs/staff-user-provisioning/spec.md` and `staff-user-management-ui`.

| Area (add row) | client | employee | manager | admin |
|----------------|--------|----------|---------|-------|
| Staff user management UI (`/admin/users` via shell) | MUST NOT | MUST NOT | MUST | MUST |

#### Scenario: Matrix documents UI access

- GIVEN MVP role verification materials
- WHEN the matrix is read
- THEN a row documents shell access to `/admin/users` with MUST for `manager` and `admin`, MUST NOT for `client` and `employee`

#### Scenario: Discoverability not URL-only

- GIVEN `admin` or `manager` using only the shell (no manual URL)
- WHEN they use primary navigation
- THEN they **SHALL** be able to open `/admin/users`

## MODIFIED Requirements

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

(Previously: CI mandate did not include UI discoverability for `/admin/users`.)

# staff-user-management-ui Specification

> **Source:** promovido desde `openspec/changes/archive/2026-04-21-ui-admin-users-discoverability/` (2026-04-21).

## Purpose

Staff who **MAY** manage users (`admin`, `manager` per `canManageUsers`) **SHALL** reach `/admin/users` from the shared authenticated shell **without** typing the URL. Others **MUST NOT** see that entry.

## Requirements

### Requirement: Primary shell navigation entry

When the main staff `AppShell` is shown, the system **SHALL** render a navigation control (e.g. button) labeled in **pt_PT** (e.g. «Utilizadores») linking to `/admin/users` **if and only if** `canManageUsers(user)`.

#### Scenario: Manager sees users nav

- GIVEN an authenticated `manager` on any `AppShell` page
- WHEN the primary nav is rendered
- THEN a control linking to `/admin/users` is visible

#### Scenario: Client does not see users nav

- GIVEN an authenticated `client` on `AppShell`
- WHEN the primary nav is rendered
- THEN no control targets `/admin/users`

### Requirement: Active navigation state

On route `/admin/users`, the shell **SHALL** mark the users nav entry as active (same visual pattern as other `activeNav` items).

#### Scenario: Active on admin users page

- GIVEN `canManageUsers` and the user is on `/admin/users`
- WHEN `AppShell` renders
- THEN the users nav control uses the active style

### Requirement: Single source of navigation truth

The repo **MUST NOT** ship duplicate, unused navigation maps for the same shell. Either `navigation.ts` (or equivalent) **SHALL** be removed, **SHALL** be the single module consumed by `AppShell`, or **SHALL** re-export only data consumed by `AppShell` (one source).

#### Scenario: No orphan staff users path

- GIVEN the codebase after this change
- WHEN searching for `/admin/users` in navigation config
- THEN at most one authoritative definition applies to `AppShell` (no dead duplicate listing the path)

### Requirement: Optional dashboard CTA

The dashboard **MAY** show a secondary link or button to `/admin/users` under the same `canManageUsers` guard.

#### Scenario: CTA optional

- GIVEN `canManageUsers` on `/dashboard`
- WHEN the dashboard implements the optional CTA
- THEN it navigates to `/admin/users` without relaxing role checks

### Requirement: Admin users page toolbar and modal provisioning

> **Merged from:** `openspec/changes/archive/2026-04-21-admin-users-app-shell-modal-parity/` (delta ADDED).

For route `/admin/users`, when the viewer **MAY** manage users (`canManageUsers`), the page **SHALL** follow the same shell pattern as other primary staff areas using `AppShell` with a **toolbar**: a page title and a primary control that opens user creation. The capture fields for provisioning **SHALL** be shown **inside a modal dialog** while that flow is open; they **SHALL NOT** be the only default main content (inline form as sole primary surface) without a toolbar CTA and dialog.

#### Scenario: Toolbar shows title and create action

- GIVEN an authenticated user with `canManageUsers` on `/admin/users`
- WHEN the page renders
- THEN a toolbar presents a page title and a control to start creating a user (e.g. «Novo utilizador»)

#### Scenario: Create action opens modal with form

- GIVEN the user has activated the create control
- WHEN the create flow is shown
- THEN a modal dialog is open and contains the provisioning fields (email, password, names, role per existing rules)

#### Scenario: Close without success leaves shell

- GIVEN the modal is open
- WHEN the user dismisses or cancels without a successful submit
- THEN the modal closes and the user remains on `/admin/users` under `AppShell`

#### Scenario: Success closes or clears and confirms

- GIVEN a successful provisioning response from the modal submit
- WHEN the client applies the result
- THEN the modal **SHALL** close or clear for another entry and a success indication **SHALL** be visible without leaving `/admin/users`

#### Scenario: Failed submit shows error in modal

- GIVEN the modal is open and submit returns failure
- WHEN the UI updates
- THEN an error message **SHALL** be visible in the modal and the user **MAY** correct input and retry

# Delta: staff-user-management-ui

## ADDED Requirements

### Requirement: Admin users page toolbar and modal provisioning

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

## MODIFIED Requirements

_None — existing navigation requirements unchanged._

## REMOVED Requirements

_None._

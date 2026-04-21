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

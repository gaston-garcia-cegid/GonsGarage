# client-auth-shell Specification

## Purpose

Define how the GonsGarage **web client** SHALL present **login** (`/auth/login`) and **registration** (`/auth/register`) as **one coherent auth experience**: shared page structure, token-driven visuals, aligned feedback patterns, and a **single** programmatic auth API surface for both flows.

## Requirements

### Requirement: Shared page shell

The product SHALL render login and registration inside a **visually equivalent page shell** (background, centred column width, card container, header area for brand/title, consistent vertical rhythm) so a user recognises the same “place” in the app on both routes.

#### Scenario: Structural parity

- GIVEN a user opens `/auth/login` and later `/auth/register` (or the reverse)
- WHEN they compare header, outer background, card width, and outer padding at the same viewport width
- THEN both routes use the same shell semantics (no divergent full-page layouts such as one route feeling like a different product skin)

#### Scenario: Token-driven surfaces

- GIVEN the design tokens in `frontend/src/styles/tokens.css`
- WHEN auth pages set backgrounds, borders, radii, shadows, and primary text colours for the shell and card
- THEN those values SHALL come from design tokens or shared utilities, not unrelated ad-hoc colour systems that contradict the rest of the app

### Requirement: Form and feedback consistency

Form labels, field density, primary submit control, and **inline alerts** (success banner after redirect query, general API error, field-level errors) SHALL follow the **same visual language** on login and registration (spacing, typography scale, success/warning/error treatment).

#### Scenario: Success message after registration

- GIVEN the user lands on login with a non-empty success `message` query parameter
- WHEN the success banner is shown
- THEN its styling SHALL match the registration flow’s alert styling (same role tokens or shared alert pattern), not a one-off colour block unique to login

#### Scenario: Validation errors

- GIVEN invalid input on either auth route
- WHEN field-level errors are displayed
- THEN they SHALL use the same text size, colour role, and proximity to fields as on the sibling route

### Requirement: Single auth consumer for both routes

Both `/auth/login` and `/auth/register` SHALL obtain authentication actions and session state through the **same documented client module** (the Zustand-based `useAuth` from `@/stores` as used elsewhere in the app), so behaviour does not silently diverge between Context-only and store-only code paths.

#### Scenario: Registration uses store auth

- GIVEN the registration page submits valid data
- WHEN `register` completes successfully
- THEN the page SHALL have used the same `useAuth` source as the login form for session-related state and redirects, without depending on a parallel legacy context API for that flow

#### Scenario: Login unchanged semantically

- GIVEN a user with valid credentials on `/auth/login`
- WHEN they submit the form
- THEN login SHALL still succeed and navigate as today (same post-login destination policy), with only presentation and wiring normalisation allowed to change

### Requirement: Cross-navigation affordances

Links or buttons that switch between “create account” and “already have an account” SHALL be styled as **secondary actions** consistent with each other and with the app’s secondary control pattern (not mismatched raw link colours between the two pages).

#### Scenario: Cross-link parity

- GIVEN the user is on login and sees the control to go to registration
- WHEN they compare it to the control on registration that returns to login
- THEN both SHALL read as the same interaction tier (secondary / link) with aligned typography and hover/focus behaviour

### Requirement: Quality gate

A release claiming this capability SHALL keep the frontend **lint-clean** and **typecheck-clean** after the refactor.

#### Scenario: Commands succeed

- GIVEN the change is ready for review
- WHEN `pnpm lint` and `pnpm typecheck` run in the frontend package
- THEN both complete successfully with no new errors introduced solely by this work

### Requirement: Registration visible copy in European Portuguese

All **user-visible** strings on `/auth/register` that are authored in the registration page module (field labels, placeholders, and client-constructed general error text shown in the auth shell banner) SHALL be in **European Portuguese** and SHALL use the **same terminology** as `/auth/login` for equivalent fields (e.g. e-mail field: label and placeholder aligned with login).

#### Scenario: Email field matches login wording

- GIVEN the user opens `/auth/register`
- WHEN they inspect the email field label and placeholder
- THEN the label and placeholder SHALL be in Portuguese and SHALL match the login form’s e-mail label and placeholder semantics (not English marketing phrasing such as “Email Address”).

#### Scenario: Confirm password label in Portuguese

- GIVEN the user opens `/auth/register`
- WHEN they inspect the confirm-password field label
- THEN the label SHALL be in Portuguese describing confirmation of the password (not English “Confirm Password”).

#### Scenario: No English “Error:” prefix from client catch path

- GIVEN the registration submit path catches an unexpected thrown error
- WHEN the general error banner is shown on the registration page
- THEN the banner text SHALL NOT begin with the English literal `Error:` as a client-added prefix.

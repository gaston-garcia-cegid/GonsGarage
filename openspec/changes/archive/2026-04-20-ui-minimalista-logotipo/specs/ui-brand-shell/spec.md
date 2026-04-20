# ui-brand-shell Specification

## Purpose

Define how the GonsGarage web UI SHALL present a **minimal, professional** look aligned with the **brand palette** (logo-derived), with **consistent light/dark** behaviour on the surfaces in scope.

## Requirements

### Requirement: Documented brand palette

The product SHALL maintain a **single documented source of truth** for brand-related colours (navy/carbon, trust blue, signal red, neutrals) used by the application shell and in-scope routes, so designers and implementers can verify alignment with the logo.

#### Scenario: Brand section exists

- GIVEN the design tokens stylesheet for the frontend
- WHEN a maintainer opens it to audit colours against the logo asset
- THEN a dedicated **Brand** (or equivalent) comment block lists the role of each brand token relative to the logo
- AND primary interaction colours are traceable to those tokens

#### Scenario: No orphan brand hex in shell

- GIVEN the authenticated application shell (header, primary nav, logout)
- WHEN the theme is toggled between light and dark
- THEN destructive and accent states remain legible without introducing **raw hex literals** for colours that already have a semantic token

### Requirement: Theme coherence on priority routes

The **vehicles** and **appointments** user-facing views that ship with this change SHALL NOT rely on hardcoded colour literals for backgrounds, borders, and text that duplicate the token palette, so light and dark themes stay visually coherent.

#### Scenario: Cars view respects theme

- GIVEN a user on the cars area with dark theme enabled
- WHEN they browse lists and forms in that area
- THEN surfaces and status colours derive from shared tokens or shared utilities, not ad-hoc hex unrelated to the token system

#### Scenario: Appointments view respects theme

- GIVEN a user on the appointments area
- WHEN status chips or alerts are shown
- THEN their colours map to semantic tokens (success, warning, error, info) or documented utilities, not standalone hex blocks that ignore dark mode

### Requirement: Non-regression quality gate

A release that claims this capability SHALL keep the frontend **lint-clean** and **buildable** after styling refactors.

#### Scenario: CI-quality commands succeed locally

- GIVEN the change is ready for review
- WHEN `pnpm lint` and `pnpm build` are run in the frontend package
- THEN both complete successfully with no new errors introduced by this work

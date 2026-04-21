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

---

## Merged additions (change `mvp-ui-visual-parity`, archived 2026-04-21)

The following sections **add** MVP-wide dark, loading, route homogeneity, toolchain spike, and Shadcn-canonical behaviour. They **do not** remove requirements above.

### Phase 1 — Dark usable + loading contract

#### Requirement: Dark theme minimum quality (MVP)

Authenticated MVP surfaces **SHALL** present legible contrast in `html[data-theme='dark']` for primary/secondary text, borders, and panel backgrounds, without relying on raw `#rrggbb` in components when an equivalent token exists.

##### Scenario: Dashboard readable in dark

- **GIVEN** `data-theme='dark'` and an authenticated user
- **WHEN** they open `/dashboard`
- **THEN** metric text and links meet at least **WCAG AA** contrast against backgrounds **OR** are justified in code with an issue link
- **AND** there are no unintended light “ghost” blocks without a surface token

##### Scenario: Accounting shell readable in dark

- **GIVEN** `data-theme='dark'`
- **WHEN** they navigate under `/accounting`
- **THEN** headers, tables, and forms use `--surface-*`, `--text-*`, or documented utilities aligned to tokens

#### Requirement: Unified loading indicator

The app **SHALL** expose a reusable loading pattern (component or documented class) with exactly **three** nominal sizes (`sm`, `md`, `lg`) with fixed proportions; MVP screens **SHALL NOT** introduce new ad hoc sizes except via documented `docs/` entry or `// UI exception:`.

##### Scenario: Size variants documented

- **GIVEN** a developer opens the unified indicator source
- **WHEN** they look for size variants
- **THEN** they find exactly three nominal variants with dimensions in `rem` in a single shared module or stylesheet

##### Scenario: Full-page load uses large variant only

- **GIVEN** progressive substitution on MVP routes
- **WHEN** a view shows full-page loading
- **THEN** it uses the `lg` variant (or named equivalent) and announces busy state with `aria-busy="true"` **OR** associated accessible text

### Phase 2 — Cross-route homogeneity

#### Requirement: MVP route coverage for theme coherence

Beyond **cars** and **appointments**, **employees**, **client**, **dashboard**, **accounting** (and sub-routes) **SHALL** follow the same policy as “Theme coherence on priority routes”: no duplicated colour literals for background/border/text on first-view components that contradict tokens.

##### Scenario: Employees list respects tokens in both themes

- **GIVEN** staff on `/employees`
- **WHEN** toggling light ↔ dark
- **THEN** rows, header, and primary actions do not end up with incompatible inherited backgrounds (e.g. pure white on white) from local literals

##### Scenario: Client home respects tokens

- **GIVEN** client role on `/client`
- **WHEN** dark theme
- **THEN** cards and empty states use tokenized surfaces and text

#### Requirement: Loader migration completeness (MVP)

For each file under `frontend/src/app/` that used `spinnerLg`, `spinnerMd`, local `.spinner` in a CSS module, or undocumented equivalent, the team **SHALL** migrate to the unified pattern **OR** document the exception in the phase PR.

##### Scenario: Inventory gate before closing Phase 2

- **GIVEN** Phase 2 closure
- **WHEN** searching the repo for agreed patterns (`spinnerLg`, `className={styles.spinner}`)
- **THEN** zero occurrences **OR** an explicit approved exception list in `verify-report.md`

### Phase 3 — Platform spike (Tailwind v4 + Next 16)

#### Requirement: Upgrade spike decision record

Before adopting **Next.js 16** and **Tailwind CSS v4** on `main`, the project **SHALL** keep a document (ADR in `docs/` or change section) with: benefit hypothesis, branches tried, `pnpm build` + `pnpm test` outcomes, and **GO** / **NO-GO** / **DEFER** decision.

##### Scenario: NO-GO preserves current stack

- **GIVEN** a NO-GO outcome
- **WHEN** the phase is archived
- **THEN** `package.json` on `main` stays on Next 15.x without mandatory Tailwind v4
- **AND** the ADR explains why

##### Scenario: GO documents migration steps

- **GIVEN** a GO outcome
- **WHEN** the spike merges
- **THEN** the ADR lists contributor migration steps (commands, known breaking changes)

#### Requirement: Spike non-regression

Any branch exercising Phase 3 **SHALL** pass the same quality commands as the base `ui-brand-shell` spec (`pnpm lint`, `pnpm build`) before proposing merge to `main`.

##### Scenario: CI parity on spike branch

- **GIVEN** a spike PR
- **WHEN** the repo CI workflow runs
- **THEN** configured frontend jobs **SHALL** finish green or the PR is not mergeable

### Phase 4 — Full Shadcn-style redesign

#### Requirement: Canonical component system (greenfield)

The product **SHALL** adopt a **Shadcn-style** component system (see `ui-component-system` spec) as the single reference for new UI and migrations; redesign **SHALL** be treated as **greenfield** on that base (no stray legacy CSS patches without an exit path).

##### Scenario: Legacy escape hatch documented

- **GIVEN** a screen not yet migrated
- **WHEN** code touching it merges
- **THEN** the PR links a Shadcn migration task **OR** extends an exception in `verify-report.md` with owner and date

##### Scenario: Post–Phase 4 MVP visual consistency

- **GIVEN** Phase 4 closed per `tasks.md`
- **WHEN** a user walks MVP routes in light and dark
- **THEN** shared interactive controls (primary button, input, modal, table) share the same Shadcn system visual language

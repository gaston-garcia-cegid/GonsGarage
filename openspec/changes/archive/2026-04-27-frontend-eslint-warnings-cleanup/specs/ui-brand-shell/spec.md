# Delta for ui-brand-shell

## ADDED Requirements

### Requirement: Lint warning budget on default frontend gate

A release or merge that touches the **Next.js frontend package** **SHALL** satisfy the repository default **`pnpm lint`** outcome **without warnings**, so React Hooks and related rules deferred as `warn` during migration **MUST NOT** accumulate unchecked.

#### Scenario: Clean lint output on review

- **GIVEN** a change ready for review that modifies files under `frontend/`
- **WHEN** a maintainer runs `pnpm lint` in the frontend package as documented in `openspec/config.yaml`
- **THEN** the command **SHALL** complete with **zero warnings** in the aggregate lint report for that package
- **AND** `pnpm build` in the same package **SHALL** still complete successfully unless the change explicitly documents a temporary build exception approved outside this requirement

#### Scenario: Capped, documented waivers

- **GIVEN** a specific warning cannot be removed in the same change without unacceptable product risk
- **WHEN** the change documents the exception in its proposal or verify notes with rationale
- **THEN** at most **three** source files **MAY** retain a targeted suppression
- **AND** each suppression **SHALL** be scoped (rule + line or narrow block), not a file-wide disable of the default hook ruleset

## MODIFIED Requirements

### Requirement: Non-regression quality gate

A release that claims this capability SHALL keep the frontend **lint-clean (no errors and no warnings on the default `pnpm lint` invocation)** and **buildable** after styling refactors and related UI refactors, subject only to the **Lint warning budget** waiver rules in this delta.

(Previously: success required no new **errors** from lint; **warnings** were not contractually capped.)

#### Scenario: CI-quality commands succeed locally

- **GIVEN** the change is ready for review
- **WHEN** `pnpm lint` and `pnpm build` are run in the frontend package
- **THEN** both complete successfully with **no errors and no warnings** introduced or left unresolved by this work, except waivers that meet the **Lint warning budget** requirement
- **AND** any waiver **SHALL** appear in the active change documentation with rationale

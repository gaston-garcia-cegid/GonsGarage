# Skill Registry

**Delegator use only.** Any agent that launches sub-agents reads this registry to resolve compact rules, then injects them directly into sub-agent prompts. Sub-agents do NOT read this registry or individual SKILL.md files.

See `_shared/skill-resolver.md` (Gentleman skills) for the full resolution protocol.

## User Skills

| Trigger | Skill | Path |
|---------|-------|------|
| When creating a pull request, opening a PR, or preparing changes for review. | branch-pr | `C:\Users\gaston.garcia\.claude\skills\branch-pr\SKILL.md` |
| When creating a GitHub issue, reporting a bug, or requesting a feature. | issue-creation | `C:\Users\gaston.garcia\.claude\skills\issue-creation\SKILL.md` |
| When writing Go tests, using teatest, or adding test coverage. | go-testing | `C:\Users\gaston.garcia\.claude\skills\go-testing\SKILL.md` |
| judgment day, judgment-day, review adversarial, dual review, doble review, juzgar, que lo juzguen | judgment-day | `C:\Users\gaston.garcia\.claude\skills\judgment-day\SKILL.md` |
| When the user asks to create a new skill, add agent instructions, or document patterns for AI. | skill-creator | `C:\Users\gaston.garcia\.claude\skills\skill-creator\SKILL.md` |
| Keep a PR merge-ready by triaging comments, resolving clear conflicts, and fixing CI in a loop. | babysit | `C:\Users\gaston.garcia\.cursor\skills-cursor\babysit\SKILL.md` |
| Standalone analytical artifacts, data-heavy deliverables, MCP tool results as UI; create/edit `.canvas.tsx`. | canvas | `C:\Users\gaston.garcia\.cursor\skills-cursor\canvas\SKILL.md` |
| Create a rule, add coding standards, `.cursor/rules/`, AGENTS.md. | create-rule | `C:\Users\gaston.garcia\.cursor\skills-cursor\create-rule\SKILL.md` |
| Create, write, or author a new skill; SKILL.md format. | create-skill | `C:\Users\gaston.garcia\.cursor\skills-cursor\create-skill\SKILL.md` |
| Split a chat, set of changes, branch, or PR into small PRs. | split-to-prs | `C:\Users\gaston.garcia\.cursor\skills-cursor\split-to-prs\SKILL.md` |
| Change editor settings, settings.json, themes, format on save. | update-cursor-settings | `C:\Users\gaston.garcia\.cursor\skills-cursor\update-cursor-settings\SKILL.md` |
| Create hooks, hooks.json, automate around agent events. | create-hook | `C:\Users\gaston.garcia\.cursor\skills-cursor\create-hook\SKILL.md` |
| status line, statusline, CLI status bar, prompt footer. | statusline | `C:\Users\gaston.garcia\.cursor\skills-cursor\statusline\SKILL.md` |
| User explicitly invokes `/shell` and wants literal command execution. | shell | `C:\Users\gaston.garcia\.cursor\skills-cursor\shell\SKILL.md` |
| Migrate `.mdc` rules / slash commands to `.cursor/skills/`. | migrate-to-skills | `C:\Users\gaston.garcia\.cursor\skills-cursor\migrate-to-skills\SKILL.md` |
| Create custom subagents, `.cursor/agents/`. | create-subagent | `C:\Users\gaston.garcia\.cursor\skills-cursor\create-subagent\SKILL.md` |
| Change CLI settings, `cli-config.json`, permissions, sandbox. | update-cli-config | `C:\Users\gaston.garcia\.cursor\skills-cursor\update-cli-config\SKILL.md` |

## Compact Rules

Pre-digested rules per skill. Delegators copy matching blocks into sub-agent prompts as `## Project Standards (auto-resolved)`.

### branch-pr
- Every PR MUST link an approved issue (`Closes`/`Fixes`/`Resolves #N`).
- Every PR MUST have exactly one `type:*` label.
- Branch names MUST match `^(feat|fix|chore|docs|style|refactor|perf|test|build|ci|revert)\/[a-z0-9._-]+$`.
- Use `.github/PULL_REQUEST_TEMPLATE.md`; automated checks must pass before merge.

### issue-creation
- Use GitHub templates only (blank issues disabled).
- New issues get `status:needs-review`; maintainer MUST add `status:approved` before any PR.
- Search duplicates first; questions go to Discussions per template workflow.

### go-testing
- Prefer table-driven tests with `t.Run` subtests.
- Skill targets Gentleman.Dots / Bubbletea patterns; for this repo use standard `testing` + `testify` as in `backend/`.
- Use `httptest` for HTTP handler tests; integration tests live under `backend/tests/integration/`.

### judgment-day
- Before judges: resolve skills via registry or Engram `mem_search(skill-registry)`; inject `## Project Standards (auto-resolved)` into both judges and fix agent.
- Launch two blind judges in parallel; synthesize Confirmed / Suspect / Contradiction; classify WARNING as real vs theoretical.
- Re-judge after fixes until pass or two iterations; orchestrator coordinates only.

### skill-creator
- Use Agent Skills layout: `skills/{name}/SKILL.md` (+ optional assets/references).
- Frontmatter: `name`, `description` with Trigger line, version metadata.
- Do not create a skill when docs or a one-off suffices.

### babysit
- Triage every PR comment (including Bugbot); fix only what you agree with, explain the rest.
- Resolve conflicts only when intent is clear; otherwise stop for clarification.
- Fix CI with small scoped fixes and re-check until green and mergeable.

### canvas
- Use for standalone analytical deliverables (audits, metrics, MCP-sourced tables); skip for normal code fixes or short answers.
- One `.canvas.tsx` per canvas; import only from `cursor/canvas`; no `fetch`, no extra modules.
- Read `canvas/sdk` typings before using uncommon components; prefer built-ins over hand-rolled UI.

### create-rule
- Rules are `.mdc` in `.cursor/rules/` with YAML frontmatter (`description`, optional `globs`, `alwaysApply`).
- Clarify scope and file patterns before writing; use AskQuestion when ambiguous.

### create-skill
- Gather purpose, location (user vs project), triggers, output format; respect user verbatim wording in SKILL.md.
- Skills are directories with `SKILL.md`; follow Cursor Agent Skills conventions.

### split-to-prs
- Do not branch/commit/push/open PRs until the user approves the split plan.
- No destructive git without explicit approval; no `git add .`; save recoverable snapshot before moving work.
- Default independent PRs off default branch; stack only for real dependencies.

### update-cursor-settings
- Read `%APPDATA%\Cursor\User\settings.json` (Windows) first; preserve unrelated keys; validate JSON.

### create-hook
- Choose project (`.cursor/hooks.json`) vs user hooks; narrowest event; decide fail-open vs fail-closed.
- Project hook paths are relative to repo root.

### statusline
- Configure via `~/.cursor/cli-config.json` `statusLine.command`; stdin JSON matches CLI `StatusLinePayload`.
- Respect `updateIntervalMs` (>=300) and `timeoutMs` defaults.

### shell
- Only when user invokes `/shell`; treat following text as literal command; no rewrite before run.

### migrate-to-skills
- Copy rule/command body verbatim when converting `.mdc`/commands to skills.
- Migrate rules that have `description` but no `globs` and not `alwaysApply: true`.

### create-subagent
- Project agents in `.cursor/agents/` override user `~/.cursor/agents/` by name.
- Markdown + frontmatter body is the system prompt.

### update-cli-config
- Home config `~/.cursor/cli-config.json`; merge project `.cursor/cli.json` along path; restart CLI to apply.

## Project Conventions

| File | Path | Notes |
|------|------|-------|
| nextjs-app-router-navigation | `d:\Repos\GonsGarage\.cursor\rules\nextjs-app-router-navigation.mdc` | `alwaysApply: true` — App Router `page.tsx` must exist for each `Link`/`router.push` target; prefer modals over invented edit routes. |

Read the convention files listed above for project-specific patterns and rules.

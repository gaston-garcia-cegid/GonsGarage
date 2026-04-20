# Skill Registry

**Delegator use only.** Any agent that launches sub-agents reads this registry to resolve compact rules, then injects them directly into sub-agent prompts. Sub-agents do NOT read this registry or individual SKILL.md files.

See `_shared/skill-resolver.md` for the full resolution protocol (if present in your tooling).

## User Skills

| Trigger | Skill | Path |
|---------|-------|------|
| When creating a pull request, opening a PR, or preparing changes for review. | branch-pr | `C:\Users\gaston.garcia\.cursor\skills\branch-pr\SKILL.md` |
| When creating a GitHub issue, reporting a bug, or requesting a feature. | issue-creation | `C:\Users\gaston.garcia\.cursor\skills\issue-creation\SKILL.md` |
| When writing Go tests, using teatest, or adding test coverage. | go-testing | `C:\Users\gaston.garcia\.cursor\skills\go-testing\SKILL.md` |
| When user says "judgment day", "judgment-day", "review adversarial", "dual review", "doble review", "juzgar", "que lo juzguen". | judgment-day | `C:\Users\gaston.garcia\.cursor\skills\judgment-day\SKILL.md` |
| When user asks to create a new skill, add agent instructions, or document patterns for AI. | skill-creator | `C:\Users\gaston.garcia\.cursor\skills\skill-creator\SKILL.md` |
| Standalone analytical artifact, data-heavy tables, MCP tool deliverables as UI; create/edit `.canvas.tsx`. | canvas | `C:\Users\gaston.garcia\.cursor\skills-cursor\canvas\SKILL.md` |
| Keep a PR merge-ready (comments, conflicts, CI). | babysit | `C:\Users\gaston.garcia\.cursor\skills-cursor\babysit\SKILL.md` |
| Create a rule, coding standards, `.cursor/rules/`, AGENTS.md. | create-rule | `C:\Users\gaston.garcia\.cursor\skills-cursor\create-rule\SKILL.md` |
| Create, write, or author a new skill; SKILL.md format. | create-skill | `C:\Users\gaston.garcia\.cursor\skills-cursor\create-skill\SKILL.md` |
| Create a hook, hooks.json, automate around agent events. | create-hook | `C:\Users\gaston.garcia\.cursor\skills-cursor\create-hook\SKILL.md` |
| Change editor settings, settings.json, themes, format on save. | update-cursor-settings | `C:\Users\gaston.garcia\.cursor\skills-cursor\update-cursor-settings\SKILL.md` |
| Status line, statusline, CLI status bar, prompt footer. | statusline | `C:\Users\gaston.garcia\.cursor\skills-cursor\statusline\SKILL.md` |
| User invokes `/shell` and wants literal command execution. | shell | `C:\Users\gaston.garcia\.cursor\skills-cursor\shell\SKILL.md` |
| Migrate .mdc rules or slash commands to Agent Skills. | migrate-to-skills | `C:\Users\gaston.garcia\.cursor\skills-cursor\migrate-to-skills\SKILL.md` |
| Create custom subagents, task-specific agents. | create-subagent | `C:\Users\gaston.garcia\.cursor\skills-cursor\create-subagent\SKILL.md` |
| Change CLI settings, cli-config.json, sandbox, approval mode. | update-cli-config | `C:\Users\gaston.garcia\.cursor\skills-cursor\update-cli-config\SKILL.md` |

## Project conventions

No `AGENTS.md`, `CLAUDE.md`, `.cursorrules`, `GEMINI.md`, or `copilot-instructions.md` in the repository root (this scan: 2026-04-20).

## Compact Rules

Pre-digested rules per skill. Delegators copy matching blocks into sub-agent prompts as `## Project Standards (auto-resolved)`.

### branch-pr
- Every PR MUST link an **approved** issue (`status:approved`); every PR MUST have exactly one `type:*` label.
- Branch names: `^(feat|fix|chore|docs|style|refactor|perf|test|build|ci|revert)\/[a-z0-9._-]+$`.
- Run shellcheck on modified scripts before PR; use the repo PR template.

### issue-creation
- Use a template only (blank issues disabled); issues start with `status:needs-review`; maintainer adds `status:approved` before any PR.

### go-testing
- Prefer table-driven `t.Run` subtests; use `httptest` for HTTP handler tests; integration tests can live under e.g. `backend/tests/integration` (match existing patterns).
- For Bubbletea/TUI: test `Model.Update` transitions directly when applicable.

### judgment-day
- Before judges: resolve skill registry (Engram `skill-registry` or `.atl/skill-registry.md`), inject matching **Compact Rules** into both judge and fix-agent prompts identically.
- Launch **two** blind reviewers in **parallel**; synthesize verdicts; fix and re-judge up to 2 iterations or escalate.

### skill-creator
- Skills use `SKILL.md` + optional `assets/` and `references/`; frontmatter needs `name`, `description` with **Trigger:** line.
- Do not create a skill when docs already cover it or the task is one-off.

### canvas
- Use a **canvas** (`.canvas.tsx`) for standalone analytical deliverables, large structured tables, or MCP-sourced data as the primary output — not for simple code fixes or short answers.
- Read this skill when creating or debugging any `.canvas.tsx` file.

### babysit
- Triage all PR comments (including bots); resolve conflicts only when intent is clear; fix CI in small scoped loops until green and mergeable.

### create-rule
- Place rules in `.cursor/rules/`; define scope (always vs glob); clarify file patterns before writing.

### create-skill
- Gather purpose, location (user vs project), triggers, domain knowledge, and output format; prefer project skills when the workflow is repo-specific.

### create-hook
- Project hooks: `.cursor/hooks.json` + `.cursor/hooks/*`; user hooks under `~/.cursor/`; choose narrowest hook event; decide fail-open vs fail-closed.

### update-cursor-settings
- Read `%APPDATA%\Cursor\User\settings.json` first; change only requested keys; validate JSON before save.

### statusline
- Configure `statusLine` in `~/.cursor/cli-config.json` pointing at a command; stdin receives session JSON; respect timeout/update interval limits.

### shell
- Only when user uses `/shell`: run the following text as a **literal** command with no rewriting.

### migrate-to-skills
- **CRITICAL:** Copy rule/command body **verbatim** — do not reformat or “improve” while migrating.

### create-subagent
- Subagents live in user or project config; isolate context with focused system prompts for repeatable tasks.

### update-cli-config
- Main file `~/.cursor/cli-config.json`; project overrides via layered `.cursor/cli.json` along path; restart CLI after changes.

# Skill Registry

**Delegator use only.** Any agent that launches sub-agents reads this registry to resolve compact rules, then injects them directly into sub-agent prompts. Sub-agents do NOT read this registry or individual SKILL.md files.

## User Skills

| Trigger | Skill | Path |
|---------|-------|------|
| When writing Go tests, using teatest, or adding test coverage. | go-testing | `C:\Users\gaston.garcia\.claude\skills\go-testing\SKILL.md` |
| When creating a pull request, opening a PR, or preparing changes for review. | branch-pr | `C:\Users\gaston.garcia\.claude\skills\branch-pr\SKILL.md` |
| When creating a GitHub issue, reporting a bug, or requesting a feature. | issue-creation | `C:\Users\gaston.garcia\.claude\skills\issue-creation\SKILL.md` |
| judgment day, judgment-day, review adversarial, dual review, doble review, juzgar, que lo juzguen | judgment-day | `C:\Users\gaston.garcia\.claude\skills\judgment-day\SKILL.md` |
| When user asks to create a new skill, add agent instructions, or document patterns for AI. | skill-creator | `C:\Users\gaston.garcia\.claude\skills\skill-creator\SKILL.md` |
| Standalone analytical artifacts, data-heavy deliverables, MCP tool results as primary output; creating/editing `.canvas.tsx` | canvas | `C:\Users\gaston.garcia\.cursor\skills-cursor\canvas\SKILL.md` |
| Create a rule, coding standards, `.cursor/rules/`, AGENTS.md | create-rule | `C:\Users\gaston.garcia\.cursor\skills-cursor\create-rule\SKILL.md` |
| Create, write, or author a new skill; SKILL.md format | create-skill | `C:\Users\gaston.garcia\.cursor\skills-cursor\create-skill\SKILL.md` |
| Keep a PR merge-ready; triage comments, conflicts, CI | babysit | `C:\Users\gaston.garcia\.cursor\skills-cursor\babysit\SKILL.md` |
| Split a chat, branch, or PR into smaller PRs | split-to-prs | `C:\Users\gaston.garcia\.cursor\skills-cursor\split-to-prs\SKILL.md` |
| Change editor settings, settings.json, themes, format on save | update-cursor-settings | `C:\Users\gaston.garcia\.cursor\skills-cursor\update-cursor-settings\SKILL.md` |
| Create hooks, hooks.json, automate around agent events | create-hook | `C:\Users\gaston.garcia\.cursor\skills-cursor\create-hook\SKILL.md` |
| status line, statusline, CLI prompt footer | statusline | `C:\Users\gaston.garcia\.cursor\skills-cursor\statusline\SKILL.md` |
| User explicitly invokes `/shell` | shell | `C:\Users\gaston.garcia\.cursor\skills-cursor\shell\SKILL.md` |
| Create custom subagents, `.cursor/agents/` | create-subagent | `C:\Users\gaston.garcia\.cursor\skills-cursor\create-subagent\SKILL.md` |
| Migrate .mdc rules / slash commands to `.cursor/skills/` | migrate-to-skills | `C:\Users\gaston.garcia\.cursor\skills-cursor\migrate-to-skills\SKILL.md` |
| Change `~/.cursor/cli-config.json`, CLI permissions, display | update-cli-config | `C:\Users\gaston.garcia\.cursor\skills-cursor\update-cli-config\SKILL.md` |

## Project Conventions

| File | Notes |
|------|-------|
| `.cursor/rules/nextjs-app-router-navigation.mdc` | App Router: every `Link`/`router.push` to a page must exist under `frontend/src/app/**/page.tsx`; prefer modal/query patterns over unimplemented dynamic routes. |

## Compact Rules

### go-testing
- Prefer table-driven tests with `t.Run` subtests; cover happy path, errors, and boundaries.
- For Bubbletea: use `teatest` patterns from the skill when testing TUI models.
- Integration: use `httptest` / real DB fixtures consistent with `backend/tests/integration/`.
- Golden files: update with intent; never silence diffs without review.

### branch-pr
- Every PR MUST link an approved issue; exactly one `type:*` label.
- Branch naming and conventional commits per team workflow; checks green before merge.

### issue-creation
- Use repo templates only; issues start `status:needs-review`; maintainer adds `status:approved` before PR work.

### judgment-day
- Before judges: resolve skills from this registry → inject identical `## Project Standards` into both judge prompts and fix agent.
- Two blind judges in parallel; synthesize; fix; re-judge until pass or escalate after two iterations.

### skill-creator
- Follow Agent Skills spec: `SKILL.md` + optional `assets/`; frontmatter with name, description, triggers.
- Do not duplicate existing docs; skills are for repeatable agent workflows.

### canvas
- Use for analytical artifacts and rich layouts when the deliverable *is* the structured view; skip for normal code fixes or short answers.
- Read this skill when editing `.canvas.tsx`; prefer canvas over large markdown tables for MCP-driven data.

### create-rule
- Scope rules (`alwaysApply` vs globs); confirm file patterns before writing `.mdc`.

### create-skill
- Gather purpose, location (user vs project), triggers, output format; match repo conventions.

### babysit
- Triage every PR comment; resolve conflicts only when intent is clear; fix CI with small scoped changes.

### split-to-prs
- No branch/commit/PR until user approves the plan; never destructive git without explicit approval; recoverable stash snapshot before moving work.

### update-cursor-settings
- Read `%APPDATA%\Cursor\User\settings.json` (Windows); preserve unrelated keys; validate JSON.

### create-hook
- Choose project vs user path; narrowest hook event; define fail-open vs fail-closed.

### statusline
- Configure via `~/.cursor/cli-config.json` `statusLine.command`; respect timeout and update interval.

### shell
- Only when user invokes `/shell`: run the following text literally; do not rewrite the command.

### create-subagent
- Project agents in `.cursor/agents/` override user; markdown + YAML frontmatter per Cursor format.

### migrate-to-skills
- Copy rule/command bodies verbatim; do not “improve” prose during migration.

### update-cli-config
- Merge project `.cursor/cli.json` overrides; home config is `~/.cursor/cli-config.json`.

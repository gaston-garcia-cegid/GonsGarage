# Skill Registry

**Delegator use only.** Any agent that launches sub-agents reads this registry to resolve compact rules, then injects them directly into sub-agent prompts. Sub-agents do NOT read this registry or individual SKILL.md files.

See `_shared/skill-resolver.md` (Gentleman skills) for the full resolution protocol.

## User Skills

| Trigger | Skill | Path |
|---------|-------|------|
| When creating a pull request, opening a PR, or preparing changes for review. | branch-pr | `C:\Users\gaston.garcia\.claude\skills\branch-pr\SKILL.md` |
| When creating a GitHub issue, reporting a bug, or requesting a feature. | issue-creation | `C:\Users\gaston.garcia\.claude\skills\issue-creation\SKILL.md` |
| When writing Go tests, using teatest, or adding test coverage. | go-testing | `C:\Users\gaston.garcia\.claude\skills\go-testing\SKILL.md` |
| When user says "judgment day", "judgment-day", "review adversarial", "dual review", "doble review", "juzgar", "que lo juzguen". | judgment-day | `C:\Users\gaston.garcia\.claude\skills\judgment-day\SKILL.md` |
| When user asks to create a new skill, add agent instructions, or document patterns for AI. | skill-creator | `C:\Users\gaston.garcia\.claude\skills\skill-creator\SKILL.md` |
| Standalone analytical artifact, data-heavy tables, MCP tool deliverables; create/edit `.canvas.tsx`. | canvas | `C:\Users\gaston.garcia\.cursor\skills-cursor\canvas\SKILL.md` |
| Keep a PR merge-ready: comments, conflicts, CI loop. | babysit | `C:\Users\gaston.garcia\.cursor\skills-cursor\babysit\SKILL.md` |
| Create hook, write hooks.json, automate around agent events. | create-hook | `C:\Users\gaston.garcia\.cursor\skills-cursor\create-hook\SKILL.md` |
| Create a rule, coding standards, `.cursor/rules/`, AGENTS.md. | create-rule | `C:\Users\gaston.garcia\.cursor\skills-cursor\create-rule\SKILL.md` |
| Create or author a new skill; SKILL.md format. | create-skill | `C:\Users\gaston.garcia\.cursor\skills-cursor\create-skill\SKILL.md` |
| Split a chat, branch, or PR into small PRs. | split-to-prs | `C:\Users\gaston.garcia\.cursor\skills-cursor\split-to-prs\SKILL.md` |
| Status line, statusline, CLI prompt footer. | statusline | `C:\Users\gaston.garcia\.cursor\skills-cursor\statusline\SKILL.md` |
| Change editor settings, settings.json, format on save. | update-cursor-settings | `C:\Users\gaston.garcia\.cursor\skills-cursor\update-cursor-settings\SKILL.md` |
| User invokes `/shell` — run following text literally. | shell | `C:\Users\gaston.garcia\.cursor\skills-cursor\shell\SKILL.md` |
| Migrate .mdc rules / slash commands to `.cursor/skills/`. | migrate-to-skills | `C:\Users\gaston.garcia\.cursor\skills-cursor\migrate-to-skills\SKILL.md` |
| Create subagents, `.cursor/agents/`. | create-subagent | `C:\Users\gaston.garcia\.cursor\skills-cursor\create-subagent\SKILL.md` |
| Change `~/.cursor/cli-config.json`, CLI permissions, sandbox. | update-cli-config | `C:\Users\gaston.garcia\.cursor\skills-cursor\update-cli-config\SKILL.md` |

## Compact Rules

### branch-pr

- Every PR MUST link an approved issue; every PR MUST have exactly one `type:*` label.
- Verify issue has `status:approved`; branch naming `type/description`; conventional commits; checks green before merge.

### issue-creation

- Use a template only (blank issues disabled); new issues get `status:needs-review`; maintainer adds `status:approved` before PR; questions → Discussions not issues.

### go-testing

- Prefer table-driven tests with `name` subtests; use `t.Run`; golden files when output is large/stable; integration: real boundaries + `httptest` where appropriate.

### judgment-day

- Before judges: resolve skills (Engram `skill-registry` or `.atl/skill-registry.md`), build identical `## Project Standards (auto-resolved)` for both judges + fix agent; two blind judges in parallel; synthesize, fix, re-judge ≤2 iterations or escalate.

### skill-creator

- Skill = SKILL.md + optional assets; include frontmatter name/description/trigger; skip creating a skill when docs already suffice or task is one-off.

### canvas

- Use canvas for new standalone analytical deliverables and wide MCP-driven tables; skip when user asked for a specific external tool deliverable, a code fix/PR, or targeted debugging.

### babysit

- Triage every PR comment before acting; resolve conflicts only when intent matches base; fix CI with small scoped fixes and re-watch until mergeable.

### create-hook

- Decide project vs user hook (`~/.cursor`); pick event + fail-open vs fail-closed; project hooks use paths from repo root.

### create-rule

- Put rules in `.cursor/rules/`; decide always-apply vs globs; infer from conversation when possible.

### create-skill

- Gather purpose, location (user vs project), triggers, output format; match repo patterns; use AskQuestion only when needed.

### split-to-prs

- No branch/commit/PR until user approves plan; no destructive git without approval; snapshot first; stage named files/hunks only, never `git add .`.

### statusline

- Configure in `~/.cursor/cli-config.json` → `statusLine.command`; stdin JSON per Cursor spec.

### update-cursor-settings

- Read `%APPDATA%\Cursor\User\settings.json` (Windows) first; change only requested keys; validate JSON.

### shell

- Only on `/shell`: treat remainder as literal command; run as-is; no repo inspection unless the command needs it.

### migrate-to-skills

- Copy rule/command body verbatim to SKILL.md; ignore `~/.cursor/worktrees` and `skills-cursor` sources.

### create-subagent

- Project `.cursor/agents/` overrides user; isolated prompts for specialized tasks.

### update-cli-config

- Edit `~/.cursor/cli-config.json`; project `.cursor/cli.json` merges from git root downward; restart CLI to apply.

## Project Conventions

| File | Path | Notes |
|------|------|-------|
| Next.js App Router navigation | `d:\Repos\GonsGarage\.cursor\rules\nextjs-app-router-navigation.mdc` | `alwaysApply: true` — Link/router targets must exist under `frontend/src/app/` |

Read the convention files listed above for project-specific patterns and rules.

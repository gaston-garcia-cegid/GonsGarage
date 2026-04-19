# Skill Registry

**Delegator use only.** Any agent that launches sub-agents reads this registry to resolve compact rules, then injects them directly into sub-agent prompts. Sub-agents do NOT read this registry or individual SKILL.md files.

See `_shared/skill-resolver.md` (en paquete de skills del usuario) para el protocolo completo.

## User Skills

| Trigger | Skill | Path |
|---------|-------|------|
| When creating a pull request, opening a PR, or preparing changes for review. | branch-pr | `C:\Users\gaston.garcia\.cursor\skills\branch-pr\SKILL.md` |
| When creating a GitHub issue, reporting a bug, or requesting a feature. | issue-creation | `C:\Users\gaston.garcia\.cursor\skills\issue-creation\SKILL.md` |
| When writing Go tests, using teatest, or adding test coverage. | go-testing | `C:\Users\gaston.garcia\.cursor\skills\go-testing\SKILL.md` |
| When user says "judgment day", "judgment-day", "review adversarial", "dual review", "doble review", "juzgar", "que lo juzguen". | judgment-day | `C:\Users\gaston.garcia\.cursor\skills\judgment-day\SKILL.md` |
| When the user asks to create a new skill, add agent instructions, or document patterns for AI. | skill-creator | `C:\Users\gaston.garcia\.cursor\skills\skill-creator\SKILL.md` |
| Standalone analytical artifacts, data-heavy deliverables, MCP tool results as primary output, or any `.canvas.tsx` work. | canvas | `C:\Users\gaston.garcia\.cursor\skills-cursor\canvas\SKILL.md` |
| Keep a PR merge-ready: triage comments, resolve conflicts, fix CI loop. | babysit | `C:\Users\gaston.garcia\.cursor\skills-cursor\babysit\SKILL.md` |
| Create, write, or author a new skill; SKILL.md format and best practices. | create-skill | `C:\Users\gaston.garcia\.cursor\skills-cursor\create-skill\SKILL.md` |
| Create a rule, `.cursor/rules/`, AGENTS.md, coding standards. | create-rule | `C:\Users\gaston.garcia\.cursor\skills-cursor\create-rule\SKILL.md` |
| Create hooks, `hooks.json`, automate around agent events. | create-hook | `C:\Users\gaston.garcia\.cursor\skills-cursor\create-hook\SKILL.md` |
| Change editor `settings.json`, themes, format on save, keybindings. | update-cursor-settings | `C:\Users\gaston.garcia\.cursor\skills-cursor\update-cursor-settings\SKILL.md` |
| Change `~/.cursor/cli-config.json`, permissions, CLI display, sandbox. | update-cli-config | `C:\Users\gaston.garcia\.cursor\skills-cursor\update-cli-config\SKILL.md` |
| Status line / prompt footer / CLI status bar. | statusline | `C:\Users\gaston.garcia\.cursor\skills-cursor\statusline\SKILL.md` |
| Migrate `.cursor/rules` (intelligent) and `.cursor/commands` to skills. | migrate-to-skills | `C:\Users\gaston.garcia\.cursor\skills-cursor\migrate-to-skills\SKILL.md` |
| Create custom subagents under `.cursor/agents/` or `~/.cursor/agents/`. | create-subagent | `C:\Users\gaston.garcia\.cursor\skills-cursor\create-subagent\SKILL.md` |
| User invokes `/shell` — run following text literally as shell. | shell | `C:\Users\gaston.garcia\.cursor\skills-cursor\shell\SKILL.md` |

## Compact Rules

### branch-pr
- Every PR MUST link an approved issue (`Closes`/`Fixes`/`Resolves #N`); exactly one `type:*` label.
- Branch names MUST match `^(feat|fix|chore|docs|style|refactor|perf|test|build|ci|revert)\/[a-z0-9._-]+$`.
- Use `.github/PULL_REQUEST_TEMPLATE.md`; run shellcheck on modified scripts; checks green before merge.

### issue-creation
- Use templates only (blank issues disabled); fill required fields; search duplicates first.
- New issues get `status:needs-review`; maintainer adds `status:approved` before any PR.
- Questions → Discussions, not issues.

### go-testing
- Prefer table-driven tests with `t.Run` per case; assert errors with `(err != nil) != wantErr` pattern.
- Bubbletea: test `Update` transitions on model state directly; use `teatest` for interactive flows when applicable.
- Golden files: commit expected output; update goldens deliberately when behavior changes.

### judgment-day
- Before judges: resolve skill registry → build identical `## Project Standards (auto-resolved)` for both judges + fix agent.
- Launch two blind judges in parallel (same target, no cross-hints); orchestrator synthesizes Confirmed / Suspect / Contradiction.
- Classify warnings as real vs theoretical; only real warnings drive fixes and re-judge (max 2 iterations or escalate).

### skill-creator
- Skills need `SKILL.md` with frontmatter (`name`, `description` + Trigger, `license`, `metadata.version`).
- Skip creating a skill when docs already cover it or the workflow is one-off trivial.

### canvas
- Use canvas only when output is a **standalone artifact** (analysis, audit, heavy tables); not for normal code fixes or short answers.
- One `.canvas.tsx` per feature; default export; imports only from `cursor/canvas`; no `fetch`, no extra files.
- Read `~/.cursor/skills-cursor/canvas/sdk/*.d.ts` for real exports; flat minimal UI (no gradients/emojis/heavy shadows).

### babysit
- Triage every PR comment (including bots); fix only what you agree with; explain disagreements.
- Resolve conflicts only when intent is clear; otherwise stop for clarification.
- Fix CI with small scoped changes and re-watch until green and mergeable.

### create-skill
- Gather purpose, location (`~/.cursor/skills` vs project), triggers, output format before writing.
- Layout: `skill-name/SKILL.md` (+ optional `reference.md`, `examples.md`, `scripts/`).

### create-rule
- Rules are `.mdc` in `.cursor/rules/` with frontmatter (`description`, optional `globs`, `alwaysApply`).
- Clarify always-on vs file-scoped; get concrete glob patterns before writing.

### create-hook
- Project: `.cursor/hooks.json` + `.cursor/hooks/*`; user: `~/.cursor/hooks.json` + `~/.cursor/hooks/*`.
- Pick narrowest event (`preToolUse`, `beforeShellExecution`, etc.); decide fail-open vs fail-closed.

### update-cursor-settings
- Windows path: `%APPDATA%\Cursor\User\settings.json`; read first, merge minimally, preserve JSON with comments.

### update-cli-config
- Global `~/.cursor/cli-config.json`; project layers via `.cursor/cli.json` merged from repo root downward.
- Restart CLI after edits; respect `permissions.allow` / `permissions.deny` shape.

### statusline
- Configure `statusLine` in `~/.cursor/cli-config.json` with `type: command`, script path, optional `padding` / `updateIntervalMs` (>=300).

### migrate-to-skills
- Copy body **verbatim** — no reformat or "improve" when converting `.mdc` / commands to `SKILL.md`.
- Migrate rules that have `description` but no `globs` and not `alwaysApply: true`; migrate all slash commands.

### create-subagent
- Project `.cursor/agents/*.md` overrides user `~/.cursor/agents/*.md` on name collision.
- Frontmatter: `name`, `description`; body is the system prompt.

### shell
- Only when user explicitly uses `/shell`: run the text after it as a literal command; do not rewrite or pre-analyze.

## Project Conventions

| File | Path | Notes |
|------|------|-------|
| — | — | No `AGENTS.md`, `CLAUDE.md`, `.cursorrules`, or `copilot-instructions.md` en la raíz del repo. |

Read stack-specific rules from `openspec/config.yaml` y documentación en `docs/` cuando apliquen.

# Skill Registry

**Delegator use only.** Any agent that launches sub-agents reads this registry to resolve compact rules, then injects them directly into sub-agent prompts. Sub-agents do NOT read this registry or individual SKILL.md files.

See `_shared/skill-resolver.md` for the full resolution protocol.

## User Skills

| Trigger | Skill | Path |
|---------|-------|------|
| When writing Go tests, using teatest, or adding test coverage | go-testing | C:\Users\gaston.garcia\.claude\skills\go-testing\SKILL.md |
| When creating a pull request, opening a PR, or preparing changes for review | branch-pr | C:\Users\gaston.garcia\.claude\skills\branch-pr\SKILL.md |
| When creating a GitHub issue, reporting a bug, or requesting a feature | issue-creation | C:\Users\gaston.garcia\.claude\skills\issue-creation\SKILL.md |
| judgment day, judgment-day, review adversarial, dual review, doble review, juzgar, que lo juzguen | judgment-day | C:\Users\gaston.garcia\.claude\skills\judgment-day\SKILL.md |
| When user asks to create a new skill, add agent instructions, or document patterns for AI | skill-creator | C:\Users\gaston.garcia\.claude\skills\skill-creator\SKILL.md |
| User explicitly wants Cursor Automations | automate | C:\Users\gaston.garcia\.cursor\skills-cursor\automate\SKILL.md |
| Keep a PR merge-ready (comments, conflicts, CI) | babysit | C:\Users\gaston.garcia\.cursor\skills-cursor\babysit\SKILL.md |
| Standalone analytical artifact, data-heavy output, MCP tool results | canvas | C:\Users\gaston.garcia\.cursor\skills-cursor\canvas\SKILL.md |
| Create Cursor rules, .cursor/rules/, AGENTS.md | create-rule | C:\Users\gaston.garcia\.cursor\skills-cursor\create-rule\SKILL.md |
| Authoring new skill or SKILL.md structure | create-skill | C:\Users\gaston.garcia\.cursor\skills-cursor\create-skill\SKILL.md |
| Split work into small reviewable PRs | split-to-prs | C:\Users\gaston.garcia\.cursor\skills-cursor\split-to-prs\SKILL.md |
| Cursor SDK (@cursor/sdk, cursor-sdk, Agent.create, run.stream) | sdk | C:\Users\gaston.garcia\.cursor\skills-cursor\sdk\SKILL.md |
| User invokes /shell | shell | C:\Users\gaston.garcia\.cursor\skills-cursor\shell\SKILL.md |
| /loop recurring prompt | loop | C:\Users\gaston.garcia\.cursor\skills-cursor\loop\SKILL.md |
| Custom CLI status line | statusline | C:\Users\gaston.garcia\.cursor\skills-cursor\statusline\SKILL.md |
| Cursor hooks (hooks.json) | create-hook | C:\Users\gaston.garcia\.cursor\skills-cursor\create-hook\SKILL.md |
| settings.json, editor preferences | update-cursor-settings | C:\Users\gaston.garcia\.cursor\skills-cursor\update-cursor-settings\SKILL.md |

## Compact Rules

Pre-digested rules per skill. Delegators copy matching blocks into sub-agent prompts as `## Project Standards (auto-resolved)`.

### go-testing
- Prefer table-driven tests with `t.Run(tt.name, ...)` for multiple cases
- Use `testify` (`assert`/`require`) already in go.mod — match existing test style in `backend/internal/**/**_test.go`
- Handler tests: `net/http/httptest` + Gin test context; see `backend/internal/handler/*_test.go`
- Integration flows: `backend/tests/integration/` with real HTTP against test server + sqlite/postgres as peers do
- Repo tests: glebarez/sqlite driver for local DB; follow `*_repository_test.go` patterns
- Run: `cd backend && go test ./... -count=1`; CI uses `-race -timeout=2m`
- Never skip error-path cases when adding domain validation tests

### branch-pr
- Every PR MUST link an approved issue with `status:approved` (Agent Teams Lite repos)
- Exactly one `type:*` label on PR
- Branch naming: `^(feat|fix|chore|docs|style|refactor|perf|test|build|ci|revert)/[a-z0-9._-]+$`
- Use `gh pr create` with template; wait for CI green before merge
- GonsGarage CI: backend go vet/test + frontend lint/typecheck/test/build

### issue-creation
- Use GitHub issue templates — blank issues disabled in ATL repos
- New issues get `status:needs-review`; maintainer adds `status:approved` before PR
- Search duplicates first; questions → Discussions not issues

### judgment-day
- Load skill registry (engram or `.atl/skill-registry.md`) before launching judges
- Launch TWO blind judge sub-agents in parallel — orchestrator coordinates only
- Inject matching compact rules into both judges and fix agent
- Re-judge after fixes; escalate after 2 iterations if both still fail

### skill-creator
- Skills live in `skills/{name}/SKILL.md` with YAML frontmatter (`name`, `description` with Trigger)
- Keep SKILL.md focused; use `references/` for long docs
- Don't create skills for trivial one-off tasks or existing docs

### automate
- Only for explicit Cursor Automations requests — not generic CI/scripts
- Plain language in user chat; no MCP/proto names
- Finish via Automations editor handoff after user approves draft table

### babysit
- Resolve merge conflicts preserving branch intent; ask if intents conflict
- Triage unresolved PR comments; validate Bugbot before acting
- Fix CI within PR scope only — never weaken workflows to pass
- Merge latest base if failures seem unrelated

### canvas
- Use for standalone analytical artifacts (tables, charts, timelines, MCP data dumps)
- Single `.canvas.tsx` beside chat — not for code fixes or draft messages
- Read canvas skill when editing `.canvas.tsx`

### create-rule
- Rules in `.cursor/rules/*.mdc` with frontmatter (`description`, `alwaysApply` or globs)
- One concern per rule; reference real bug examples when enforcing patterns

### create-skill
- Personal: `~/.cursor/skills/`; project: `.cursor/skills/` or repo `skills/`
- Frontmatter `description` must include Trigger line
- Respect user verbatim wording when provided

### split-to-prs
- Never branch/commit/push until user approves split plan
- Snapshot before moving work; stage named files only — no `git add .`
- Split by reviewer/ownership boundaries; stack only when dependency is real

### sdk
- TypeScript `@cursor/sdk` or Python `cursor-sdk` — read skill for current API, don't guess
- Agent → Run model; local vs cloud runtime choice matters for cwd/repo
- Handle streaming, cancellation, CursorAgentError explicitly

### shell
- Only when user invokes `/shell` — run following text literally, no rewriting
- Report exit status and key stdout/stderr after run

### loop
- Recurring prompt execution at user-specified interval
- Read loop skill for /loop syntax and persistence

### statusline
- Customize Cursor CLI status line via skill instructions
- Read skill for config file locations

### create-hook
- Cursor hooks in hooks.json; read skill for event types and script layout

### update-cursor-settings
- Modify settings.json values only user requested
- Preserve existing keys; validate JSON before write

## Project Conventions

| File | Path | Notes |
|------|------|-------|
| Next.js navigation rule | `.cursor/rules/nextjs-app-router-navigation.mdc` | Always apply — verify `page.tsx` exists before Link/router.push |
| SDD config | `openspec/config.yaml` | strict_tdd, testing commands, SDD rules |
| CI | `.github/workflows/ci.yml` | Backend + frontend gates |

Read the convention files listed above for project-specific patterns and rules.

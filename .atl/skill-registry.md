# Skill Registry

**Delegator use only.** Any agent that launches sub-agents reads this registry to resolve compact rules, then injects them directly into sub-agent prompts. Sub-agents do NOT read this registry or individual SKILL.md files.

See `_shared/skill-resolver.md` in the user skills tree for the full resolution protocol.

Skills `sdd-*`, `_shared`, and `skill-registry` are omitted from the table (SDD vía `/sdd-*`; el catálogo de skill-registry es esta propia tarea).

## User Skills

| Trigger | Skill | Path |
|---------|-------|------|
| PR, pull request, branch-pr | branch-pr | `~/.cursor/skills/branch-pr/SKILL.md` |
| GitHub issue, bug, feature request | issue-creation | `~/.cursor/skills/issue-creation/SKILL.md` |
| Go tests, teatest, Bubbletea TUI tests | go-testing | `~/.cursor/skills/go-testing/SKILL.md` |
| judgment day, dual review, adversarial review | judgment-day | `~/.cursor/skills/judgment-day/SKILL.md` |
| create skill, Agent Skills, SKILL.md | skill-creator | `~/.cursor/skills/skill-creator/SKILL.md` |
| Canvas, datos pesados, MCP tables | canvas | `~/.cursor/skills-cursor/canvas/SKILL.md` |
| Cursor hooks, hooks.json | create-hook | `~/.cursor/skills-cursor/create-hook/SKILL.md` |
| Cursor rules, RULE.md, .cursor/rules | create-rule | `~/.cursor/skills-cursor/create-rule/SKILL.md` |
| crear skill Cursor | create-skill | `~/.cursor/skills-cursor/create-skill/SKILL.md` |
| subagent, custom subagent | create-subagent | `~/.cursor/skills-cursor/create-subagent/SKILL.md` |
| migrate rules, .mdc to skills | migrate-to-skills | `~/.cursor/skills-cursor/migrate-to-skills/SKILL.md` |
| /shell | shell | `~/.cursor/skills-cursor/shell/SKILL.md` |
| statusline, CLI prompt footer | statusline | `~/.cursor/skills-cursor/statusline/SKILL.md` |
| settings.json, editor preferences | update-cursor-settings | `~/.cursor/skills-cursor/update-cursor-settings/SKILL.md` |
| cli-config.json, CLI permissions | update-cli-config | `~/.cursor/skills-cursor/update-cli-config/SKILL.md` |
| babysit, PR merge-ready, CI loop | babysit | `~/.cursor/skills-cursor/babysit/SKILL.md` |

## Compact Rules

Pre-digested rules per skill. Delegators copy matching blocks into sub-agent prompts as `## Project Standards (auto-resolved)`.

### branch-pr
- Issue-first: enlazar el issue en la descripción del PR; seguir el flujo del SKILL para Agent Teams Lite.

### issue-creation
- Issue con contexto reproducible, pasos, resultado esperado vs actual; encaje con enforcement issue-first del SKILL.

### go-testing
- Preferir table-driven en servicios; testify; desde `backend/` usar `go test ./...`; ver SKILL para patrones Gin/teatest si aplica.

### judgment-day
- Dos revisores ciegos en paralelo; sintetizar, corregir, re-juzgar hasta pass o escalar según el SKILL.

### skill-creator
- Seguir Agent Skills spec: frontmatter con `name`, `description` y triggers en `Trigger:`.

### canvas
- Para análisis cuantitativos, tablas densas o herramientas interactivas: canvas en lugar de volcar solo markdown.

### create-hook / create-rule / create-skill
- Respetar rutas y convenciones Cursor del SKILL; triggers explícitos en descripción.

### create-subagent
- Subagentes en rutas documentadas en el SKILL; prompts enfocados y reutilizables.

### migrate-to-skills
- Copiar cuerpo de reglas/comandos verbatim al migrar; no reescribir ni “mejorar” el contenido fuente.

### shell
- Solo con invocación explícita `/shell`: ejecutar el texto subsiguiente literalmente, sin reinterpretar.

### statusline / update-cursor-settings / update-cli-config
- Editar solo los archivos indicados en el SKILL; reinicio CLI si aplica tras cambiar `cli-config.json`.

### babysit
- Bucle: triage de comentarios del PR, conflictos claros, CI verde, hasta estado merge-ready.

## Project Conventions

| File | Path | Notes |
|------|------|-------|
| README | `README.md` | Arranque, API, Swagger |
| Development guide | `docs/development-guide.md` | pnpm, Go, compose |
| MVP phases | `docs/mvp-minimum-phases.md` | Alcance MVP |
| i18n reminder | `docs/i18n-reminder.md` | Futuro i18n |
| Contributing | `CONTRIBUTING.md` | go vet, tests |

No hay `AGENTS.md` / `.cursorrules` en la raíz del repo; convenciones principales en `docs/` y README.

---

*Engram (`mem_save`): no disponible en esta sesión. Re-ejecutar sdd-init o skill-registry tras instalar skills nuevas.*

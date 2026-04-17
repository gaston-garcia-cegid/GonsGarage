# Skill registry — GonsGarage

Generado por **sdd-init** (2026-04-16). Las skills `sdd-*` y `_shared` se omiten del catálogo de trabajo; siguen disponibles vía comandos `/sdd-*`.

## Convenciones del proyecto

| Fuente | Notas |
|--------|--------|
| `README.md` | Arranque, API, Swagger |
| `docs/development-guide.md` | pnpm, Go, compose |
| `docs/mvp-minimum-phases.md` | Alcance MVP |
| `docs/i18n-reminder.md` | Futuro i18n |

## Skills detectadas (Cursor user + skills-cursor)

| Trigger (resumen) | Nombre | Ruta |
|--------------------|--------|------|
| PR, branch-pr | branch-pr | `~/.cursor/skills/branch-pr/` |
| GitHub issue | issue-creation | `~/.cursor/skills/issue-creation/` |
| Go tests, teatest | go-testing | `~/.cursor/skills/go-testing/` |
| judgment day, dual review | judgment-day | `~/.cursor/skills/judgment-day/` |
| nueva skill | skill-creator | `~/.cursor/skills/skill-creator/` |
| update skills, registry | skill-registry | `~/.cursor/skills/skill-registry/` |
| Canvas / datos pesados | canvas | `~/.cursor/skills-cursor/canvas/` |
| hooks | create-hook | `~/.cursor/skills-cursor/create-hook/` |
| rules, AGENTS.md | create-rule | `~/.cursor/skills-cursor/create-rule/` |
| crear skill | create-skill | `~/.cursor/skills-cursor/create-skill/` |
| settings.json | update-cursor-settings | `~/.cursor/skills-cursor/update-cursor-settings/` |
| PR merge-ready | babysit | `~/.cursor/skills-cursor/babysit/` |
| statusline CLI | statusline | `~/.cursor/skills-cursor/statusline/` |

## Reglas compactas (sub-agentes)

### branch-pr
- Flujo centrado en issue antes del PR; enlazar issue en descripción del PR.

### issue-creation
- Issue con contexto reproducible, pasos, resultado esperado vs actual.

### go-testing
- Tests table-driven en services; testify; `go test ./...` desde `backend/`; seguir SKILL para Gin/teatest.

### judgment-day
- Lanzar dos revisores ciegos en paralelo; iterar o escalar según el SKILL.

### canvas
- Usar canvas para artefactos analíticos/tablas interactivas; no volcar solo markdown masivo.

### create-rule / create-skill
- Seguir formato Agent Skills; triggers claros en frontmatter.

### babysit
- Bucle: comentarios CI, conflictos claros, hasta merge-ready.

---

*Engram: no disponible en esta sesión; no se llamó `mem_save`. Re-ejecutar `skill-registry` tras instalar skills nuevas.*

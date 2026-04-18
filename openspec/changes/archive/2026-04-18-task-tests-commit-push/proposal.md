# Proposal: Task «tests completos → comentario → commit → push»

## Intent

Hoy no hay un único disparador reproducible que ejecute **toda** la verificación de tests (frontend + backend), deje constancia en un **comentario** (trazabilidad / PR), y encadene **commit + push**. Los desarrolladores repiten comandos a mano o confían solo en CI remota; se quiere una **Task** (VS Code / Cursor) o equivalente documentado en el repo que automatice esa secuencia con riesgo acotado.

## Scope

### In Scope

- Añadir una **VS Code Task** compuesta (`.vscode/tasks.json`) que, en orden: ejecute tests del **backend** (`go test ./...` desde `backend/`) y del **frontend** (`pnpm test` desde `frontend/`); si alguno falla, **detener** la cadena sin commit ni push.
- Tras tests OK: generar el **comentario** (Markdown: fecha, comandos, OK) e inyectarlo en el **mensaje de commit**; opcionalmente el mismo texto vía **`gh pr comment`** si aplica.
- Tras generar el mensaje: **`git add`**, **`git commit`** (con cuerpo multi-línea que incluye el comentario) y **`git push`** — vía tasks encadenadas y script (PowerShell en Windows + notas para POSIX en `design`).

### Out of Scope

- Cambiar lógica de producto ni tests existentes.
- Forzar comentario en GitHub sin `gh` ni contexto de PR (solo documentar cómo hacerlo).
- Ejecutar en agentes cloud sin acceso git: la Task es **local** al workspace.

## Capabilities

### New Capabilities

- None — no hay requisitos de producto nuevos en `openspec/specs/`; solo flujo de desarrollo en el repo.

### Modified Capabilities

- None

## Approach

1. **`.vscode/tasks.json`**: task compuesta en secuencia — backend `go test`, frontend `pnpm test`; propagar fallo y no ejecutar pasos posteriores si algo rompe.
2. **«Comentario»**: generar texto Markdown (timestamp, suites, resultado) y usarlo como (a) **cuerpo extendido del mensaje de `git commit`** (p. ej. `chore: verify all tests` + bloque *Test report*), y (b) opcionalmente el mismo texto con **`gh pr comment`** si hay `gh` y rama con PR (`design` fija variables/env).
3. **Commit / push**: tras tests OK, `git add -A` (o lista acotada en design) + `git commit` con el mensaje anterior + `git push` al remoto por defecto; **input** de VS Code para confirmar mensaje base y desactivar push en un perfil «solo verify».

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `.vscode/tasks.json` | New | Task compuesta verify + opcional commit/push |
| `scripts/` | New | PowerShell o script verify |
| `.gitignore` | Modified | Ignorar `.local/` si se usa |
| `CONTRIBUTING.md` | Modified | Una sección «Verify task» |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Push accidental | Med | Task separada «solo verify»; push en task opt-in con `input` confirmación |
| `gh` no instalado | Med | Documentar; fallar en silencio en paso PR comment |

## Rollback Plan

Borrar `.vscode/tasks.json` entradas añadidas y scripts; revert commit único.

## Dependencies

- Node/pnpm, Go, git en PATH; opcional `gh` para comentario en PR.

## Success Criteria

- [ ] Desde VS Code/Cursor se ejecuta una Task que corre **todos** los tests front y back y falla en el primer error.
- [ ] Queda un **comentario** generado (en el cuerpo del commit y, si se configura, en el PR) con resumen verificable.
- [ ] Documentado cómo encadenar commit + push sin sobrescribir trabajo no deseado.

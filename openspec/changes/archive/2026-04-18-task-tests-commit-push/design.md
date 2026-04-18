# Design: Task «tests completos → comentario → commit → push»

## Resolución de preguntas abiertas

| Pregunta | Decisión |
|----------|----------|
| ¿Incluir `go vet`, `pnpm lint`, `pnpm typecheck` en la misma cadena? | **No.** La cadena de esta Task incluye **solo tests** (backend `go test`, frontend `pnpm test`), como en la propuesta literal. Vet/lint/typecheck siguen en CI y en la sección manual de `CONTRIBUTING.md`. |
| ¿Input de VS Code para el asunto? | **Sí, ahora.** Dos `promptString`: **asunto** (primera línea del commit) y **explicación** (cuerpo humano). El script añade automáticamente bloques **Verificación** (tests pasados, timestamp UTC) y **Análisis git** (`git status -sb`, `diff --cached --stat`, `diff --stat` no staged como referencia). |

## Technical Approach

1. **`dev: verify (tests only)`** — tareas shell en secuencia: backend (`CGO_ENABLED=1`, `go test ./... -count=1 -race -timeout=2m`) y frontend (`pnpm test -- --passWithNoTests`). Sin git.

2. **`dev: verify + commit + push`** — depende en secuencia de `dev: verify (tests only)`; ejecuta script PowerShell (Windows) / sh (POSIX) con **`-SkipTests`** para no duplicar tests. El script: compone el mensaje con inputs (`GONS_COMMIT_SUBJECT`, `GONS_COMMIT_EXPLANATION`) + análisis git; **`git commit -F`** si hay **cambios en stage**; si el índice está vacío, error claro salvo `GONS_ALLOW_EMPTY_COMMIT=1` (commit vacío de marcador); luego `git push` si no `GONS_SKIP_PUSH=1`. Opcional: `gh pr comment --body-file` si `GONS_GH_PR=1`.

## Architecture Decisions

| Decision | Alternatives | Choice & rationale |
|----------|--------------|-------------------|
| Alcance verify | Incluir lint/vet | **Solo tests** — respuesta explícita del product owner. |
| Staging | `git add -A` | **Solo lo ya en stage** — el mensaje describe exactamente `diff --cached`; el usuario hace `git add` antes de la Task. |
| Commit vacío | Prohibido | Permitido solo con **`GONS_ALLOW_EMPTY_COMMIT=1`** (documentado). |
| Duplicar tests en task larga | Siempre en script | **`dependsOn` verify** + script `-SkipTests` — una sola ejecución de tests por invocación. |

## Sequence

```
dev: verify + commit + push
  │
  ├─► dev: verify (tests only)  ──fail──► STOP
  │
  └─► script (-SkipTests)
         ├─► compose message (subject + explanation + tests OK + git analysis)
         ├─► staged?  git commit -F
         │      else  allow-empty env?  else  ERROR
         ├─► (opt) gh pr comment
         └─► git push
```

## File Changes

| File | Action |
|------|--------|
| `.vscode/tasks.json` | Create |
| `scripts/dev-verify-and-commit.ps1` | Create |
| `scripts/dev-verify-and-commit.sh` | Create |
| `CONTRIBUTING.md` | Modify — sección Tasks VS Code |

## Variables de entorno

| Variable | Efecto |
|----------|--------|
| `GONS_COMMIT_SUBJECT` | Primera línea (VS Code `input`) |
| `GONS_COMMIT_EXPLANATION` | Cuerpo explicativo (VS Code `input`) |
| `GONS_ALLOW_EMPTY_COMMIT` | `1` → permite `git commit --allow-empty` si no hay stage |
| `GONS_SKIP_PUSH` | `1` → no push |
| `GONS_GH_PR` | `1` → intenta `gh pr comment --body-file` |

## Testing Strategy

Manual: ejecutar verify only con test roto; ejecutar pipeline con/sin staged; revisar `git log -1 --format=full`.

## Migration / Rollout

Documentar en `CONTRIBUTING.md` que hay que **stagear** antes del commit automatizado.

## Open Questions

- None (cerradas arriba).

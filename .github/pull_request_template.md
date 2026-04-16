## Qué cambia

<!-- Breve descripción para revisores. -->

## Tipo de cambio

- [ ] Corrección
- [ ] Funcionalidad
- [ ] Documentación
- [ ] Chore / refactor / CI

## Issue

<!-- Si aplica: Closes #123 -->

## Verificación

- [ ] Backend: `go test ./...` y `go vet ./...` (desde `backend/`)
- [ ] Frontend: `pnpm lint`, `pnpm typecheck`, `pnpm test` (desde `frontend/`) si toca UI o cliente
- [ ] OpenAPI: `swag@v1.8.12` regenerado si cambiaron anotaciones de handlers

# Contribuir a GonsGarage

## Herramientas

- **Go**: versión indicada en `backend/go.mod`.
- **Node.js**: 22+ recomendado.
- **pnpm**: 9.x (activar con `corepack enable` si usas `packageManager` del `frontend/package.json`).
- **Docker**: para Postgres + Redis locales (`docker-compose.yml` en la raíz).

## Instalación rápida

```powershell
# Raíz del repo
docker compose up -d

cd backend
copy .env.example .env   # o cp en bash
go run ./cmd/api

cd ../frontend
copy .env.local.example .env.local
pnpm install
pnpm dev
```

## TDD obligatorio

Antes de implementar una funcionalidad nueva o corregir un bug con impacto en comportamiento:

1. Escribir o actualizar un **test automatizado** que describa el comportamiento deseado.
2. Verificar que falla por la razón correcta.
3. Implementar el cambio mínimo para que pase.
4. Refactorizar manteniendo los tests verdes.

Detalle y convenciones: [docs/testing-tdd.md](docs/testing-tdd.md).

## Comprobaciones antes de abrir PR

```powershell
cd backend && go vet ./... && go test ./... -count=1
cd ../frontend && pnpm lint && pnpm typecheck && pnpm test
```

## Commits

Se recomienda [Conventional Commits](https://www.conventionalcommits.org/) (`feat:`, `fix:`, `test:`, `docs:`, etc.), alineado con el flujo descrito en [Agent.md](Agent.md).

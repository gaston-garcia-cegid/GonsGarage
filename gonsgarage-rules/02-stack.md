# §2 — Stack tecnológico

| Componente | Elección plantilla | Notas repo |
|--------------|-------------------|------------|
| Go | **1.25** (`go.mod`) | Alineado. |
| HTTP | **Gin** | Alineado. |
| BD | **PostgreSQL 16**, **sqlx**, driver **pq** | **Deuda:** acceso principal sigue siendo GORM; **Fase 2 en curso:** dependencias `sqlx` + `lib/pq` y paquete `internal/platform/sqlxdb` (`Open`) para nuevas rutas (ver `docs/template-adoption-plan.md` Fase 2). |
| Migraciones | **golang-migrate**, SQL en `migrations/` | Alineado (rutas bajo `backend/db/migrations` o convención acordada). |
| Redis | **go-redis v8** + miniredis | **Desviación aprobada:** repo usa **go-redis v9**; mantener v9 salvo requisito estricto de paridad Arnela. |
| Node | **22+**, **pnpm**, CI `--frozen-lockfile` | Alineado. |
| Frontend | **Next.js 16**, **React 19**, **TypeScript 5.9+** | Subir `typescript` a 5.9+ en `package.json` cuando CI lo exija. |
| CSS | **Tailwind v4** + `@tailwindcss/postcss` | Alineado. |
| Logging | **Zerolog** | **Deuda:** migrar desde slog (plan Fase 4). |
| Config | **godotenv** en dev; sin secretos en repo | Alineado en espíritu. |
| Probes | `GET /health`, `GET /readiness` | Deben existir en el binario API. |

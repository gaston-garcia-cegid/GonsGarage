# §3 — Librerías y herramientas

## Backend (Go)

- Auth: **JWT** `golang-jwt/jwt/v5`; passwords **golang.org/x/crypto**.
- HTTP/CORS: **Gin** + `gin-contrib/cors`.
- Validación: cadena Gin / `go-playground/validator`.
- IDs: `google/uuid`.
- OpenAPI: `swag`, `gin-swagger`, `files`.
- Tests: **testify** (assert + mock); tests **table-driven** en **services**.

## Frontend (TypeScript)

- UI: **Radix** + **Shadcn** (`components/ui`).
- CSS utils: `class-variance-authority`, `clsx`, `tailwind-merge`, `tailwindcss-animate`.
- Iconos: `lucide-react`.
- Tema: `next-themes`.
- Formularios: `react-hook-form` + `zod` + `@hookform/resolvers`.
- Fechas: `date-fns` + `react-day-picker`.
- Toasts: `sonner`.
- Opcional (no añadir sin decisión explícita): `xlsx`, `@sentry/browser`.
- Lint/types: ESLint 9 + `eslint-config-next` + `typescript-eslint`.
- Tests: **Vitest** + Testing Library + jsdom.

## Infra / calidad

- CI: GitHub Actions (`backend` / `frontend`, `working-directory`).
- Contenedores: Docker + Compose; `.env.example` en backend y frontend.

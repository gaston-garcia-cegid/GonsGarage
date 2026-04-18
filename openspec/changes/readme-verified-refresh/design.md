# Design: readme-verified-refresh

## Language

- **Root `README.md`:** English (GitHub-facing, single language; old file mixed EN/ES).

## Structure (new README)

1. Title, one-line description, badges tied to **pinned or directive** versions from `backend/go.mod` and `frontend/package.json`.
2. **Docs hub** → `docs/README.md` for architecture, TDD, roadmap.
3. **Verified stack** table (Go, Gin, Postgres/Redis images from `docker-compose.yml`, Next, React, pnpm, Vitest).
4. **Prerequisites** aligned with CI (`Node 22`, `pnpm 9`, Go from `go.mod`).
5. **Quick start:** `docker compose`, `go run ./cmd/api` from `backend`, `pnpm dev` from `frontend`; `.env` copy commands. **No** `generate-types` unless script exists.
6. **Local URLs** + Swagger regen command (one block, same as `docs/development-guide.md`).
7. **Test commands** mirroring `.github/workflows/ci.yml`.
8. **Sample users / seed** — short, link to `docs/development-guide.md` for detail if needed.
9. **Contributing** + **License** — links only.

## Styling claim

- Do **not** claim Tailwind unless listed in `frontend/package.json`. Document **CSS Modules** + shared styles under `frontend/src/styles/` (tokens/utilities).

## Follow-up (post-README)

- `docs/development-guide.md` y `Agent.md`: versión Go y Postgres/Redis alineadas a `go.mod` y `docker-compose.yml`.

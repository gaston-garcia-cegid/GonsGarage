# GonsGarage frontend

Next.js **15.5** (App Router), React **19**, TypeScript, Zustand. Package manager: **pnpm 9** (see `packageManager` in [`package.json`](./package.json)).

- **Monorepo overview & setup:** [../README.md](../README.md)
- **Detailed dev guide:** [../docs/development-guide.md](../docs/development-guide.md)
- **API client notes:** [docs/api-client.md](./docs/api-client.md)

## Commands

```bash
pnpm install
pnpm dev          # http://localhost:3000
pnpm lint
pnpm typecheck
pnpm test         # Vitest (default; same as CI)
pnpm build
```

Copy [`./.env.local.example`](./.env.local.example) to `.env.local` and set `NEXT_PUBLIC_API_URL` if your API is not at `http://localhost:8080`.

## Create Next App

This app was originally bootstrapped with [`create-next-app`](https://nextjs.org/docs/app/api-reference/cli/create-next-app). For framework docs see [Next.js Documentation](https://nextjs.org/docs).

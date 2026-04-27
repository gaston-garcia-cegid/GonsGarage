# Proposal: UI homogeneity — modal parity (peças + taller)

## Intent

Unify **layout/CSS/UX** with the rest of the staff app: **list + toolbar + modal** for creates (same pattern as accounting `?create=1` redirects and appointments). Fix **Taller → visita** flows: **Nova visita** should not feel like a one-off full-page divergent flow; **detalle/edición de visita** (`/workshop/[id]`) must **load and show** job data or a clear error (today users report empty state).

## Scope

### In Scope

- **Nova peça**: replace dedicated `admin/parts/new` full-page form with **list-driven flow** — toolbar CTA opens **Dialog** (or `?create=1` → modal) on `admin/parts`; optional **redirect** from `/admin/parts/new` for bookmarks (mirror suppliers pattern).
- **Taller → Nova visita**: avoid raw `router.push` to a sparse detail page as the only path; align with **modal or inline panel** on `/workshop` (confirm car, optional notes) then open detail — exact UX to be fixed in design/spec.
- **Taller → visita (editar)**: **diagnose and fix** why detail does not populate (params vs auth hydration, API response handling, type mismatch, silent errors); add **loading/error** UX consistent with `AppShell` routes.
- **Sweep**: same pass for obvious **legacy-only CSS** on touched routes — prefer shared `components/ui` + tokens per `ui-component-system`; keep scope to **admin parts + workshop** unless trivial shared fixes.

### Out of Scope

- Full-app visual redesign; i18n copy overhaul; backend schema changes unless a **confirmed bug** requires it; E2E Playwright (not in repo).

## Capabilities

### New Capabilities

- None

### Modified Capabilities

- **`ui-brand-shell`**: extend **route/shell homogeneity** to **admin parts** and **workshop** (toolbar + modal/create parity, theme tokens over ad-hoc module-only forms where touched).
- **`ui-component-system`**: **MUST** use **Dialog/Button/Input** primitives for new create surfaces on these routes (per migration requirement).
- **`parts-inventory`**: **SHALL** document that **create** is initiated from the **list** context (modal/query), not a standalone page as primary UX; **redirect** MAY preserve old URL.
- **`workshop-repair-execution`**: **SHALL** require **readable visit detail** after navigation (data or explicit failure); **SHALL** align **nova visita** UX with list-primary pattern (modal/confirm), not only deep-link detail.

## Approach

1. **Audit** `appointments/page.tsx`, `accounting/*/new` redirects, and `admin-users` modal parity (archived change) as **templates**.
2. **Parts**: extract create form → modal component; wire `admin/parts`; `new/page.tsx` → `router.replace('/admin/parts?create=1')`.
3. **Workshop**: add create modal (car already selected on list) or confirm step; adjust navigation after create.
4. **Detail load**: trace `getServiceJob` + `useParams`; add tests (Vitest RTL or unit for hook) if logic fix; verify 403/404 surfaces in UI.
5. **Verify**: `pnpm lint`, `pnpm typecheck`, `pnpm test`, `go test` unchanged paths unless API fix.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `frontend/src/app/admin/parts/page.tsx` | Modified | Toolbar; modal state; create flow |
| `frontend/src/app/admin/parts/new/page.tsx` | Modified | Redirect legacy URL |
| `frontend/src/app/admin/parts/*.module.css` | Modified | Reduce one-off layout |
| `frontend/src/app/workshop/page.tsx` | Modified | Nova visita → modal/aligned flow |
| `frontend/src/app/workshop/[id]/page.tsx` | Modified | Fix load + UX states |
| `frontend/src/lib/api.ts` | Maybe | Only if client/server contract bug |
| `openspec/specs/{ui-brand-shell,ui-component-system,parts-inventory,workshop-repair-execution}` | Delta in change | Spec scenarios |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Detail bug is backend/auth | Med | Reproduce with network trace; scope backend fix in same change if small |
| Modal + SSR/hydration quirks | Low | Client-only dialogs; match existing accounting patterns |
| Bookmark `/admin/parts/new` | Low | Keep redirect |

## Rollback Plan

Revert the frontend PR (and any spec delta merge). Restore prior `new/page.tsx` behaviour if redirect is insufficient alone.

## Dependencies

- None external.

## Success Criteria

- [ ] Nova peça desde lista usa **modal** (o query equivalente); `/admin/parts/new` redirige sin UX rota.
- [ ] Nova visita alineada con patrón **lista + flujo coherente** (sin parche solo cosmético).
- [ ] `/workshop/[id]` muestra **datos de visita** tras carga correcta o **error accionable**; caso repro documentado en verify.
- [ ] CI verde: lint, typecheck, tests frontend.

# §4 — Gobernanza y convenciones

1. **Ramas:** `main` protegida; features `feat/...`, fixes `fix/...`.
2. **Commits:** Conventional Commits (`feat:`, `fix:`, `chore:`, `docs:`, `test:`, `refactor:`).
3. **PRs:** desde `main`, descripción clara, **CI verde**; API reflejada en Swagger.
4. **API JSON:** claves **camelCase**; en Go, tags `json:"camelCase"` alineados con TypeScript.
5. **Naming:** Go — `PascalCase` exportado, `camelCase` interno; TS — `PascalCase` componentes/tipos, `camelCase` props/funciones.
6. **Testing:** TDD preferente en **servicios** backend; mocks de repos con **testify/mock**; frontend: hooks, cliente API, validadores, componentes críticos (**Vitest**).
7. **Estructura:** `internal/domain|service|repository|handler|middleware` y `frontend/src/{app,components,hooks,lib,stores,types}`.
8. **Validación de dominio:** reglas por mercado — definir en `{{VALIDATION_RULES}}` (p. ej. fiscal/teléfono/documento por país); cuando exista, archivo en `docs/`.
9. **Docs / agentes:** `docs/DOCUMENTATION_INDEX.md`; reglas en **`gonsgarage-rules/`**.
10. **Seguridad:** rate limiting en auth pública; JWT `Authorization: Bearer`; roles explícitos según producto.

**Locales (`{{LOCALE}}`):** `pt_PT`, `es_ES`, `en_GB` para UI/docs principal según tabla del proyecto.

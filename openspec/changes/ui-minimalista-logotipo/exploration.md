## Exploration: UI/UX minimalista alineada al logotipo

### Current State

- **Tokens centralizados** en `frontend/src/styles/tokens.css`: marca documentada (navy/carbón, azul confianza, rojo señalética), escalas de gris, superficies, tipografía Geist, radios y sombras suaves. Modo claro/oscuro vía `html[data-theme]`.
- **Base global** en `frontend/src/app/globals.css`: reset, jerarquía tipográfica, `.container`, enlaces con `--color-primary`.
- **Shell autenticado** (`AppShell.module.css`): cabecera 1200px, nav con pestañas, gradiente en “logoIcon” (navy → primary), **logout** en rojo fijo `#b91c1c` en hover (duplica intención de `--brand-signal` pero hardcoded).
- **Landing** (`landing.module.css`): uso mayoritario de tokens; cabecera fija con blur; logo en caja con borde/sombra.
- **Inconsistencia**: varios `*.module.css` usan **hex sueltos** (p. ej. `cars.module.css`, `appointments.module.css`) en lugar de variables, lo que rompe dark mode y diluye el look unificado.
- **Logo real**: asset `frontend/public/images/LogoGonsGarage.jpg`; los tokens ya nombran intención de marca pero **no hay extracción formal** de muestras del JPG en código (riesgo de deriva cromática vs. logo).

### Affected Areas

- `frontend/src/styles/tokens.css` — fuente de verdad de color/espaciado; posible refinado tras muestreo del logo.
- `frontend/src/app/globals.css` — tipografía, ritmo vertical, enlaces.
- `frontend/src/components/layout/AppShell.module.css` — navegación staff; hardcoded hover logout.
- `frontend/src/app/landing.module.css` — marketing; coherencia con shell app.
- `frontend/src/components/ui/*` — Button, Card, Input, Modal (patrones reutilizables).
- `frontend/src/app/**/*.module.css` — reemplazo gradual de hex por tokens / utilidades.
- `docs/` (opcional) — guía breve de uso de marca si se formaliza.

### Approaches

1. **Token-first + barrido de módulos** — Auditar hex en CSS modules; mapear a tokens existentes o nuevos (`--semantic-*`); mantener una sola paleta.
   - Pros: dark mode consistente, mantenimiento bajo.
   - Cons: trabajo repetitivo en muchos archivos.
   - Effort: **Medium**

2. **Design system doc + solo shell crítico** — Documentar reglas mínimas (espaciado, pesos, uso de rojo) y tocar solo AppShell + landing.
   - Pros: rápido, bajo riesgo.
   - Cons: pantallas internas siguen inconsistentes.
   - Effort: **Low**

3. **Capa utilitaria (p. ej. variantes “subtle / solid”)** — Añadir clases en `utilities.css` para badges/estados y migrar appointments/cars.
   - Pros: menos duplicación semántica.
   - Cons: requiere convención de nombres y disciplina en PRs.
   - Effort: **Medium–High**

### Recommendation

Combinar **(1)** con un subconjunto de **(3)**: fijar **muestras del logo** (3–5 swatches) en `tokens.css` como `--brand-*` verificados, unificar **logout y estados** en tokens, luego **barrido dirigido** de los módulos con más hex (`cars`, `appointments`, `dashboard` si aplica).

### Risks

- Sobrecorregir contraste y perder jerarquía en tablas densas.
- `color-mix` / backdrop-filter: comprobar en navegadores objetivo del MVP.

### Ready for Proposal

**Sí.** El orquestador puede pasar a **sdd-propose** con alcance: look profesional minimalista, paleta anclada al logo, reducción de hardcoded y alineación shell/landing.

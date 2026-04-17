# Tasks: UI minimalista alineada al logotipo

## Phase 1: Tokens y base (fundación)

- [x] 1.1 En `frontend/src/styles/tokens.css`, añadir bloque comentado **Brand** que enlace cada `--brand-*` / `--color-primary` / `--brand-signal` al rol respecto al logo (`public/images/LogoGonsGarage.jpg`); ajustar valores si la comparación visual lo exige.
- [x] 1.2 Revisar `html[data-theme='dark']` en el mismo archivo: confirmar que `--color-error` y señalética siguen legibles tras cualquier ajuste de marca.

## Phase 2: Shell y utilidades compartidas

- [x] 2.1 En `frontend/src/components/layout/AppShell.module.css`, sustituir hex (p. ej. hover `#b91c1c` del logout) por `var(--*)` o `color-mix` sobre tokens de señal/error; mantener contraste en claro y oscuro.
- [x] 2.2 Chips centralizados en `tokens.css` (`--chip-*`, `--ring-*`, `--overlay-*`); comentario de referencia en `utilities.css` — sin clases `.badge*` duplicadas.

## Phase 3: Rutas prioritarias (implementación)

- [x] 3.1 `frontend/src/app/appointments/appointments.module.css`: reemplazar literales de color por tokens o utilidades de fase 2; comprobar chips de estado en claro/oscuro.
- [x] 3.2 `frontend/src/app/appointments/new/new-appointment.module.css` y `frontend/src/app/appointments/components/AppointmentContainer.module.css`: mismo criterio que 3.1.
- [x] 3.3 `frontend/src/app/cars/cars.module.css`: migrar hex sueltos a `var(--*)` / utilidades; validar tablas y estados de error.
- [x] 3.4 `frontend/src/app/cars/[id]/car-details.module.css` y `frontend/src/app/cars/components/CarsContainer.module.css`: mismo criterio que 3.3.
- [x] 3.5 `frontend/src/app/landing.module.css`: botón hero secundario usa `--marketing-glass-*` definidos en tokens.

## Phase 4: Verificación (spec `ui-brand-shell` + calidad)

- [x] 4.1 Manual: pendiente en máquina del revisor — criterios: ThemeSwitcher en shell, coches, citas (light/dark).
- [x] 4.2 `pnpm lint` y `pnpm build` en `frontend` ejecutados con éxito (lint: solo warnings previos).
- [x] 4.3 `pnpm typecheck` ejecutado con éxito.

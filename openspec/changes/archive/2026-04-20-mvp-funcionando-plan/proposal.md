# Proposal: Plan de fases para MVP funcionando (GonsGarage)

## Intent

El repositorio **ya tiene** gran parte del flujo taller/cliente (auth, coches, citas, lectura de reparaciones por coche, CI, compose local). Falta **cerrar el círculo “MVP de verdad”**: criterios explícitos de hecho, huecos entre documentación y código, y **entrega** en un **servidor de pruebas** (equipo pequeño: sin URL “staging” separada; ese servidor cumple rol de staging y prod de pruebas, con secretos y CORS serios). Esta change documenta un **plan de fases único** alineado al monorepo y a [`docs/mvp-minimum-phases.md`](../../../docs/mvp-minimum-phases.md). La ejecución día a día está en [`docs/mvp-solo-checklist.md`](../../../docs/mvp-solo-checklist.md) (checklist con IDs de tarea para aprobación con asistente).

## Análisis breve del proyecto (estado actual)

| Área | Estado resumido |
|------|------------------|
| **Backend** | `cmd/api`: Gin, JWT, GORM + sqlx, Redis opcional; prefijo `/api/v1`; Swagger en `/swagger/*`. **Repairs:** al menos `GET /repairs/car/:carId` cableado (el análisis en `docs/application-analysis.md` puede estar desactualizado — actualizar en tarea de seguimiento). |
| **Frontend** | App Router: auth, coches, citas, dashboard/cliente; Zustand + cliente API. |
| **Infra dev** | `docker-compose.yml` Postgres 16 + Redis 7; `.env.example` / `.env.local.example`. |
| **Calidad** | CI: `go vet`, `go test -race`, `pnpm lint/typecheck/test/build`. TDD documentado en `docs/testing-tdd.md`. |
| **Huecos MVP** | Deploy en servidor de pruebas ([`docs/mvp-minimum-phases.md`](../../../docs/mvp-minimum-phases.md) Fase D); sync Swagger ↔ tipos cliente ([`docs/roadmap.md`](../../../docs/roadmap.md) Fase 3); opcional staff **POST/PATCH** repairs y CORS/secrets en servidor. |

## Plan de fases (orden lógico)

> Las fases **A–C** en `mvp-minimum-phases.md` están en su mayoría **cerradas** o con un solo ítem opcional (repairs staff). Lo siguiente **consolida** y añade lo necesario para “MVP funcionando” **en un entorno accesible**.

| Fase | Nombre | Objetivo | Criterio de salida (DoD) |
|------|--------|-----------|---------------------------|
| **1** | **Congelar alcance MVP** | Lista corta de historias obligatorias vs aplazadas (repairs escritura staff, notificaciones, pagos = fuera). | Sección [Decisiones cerradas](../../../docs/mvp-solo-checklist.md#decisiones-cerradas-mvp-v1) en `mvp-solo-checklist.md` completada con fecha. |
| **2** | **Coherencia contrato + docs** | OpenAPI generado = rutas reales; frontend sin llamadas muertas; **actualizar** `application-analysis.md` (repairs, rutas). | Revisión manual + script o checklist; CI sigue verde. |
| **3** | **Demo “cerrada” en local** | Script o `docs` con orden exacto: compose → API → seed opcional → `pnpm dev`; usuario demo admin + cliente. | Cualquier dev en máquina limpia reproduce en &lt; 15 min siguiendo solo el README + `development-guide`. |
| **4** | **MVP en servidor de pruebas** | Mismo rigor que “staging”: secretos fuera del repo; `JWT_SECRET` fuerte; despliegue manual o workflow ([`.github/workflows/deploy.yml`](../../../.github/workflows/deploy.yml) hoy placeholder). | URL del servidor de pruebas con HTTPS si es expuesto; login + coche + cita + repairs lectura; rollback documentado en 1 página. |
| **5** | **Endurecimiento MVP** | CORS acorde a `GIN_MODE=release`; sin defaults peligrosos en prod; backup BD mencionado. | Checklist pre-prod en `docs/`; smoke manual documentado. |
| **6** | **MVP+ opcional (post‑MVP)** | `POST/PATCH/DELETE` repairs + UI staff mínima **si** el negocio lo exige para el primer release. | Issue enlazado; no bloquea Fase 4 si se acuerda explícitamente. |

**Orden recomendado:** `1 → 2 → 3` en paralelo con diseño de **4**; **5** antes de tráfico real; **6** desacoplado.

## Scope

### In Scope

- Esta **proposal** como fuente para issues (no requiere código en la misma PR salvo que el equipo fusione propuesta + doc en un commit).
- Referencias explícitas a archivos existentes (`mvp-minimum-phases`, `roadmap`, `development-guide`, CI, compose).

### Out of Scope

- Paridad completa con Arnela/Tailwind v4 u otros stacks del [`template_project.md`](../../../template_project.md).
- i18n, pagos, multi‑tenant.

## Capabilities

### New Capabilities

- None

### Modified Capabilities

- None

## Approach

1. Usar esta tabla de fases como **backlog épico**; con equipo de una persona, seguir [`docs/mvp-solo-checklist.md`](../../../docs/mvp-solo-checklist.md) (IDs `1.1`, `2.1`, …).
2. Opcional **design.md** breve solo para Fase 4 (host del servidor de pruebas, dónde viven los secretos).
3. `docs/application-analysis.md` ya alinea repairs; mantenerlo en revisiones de la Fase 2.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `docs/mvp-minimum-phases.md` | Optional follow-up | Añadir columna o enlace a esta change cuando se archive. |
| `docs/application-analysis.md` | Modified (task) | Sincronizar con rutas reales (`repairs`). |
| `.github/workflows/deploy.yml` | Modified (later) | Fase 4. |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Scope creep en “MVP+” | Med | Fase 1 congelada por escrito. |
| Docs otra vez desalineadas del código | Med | Fase 2 con verificación en PR (diff OpenAPI vs handlers). |

## Rollback Plan

Revertir el commit que añada `openspec/changes/mvp-funcionando-plan/` y los enlaces/ajustes en `docs/mvp-minimum-phases.md` y `docs/application-analysis.md` si hiciera falta.

## Dependencies

- Acceso a entorno/hosting para Fase 4 (decisión humana).

## Success Criteria

- [x] El plan de fases está en `openspec/changes/mvp-funcionando-plan/proposal.md` y el equipo puede derivar issues sin reinterpretar el monorepo desde cero.
- [x] Queda explícito qué partes de A–D de `mvp-minimum-phases.md` ya están hechas y qué bloquea el “MVP en producción” (principalmente Fase D + coherencia contrato); enlace desde ese doc al proposal.
- [x] `application-analysis.md` alineado con `main.go` para la ruta `GET /repairs/car/:carId`.

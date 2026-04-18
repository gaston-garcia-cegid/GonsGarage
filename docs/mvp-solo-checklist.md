# Checklist MVP — equipo de una persona

**Contexto**

- **Equipo:** una sola persona (vos); las decisiones de alcance las cerrás vos en la **Fase 1** y podés validar el resto con el asistente tarea por tarea.
- **Entorno:** no hay URL “staging” separada de “producción”. Tu **servidor de pruebas** es hoy el único entorno remoto: tratálo como **staging = prod de pruebas** (mismas reglas de secretos y CORS que usarías en prod real).

**Plan maestro:** [openspec/changes/mvp-funcionando-plan/proposal.md](../openspec/changes/mvp-funcionando-plan/proposal.md) · Fases técnicas históricas: [mvp-minimum-phases.md](./mvp-minimum-phases.md).

---

## Cómo avanzar con el asistente (aprobado)

1. Seguí el orden **1 → 2 → 3** (podés adelantar partes de **4** si ya tenés servidor listo).
2. Para cada ítem, en el chat: **“Fase X.Y — aprobado”** o **“Fase X.Y — no, porque …”** (el asistente puede ejecutar o ajustar el checklist en el repo).
3. Marcá vos las casillas `- [x]` en este archivo al cerrar cada ítem, o pedí al asistente que lo haga tras tu “aprobado”.

**IDs de tarea** (para referencia en chat): `1.1`, `1.2`, … como en la tabla abajo.

---

## Fase 1 — Congelar alcance MVP v1

| ID | Tarea | Criterio de hecho |
|----|--------|-------------------|
| 1.1 | Escribir la lista **«Entra»** en MVP v1 | **Hecho** 2026-04-17 — bloque [Entra](#decisiones-cerradas-mvp-v1) abajo |
| 1.2 | Escribir **3 bullets** “fuera de MVP v1” (ej. pagos, i18n, multi-tenant) | **Hecho** 2026-04-17 — tres ítems bajo «Fuera» + nota de propuesta |
| 1.3 | Cerrar explícitamente **repairs staff** (API/UI escritura) vs solo lectura | **Hecho** 2026-04-17 — línea «Repairs staff» abajo (**incluido**) |

### Decisiones cerradas MVP v1

**Fecha cierre Fase 1 (alcance):** 2026-04-17 — *1.1, 1.2 y 1.3 cerradas en checklist.*  
**Reconfirmación 1.1 (chat):** lista Entra con desglose de **invoices** (clientes, proveedores, recibos de sueldos) y **billing** solo clientes.  
**Propuesta 1.2 (chat) — «Fuera»:** tres límites claros compatibles con **CRUD billing** (registro contable / facturación hacia clientes) sin convertir el MVP en un PSP.

- Entra (MVP v1):
  - **CRUD users** — registro / login JWT y perfil (`/auth/me`).
  - **CRUD cars**
  - **CRUD appointments** (citas)
  - **CRUD repairs** (reparaciones)
  - **CRUD invoices** (facturas clientes, facturas proveedores, recibos de sueldos de empleados)
  - **CRUD billing** (cobros / facturación exclusiva de clientes)
- Fuera (MVP v1):
  - **Cobro electrónico con terceros** — pasarelas (Stripe, Adyen, …), TPV con adquirente, **MB Way / transferencia como canal conectado** (API, webhooks, firma, devoluciones), tokenización de tarjeta y **conciliación bancaria automática**. Sí podés modelar *medio* e *importe* como **datos** en facturas/billing (efectivo, tarjeta, MB Way, transferencia como etiqueta o importe registrado); **no** integración ni cobro online end‑to‑end.
  - **i18n** — una sola lengua de interfaz y de documentación operativa en MVP v1.
  - **Multi‑tenant / varios talleres** — una instancia = un taller (una marca, un conjunto de datos); sin subdominios por cliente corporativo ni aislamiento duro entre empresas.
- Repairs staff (POST/PATCH/DELETE + UI): **incluido** en MVP v1 — coherente con «CRUD repairs» en Entra (1.1 aprobado).

---

## Fase 2 — Coherencia contrato + docs

| ID | Tarea | Criterio de hecho |
|----|--------|-------------------|
| 2.1 | Revisar que **Swagger generado** (`backend/docs/`) refleje rutas reales en `cmd/api` (grupos `/api/v1/...`) | **Hecho** 2026-04-18 — faltaba `GET /api/v1/repairs/car/{carId}` en Swagger; añadidas anotaciones swag en `repair_handler_gin.go` y regenerado `swagger.json` / `swagger.yaml` / `docs.go`. Resto de paths alineados con `setupRoutes`. |
| 2.2 | **Frontend:** buscar llamadas a endpoints que no existan o estén mal prefijados | **Hecho** 2026-04-18 — citas: `PATCH …/cancel|confirm|complete` → `PUT /appointments/{id}` con `{ status }`. Coches: `GET /cars/owner/:id` → `GET /cars?ownerId&limit&offset`. Repairs: `getCarWithRepairs` usa `/repairs/car/{id}`; `api.ts` `getRepairs` exige `carId`; `getProfile`→`/auth/me`; logout sin `POST /auth/logout`. Constantes `API_ENDPOINTS` alineadas. Tests Jest `auth.service` actualizados. |
| 2.3 | Mantener [application-analysis.md](./application-analysis.md) alineado cuando agregues rutas | **Hecho** 2026-04-18 — citas `PUT`+`status`, coches `ownerId` query, nota cliente HTTP. |

---

## Fase 3 — Demo local reproducible (una máquina limpia)

| ID | Tarea | Criterio de hecho |
|----|--------|-------------------|
| 3.1 | Seguir solo [README.md](../README.md) + [development-guide.md](./development-guide.md): compose → API → `pnpm dev` | Anotar cualquier paso faltante y corregir doc |
| 3.2 | Flujo mínimo manual: **login** → **coche** → **cita** → **ver repairs** en detalle coche (si aplica a tu MVP v1) | Checklist mental OK; opcional: capturas en `docs/` |
| 3.3 | Seed cliente demo (`go run ./cmd/seed-test-client`) documentado si lo usás en demos | Comando probado desde `backend/` |

---

## Fase 4 — Servidor de pruebas (= tu staging / prod de pruebas)

| ID | Tarea | Criterio de hecho |
|----|--------|-------------------|
| 4.1 | Definir **URL base** del front y del API (aunque sea IP + puerto) y anotarlas abajo | Texto en [Entorno remoto](#entorno-remoto-servidor-de-pruebas) |
| 4.2 | **Secretos:** `JWT_SECRET` fuerte; `DATABASE_URL` / Redis solo en el servidor (no en git) | `.env` en servidor o secret manager; nada sensible en repo |
| 4.3 | **Backend:** build (`go build -o … ./cmd/api`) o imagen Docker; proceso bajo systemd/supervisor o compose en el servidor | API responde `/health` y `/ready` |
| 4.4 | **Frontend:** `pnpm build` con `NEXT_PUBLIC_API_URL` apuntando al API del servidor; servir con lo que elijas (nginx, Node, etc.) | Login usable contra ese API |
| 4.5 | **HTTPS** si el servidor es expuesto (certificado Let’s Encrypt o TLS detrás de proxy) | Navegador sin warning crítico (o documentar excepción solo LAN) |
| 4.6 | **Rollback:** una página o sección “cómo volver atrás” (versión anterior binario + migración si aplica) | Enlace o párrafo en `docs/` o nota privada |

### Entorno remoto (servidor de pruebas)

- URL API: ____________________
- URL frontend: ____________________
- Notas (proveedor, SSH, etc.): ____________________

---

## Fase 5 — Endurecimiento antes de dar por cerrado el MVP remoto

| ID | Tarea | Criterio de hecho |
|----|--------|-------------------|
| 5.1 | **CORS:** origen permitido = URL real del front de pruebas (no `*` en `release` si podés evitarlo) | Revisión `main.go` / middleware |
| 5.2 | **`GIN_MODE=release`** en el servidor de pruebas si simulás prod | Variable documentada |
| 5.3 | **Backup BD:** comando o política mínima (pg_dump semanal, etc.) | Párrafo en `docs/` o runbook propio |

---

## Fase 6 — MVP+ (opcional; no bloquea Fase 4–5)

| ID | Tarea | Criterio de hecho |
|----|--------|-------------------|
| 6.1 | API + UI staff para **crear/editar** repairs | Solo si lo marcaste “incluido” en Fase 1 |
| 6.2 | Automatizar **deploy** (GitHub Actions → tu servidor) | Solo si querés; [`.github/workflows/deploy.yml`](../.github/workflows/deploy.yml) hoy es placeholder |

---

## Estado rápido (rellená al ir cerrando)

| Fase | Estado |
|------|--------|
| 1 Congelar alcance | **hecha** (1.1–1.3 cerradas 2026-04-17) |
| 2 Contrato + docs | **hecha** (2.1–2.3, 2026-04-18) |
| 3 Demo local | … |
| 4 Servidor pruebas | … |
| 5 Endurecimiento | … |
| 6 MVP+ | N/A / pendiente |

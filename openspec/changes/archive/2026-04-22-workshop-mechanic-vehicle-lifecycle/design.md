# Design: workshop-mechanic-vehicle-lifecycle (modelo)

Técnica: nuevo agregado **ServiceJob** (visita), persistencia 1:1 de **recepción** y **cierre** en fase MVP1, enlaces suaves a `repairs` y a **appointments** opcional; OBD/estimate = API *stub* acorde al spec.

## Architecture decisions

| # | Options | **Choice** | Rationale |
|---|-----------|--------------|-----------|
| 1 | Extender `repairs` con JSON / nuevo agregado | **Nuevo `ServiceJob`** (tabla `service_jobs`) | Trazabilidad de visita sin colapsar reparación lineal; compatible con reparations legadas. |
| 2 | 1 fila con JSON / tablas hijas p/ checklists | **2 tablas hijas 1:1** (recepción, cierre) o columnas estructuradas con `schema_version` | Mejor que JSON opaco: consultas, migraciones, tests. Se priorizan filas; si pesa, JSON versionado mín. |
| 3 | Visita = N `repairs` | **Visita 1, `repairs` 0..N** con `repairs.service_job_id` NULL legado, NOT NULL cando aplique (post–MVP1 ou rule-based) | Preserva API actual; `Repair` segue cando; acople gradual. |
| 4 | Prefijo API | **`/api/v1/service-jobs`**, sub-rutas `.../reception`, `.../handover` (nomes afinables) | Coherente con REST actual; fácil anexar OBD/estimate. |
| 5 | Gating auth | Misma capa `protected` + reglas en servicio (estilo `RepairService`); opción `RequireWorkshopStaff` si todas las mutaciones son *staff*; el spec exige **cliente MUST NOT mutar** con cobertura en tests | Mantiene flexibilidad para lectura cliente P1+ sin duplicar middleware. |

## Modelo de datos (MVP1)

```text
service_jobs
  id, car_id, status (open|in_progress|closed|cancelled),
  opened_by, opened_at, closed_at?,
  optional: appointment_id NULL, schema_version
```

```text
service_job_reception   -- 0..1 por job en MVP1 (1:1 tras primeiro save)
  service_job_id PK/FK, km, -- campos puntuales ou JSON versionado
  created_at, created_by, payload_version

service_job_handover     -- 0..1, semántica: entrega
  (mesma idea)
```

**Estado mínimo:** *open* á creación; transición a *closed* só cando pase agregada `handover` (ou regra de servicio explícita).

OBD/estimate: tablas baleiras ou **non creadas**; endpoints devolven 501/404 *documentado* (spec *stub*). Sen hardware.

## Flujo (MVP1)

```text
Employee → POST /service-jobs { car_id } → open
         → PUT  /service-jobs/:id/reception
         → PUT  /service-jobs/:id/handover   → status closed
(optional) POST /repairs con service_job_id (fase D)
```

## File changes (plan)

| Ruta / área | Acción |
|-------------|--------|
| `backend/internal/domain/service_job.go` (ou subpaquete `workshop`) | Crear agregado + DTOs |
| `.../repository/postgres/..._repository.go` | GORM + migracións |
| `.../service/.../service.go` | Regras, permisos por coche (misma idea que reparos) |
| `.../handler/..._gin.go` + `main.go` | Ruta e inyección |
| `backend/.../authorization_*` o tests de rol | Cumplir delta *mvp-role-access* |
| `frontend/src/app/workshop/...` (o bajo *employees*) | Listado, detalle, formularios |
| `openspec/.../specs/**` | Creado en este change |

## Testing

- Go: *httptest* (client 403/401 vs staff 2xx) para creación visita e *PUT* checklist.
- Criterio: `go test` como en `openspec/config.yaml`.

## Migration / Rollout

- *Add only*: tablas novas, columna `repairs.service_job_id` nullable cando se implemente enlace.
- Rollback: *feature flag* ou non rexistrar ruta; datos conservados.

## Open questions

- [ ] Nomenclatura pública: `service-jobs` vs. `taller/visitas` (UI pt_PT).
- [ ] Lectura *client* read-only de visita: cal endpoint e `car` *ownership* (P1+).

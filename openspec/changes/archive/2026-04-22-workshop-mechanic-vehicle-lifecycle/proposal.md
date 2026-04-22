# Proposal: Taller — ciclo de servicio mecánico (empleado)

## Intent

Dar a **`employee` (mecánico)** un flujo estructurado alrededor de un **vehículo en servicio**: desde la **recepción documentada** (estado, km, fluidos, ruedas) hasta la **entrega con verificación**, pasando por **OBD / diagnóstico**, **presupuesto aprobable**, **ejecución bajo plazo pactado** y trazabilidad. Hoy el dominio `repairs` y la matriz de roles (`mvp-role-access`) cubren poco de este ciclo; Falta alinear producto, API y UI.

## Scope (épica — desglosar en fases)

### In Scope (visión v1+)

- **Fase A — Recepción (checklist):** km, estado general, niveles (aceite, refrigerante, etc.), neumáticos, observaciones visuales; opcional **fotos** y registro de **quién** entrega el vehículo.
- **Fase B — Diagnóstico OBD:** enlace/lectura (según dispositivo), códigos DTC, snapshot almacenado, correlación con “posibles tareas” (no sustituir criterio humano).
- **Fase C — Presupuesto y aprobación:** tareas **solicitadas** + **descubiertas**; importes/horas; **aprobación explícita del cliente** o “trabajo interno / garantía” según regla; pase a cola/orden de **gestión** (manager/orden de trabajo).
- **Fase D — Ejecución:** repair/work order, partes, tiempos; respeto del **plazo/ventana pactada con el manager** (o desviación registrada con causa).
- **Fase E — Entrega (checklist cierre):** re-verificar km, presión neumáticos, fluidos, firma o constancia; opcional comprobante al cliente.

### Out of Scope (inicial)

- Conector OBD real end-to-end (hardware) sin spike; contabilidad completa de repuestos (P1); multi-taller; app móvil off-line; integración con proveedores de piezas en tiempo real.

## Capabilities

| Tipo | Nome | Notas |
|------|------|--------|
| **New (delta)** | `workshop-repair-execution` | Requisitos y escenarios: checklist recepción/cierre, OBD, presupuesto, aprobación, work order, SLA/desviación, entrega. Puede partir como **delta** bajo un único `specs/workshop-repair-execution/` ou partir **modular** (p.ej. `vehicle-intake`, `obd-snapshot`, `work-order-approval`) en fases. |
| **Modified** | `mvp-role-access` | Extender qué hace el rol `employee` (y quizá `manager`) sobre reparaciones / work orders, lecturas vs mutaciones, sin romper *client* MUST NOT mutar. |

*Si a primeira entrega for só *documento de requisitos* e spike técnico, pódese marcar *Modified* como “ningún requisito normativo aínda” e facer só boceto en spec.*

## Mapeo ás túas 5 fases (resumo)

| Súa frase | Épica | Recordatorio UC |
|----------|-------|------------------|
| Checklist (km, estado, niveis, neumáticos) | A | + fotos, dano preexistente, asinatura, vehículo bloqueado/bloqueo baúl. |
| OBD e análise de posibles problemas | B | + limpar códigos tras reparar? historial de lecturas, export. |
| Presupuesto tareas + novas tareas → xestión | C | aprobación cliente, rexeitar presuposto, pago a conta, orde interna. |
| Reparación, pezas, tempo pactado | D | recepción pezas, substitúe plan, bloqueo por agarda peza. |
| Entregar con reverificación | E | proba de ruta? checklist impreso/PDF, queixas pós-entrega 24/48h. |

## Outros UC que convén (ademais dos que xa listaches)

- **Trazabilidade e auditoría:** *quién* cambia estado, hora, orixe do dato (manual vs OBD).
- **Foto / dano / inventario** antes de tocar o vehículo (disputa de responsabilidade).
- **Orde de traballo (OT)** e **asignación a baía/mecánico**; reasignar se baja.
- **Atraso e causas** (taller, cliente, prov.) notificable ao manager.
- **Garantía e retraballos** ligados á reparación orixinal.
- **Pezas:** pedido, recepción, devolución, números de serie; stock mínimo (MVP+).
- **Cliente:** ver estado aprobado (“en diagnóstico / agarda aprobación / en reparación / listo”); aprobar presuposto (canle segura); non bloquea UC mecánico, pero alimenta a Fase C.
- **Cumplimento normativo/seguridade:** dixitalizar o que toque (MOT/ITV, residuos) só como requisito opcional por xurisdicción.
- **Métricas do taller:** TTM (tiempo en taller), tasa de rotraballos, cumprimento de citas.
- **Integración futura** con **appointments** (a cita acaba en “vehículo recibido en taller” e abre a recepción mecánica).

## Aproximación técnica (alto nivel)

- Modelar un **agregado** (p.ej. `ServiceJob` / `WorkOrder`) ligado a `car` + (opcional) `appointment`/`repair` existente, con estados e historial.
- API: evolución por etapas ou sub-recursos (checklist, OBD, estimate, work_log, handover).
- UI: área de taller para `employee` (navegación separada do cliente) ou extensión de *employees* — decidir en `design`.

## Affected areas (direccional)

| Area | Impact |
|------|--------|
| `backend/internal/.../repair` ou novo paquete `workorder` / `service_job` | Alto |
| `backend/internal/handler` | Nuevo rutas ou extensión |
| `openspec/specs/*` | Novo ou delta + `mvp-role-access` |
| `frontend` (ruta staff tipo `/workshop/...` ou semellante) | Alto |

## Risks

| Risco | Mitigación |
|--------|------------|
| Alcance épico inmanexable de golpe | Mínimo viable: só **recepción + 1 reparación simple** ou spike OBD mock. |
| OBD real varía | Comezar con *upload de informe* ou mock API. |
| Rol *employee* con demasiados permisos | Afinar con `manager` aprobación e tests en `mvp-role-access`. |

## Rollback

Feature flag ou rama: deshabilitar rutas novas; datos novos con migración **add only**; rollback desplegable sen borrar agregado se xa hai datos reais (marcar *deprecated*).  

## Dependencias

- Criterio de produto: **obriga** de firma de cliente en recepción? canal de aprobación de presuposto?
- `repairs` existente: compatibilidade ou migración desde histórico simple.

## Success criteria

- [ ] Listaxe priorizada de UC (Must / Should) acordada co anfitrión do taller.
- [ ] Un documento de **delta spec** (ou specs modulares) con escenarios *Given/When/Then* para a primeira entrega.
- [ ] Prototipo API + 1 test de rol ou spike OBD aceptado.
- [ ] Criterio de “listo” para corte v1: recepción **persistida** + cierre básico **reproducible** en entorno de probas.

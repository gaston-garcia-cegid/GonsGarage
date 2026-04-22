# Delta for workshop-repair-execution

> Propuesta: `openspec/changes/taller-benchmark-mvp-priorities/proposal.md`. `mvp-role-access` **no** se modifica (sin rol ni ruta de matriz adicional).

Solo se **añaden** requisitos; el catálogo principal de `openspec/specs/workshop-repair-execution/spec.md` se mantiene para fusión vía sdd-archive.

## ADDED Requirements

### Requirement: Superficie de recepción prioritaria (flujo “app”)

El personal de taller con permiso para mutar checklists de la visita **SHALL** poder completar la **recepción** con el **mismo** contrato de validación y persistencia del requisito *Recepción y cierre mínimos* desde una **superficie** orientada a captura rápida: ruta o **vista dedicada** a recepción, o **diseño responsive** equivalente. La API **MUST NOT** exigir un canal, cabecera o formato de cuerpo distinto al ya definido para la recepción; solo cambia la experiencia de producto (navegación y prioridad de tareas en UI).

#### Scenario: Recepción válida desde la superficie prioritaria

- GIVEN un miembro de staff con permiso de mutación de visita según `mvp-role-access` y visita en estado abierto
- WHEN usa la superficie prioritaria y envía un checklist de recepción válido
- THEN la recepción queda persistida

#### Scenario: Payload inválido o incompleto

- GIVEN la misma visita
- WHEN el cuerpo es inválido o incompleto
- THEN la respuesta **MUST NOT** ser 2xx; la semántica de error **SHALL** ser coherente con *Recepción incompleta*

### Requirement: Listado de visitas del “día” operativo

El sistema **SHALL** ofrecer a personal de staff (matriz de `mvp-role-access`) un medio de **listar o filtrar** visitas (service jobs) según criterio de “día” **documentado** (p. ej. fecha/hora de apertura en UTC, o regla de calendario local fijada en el producto). La visita **SHALL** seguir siendo la unidad de trazabilidad; el listado **MUST NOT** reemplazar la facturación ni usarse como comprobante de cobro.

#### Scenario: Día sin visitas

- GIVEN criterio “día” tal que ninguna visita califica
- WHEN se pide el listado o el recurso de filtrado
- THEN el resultado **SHALL** ser vacío o equivalente explícito; **MUST NOT** 5xx

#### Scenario: Día con al menos una visita

- GIVEN al menos una visita abierta que califica bajo el criterio
- WHEN se pide el listado
- THEN el resultado **SHALL** incluir al menos un identificador de visita accesible al staff autorizado

### Requirement: Lectura: enlace visita–reparaciones

Donde el modelo de datos permita vincular `repairs` a la visita, la lectura de la visita (mismo `GET` principal, o sub-recurso **coherente y documentado** en un solo criterio de producto) **SHOULD** permitir ver si la visita tiene **cero, una o varias** reparaciones vinculadas, **MUST NOT** implicar trabajo inexistente.

#### Scenario: Visita sin reparación

- GIVEN la visita sin reparación asociada
- WHEN se lee el recurso de la visita con el enfoque de enlace
- THEN la carga **MUST NOT** sugerir reparaciones inexistentes

#### Scenario: Visita con reparación vinculada

- GIVEN una reparación persistida vinculada a esa visita
- WHEN se lee de la misma forma
- THEN la carga **SHALL** permitir identificar la vinculación (identificador o lista, según el contrato elegido)

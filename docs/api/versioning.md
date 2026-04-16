# Versionado de la API HTTP (GonsGarage)

## Prefijo de ruta

- La versión **mayor** de la API va en el path: **`/api/v1/...`**.
- Los clientes deben usar siempre el prefijo explícito; no hay redirecciones automáticas entre versiones.

## Versión de contrato vs OpenAPI

- **`apiVersion`** (p. ej. `1.0.0`): versión **semántica del contrato** que documentamos en [CHANGELOG.md](../../CHANGELOG.md) y exponemos en **`GET /health`** (`apiVersion`).
- El artefacto **Swagger** (`backend/docs/`) puede seguir anotaciones `@version` alineadas con ese número cuando haya cambios visibles en la doc.

## Cambios compatibles (mismo `/api/v1`)

- Nuevos endpoints, campos JSON opcionales nuevos, validaciones más estrictas que no rompan clientes conformes.
- Registrar en **CHANGELOG** bajo `[Unreleased]` → sección **Added** o **Changed**.

## Cambios incompatibles

- Renombrar o eliminar campos obligatorios, cambiar semántica de códigos HTTP o de autenticación en rutas existentes.
- Plan: documentar en CHANGELOG **Removed** / **Changed**, incrementar **`apiVersion` major** (p. ej. `2.0.0`), y en el futuro exponer **`/api/v2`** en paralelo durante un periodo de deprecación acordado con el equipo.

## Deprecación

- Anotar en CHANGELOG con fecha o versión objetivo de retirada.
- Donde sea posible, mantener el comportamiento antiguo con flag o respuesta dual durante una ventana corta.

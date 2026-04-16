# Documentación GonsGarage

Índice de la documentación técnica del repositorio. El README principal del proyecto sigue siendo la puerta de entrada general; aquí se concentra el material para **desarrollo, arquitectura y planificación**.

| Documento | Contenido |
|-----------|-----------|
| [application-analysis.md](./application-analysis.md) | Análisis de la aplicación: stack, módulos, rutas API, frontend, infraestructura. |
| [development-guide.md](./development-guide.md) | Cómo levantar backend y frontend, variables de entorno, Docker y comprobaciones rápidas. |
| [arnela-specs.md](./arnela-specs.md) | Arnela en **`D:\Repos\Arnela`**, matriz vs GonsGarage y enlaces al resumen en `specs/arnela/`. |
| [specs/arnela/ARNELA_SYNOPSIS.md](./specs/arnela/ARNELA_SYNOPSIS.md) | Resumen extraído del repo Arnela (stack, estructura, convenciones, diferencias). |
| [roadmap.md](./roadmap.md) | Roadmap para documentación, configuración y alineación con el enfoque Arnela (cuando existan specs). |
| [testing-tdd.md](./testing-tdd.md) | TDD obligatorio, backend/frontend, CI y exclusiones Jest documentadas. |
| [mvp-minimum-phases.md](./mvp-minimum-phases.md) | Fases mínimas sugeridas para un primer MVP. |
| [api/versioning.md](./api/versioning.md) | Versionado del contrato HTTP (`/api/v1`, `apiVersion`, CHANGELOG). |
| [deployment-staging.md](./deployment-staging.md) | Docker compose prod/smoke, secretos, CORS, probes `/health` y `/ready`, backup. |
| [../CHANGELOG.md](../CHANGELOG.md) | Cambios notables del contrato y del producto (Keep a Changelog). |

## Documentación ya existente en el repo

- [../README.md](../README.md): visión del producto, stack y guía amplia.
- [../Agent.md](../Agent.md): estándares de código y arquitectura para el equipo.
- [../frontend/README.md](../frontend/README.md): enlaces a `docs/` y `.env.local.example`.
- [../docker-compose.yml](../docker-compose.yml): Postgres + Redis para desarrollo local.
- [../frontend/docs/api-client.md](../frontend/docs/api-client.md): cliente API del frontend.
- Archivos históricos de migración en la raíz y `frontend/`: ver enlaces en [application-analysis.md](./application-analysis.md).

## Convención sugerida

- Especificaciones nuevas o importadas (p. ej. desde Arnela): colocar bajo `docs/specs/` o documentar la ruta en [arnela-specs.md](./arnela-specs.md).
- Cambios de producto grandes: actualizar `roadmap.md` y, si aplica, el README raíz.

# Especificaciones del proyecto Arnela

## Fuente local (canónica)

El proyecto **Arnela** está disponible en esta máquina en:

**`D:\Repos\Arnela`**

No forma parte del workspace de GonsGarage; la información se consulta **leyendo ese repo directamente**.

## Resumen en GonsGarage

Para no duplicar mantenimiento, el equipo puede usar:

- **[docs/specs/arnela/README.md](./specs/arnela/README.md)** — Índice y punteros.  
- **[docs/specs/arnela/ARNELA_SYNOPSIS.md](./specs/arnela/ARNELA_SYNOPSIS.md)** — Resumen de producto, stack, estructura, convenciones y diferencias frente a GonsGarage (última extracción: 2026-04-16).

## Matriz de comparación (borrador)

| Tema | Arnela | GonsGarage hoy | Acción propuesta |
|------|--------|----------------|------------------|
| **Compose / DB** | `docker-compose.yml` en raíz: PostgreSQL + Redis; `.env.example` | **Hecho en GonsGarage:** `docker-compose.yml` raíz + `backend/.env.example` alineados con `cmd/api` | Perfiles opcionales para app en contenedor; endurecer secretos en prod |
| **Migraciones** | `golang-migrate`, SQL en `migrations/` | GORM `AutoMigrate` en `cmd/api` | Decidir: mantener GORM o acercarse al modelo Arnela si se exige trazabilidad SQL |
| **CI** | GitHub Actions (README) | Sin workflows `.github` | Añadir pipeline similar (Go test + frontend lint/test) |
| **Estructura docs** | `docs/DOCUMENTATION_INDEX.md` + muchos temas + `arnela-rules/` | `docs/` reciente, sin índice tan grande | Opcional: `docs/DOCUMENTATION_INDEX.md` propio enlazando análisis, dev, roadmap |
| **Auth / API** | JWT + roles; `GET /auth/me`; rate limiting; CRUD users/clients documentado | JWT; register/login; sin `auth/me` en el listado revisado; cars/employees/appointments | Evaluar `auth/me`, rate limit, políticas por rol según necesidad del taller |
| **Frontend** | Next 16, pnpm, Tailwind v4, Shadcn, rutas por `(auth)`/`(client)`/`(backoffice)` | Next 15, npm, estructura `app/` plana en parte | Roadmap: subir versión/herramientas solo si hay beneficio claro |
| **Observabilidad** | `/health` y `/readiness` (README Arnela) | `/health` | Añadir readiness si se despliega en orquestador |

## Si la ruta cambia de máquina

Actualizar la ruta absoluta en este archivo y en `docs/specs/arnela/ARNELA_SYNOPSIS.md`. Alternativa estable: clonar Arnela como submódulo bajo `docs/external/arnela` y enlazar desde aquí.

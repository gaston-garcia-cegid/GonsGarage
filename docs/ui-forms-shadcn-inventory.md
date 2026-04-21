# Inventario — forms y CTAs vs `components/ui` (Shadcn)

**Change:** `mvp-ui-visual-parity`  
**Fecha:** 2026-04-21  
**Propietario:** equipo frontend (actualizar al migrar).

## Hecho (primitives Shadcn en flujo)

| Área | Controles |
|------|-----------|
| Auth | `LoginForm`, `RegisterPage` — `Button`, `Input`, `Label` |
| `AppShell` | Navegación y logout — `Button` |
| Dashboard | CTAs de tarjetas — `Button` |
| Cars | Cabecera “adicionar” — `Button` en `CarsContainer` |
| Appointments | Toolbar, errores, empty state — `Button` |
| Accounting (detalle + new) | `Guardar` / `Criar` / `Eliminar` — `Button` (`variant="destructive"` en eliminar) |
| `employees/page.tsx` | Pesquisa `Input`/`Label`, CTAs e paginação `Button`, modal colaborador `Input`/`Label`/`Button` |
| Área cliente | `ClientDashboard` (CTAs) + `DashboardLayout` (logout e tabs) — `Button` |

## Pendiente (inputs nativos o `<button>` sin Shadcn)

| Ruta / componente | Notas |
|-------------------|--------|
| `my-invoices/**`, resto de `cars/**` detalle | Prioridad media |
| `ConfirmModal`, `CarModal`, otros modales | Sustituir por `Dialog` de Shadcn cuando se toque cada flujo |

**Meta:** cero filas en “Pendiente” o lista fechada por issue en el próximo verify.

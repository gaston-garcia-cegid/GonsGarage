# Tema Shadcn ↔ marca GonsGarage

Variables HSL (sin función `hsl()`) definidas en `frontend/src/styles/shadcn-theme.css`. Tailwind las expone como `bg-primary`, `text-foreground`, etc. vía `frontend/tailwind.config.ts`.

| Token marca / app (`tokens.css`) | Variable Shadcn | Uso típico |
|----------------------------------|-----------------|------------|
| `--color-primary` / `--brand-accent` (claro) | `--primary` | Botón primario, anillos focus |
| `--text-on-primary` | `--primary-foreground` | Texto sobre primario |
| `--surface-page` (sensación visual) | `--background` | Fondo base en utilidades Shadcn |
| `--text-primary` | `--foreground` | Texto principal |
| `--color-error` / `--brand-signal` | `--destructive` | Eliminar / acciones destructivas |
| `--border-default` | `--border`, `--input` | Bordes de campos y tarjetas |
| `--surface-muted` | `--muted`, `--accent` | Superficies secundarias / hover |
| `--text-muted` | `--muted-foreground` | Texto secundario en forms |

**Dark:** `html[data-theme='dark']` sobrescribe las mismas variables Shadcn para alinearlas con la rampa oscura de `tokens.css` (paneles `#0e1522`, texto claro, primario `#60a5fa`).

**Mantenimiento:** al cambiar una marca en `tokens.css`, revisar si hace falta ajustar el triple HSL correspondiente en `shadcn-theme.css` y esta tabla.

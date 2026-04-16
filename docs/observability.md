# Observabilidad (API)

## Logs estructurados (`log/slog`)

- **Salida JSON:** automática con `GIN_MODE=release`, o forzada con `LOG_FORMAT=json`.
- **Texto (dev):** por defecto en local si no aplicas las condiciones anteriores.
- **Nivel:** `LOG_LEVEL` = `debug` | `info` | `warn` | `error` (por defecto `info`).
- El paquete estándar **`log`** del arranque (`cmd/api`) se reenvía al mismo handler slog (una línea = un evento `INFO`).
- Cada petición HTTP se registra con el middleware `SlogRequestLogger`: `method`, `path`, `status`, `latency_ms`. Las rutas `/health`, `/ready` y `/metrics` se registran en nivel **debug** para no inundar los logs.

## Métricas Prometheus

- **`GET /metrics`**: handler estándar de `prometheus/client_golang` (métricas de runtime Go y del cliente HTTP, etc.).
- **Seguridad:** no debe quedar expuesto a Internet abierto; usar red interna, `NetworkPolicy`, o que el scrape vaya solo desde el monitor (VictoriaMetrics, Grafana Agent, etc.).

## OpenTelemetry

- No integrado aún; si se adopta, enlazar trazas con el `request_id` (cabecera opcional `X-Request-Id` ya contemplada en el middleware de logs).

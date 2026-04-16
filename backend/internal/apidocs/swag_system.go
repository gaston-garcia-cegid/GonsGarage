// Package apidocs contiene anclas solo para swag (sin lógica en runtime).
// Los handlers reales están en cmd/api (main) y en handlers.
//
// Incluir este directorio en: swag init -d ...,./internal/apidocs,...
package apidocs

// healthCheckDoc ancla /health para swag (el handler real está en setupRoutes de main).
//
// @Summary Comprobación de salud
// @Description Estado del proceso API (sin autenticación).
// @Tags system
// @Produce json
// @Success 200 {object} map[string]string "status, message"
// @Router /health [get]
func healthCheckDoc() {}

// readyCheckDoc ancla /ready para swag.
//
// @Summary Comprobación de readiness
// @Description PostgreSQL accesible (sin autenticación; usar en probes de despliegue).
// @Tags system
// @Produce json
// @Success 200 {object} map[string]interface{} "status, checks"
// @Failure 503 {object} map[string]interface{} "not_ready"
// @Router /ready [get]
func readyCheckDoc() {}

// metricsDoc ancla /metrics para swag.
//
// @Summary Métricas Prometheus
// @Description Métricas de proceso y runtime (sin autenticación; restringir en producción).
// @Tags system
// @Produce plain
// @Success 200 {string} string "texto formato Prometheus"
// @Router /metrics [get]
func metricsDoc() {}

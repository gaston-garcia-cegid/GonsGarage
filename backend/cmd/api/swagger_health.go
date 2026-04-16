package main

// healthCheckDoc ancla /health para swag (el handler real está en setupRoutes).
//
// @Summary Comprobación de salud
// @Description Estado del proceso API (sin autenticación).
// @Tags system
// @Produce json
// @Success 200 {object} map[string]string "status, message"
// @Router /health [get]
func healthCheckDoc() {}

// readyCheckDoc ancla /ready para swag (el handler real está en setupRoutes).
//
// @Summary Comprobación de readiness
// @Description PostgreSQL accesible (sin autenticación; usar en probes de despliegue).
// @Tags system
// @Produce json
// @Success 200 {object} map[string]interface{} "status, checks"
// @Failure 503 {object} map[string]interface{} "not_ready"
// @Router /ready [get]
func readyCheckDoc() {}

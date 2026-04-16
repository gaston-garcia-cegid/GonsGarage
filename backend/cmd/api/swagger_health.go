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

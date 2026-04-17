package main

// readyCheckDoc ancla /ready para swag (el handler real está en setupRoutes).
//
// @Summary Comprobación de preparación
// @Description PostgreSQL alcanzable vía pool compartido (sqlx sobre *sql.DB de GORM).
// @Tags system
// @Produce json
// @Success 200 {object} map[string]string "status ready"
// @Failure 503 {object} map[string]string "not_ready"
// @Router /ready [get]
func readyCheckDoc() {}

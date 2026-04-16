package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserHandler es un stub de gestión de usuarios (sin rutas registradas en cmd/api).
// Firma alineada con Gin para no arrastrar gorilla/mux.
type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// ListUsers lista usuarios (placeholder).
func (h *UserHandler) ListUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "List users"})
}

// GetUser obtiene un usuario por id (placeholder).
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Get user", "id": id})
}

// UpdateUser actualiza un usuario (placeholder).
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Update user", "id": id})
}

// DeleteUser elimina un usuario (placeholder).
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Delete user", "id": id})
}

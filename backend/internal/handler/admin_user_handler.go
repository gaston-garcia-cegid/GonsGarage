package handler

import (
	"errors"
	"net/http"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/ports"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdminUserHandler struct {
	authService ports.AuthService
}

func NewAdminUserHandler(authService ports.AuthService) *AdminUserHandler {
	return &AdminUserHandler{authService: authService}
}

// ProvisionUser creates a user (roles manager, employee, or client only). Requires JWT; only admin and manager reach this handler.
// @Summary     Aprovisionar utilizador (staff)
// @Description Cria utilizador com papel manager, employee ou client. Admin pode todos; manager não pode criar manager. Nunca cria admin por este fluxo.
// @Tags        admin
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       body body ports.ProvisionUserRequest true "Dados do novo utilizador"
// @Success     201 {object} SwaggerProvisionUserOK
// @Failure     400 {object} SwaggerMessage
// @Failure     403 {object} SwaggerMessage
// @Failure     409 {object} SwaggerMessage
// @Failure     500 {object} SwaggerMessage
// @Router      /api/v1/admin/users [post]
func (h *AdminUserHandler) ProvisionUser(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	uidStr, ok := userIDStr.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
		return
	}
	callerID, err := uuid.Parse(uidStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
		return
	}

	roleVal, ok := c.Get("userRole")
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	callerRole, _ := roleVal.(string)

	var req ports.ProvisionUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	user, err := h.authService.ProvisionUser(c.Request.Context(), callerID, callerRole, req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case errors.Is(err, domain.ErrPermissionDenied):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case errors.Is(err, domain.ErrInvalidRole):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user})
}

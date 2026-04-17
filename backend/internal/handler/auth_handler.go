package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/ports"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authService ports.AuthService
}

func NewAuthHandler(authService ports.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Login autentica con email y contraseña.
// @Summary     Iniciar sesión
// @Description Devuelve un JWT (campo token). El cliente debe llamar después a GET /auth/me para el perfil completo.
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body LoginRequest true "Credenciales"
// @Success     200 {object} SwaggerLoginOK
// @Failure     400 {object} SwaggerMessage
// @Failure     401 {object} SwaggerMessage
// @Router      /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	token, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}

// Register crea una cuenta (rol por defecto client si no se envía role).
// @Summary     Registro
// @Description Registro público. Email único; 409 si el usuario ya existe.
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body ports.RegisterRequest true "Datos de registro"
// @Success     201 {object} SwaggerRegisterOK
// @Failure     400 {object} SwaggerMessage
// @Failure     409 {object} SwaggerMessage
// @Router      /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req ports.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Convert to service request
	serviceReq := ports.RegisterRequest{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      req.Role,
	}

	user, err := h.authService.Register(c.Request.Context(), serviceReq)
	if err != nil {
		log.Printf("Error when try to register user: %s", err.Error())
		log.Printf("/******************************************/")

		statusCode := http.StatusBadRequest
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			statusCode = http.StatusConflict
		}

		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user": gin.H{
			"id":        user.ID,
			"email":     user.Email,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"role":      user.Role,
			"isActive":  user.IsActive,
			"createdAt": user.CreatedAt,
			"updatedAt": user.UpdatedAt,
		},
	})
}

// Me devuelve el usuario asociado al JWT.
// @Summary     Perfil actual
// @Description Requiere cabecera Authorization con esquema Bearer y el JWT.
// @Tags        auth
// @Produce     json
// @Security    BearerAuth
// @Success     200 {object} SwaggerMeOK
// @Failure     401 {object} SwaggerMessage
// @Failure     404 {object} SwaggerMessage
// @Router      /api/v1/auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID"})
		return
	}

	user, err := h.authService.CurrentUser(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":        user.ID,
			"email":     user.Email,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"role":      user.Role,
			"isActive":  user.IsActive,
			"createdAt": user.CreatedAt,
			"updatedAt": user.UpdatedAt,
		},
	})
}

package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/ports"
)

type AppointmentHandler struct {
	appointmentService ports.AppointmentService
}

func NewAppointmentHandler(appointmentService ports.AppointmentService) *AppointmentHandler {
	return &AppointmentHandler{
		appointmentService: appointmentService,
	}
}

// CreateAppointmentRequest represents the request payload for creating an appointment
type CreateAppointmentRequest struct {
	CustomerID    string `json:"customerID"`    // ✅ camelCase
	EmployeeID    string `json:"employeeID"`    // ✅ camelCase
	CarID         string `json:"carID"`         // ✅ camelCase
	ScheduledTime string `json:"scheduledTime"` // ✅ camelCase
	ScheduledAt   string `json:"scheduledAt"`   // ✅ camelCase
	Reason        string `json:"reason"`
	Notes         string `json:"notes"`
	Status        string `json:"status"`
	ServiceType   string `json:"serviceType"` // ✅ camelCase
}

// UpdateAppointmentRequest represents the request payload for updating an appointment
type UpdateAppointmentRequest struct {
	CustomerID    string `json:"customerID"`
	EmployeeID    string `json:"employeeID"`
	CarID         string `json:"carId"`         // ✅ camelCase
	ScheduledTime string `json:"scheduledTime"` // ✅ camelCase
	ScheduledAt   string `json:"scheduledAt"`   // ✅ camelCase
	Reason        string `json:"reason"`
	Notes         string `json:"notes"`
	Status        string `json:"status"`
	ServiceType   string `json:"serviceType"` // ✅ camelCase
}

// AppointmentResponse represents the response payload for an appointment
type AppointmentResponse struct {
	ID         string  `json:"id"`
	ClientName string  `json:"clientName"`
	CarID      string  `json:"carId"`
	Service    string  `json:"service"`
	Date       string  `json:"date"`
	Time       string  `json:"time"`
	Status     string  `json:"status"`
	Notes      string  `json:"notes"`
	CreatedAt  string  `json:"createdAt"`
	UpdatedAt  string  `json:"updatedAt"`
	DeletedAt  *string `json:"deletedAt,omitempty"`
}

// CreateAppointment handles POST /api/v1/appointments
func (h *AppointmentHandler) CreateAppointment(c *gin.Context) {
	// Get user from Gin context
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
	// Parse request
	var req CreateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Parse campos
	carID, err := uuid.Parse(req.CarID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid car ID"})
		return
	}
	scheduledAt, err := time.Parse("2006-01-02T15:04:05Z07:00", req.ScheduledAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid scheduled time format"})
		return
	}

	// Convert to domain object
	appointment := &domain.Appointment{
		CustomerID:  userID,
		CarID:       carID,
		ScheduledAt: scheduledAt,
		Status:      domain.AppointmentStatus(req.Status),
		Notes:       req.Notes,
		ServiceType: req.ServiceType,
	}

	// Create appointment
	createdAppointment, err := h.appointmentService.CreateAppointment(c.Request.Context(), appointment, userID)
	if err != nil {
		if err == domain.ErrUnauthorizedAccess {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		if err == domain.ErrAppointmentAlreadyExists {
			c.JSON(http.StatusConflict, gin.H{"error": "appointment with this ID already exists"})
			return
		}
		if err == domain.ErrInvalidAppointmentData {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid appointment data"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Convert to response
	response := h.toAppointmentResponse(createdAppointment)

	c.JSON(http.StatusCreated, response)
}

// GetAppointment handles GET /api/v1/appointments/{id}
func (h *AppointmentHandler) GetAppointment(c *gin.Context) {
	// Get user from Gin context
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
	// Parse appointment ID from URL parameter
	appointmentIDStr := c.Param("id")
	appointmentID, err := uuid.Parse(appointmentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid appointment ID"})
		return
	}
	// Get appointment
	appointment, err := h.appointmentService.GetAppointment(c.Request.Context(), appointmentID, userID)
	if err != nil {
		if err == domain.ErrAppointmentNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "appointment not found"})
			return
		}
		if err == domain.ErrUnauthorizedAccess {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Convert to response
	response := h.toAppointmentResponse(appointment)

	c.JSON(http.StatusOK, response)
}

// ListAppointments handles GET /api/v1/appointments
func (h *AppointmentHandler) ListAppointments(c *gin.Context) {
	// Get user from Gin context
	// userIDStr, exists := c.Get("userID")
	// if !exists {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	// 	return
	// }
	// userID, err := uuid.Parse(userIDStr.(string))
	// if err != nil {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID"})
	// 	return
	// }
	// Parse query parameters for filters (if any)
	// For simplicity, we will not implement filters in this example
	var filters *ports.AppointmentFilters = nil
	// List appointments
	appointments, _, err := h.appointmentService.ListAppointments(c.Request.Context(), filters)
	if err != nil {
		if err == domain.ErrUnauthorizedAccess {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	// Convert to response
	var responses []AppointmentResponse
	for _, appointment := range appointments {
		responses = append(responses, h.toAppointmentResponse(appointment))
	}
	c.JSON(http.StatusOK, responses)
}

// UpdateAppointment handles PUT /api/v1/appointments/{id}
func (h *AppointmentHandler) UpdateAppointment(c *gin.Context) {
	// Get user from Gin context
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

	// Parse appointment ID from URL parameter
	appointmentIDStr := c.Param("id")
	appointmentID, err := uuid.Parse(appointmentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid appointment ID"})
		return
	}

	// Parse request
	var req UpdateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Convert to domain object
	// customerUUID, err := uuid.Parse(req.CustomerID)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid customer ID"})
	// 	return
	// }

	scheduledAt, err := time.Parse("2006-01-02T15:04:05Z07:00", req.ScheduledAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid scheduled time format"})
		return
	}

	carUUID, err := uuid.Parse(req.CarID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid car ID"})
		return
	}

	appointment := &domain.Appointment{
		ID:          appointmentID,
		CustomerID:  userID,
		CarID:       carUUID,
		ScheduledAt: scheduledAt,
		Status:      domain.AppointmentStatus(req.Status),
		Notes:       req.Notes,
		ServiceType: req.ServiceType,
	}

	// Update appointment
	appointment, err = h.appointmentService.UpdateAppointment(c.Request.Context(), appointment, userID)
	if err != nil {
		if err == domain.ErrAppointmentNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "appointment not found"})
			return
		}
		if err == domain.ErrUnauthorizedAccess {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Convert to response
	response := h.toAppointmentResponse(appointment)
	c.JSON(http.StatusOK, response)
}

// DeleteAppointment handles DELETE /api/v1/appointments/{id}
func (h *AppointmentHandler) DeleteAppointment(c *gin.Context) {
	// Get user from Gin context
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

	// Parse appointment ID from URL parameter
	appointmentIDStr := c.Param("id")
	appointmentID, err := uuid.Parse(appointmentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid appointment ID"})
		return
	}

	// Delete appointment
	err = h.appointmentService.DeleteAppointment(c.Request.Context(), appointmentID, userID)
	if err != nil {
		if err == domain.ErrAppointmentNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "appointment not found"})
			return
		}
		if err == domain.ErrUnauthorizedAccess {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.Status(http.StatusNoContent)
}

// Helper methods

func (h *AppointmentHandler) toAppointmentResponse(appointment *domain.Appointment) AppointmentResponse {
	// Extrai nome do cliente (podes buscar do domínio ou da DB, ou deixar vazio se não existir)
	clientName := "" // TODO: buscar nome do cliente se necessário

	// Extrai service (pode ser appointment.ServiceType ou outro campo)
	service := appointment.ServiceType

	// Divide data/hora
	date := appointment.ScheduledAt.Format("2006-01-02")
	time := appointment.ScheduledAt.Format("15:04")

	var deletedAt *string
	if appointment.DeletedAt != nil {
		s := appointment.DeletedAt.Format("2006-01-02T15:04:05Z07:00")
		deletedAt = &s
	}

	return AppointmentResponse{
		ID:         appointment.ID.String(),
		ClientName: clientName,
		CarID:      appointment.CarID.String(),
		Service:    service,
		Date:       date,
		Time:       time,
		Status:     string(appointment.Status),
		Notes:      appointment.Notes,
		CreatedAt:  appointment.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  appointment.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		DeletedAt:  deletedAt,
	}
}

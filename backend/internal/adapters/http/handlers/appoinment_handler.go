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
	Reason        string `json:"reason"`
	Status        string `json:"status"`
}

// UpdateAppointmentRequest represents the request payload for updating an appointment
type UpdateAppointmentRequest struct {
	CustomerID    string `json:"customerID"`    // ✅ camelCase
	EmployeeID    string `json:"employeeID"`    // ✅ camelCase
	CarID         string `json:"carID"`         // ✅ camelCase
	ScheduledTime string `json:"scheduledTime"` // ✅ camelCase
	Reason        string `json:"reason"`
	Status        string `json:"status"`
}

// AppointmentResponse represents the response payload for an appointment
type AppointmentResponse struct {
	ID          string `json:"id"`
	CustomerID  string `json:"customerID"`    // ✅ camelCase
	EmployeeID  string `json:"employeeID"`    // ✅ camelCase
	CarID       string `json:"carID"`         // ✅ camelCase
	ScheduledAt string `json:"scheduledTime"` // ✅ camelCase
	Reason      string `json:"reason"`
	Status      string `json:"status"`
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

	// Convert to domain object
	appointment := &domain.Appointment{
		CustomerID: userID,
		CarID:      uuid.MustParse(req.CarID),
		Status:     domain.AppointmentStatus(req.Status),
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
	customerUUID, err := uuid.Parse(req.CustomerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid customer ID"})
		return
	}
	scheduledAt, err := time.Parse("2006-01-02T15:04:05Z07:00", req.ScheduledTime)
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
		CustomerID:  customerUUID,
		CarID:       carUUID,
		ScheduledAt: scheduledAt,
		Status:      domain.AppointmentStatus(req.Status),
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
	return AppointmentResponse{
		ID:         appointment.ID.String(),
		CustomerID: appointment.CustomerID.String(),
		//EmployeeID:  appointment.EmployeeID.String(),
		CarID:       appointment.CarID.String(),
		ScheduledAt: appointment.ScheduledAt.Format("2006-01-02T15:04:05Z07:00"),
		//Reason:      appointment.Reason,
		Status: string(appointment.Status),
	}
}

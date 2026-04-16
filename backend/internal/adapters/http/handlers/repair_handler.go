package handlers

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/ports"
)

func parseOptionalRFC3339OrDate(s *string) (*time.Time, error) {
	if s == nil {
		return nil, nil
	}
	raw := strings.TrimSpace(*s)
	if raw == "" {
		return nil, nil
	}
	if t, err := time.Parse(time.RFC3339, raw); err == nil {
		return &t, nil
	}
	if t, err := time.Parse("2006-01-02T15:04:05Z07:00", raw); err == nil {
		return &t, nil
	}
	t, err := time.Parse("2006-01-02", raw)
	if err != nil {
		return nil, err
	}
	utc := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	return &utc, nil
}

// RepairHandler expone reparaciones (taller CRUD parcial; cliente lectura por coche).
type RepairHandler struct {
	repairService ports.RepairService
}

func NewRepairHandler(repairService ports.RepairService) *RepairHandler {
	return &RepairHandler{repairService: repairService}
}

// CreateRepairRequest cuerpo POST /repairs.
type CreateRepairRequest struct {
	CarID       string  `json:"carId" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Status      string  `json:"status,omitempty"`
	StartedAt   *string `json:"startedAt,omitempty"`
	Cost        float64 `json:"cost"`
}

// UpdateRepairRequest cuerpo PUT /repairs/:id (campos vacíos no sustituyen descripción/estado).
type UpdateRepairRequest struct {
	Description string  `json:"description"`
	Status      string  `json:"status"`
	Cost        float64 `json:"cost"`
	StartedAt   *string `json:"startedAt,omitempty"`
	CompletedAt *string `json:"completedAt,omitempty"`
}

// RepairResponse respuesta JSON camelCase.
type RepairResponse struct {
	ID             string  `json:"id"`
	CarID          string  `json:"carId"`
	TechnicianID   string  `json:"technicianId"`
	Description    string  `json:"description"`
	Status         string  `json:"status"`
	Cost           float64 `json:"cost"`
	StartedAt      *string `json:"startedAt,omitempty"`
	CompletedAt    *string `json:"completedAt,omitempty"`
	CreatedAt      string  `json:"createdAt"`
	UpdatedAt      string  `json:"updatedAt"`
}

func repairToResponse(r *domain.Repair) RepairResponse {
	out := RepairResponse{
		ID:           r.ID.String(),
		CarID:        r.CarID.String(),
		TechnicianID: r.TechnicianID.String(),
		Description:  r.Description,
		Status:       string(r.Status),
		Cost:         r.Cost,
		CreatedAt:    r.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:    r.UpdatedAt.UTC().Format(time.RFC3339),
	}
	if r.StartedAt != nil {
		s := r.StartedAt.UTC().Format(time.RFC3339)
		out.StartedAt = &s
	}
	if r.CompletedAt != nil {
		s := r.CompletedAt.UTC().Format(time.RFC3339)
		out.CompletedAt = &s
	}
	return out
}

func writeRepairError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrUnauthorizedAccess):
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
	case errors.Is(err, domain.ErrRepairNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "repair not found"})
	case errors.Is(err, domain.ErrCarNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "car not found"})
	case errors.Is(err, domain.ErrInvalidRepairData):
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid repair data"})
	case errors.Is(err, domain.ErrUserNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

// CreateRepair registra una reparación (solo personal del taller).
// @Summary     Crear reparación
// @Tags        repairs
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       body body CreateRepairRequest true "carId, description, cost; startedAt opcional (RFC3339)"
// @Success     201 {object} RepairResponse
// @Failure     400 {object} SwaggerMessage
// @Failure     401 {object} SwaggerMessage
// @Failure     403 {object} SwaggerMessage
// @Failure     404 {object} SwaggerMessage
// @Router      /api/v1/repairs [post]
func (h *RepairHandler) CreateRepair(c *gin.Context) {
	userID, err := uuid.Parse(c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID"})
		return
	}

	var req CreateRepairRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	carID, err := uuid.Parse(strings.TrimSpace(req.CarID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid carId"})
		return
	}

	startedAt, err := parseOptionalRFC3339OrDate(req.StartedAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid startedAt format"})
		return
	}

	repair := &domain.Repair{
		CarID:       carID,
		Description: req.Description,
		Status:      domain.RepairStatus(req.Status),
		StartedAt:   startedAt,
		Cost:        req.Cost,
	}

	created, err := h.repairService.CreateRepair(c.Request.Context(), repair, userID)
	if err != nil {
		writeRepairError(c, err)
		return
	}

	c.JSON(http.StatusCreated, repairToResponse(created))
}

// ListRepairsByCar lista reparaciones de un coche (cliente: solo su coche; staff: cualquier coche).
// @Summary     Listar reparaciones por coche
// @Tags        repairs
// @Security    BearerAuth
// @Produce     json
// @Param       id path string true "ID del coche"
// @Success     200 {array} RepairResponse
// @Failure     401 {object} SwaggerMessage
// @Failure     403 {object} SwaggerMessage
// @Failure     404 {object} SwaggerMessage
// @Router      /api/v1/cars/{id}/repairs [get]
func (h *RepairHandler) ListRepairsByCar(c *gin.Context) {
	userID, err := uuid.Parse(c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID"})
		return
	}

	carID, err := uuid.Parse(strings.TrimSpace(c.Param("id")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid car id"})
		return
	}

	repairs, err := h.repairService.GetRepairsByCarID(c.Request.Context(), carID, userID)
	if err != nil {
		writeRepairError(c, err)
		return
	}

	out := make([]RepairResponse, 0, len(repairs))
	for _, r := range repairs {
		out = append(out, repairToResponse(r))
	}
	c.JSON(http.StatusOK, out)
}

// GetRepair obtiene una reparación por id.
// @Summary     Obtener reparación
// @Tags        repairs
// @Security    BearerAuth
// @Produce     json
// @Param       id path string true "ID de la reparación"
// @Success     200 {object} RepairResponse
// @Failure     401 {object} SwaggerMessage
// @Failure     403 {object} SwaggerMessage
// @Failure     404 {object} SwaggerMessage
// @Router      /api/v1/repairs/{id} [get]
func (h *RepairHandler) GetRepair(c *gin.Context) {
	userID, err := uuid.Parse(c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID"})
		return
	}

	repairID, err := uuid.Parse(strings.TrimSpace(c.Param("id")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid repair id"})
		return
	}

	repair, err := h.repairService.GetRepair(c.Request.Context(), repairID, userID)
	if err != nil {
		writeRepairError(c, err)
		return
	}

	c.JSON(http.StatusOK, repairToResponse(repair))
}

// UpdateRepair actualiza una reparación (solo personal del taller).
// @Summary     Actualizar reparación
// @Tags        repairs
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       id path string true "ID de la reparación"
// @Param       body body UpdateRepairRequest true "Campos a actualizar"
// @Success     200 {object} RepairResponse
// @Failure     400 {object} SwaggerMessage
// @Failure     401 {object} SwaggerMessage
// @Failure     403 {object} SwaggerMessage
// @Failure     404 {object} SwaggerMessage
// @Router      /api/v1/repairs/{id} [put]
func (h *RepairHandler) UpdateRepair(c *gin.Context) {
	userID, err := uuid.Parse(c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID"})
		return
	}

	repairID, err := uuid.Parse(strings.TrimSpace(c.Param("id")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid repair id"})
		return
	}

	var req UpdateRepairRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	startedAt, err := parseOptionalRFC3339OrDate(req.StartedAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid startedAt format"})
		return
	}
	completedAt, err := parseOptionalRFC3339OrDate(req.CompletedAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid completedAt format"})
		return
	}

	patch := &domain.Repair{
		ID:          repairID,
		Description: req.Description,
		Status:      domain.RepairStatus(req.Status),
		Cost:        req.Cost,
		StartedAt:   startedAt,
		CompletedAt: completedAt,
	}

	updated, err := h.repairService.UpdateRepair(c.Request.Context(), patch, userID)
	if err != nil {
		writeRepairError(c, err)
		return
	}

	c.JSON(http.StatusOK, repairToResponse(updated))
}

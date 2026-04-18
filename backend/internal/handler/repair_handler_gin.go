package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/domain"
)

// ListRepairsByCar returns repairs for a car (client: own car only).
// @Summary     Listar reparaciones por coche
// @Description Cliente: solo su coche. Personal del taller según reglas del dominio.
// @Tags        repairs
// @Security    BearerAuth
// @Produce     json
// @Param       carId path string true "UUID del coche"
// @Success     200 {array} RepairResponse
// @Failure     400 {object} SwaggerMessage
// @Failure     401 {object} SwaggerMessage
// @Failure     403 {object} SwaggerMessage
// @Failure     500 {object} SwaggerMessage
// @Router      /api/v1/repairs/car/{carId} [get]
func (h *RepairHandler) ListRepairsByCar(c *gin.Context) {
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

	carIDStr := c.Param("carId")
	if carIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "car ID is required"})
		return
	}

	carID, err := uuid.Parse(carIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid car ID"})
		return
	}

	repairs, err := h.repairService.GetRepairsByCarID(c.Request.Context(), carID, userID)
	if err != nil {
		if err == domain.ErrUnauthorizedAccess {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if repairs == nil {
		c.JSON(http.StatusOK, []domain.Repair{})
		return
	}

	out := make([]domain.Repair, 0, len(repairs))
	for _, r := range repairs {
		if r != nil {
			out = append(out, *r)
		}
	}
	c.JSON(http.StatusOK, out)
}

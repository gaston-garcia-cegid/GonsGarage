package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/ports"
)

type CarHandler struct {
	carService ports.CarService
}

func NewCarHandler(carService ports.CarService) *CarHandler {
	return &CarHandler{
		carService: carService,
	}
}

// CreateCarRequest represents the request payload for creating a car
type CreateCarRequest struct {
	Make         string `json:"make"`
	Model        string `json:"model"`
	Year         int    `json:"year"`
	LicensePlate string `json:"licensePlate"`
	VIN          string `json:"vin"`
	Color        string `json:"color"`
	Mileage      int    `json:"mileage"`
}

// UpdateCarRequest represents the request payload for updating a car
type UpdateCarRequest struct {
	Make         string `json:"make"`
	Model        string `json:"model"`
	Year         int    `json:"year"`
	LicensePlate string `json:"licensePlate"` // ✅ camelCase
	VIN          string `json:"vin"`
	Color        string `json:"color"`
	Mileage      int    `json:"mileage"`
}

// CarResponse represents the response payload for a car
type CarResponse struct {
	ID           string `json:"id"`
	Make         string `json:"make"`
	Model        string `json:"model"`
	Year         int    `json:"year"`
	LicensePlate string `json:"licensePlate"` // ✅ camelCase
	VIN          string `json:"vin"`
	Color        string `json:"color"`
	Mileage      int    `json:"mileage"`
	OwnerID      string `json:"ownerID"`   // ✅ camelCase
	CreatedAt    string `json:"createdAt"` // ✅ camelCase
	UpdatedAt    string `json:"updatedAt"` // ✅ camelCase
}

// CreateCar handles POST /api/v1/cars
func (h *CarHandler) CreateCar(c *gin.Context) {
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
	var req CreateCarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Convert to domain object
	car := &domain.Car{
		Make:         req.Make,
		Model:        req.Model,
		Year:         req.Year,
		LicensePlate: req.LicensePlate,
		VIN:          req.VIN,
		Color:        req.Color,
		Mileage:      req.Mileage,
	}

	// Create car
	createdCar, err := h.carService.CreateCar(c.Request.Context(), car, userID)
	if err != nil {
		if err == domain.ErrUnauthorizedAccess {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		if err == domain.ErrCarAlreadyExists {
			c.JSON(http.StatusConflict, gin.H{"error": "car with this license plate already exists"})
			return
		}
		if err == domain.ErrInvalidCarData {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid car data"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Convert to response
	response := h.toCarResponse(createdCar)

	c.JSON(http.StatusCreated, response)
}

// GetCar handles GET /api/v1/cars/{id}
func (h *CarHandler) GetCar(c *gin.Context) {
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

	// Parse car ID from URL parameter
	carIDStr := c.Param("id")
	carID, err := uuid.Parse(carIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid car ID"})
		return
	}

	car, err := h.carService.GetCar(c.Request.Context(), carID, userID)
	if err != nil {
		if err == domain.ErrCarNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "car not found"})
			return
		}
		if err == domain.ErrUnauthorizedAccess {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	response := h.toCarResponse(car)
	c.JSON(http.StatusOK, response)
}

// ListCars handles GET /api/v1/cars
func (h *CarHandler) ListCars(c *gin.Context) {
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

	// For clients, list their own cars
	// For admins/managers, this could list all cars with pagination
	cars, err := h.carService.GetCarsByOwner(c.Request.Context(), userID, userID)
	if err != nil {
		if err == domain.ErrUnauthorizedAccess {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Convert to response
	responses := make([]CarResponse, len(cars))
	for i, car := range cars {
		responses[i] = h.toCarResponse(car)
	}

	c.JSON(http.StatusOK, responses)
}

// UpdateCar handles PUT /api/v1/cars/{id}
func (h *CarHandler) UpdateCar(c *gin.Context) {
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

	// Parse car ID from URL parameter
	carIDStr := c.Param("id")
	carID, err := uuid.Parse(carIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid car ID"})
		return
	}

	// Parse request
	var req UpdateCarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Convert to domain object
	car := &domain.Car{
		ID:           carID,
		Make:         req.Make,
		Model:        req.Model,
		Year:         req.Year,
		LicensePlate: req.LicensePlate,
		VIN:          req.VIN,
		Color:        req.Color,
		Mileage:      req.Mileage,
	}

	// Update car
	updatedCar, err := h.carService.UpdateCar(c.Request.Context(), car, userID)
	if err != nil {
		if err == domain.ErrCarNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "car not found"})
			return
		}
		if err == domain.ErrUnauthorizedAccess {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		if err == domain.ErrInvalidCarData {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid car data"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Convert to response
	response := h.toCarResponse(updatedCar)

	c.JSON(http.StatusOK, response)
}

// DeleteCar handles DELETE /api/v1/cars/{id}
func (h *CarHandler) DeleteCar(c *gin.Context) {
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

	// Parse car ID from URL parameter
	carIDStr := c.Param("id")
	carID, err := uuid.Parse(carIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid car ID"})
		return
	}

	// Delete car
	err = h.carService.DeleteCar(c.Request.Context(), carID, userID)
	if err != nil {
		if err == domain.ErrCarNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "car not found"})
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

func (h *CarHandler) toCarResponse(car *domain.Car) CarResponse {
	return CarResponse{
		ID:           car.ID.String(),
		Make:         car.Make,
		Model:        car.Model,
		Year:         car.Year,
		LicensePlate: car.LicensePlate,
		VIN:          car.VIN,
		Color:        car.Color,
		Mileage:      car.Mileage,
		OwnerID:      car.OwnerID.String(),
		CreatedAt:    car.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    car.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

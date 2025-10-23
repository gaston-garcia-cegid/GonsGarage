package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

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
func (h *CarHandler) CreateCar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get user from middleware
	userID, ok := getUserIDFromContext(ctx)
	if !ok {
		h.respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Parse request
	var req CreateCarRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body")
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
	createdCar, err := h.carService.CreateCar(ctx, car, userID)
	if err != nil {
		if err == domain.ErrUnauthorizedAccess {
			h.respondError(w, http.StatusForbidden, "forbidden")
			return
		}
		if err == domain.ErrCarAlreadyExists {
			h.respondError(w, http.StatusConflict, "car with this license plate already exists")
			return
		}
		if err == domain.ErrInvalidCarData {
			h.respondError(w, http.StatusBadRequest, "invalid car data")
			return
		}

		//h.logger.Error("failed to create car", "error", err, "userID", userID)
		h.respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	// Convert to response
	response := h.toCarResponse(createdCar)

	h.respondJSON(w, http.StatusCreated, response)
}

// GetCar handles GET /api/v1/cars/{id}
func (h *CarHandler) GetCar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get user from middleware
	userID, ok := getUserIDFromContext(ctx)
	if !ok {
		h.respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Parse car ID
	vars := mux.Vars(r)
	carID, err := uuid.Parse(vars["id"])
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid car ID")
		return
	}

	// Get car
	car, err := h.carService.GetCar(ctx, carID, userID)
	if err != nil {
		if err == domain.ErrCarNotFound {
			h.respondError(w, http.StatusNotFound, "car not found")
			return
		}
		if err == domain.ErrUnauthorizedAccess {
			h.respondError(w, http.StatusForbidden, "forbidden")
			return
		}

		//h.logger.Error("failed to get car", "error", err, "carID", carID, "userID", userID)
		h.respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	// Convert to response
	response := h.toCarResponse(car)

	h.respondJSON(w, http.StatusOK, response)
}

// ListCars handles GET /api/v1/cars
func (h *CarHandler) ListCars(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get user from middleware
	userID, ok := getUserIDFromContext(ctx)
	if !ok {
		h.respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// For clients, list their own cars
	// For admins/managers, this could list all cars with pagination
	cars, err := h.carService.GetCarsByOwner(ctx, userID, userID)
	if err != nil {
		if err == domain.ErrUnauthorizedAccess {
			h.respondError(w, http.StatusForbidden, "forbidden")
			return
		}

		//h.logger.Error("failed to list cars", "error", err, "userID", userID)
		h.respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	// Convert to response
	responses := make([]CarResponse, len(cars))
	for i, car := range cars {
		responses[i] = h.toCarResponse(car)
	}

	h.respondJSON(w, http.StatusOK, responses)
}

// UpdateCar handles PUT /api/v1/cars/{id}
func (h *CarHandler) UpdateCar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get user from middleware
	userID, ok := getUserIDFromContext(ctx)
	if !ok {
		h.respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Parse car ID
	vars := mux.Vars(r)
	carID, err := uuid.Parse(vars["id"])
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid car ID")
		return
	}

	// Parse request
	var req UpdateCarRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body")
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
	updatedCar, err := h.carService.UpdateCar(ctx, car, userID)
	if err != nil {
		if err == domain.ErrCarNotFound {
			h.respondError(w, http.StatusNotFound, "car not found")
			return
		}
		if err == domain.ErrUnauthorizedAccess {
			h.respondError(w, http.StatusForbidden, "forbidden")
			return
		}
		if err == domain.ErrInvalidCarData {
			h.respondError(w, http.StatusBadRequest, "invalid car data")
			return
		}

		//h.logger.Error("failed to update car", "error", err, "carID", carID, "userID", userID)
		h.respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	// Convert to response
	response := h.toCarResponse(updatedCar)

	h.respondJSON(w, http.StatusOK, response)
}

// DeleteCar handles DELETE /api/v1/cars/{id}
func (h *CarHandler) DeleteCar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get user from middleware
	userID, ok := getUserIDFromContext(ctx)
	if !ok {
		h.respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Parse car ID
	vars := mux.Vars(r)
	carID, err := uuid.Parse(vars["id"])
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid car ID")
		return
	}

	// Delete car
	err = h.carService.DeleteCar(ctx, carID, userID)
	if err != nil {
		if err == domain.ErrCarNotFound {
			h.respondError(w, http.StatusNotFound, "car not found")
			return
		}
		if err == domain.ErrUnauthorizedAccess {
			h.respondError(w, http.StatusForbidden, "forbidden")
			return
		}

		//h.logger.Error("failed to delete car", "error", err, "carID", carID, "userID", userID)
		h.respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
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

func (h *CarHandler) respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func (h *CarHandler) respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// getUserIDFromContext extracts user ID from request context
func getUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return uuid.Nil, false
	}

	id, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, false
	}

	return id, true
}

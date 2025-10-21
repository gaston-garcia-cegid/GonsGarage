package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/adapters/http/middleware"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
	carUseCase "github.com/gaston-garcia-cegid/gonsgarage/internal/core/usecases/car"
)

type CarHandler struct {
	carUseCase *carUseCase.CarUseCase
}

func NewCarHandler(carUseCase *carUseCase.CarUseCase) *CarHandler {
	return &CarHandler{
		carUseCase: carUseCase,
	}
}

// CreateCarRequest represents the request payload for creating a car
type CreateCarRequest struct {
	ClientID     string `json:"client_id"`
	Make         string `json:"make"`
	Model        string `json:"model"`
	Year         int    `json:"year"`
	LicensePlate string `json:"license_plate"`
	VIN          string `json:"vin"`
	Color        string `json:"color"`
	Mileage      int    `json:"mileage"`
}

// UpdateCarRequest represents the request payload for updating a car
type UpdateCarRequest struct {
	Make         string `json:"make"`
	Model        string `json:"model"`
	Year         int    `json:"year"`
	LicensePlate string `json:"license_plate"`
	VIN          string `json:"vin"`
	Color        string `json:"color"`
	Mileage      int    `json:"mileage"`
}

// CarResponse represents the response payload for car data
type CarResponse struct {
	ID           string        `json:"id"`
	ClientID     string        `json:"client_id"`
	Make         string        `json:"make"`
	Model        string        `json:"model"`
	Year         int           `json:"year"`
	LicensePlate string        `json:"license_plate"`
	VIN          string        `json:"vin"`
	Color        string        `json:"color"`
	Mileage      int           `json:"mileage"`
	CreatedAt    string        `json:"created_at"`
	UpdatedAt    string        `json:"updated_at"`
	Client       *UserResponse `json:"client,omitempty"`
}

// UserResponse represents the response payload for user data
type UserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
	IsActive  bool   `json:"is_active"`
}

func (h *CarHandler) CreateCar(w http.ResponseWriter, r *http.Request) {
	// Get user from middleware
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateCarRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Parse client ID
	clientID, err := uuid.Parse(req.ClientID)
	if err != nil {
		http.Error(w, "Invalid client ID", http.StatusBadRequest)
		return
	}

	// Create domain car
	car := &domain.Car{
		ClientID:     clientID,
		Make:         req.Make,
		Model:        req.Model,
		Year:         req.Year,
		LicensePlate: req.LicensePlate,
		VIN:          req.VIN,
		Color:        req.Color,
		Mileage:      req.Mileage,
	}

	// Create car using use case
	createdCar, err := h.carUseCase.CreateCar(r.Context(), car, userID)
	if err != nil {
		if err == domain.ErrUnauthorizedAccess {
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert to response
	response := h.toCarResponse(createdCar)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *CarHandler) GetCar(w http.ResponseWriter, r *http.Request) {
	// Get user from middleware
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get car ID from URL
	vars := mux.Vars(r)
	carIDStr, exists := vars["id"]
	if !exists {
		http.Error(w, "Car ID is required", http.StatusBadRequest)
		return
	}

	carID, err := uuid.Parse(carIDStr)
	if err != nil {
		http.Error(w, "Invalid car ID", http.StatusBadRequest)
		return
	}

	// Get car using use case
	car, err := h.carUseCase.GetCar(r.Context(), carID, userID)
	if err != nil {
		if err == domain.ErrCarNotFound {
			http.Error(w, "Car not found", http.StatusNotFound)
			return
		}
		if err == domain.ErrUnauthorizedAccess {
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to response
	response := h.toCarResponse(car)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *CarHandler) ListCars(w http.ResponseWriter, r *http.Request) {
	// Get user from middleware
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 50 // default
	offset := 0 // default

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// Get cars using use case
	cars, err := h.carUseCase.ListCars(r.Context(), userID, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to response
	response := make([]*CarResponse, len(cars))
	for i, car := range cars {
		response[i] = h.toCarResponse(car)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"cars":   response,
		"total":  len(response),
		"limit":  limit,
		"offset": offset,
	})
}

func (h *CarHandler) UpdateCar(w http.ResponseWriter, r *http.Request) {
	// Get user from middleware
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get car ID from URL
	vars := mux.Vars(r)
	carIDStr, exists := vars["id"]
	if !exists {
		http.Error(w, "Car ID is required", http.StatusBadRequest)
		return
	}

	carID, err := uuid.Parse(carIDStr)
	if err != nil {
		http.Error(w, "Invalid car ID", http.StatusBadRequest)
		return
	}

	var req UpdateCarRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get existing car to preserve client_id
	existingCar, err := h.carUseCase.GetCar(r.Context(), carID, userID)
	if err != nil {
		if err == domain.ErrCarNotFound {
			http.Error(w, "Car not found", http.StatusNotFound)
			return
		}
		if err == domain.ErrUnauthorizedAccess {
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create updated car
	car := &domain.Car{
		ID:           carID,
		ClientID:     existingCar.ClientID, // Preserve client ID
		Make:         req.Make,
		Model:        req.Model,
		Year:         req.Year,
		LicensePlate: req.LicensePlate,
		VIN:          req.VIN,
		Color:        req.Color,
		Mileage:      req.Mileage,
	}

	// Update car using use case
	updatedCar, err := h.carUseCase.UpdateCar(r.Context(), car, userID)
	if err != nil {
		if err == domain.ErrUnauthorizedAccess {
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert to response
	response := h.toCarResponse(updatedCar)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *CarHandler) DeleteCar(w http.ResponseWriter, r *http.Request) {
	// Get user from middleware
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get car ID from URL
	vars := mux.Vars(r)
	carIDStr, exists := vars["id"]
	if !exists {
		http.Error(w, "Car ID is required", http.StatusBadRequest)
		return
	}

	carID, err := uuid.Parse(carIDStr)
	if err != nil {
		http.Error(w, "Invalid car ID", http.StatusBadRequest)
		return
	}

	// Delete car using use case
	err = h.carUseCase.DeleteCar(r.Context(), carID, userID)
	if err != nil {
		if err == domain.ErrCarNotFound {
			http.Error(w, "Car not found", http.StatusNotFound)
			return
		}
		if err == domain.ErrUnauthorizedAccess {
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *CarHandler) toCarResponse(car *domain.Car) *CarResponse {
	response := &CarResponse{
		ID:           car.ID.String(),
		ClientID:     car.ClientID.String(),
		Make:         car.Make,
		Model:        car.Model,
		Year:         car.Year,
		LicensePlate: car.LicensePlate,
		VIN:          car.VIN,
		Color:        car.Color,
		Mileage:      car.Mileage,
		CreatedAt:    car.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:    car.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	if car.Client != nil {
		response.Client = &UserResponse{
			ID:        car.Client.ID.String(),
			Email:     car.Client.Email,
			FirstName: car.Client.FirstName,
			LastName:  car.Client.LastName,
			Role:      car.Client.Role,
			IsActive:  car.Client.IsActive,
		}
	}

	return response
}

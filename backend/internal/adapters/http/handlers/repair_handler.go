package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/adapters/http/middleware"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
	repairUseCase "github.com/gaston-garcia-cegid/gonsgarage/internal/core/usecases/repair"
)

type RepairHandler struct {
	repairUseCase *repairUseCase.RepairUseCase
}

func NewRepairHandler(repairUseCase *repairUseCase.RepairUseCase) *RepairHandler {
	return &RepairHandler{
		repairUseCase: repairUseCase,
	}
}

// CreateRepairRequest represents the request payload for creating a repair
type CreateRepairRequest struct {
	CarID       string  `json:"car_id"`
	Description string  `json:"description"`
	Status      string  `json:"status,omitempty"`
	StartDate   string  `json:"start_date"`
	Cost        float64 `json:"cost"`
}

// UpdateRepairRequest represents the request payload for updating a repair
type UpdateRepairRequest struct {
	Description string  `json:"description"`
	Status      string  `json:"status"`
	EndDate     *string `json:"end_date,omitempty"`
	Cost        float64 `json:"cost"`
}

// RepairResponse represents the response payload for repair data
type RepairResponse struct {
	ID          string        `json:"id"`
	CarID       string        `json:"car_id"`
	EmployeeID  string        `json:"employee_id"`
	Description string        `json:"description"`
	Status      string        `json:"status"`
	StartDate   string        `json:"start_date"`
	EndDate     *string       `json:"end_date,omitempty"`
	Cost        float64       `json:"cost"`
	CreatedAt   string        `json:"created_at"`
	UpdatedAt   string        `json:"updated_at"`
	Car         *CarResponse  `json:"car,omitempty"`
	Employee    *UserResponse `json:"employee,omitempty"`
}

// UserResponse represents the response payload for user data (employee)
type UserResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func (h *RepairHandler) CreateRepair(w http.ResponseWriter, r *http.Request) {
	// Get user from middleware
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateRepairRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Parse car ID
	carID, err := uuid.Parse(req.CarID)
	if err != nil {
		http.Error(w, "Invalid car ID", http.StatusBadRequest)
		return
	}

	// Parse start date
	startDate, err := time.Parse("2006-01-02T15:04:05Z", req.StartDate)
	if err != nil {
		// Try alternative format
		startDate, err = time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			http.Error(w, "Invalid start date format", http.StatusBadRequest)
			return
		}
	}

	// Create domain repair
	repair := &domain.Repair{
		CarID:       carID,
		Description: req.Description,
		Status:      domain.RepairStatus(req.Status),
		StartedAt:   &startDate,
		Cost:        req.Cost,
	}

	// Create repair using use case
	createdRepair, err := h.repairUseCase.CreateRepair(r.Context(), repair, userID)
	if err != nil {
		if err == domain.ErrUnauthorizedAccess {
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert to response
	response := h.toRepairResponse(createdRepair)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *RepairHandler) GetRepair(w http.ResponseWriter, r *http.Request) {
	// Get user from middleware
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get repair ID from URL
	vars := mux.Vars(r)
	repairIDStr, exists := vars["id"]
	if !exists {
		http.Error(w, "Repair ID is required", http.StatusBadRequest)
		return
	}

	repairID, err := uuid.Parse(repairIDStr)
	if err != nil {
		http.Error(w, "Invalid repair ID", http.StatusBadRequest)
		return
	}

	// Get repair using use case
	repair, err := h.repairUseCase.GetRepair(r.Context(), repairID, userID)
	if err != nil {
		if err == domain.ErrRepairNotFound {
			http.Error(w, "Repair not found", http.StatusNotFound)
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
	response := h.toRepairResponse(repair)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *RepairHandler) GetRepairsByCarID(w http.ResponseWriter, r *http.Request) {
	// Get user from middleware
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get car ID from URL
	vars := mux.Vars(r)
	carIDStr, exists := vars["carId"]
	if !exists {
		http.Error(w, "Car ID is required", http.StatusBadRequest)
		return
	}

	carID, err := uuid.Parse(carIDStr)
	if err != nil {
		http.Error(w, "Invalid car ID", http.StatusBadRequest)
		return
	}

	// Get repairs using use case
	repairs, err := h.repairUseCase.GetRepairsByCarID(r.Context(), carID, userID)
	if err != nil {
		if err == domain.ErrUnauthorizedAccess {
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to response
	response := make([]*RepairResponse, len(repairs))
	for i, repair := range repairs {
		response[i] = h.toRepairResponse(repair)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"repairs": response,
		"total":   len(response),
	})
}

func (h *RepairHandler) UpdateRepair(w http.ResponseWriter, r *http.Request) {
	// Get user from middleware
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get repair ID from URL
	vars := mux.Vars(r)
	repairIDStr, exists := vars["id"]
	if !exists {
		http.Error(w, "Repair ID is required", http.StatusBadRequest)
		return
	}

	repairID, err := uuid.Parse(repairIDStr)
	if err != nil {
		http.Error(w, "Invalid repair ID", http.StatusBadRequest)
		return
	}

	var req UpdateRepairRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get existing repair
	existingRepair, err := h.repairUseCase.GetRepair(r.Context(), repairID, userID)
	if err != nil {
		if err == domain.ErrRepairNotFound {
			http.Error(w, "Repair not found", http.StatusNotFound)
			return
		}
		if err == domain.ErrUnauthorizedAccess {
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create updated repair
	repair := &domain.Repair{
		ID:           repairID,
		CarID:        existingRepair.CarID,
		TechnicianID: existingRepair.TechnicianID,
		Description:  req.Description,
		Status:       domain.RepairStatus(req.Status),
		StartedAt:    existingRepair.StartedAt,
		Cost:         req.Cost,
	}

	// Parse end date if provided
	if req.EndDate != nil {
		endDate, err := time.Parse("2006-01-02T15:04:05Z", *req.EndDate)
		if err != nil {
			// Try alternative format
			endDate, err = time.Parse("2006-01-02", *req.EndDate)
			if err != nil {
				http.Error(w, "Invalid end date format", http.StatusBadRequest)
				return
			}
		}
		repair.CompletedAt = &endDate
	}

	// Update repair using use case
	updatedRepair, err := h.repairUseCase.UpdateRepair(r.Context(), repair, userID)
	if err != nil {
		if err == domain.ErrUnauthorizedAccess {
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert to response
	response := h.toRepairResponse(updatedRepair)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *RepairHandler) toRepairResponse(repair *domain.Repair) *RepairResponse {
	response := &RepairResponse{
		ID:          repair.ID.String(),
		CarID:       repair.CarID.String(),
		EmployeeID:  repair.TechnicianID.String(),
		Description: repair.Description,
		Status:      string(repair.Status),
		StartDate:   repair.StartedAt.Format("2006-01-02T15:04:05Z"),
		Cost:        repair.Cost,
		CreatedAt:   repair.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   repair.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	if repair.CompletedAt != nil {
		endDate := repair.CompletedAt.Format("2006-01-02T15:04:05Z")
		response.EndDate = &endDate
	}

	if repair.Car.ID != uuid.Nil {
		response.Car = &CarResponse{
			ID:           repair.Car.ID.String(),
			Make:         repair.Car.Make,
			Model:        repair.Car.Model,
			Year:         repair.Car.Year,
			LicensePlate: repair.Car.LicensePlate,
			Color:        repair.Car.Color,
		}
	}

	if repair.Technician.ID != uuid.Nil {
		response.Employee = &UserResponse{
			ID:        repair.Technician.ID.String(),
			FirstName: repair.Technician.FirstName,
			LastName:  repair.Technician.LastName,
			Email:     repair.Technician.Email,
		}
	}

	return response
}

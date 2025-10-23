package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/ports"
)

type ClientHandler struct {
	clientUseCase ports.ClientUseCase
}

func NewClientHandler(clientUseCase ports.ClientUseCase) *ClientHandler {
	return &ClientHandler{
		clientUseCase: clientUseCase,
	}
}

func (h *ClientHandler) CreateClient(w http.ResponseWriter, r *http.Request) {
	var client ports.ClientUseCase
	if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.clientUseCase.CreateClient(r.Context(), &domain.Client{}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(client)
}

func (h *ClientHandler) ListClients(w http.ResponseWriter, r *http.Request) {
	clients, err := h.clientUseCase.ListClients(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(clients)
}

func (h *ClientHandler) GetClient(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/clients/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	client, err := h.clientUseCase.GetClient(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(client)
}

func (h *ClientHandler) UpdateClient(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/clients/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var client ports.ClientUseCase
	if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.clientUseCase.UpdateClient(r.Context(), id, &domain.Client{}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(client)
}

func (h *ClientHandler) DeleteClient(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/clients/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.clientUseCase.DeleteClient(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// func (h *CarHandler) ListCarsPaginated(w http.ResponseWriter, r *http.Request) {
// 	// Get user from middleware
// 	userID, ok := middleware.GetUserIDFromContext(r.Context())
// 	if !ok {
// 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	// Get pagination parameters
// 	pageStr := r.URL.Query().Get("page")
// 	page, err := strconv.Atoi(pageStr)
// 	if err != nil || page < 1 {
// 		page = 1
// 	}

// 	pageSizeStr := r.URL.Query().Get("page_size")
// 	pageSize, err := strconv.Atoi(pageSizeStr)
// 	if err != nil || pageSize < 1 {
// 		pageSize = 10
// 	}

// 	cars, err := h.carUseCase.ListCarsPaginated(r.Context(), userID, page, pageSize)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	json.NewEncoder(w).Encode(cars)
// }
// func (h *CarHandler) GetCar(w http.ResponseWriter, r *http.Request) {
// 	// Get user from middleware
// 	userID, ok := middleware.GetUserIDFromContext(r.Context())
// 	if !ok {
// 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	carIDStr := r.URL.Path[len("/cars/"):]
// 	carID, err := strconv.Atoi(carIDStr)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	car, err := h.carUseCase.GetCar(r.Context(), userID, carID)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	json.NewEncoder(w).Encode(car)
// }

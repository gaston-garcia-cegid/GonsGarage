package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/ports"
	"github.com/google/uuid"
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
	var client domain.Client
	if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdClient, err := h.clientUseCase.CreateClient(r.Context(), &client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdClient)
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
	// Example: /clients/{userID}/{clientID}
	// Split the path to get both UUIDs
	pathParts := len("/clients/")
	idsStr := r.URL.Path[pathParts:]
	ids := make([]string, 0)
	for _, part := range split(idsStr, "/") {
		if part != "" {
			ids = append(ids, part)
		}
	}
	if len(ids) != 2 {
		http.Error(w, "Invalid path, expected /clients/{userID}/{clientID}", http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(ids[0])
	if err != nil {
		http.Error(w, "Invalid userID: "+err.Error(), http.StatusBadRequest)
		return
	}
	clientID, err := uuid.Parse(ids[1])
	if err != nil {
		http.Error(w, "Invalid clientID: "+err.Error(), http.StatusBadRequest)
		return
	}

	client, err := h.clientUseCase.GetClient(r.Context(), userID, clientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(client)
}

// Helper function to split string by separator
func split(s, sep string) []string {
	var result []string
	start := 0
	for i := range s {
		if string(s[i]) == sep {
			result = append(result, s[start:i])
			start = i + 1
		}
	}
	result = append(result, s[start:])
	return result
}

func (h *ClientHandler) UpdateClient(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/clients/"):]
	clientID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid clientID: "+err.Error(), http.StatusBadRequest)
		return
	}

	var client domain.Client
	if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedClient, err := h.clientUseCase.UpdateClient(r.Context(), &client, clientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedClient)
}

func (h *ClientHandler) DeleteClient(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/clients/"):]
	clientID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.clientUseCase.DeleteClient(r.Context(), clientID); err != nil {
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

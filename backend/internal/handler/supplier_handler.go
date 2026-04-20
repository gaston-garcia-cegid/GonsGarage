package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/ports"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/domain"
)

type SupplierHandler struct {
	svc ports.SupplierService
}

func NewSupplierHandler(svc ports.SupplierService) *SupplierHandler {
	return &SupplierHandler{svc: svc}
}

// CreateSupplierRequest body for POST /suppliers.
type CreateSupplierRequest struct {
	Name          string `json:"name"`
	ContactEmail  string `json:"contactEmail"`
	ContactPhone  string `json:"contactPhone"`
	TaxID         string `json:"taxId"`
	Notes         string `json:"notes"`
	IsActive      *bool  `json:"isActive"`
}

// UpdateSupplierRequest body for PUT /suppliers/:id.
type UpdateSupplierRequest struct {
	Name          string `json:"name"`
	ContactEmail  string `json:"contactEmail"`
	ContactPhone  string `json:"contactPhone"`
	TaxID         string `json:"taxId"`
	Notes         string `json:"notes"`
	IsActive      *bool  `json:"isActive"`
}

// SupplierResponse JSON camelCase for API.
type SupplierResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	ContactEmail string `json:"contactEmail"`
	ContactPhone string `json:"contactPhone"`
	TaxID        string `json:"taxId"`
	Notes        string `json:"notes"`
	IsActive     bool   `json:"isActive"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

func supplierToResponse(s *domain.Supplier) SupplierResponse {
	if s == nil {
		return SupplierResponse{}
	}
	return SupplierResponse{
		ID:           s.ID.String(),
		Name:         s.Name,
		ContactEmail: s.ContactEmail,
		ContactPhone: s.ContactPhone,
		TaxID:        s.TaxID,
		Notes:        s.Notes,
		IsActive:     s.IsActive,
		CreatedAt:    s.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:    s.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func writeSupplierServiceError(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, domain.ErrUnauthorizedAccess) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return true
	}
	if errors.Is(err, domain.ErrUserNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return true
	}
	if errors.Is(err, domain.ErrSupplierNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "supplier not found"})
		return true
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	return true
}

// CreateSupplier POST /api/v1/suppliers
// @Summary     Crear proveedor
// @Tags        suppliers
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       body body CreateSupplierRequest true "Datos"
// @Success     201 {object} SupplierResponse
// @Failure     400 {object} SwaggerMessage
// @Failure     401 {object} SwaggerMessage
// @Failure     403 {object} SwaggerMessage
// @Router      /api/v1/suppliers [post]
func (h *SupplierHandler) CreateSupplier(c *gin.Context) {
	uid, err := ContextUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var req CreateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	active := true
	if req.IsActive != nil {
		active = *req.IsActive
	}
	row := &domain.Supplier{
		Name: req.Name, ContactEmail: req.ContactEmail, ContactPhone: req.ContactPhone,
		TaxID: req.TaxID, Notes: req.Notes, IsActive: active,
	}
	out, err := h.svc.Create(c.Request.Context(), row, uid)
	if err != nil {
		writeSupplierServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, supplierToResponse(out))
}

// GetSupplier GET /api/v1/suppliers/:id
// @Summary     Obtener proveedor
// @Tags        suppliers
// @Security    BearerAuth
// @Produce     json
// @Param       id path string true "UUID"
// @Success     200 {object} SupplierResponse
// @Failure     404 {object} SwaggerMessage
// @Router      /api/v1/suppliers/{id} [get]
func (h *SupplierHandler) GetSupplier(c *gin.Context) {
	uid, err := ContextUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	out, err := h.svc.Get(c.Request.Context(), id, uid)
	if err != nil {
		writeSupplierServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, supplierToResponse(out))
}

// ListSuppliers GET /api/v1/suppliers
// @Summary     Listar proveedores
// @Tags        suppliers
// @Security    BearerAuth
// @Produce     json
// @Param       limit query int false "Límite (default 50, max 500)"
// @Param       offset query int false "Offset"
// @Success     200 {object} map[string]interface{} "{items,total}"
// @Router      /api/v1/suppliers [get]
func (h *SupplierHandler) ListSuppliers(c *gin.Context) {
	uid, err := ContextUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	limit, offset := QueryLimitOffset(c, 50, 500)
	list, total, err := h.svc.List(c.Request.Context(), uid, limit, offset)
	if err != nil {
		writeSupplierServiceError(c, err)
		return
	}
	items := make([]SupplierResponse, 0, len(list))
	for _, s := range list {
		items = append(items, supplierToResponse(s))
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "total": total})
}

// UpdateSupplier PUT /api/v1/suppliers/:id
// @Summary     Actualizar proveedor
// @Tags        suppliers
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       id path string true "UUID"
// @Param       body body UpdateSupplierRequest true "Datos"
// @Success     200 {object} SupplierResponse
// @Router      /api/v1/suppliers/{id} [put]
func (h *SupplierHandler) UpdateSupplier(c *gin.Context) {
	uid, err := ContextUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req UpdateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	active := true
	if req.IsActive != nil {
		active = *req.IsActive
	}
	row := &domain.Supplier{
		ID: id, Name: req.Name, ContactEmail: req.ContactEmail, ContactPhone: req.ContactPhone,
		TaxID: req.TaxID, Notes: req.Notes, IsActive: active,
	}
	out, err := h.svc.Update(c.Request.Context(), row, uid)
	if err != nil {
		writeSupplierServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, supplierToResponse(out))
}

// DeleteSupplier DELETE /api/v1/suppliers/:id
// @Summary     Eliminar proveedor (baja lógica)
// @Tags        suppliers
// @Security    BearerAuth
// @Produce     json
// @Param       id path string true "UUID"
// @Success     204
// @Router      /api/v1/suppliers/{id} [delete]
func (h *SupplierHandler) DeleteSupplier(c *gin.Context) {
	uid, err := ContextUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id, uid); err != nil {
		writeSupplierServiceError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

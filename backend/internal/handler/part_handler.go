package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/ports"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PartHandler exposes spare-parts inventory HTTP (manager/admin via middleware + service).
type PartHandler struct {
	svc ports.PartService
}

func NewPartHandler(svc ports.PartService) *PartHandler {
	return &PartHandler{svc: svc}
}

// CreatePartItemRequest body for POST /parts.
type CreatePartItemRequest struct {
	Reference       string   `json:"reference"`
	Brand           string   `json:"brand"`
	Name            string   `json:"name"`
	Barcode         string   `json:"barcode"`
	Quantity        float64  `json:"quantity"`
	UOM             string   `json:"uom"`
	MinimumQuantity *float64 `json:"minimumQuantity"`
}

// UpdatePartItemRequest body for PATCH /parts/:id (full replace of mutable fields).
type UpdatePartItemRequest struct {
	Reference       string   `json:"reference"`
	Brand           string   `json:"brand"`
	Name            string   `json:"name"`
	Barcode         string   `json:"barcode"`
	Quantity        float64  `json:"quantity"`
	UOM             string   `json:"uom"`
	MinimumQuantity *float64 `json:"minimumQuantity"`
}

// PartItemResponse JSON camelCase for API.
type PartItemResponse struct {
	ID               string   `json:"id"`
	Reference        string   `json:"reference"`
	Brand            string   `json:"brand"`
	Name             string   `json:"name"`
	Barcode          string   `json:"barcode"`
	Quantity         float64  `json:"quantity"`
	UOM              string   `json:"uom"`
	MinimumQuantity  *float64 `json:"minimumQuantity,omitempty"`
	CreatedAtRFC3339 string   `json:"createdAt"`
	UpdatedAtRFC3339 string   `json:"updatedAt"`
	DeletedAtRFC3339 *string  `json:"deletedAt,omitempty"`
}

func partItemToResponse(p *domain.PartItem) PartItemResponse {
	if p == nil {
		return PartItemResponse{}
	}
	var del *string
	if p.DeletedAt != nil {
		s := p.DeletedAt.UTC().Format(time.RFC3339)
		del = &s
	}
	return PartItemResponse{
		ID:               p.ID.String(),
		Reference:        p.Reference,
		Brand:            p.Brand,
		Name:             p.Name,
		Barcode:          p.Barcode,
		Quantity:         p.Quantity,
		UOM:              p.UOM,
		MinimumQuantity:  p.MinimumQuantity,
		CreatedAtRFC3339: p.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAtRFC3339: p.UpdatedAt.UTC().Format(time.RFC3339),
		DeletedAtRFC3339: del,
	}
}

func writePartServiceError(c *gin.Context, err error) bool {
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
	if errors.Is(err, domain.ErrPartItemNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "part item not found"})
		return true
	}
	if errors.Is(err, domain.ErrPartItemDuplicateBarcode) {
		c.JSON(http.StatusConflict, gin.H{"error": "duplicate barcode"})
		return true
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	return true
}

// CreatePartItem POST /api/v1/parts
func (h *PartHandler) CreatePartItem(c *gin.Context) {
	uid, err := ContextUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var req CreatePartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	row := &domain.PartItem{
		Reference: req.Reference, Brand: req.Brand, Name: req.Name,
		Barcode: req.Barcode, Quantity: req.Quantity, UOM: req.UOM,
		MinimumQuantity: req.MinimumQuantity,
	}
	out, err := h.svc.Create(c.Request.Context(), row, uid)
	if err != nil {
		writePartServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, partItemToResponse(out))
}

// GetPartItem GET /api/v1/parts/:id
func (h *PartHandler) GetPartItem(c *gin.Context) {
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
		writePartServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, partItemToResponse(out))
}

// ListParts GET /api/v1/parts
func (h *PartHandler) ListParts(c *gin.Context) {
	uid, err := ContextUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	limit, offset := QueryLimitOffset(c, 50, 500)
	f := ports.PartItemListFilters{Limit: limit, Offset: offset}
	if bc := c.Query("barcode"); bc != "" {
		f.Barcode = &bc
	}
	if s := c.Query("search"); s != "" {
		f.Search = &s
	}
	list, total, err := h.svc.List(c.Request.Context(), f, uid)
	if err != nil {
		writePartServiceError(c, err)
		return
	}
	items := make([]PartItemResponse, 0, len(list))
	for _, p := range list {
		items = append(items, partItemToResponse(p))
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "total": total})
}

// UpdatePartItem PATCH /api/v1/parts/:id
func (h *PartHandler) UpdatePartItem(c *gin.Context) {
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
	var req UpdatePartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	row := &domain.PartItem{
		ID: id, Reference: req.Reference, Brand: req.Brand, Name: req.Name,
		Barcode: req.Barcode, Quantity: req.Quantity, UOM: req.UOM,
		MinimumQuantity: req.MinimumQuantity,
	}
	out, err := h.svc.Update(c.Request.Context(), row, uid)
	if err != nil {
		writePartServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, partItemToResponse(out))
}

// DeletePartItem DELETE /api/v1/parts/:id
func (h *PartHandler) DeletePartItem(c *gin.Context) {
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
		writePartServiceError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

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

type BillingDocumentHandler struct {
	svc ports.BillingDocumentService
}

func NewBillingDocumentHandler(svc ports.BillingDocumentService) *BillingDocumentHandler {
	return &BillingDocumentHandler{svc: svc}
}

// CreateBillingDocumentRequest body for POST /billing-documents.
type CreateBillingDocumentRequest struct {
	Kind       string  `json:"kind"`
	Title      string  `json:"title"`
	Amount     float64 `json:"amount"`
	CustomerID *string `json:"customerId,omitempty"`
	Reference  string  `json:"reference"`
	Notes      string  `json:"notes"`
}

// UpdateBillingDocumentRequest body for PUT /billing-documents/:id.
type UpdateBillingDocumentRequest struct {
	Kind       string  `json:"kind"`
	Title      string  `json:"title"`
	Amount     float64 `json:"amount"`
	CustomerID *string `json:"customerId,omitempty"`
	Reference  string  `json:"reference"`
	Notes      string  `json:"notes"`
}

// BillingDocumentResponse JSON camelCase.
type BillingDocumentResponse struct {
	ID          string  `json:"id"`
	Kind        string  `json:"kind"`
	Title       string  `json:"title"`
	Amount      float64 `json:"amount"`
	CustomerID  *string `json:"customerId,omitempty"`
	Reference   string  `json:"reference"`
	Notes       string  `json:"notes"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
}

func billingDocToResponse(d *domain.BillingDocument) BillingDocumentResponse {
	if d == nil {
		return BillingDocumentResponse{}
	}
	var cust *string
	if d.CustomerID != nil {
		s := d.CustomerID.String()
		cust = &s
	}
	return BillingDocumentResponse{
		ID:         d.ID.String(),
		Kind:       string(d.Kind),
		Title:      d.Title,
		Amount:     d.Amount,
		CustomerID: cust,
		Reference:  d.Reference,
		Notes:      d.Notes,
		CreatedAt:  d.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:  d.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func parseOptionalUUIDPtr(s *string) (*uuid.UUID, error) {
	if s == nil || *s == "" {
		return nil, nil
	}
	id, err := uuid.Parse(*s)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func writeBillingDocServiceError(c *gin.Context, err error) bool {
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
	if errors.Is(err, domain.ErrBillingDocumentNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "billing document not found"})
		return true
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	return true
}

// CreateBillingDocument POST /api/v1/billing-documents
// @Summary     Crear documento de facturación emitido
// @Tags        billing-documents
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       body body CreateBillingDocumentRequest true "Datos"
// @Success     201 {object} BillingDocumentResponse
// @Router      /api/v1/billing-documents [post]
func (h *BillingDocumentHandler) CreateBillingDocument(c *gin.Context) {
	uid, err := ContextUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var req CreateBillingDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	custID, err := parseOptionalUUIDPtr(req.CustomerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid customerId"})
		return
	}
	doc := &domain.BillingDocument{
		Kind:       domain.BillingDocumentKind(req.Kind),
		Title:      req.Title,
		Amount:     req.Amount,
		CustomerID: custID,
		Reference:  req.Reference,
		Notes:      req.Notes,
	}
	out, err := h.svc.Create(c.Request.Context(), doc, uid)
	if err != nil {
		writeBillingDocServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, billingDocToResponse(out))
}

// GetBillingDocument GET /api/v1/billing-documents/:id
// @Summary     Obtener documento
// @Tags        billing-documents
// @Security    BearerAuth
// @Produce     json
// @Param       id path string true "UUID"
// @Success     200 {object} BillingDocumentResponse
// @Router      /api/v1/billing-documents/{id} [get]
func (h *BillingDocumentHandler) GetBillingDocument(c *gin.Context) {
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
		writeBillingDocServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, billingDocToResponse(out))
}

// ListBillingDocuments GET /api/v1/billing-documents
// @Summary     Listar documentos
// @Tags        billing-documents
// @Security    BearerAuth
// @Produce     json
// @Success     200 {object} map[string]interface{} "{items,total}"
// @Router      /api/v1/billing-documents [get]
func (h *BillingDocumentHandler) ListBillingDocuments(c *gin.Context) {
	uid, err := ContextUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	limit, offset := QueryLimitOffset(c, 50, 500)
	list, total, err := h.svc.List(c.Request.Context(), uid, limit, offset)
	if err != nil {
		writeBillingDocServiceError(c, err)
		return
	}
	items := make([]BillingDocumentResponse, 0, len(list))
	for _, d := range list {
		items = append(items, billingDocToResponse(d))
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "total": total})
}

// UpdateBillingDocument PUT /api/v1/billing-documents/:id
// @Summary     Actualizar documento
// @Tags        billing-documents
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       id path string true "UUID"
// @Param       body body UpdateBillingDocumentRequest true "Datos"
// @Success     200 {object} BillingDocumentResponse
// @Router      /api/v1/billing-documents/{id} [put]
func (h *BillingDocumentHandler) UpdateBillingDocument(c *gin.Context) {
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
	var req UpdateBillingDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	custID, err := parseOptionalUUIDPtr(req.CustomerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid customerId"})
		return
	}
	doc := &domain.BillingDocument{
		ID:         id,
		Kind:       domain.BillingDocumentKind(req.Kind),
		Title:      req.Title,
		Amount:     req.Amount,
		CustomerID: custID,
		Reference:  req.Reference,
		Notes:      req.Notes,
	}
	out, err := h.svc.Update(c.Request.Context(), doc, uid)
	if err != nil {
		writeBillingDocServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, billingDocToResponse(out))
}

// DeleteBillingDocument DELETE /api/v1/billing-documents/:id
// @Summary     Eliminar documento (baja lógica)
// @Tags        billing-documents
// @Security    BearerAuth
// @Param       id path string true "UUID"
// @Success     204
// @Router      /api/v1/billing-documents/{id} [delete]
func (h *BillingDocumentHandler) DeleteBillingDocument(c *gin.Context) {
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
		writeBillingDocServiceError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

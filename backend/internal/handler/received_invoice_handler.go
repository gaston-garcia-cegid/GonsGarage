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

type ReceivedInvoiceHandler struct {
	svc ports.ReceivedInvoiceService
}

func NewReceivedInvoiceHandler(svc ports.ReceivedInvoiceService) *ReceivedInvoiceHandler {
	return &ReceivedInvoiceHandler{svc: svc}
}

// CreateReceivedInvoiceRequest body for POST /received-invoices.
type CreateReceivedInvoiceRequest struct {
	SupplierID  *string `json:"supplierId,omitempty"`
	VendorName  string  `json:"vendorName"`
	Category    string  `json:"category"`
	Amount      float64 `json:"amount"`
	InvoiceDate string  `json:"invoiceDate"`
	Notes       string  `json:"notes"`
}

// UpdateReceivedInvoiceRequest body for PUT /received-invoices/:id.
type UpdateReceivedInvoiceRequest struct {
	SupplierID  *string `json:"supplierId,omitempty"`
	VendorName  string  `json:"vendorName"`
	Category    string  `json:"category"`
	Amount      float64 `json:"amount"`
	InvoiceDate string  `json:"invoiceDate"`
	Notes       string  `json:"notes"`
}

// ReceivedInvoiceResponse JSON camelCase.
type ReceivedInvoiceResponse struct {
	ID          string  `json:"id"`
	SupplierID  *string `json:"supplierId,omitempty"`
	VendorName  string  `json:"vendorName"`
	Category    string  `json:"category"`
	Amount      float64 `json:"amount"`
	InvoiceDate string  `json:"invoiceDate"`
	Notes       string  `json:"notes"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
}

func receivedInvoiceToResponse(r *domain.ReceivedInvoice) ReceivedInvoiceResponse {
	if r == nil {
		return ReceivedInvoiceResponse{}
	}
	var sid *string
	if r.SupplierID != nil {
		s := r.SupplierID.String()
		sid = &s
	}
	return ReceivedInvoiceResponse{
		ID:          r.ID.String(),
		SupplierID:  sid,
		VendorName:  r.VendorName,
		Category:    r.Category,
		Amount:      r.Amount,
		InvoiceDate: r.InvoiceDate.UTC().Format(time.RFC3339),
		Notes:       r.Notes,
		CreatedAt:   r.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:   r.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func parseInvoiceDate(raw string) (time.Time, error) {
	if raw == "" {
		return time.Time{}, errors.New("invoice date required")
	}
	if t, err := time.Parse(time.RFC3339, raw); err == nil {
		return t, nil
	}
	return time.ParseInLocation("2006-01-02", raw, time.UTC)
}

func writeReceivedInvoiceServiceError(c *gin.Context, err error) bool {
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
	if errors.Is(err, domain.ErrReceivedInvoiceNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "received invoice not found"})
		return true
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	return true
}

// CreateReceivedInvoice POST /api/v1/received-invoices
// @Summary     Registrar factura recibida
// @Tags        received-invoices
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       body body CreateReceivedInvoiceRequest true "Datos"
// @Success     201 {object} ReceivedInvoiceResponse
// @Router      /api/v1/received-invoices [post]
func (h *ReceivedInvoiceHandler) CreateReceivedInvoice(c *gin.Context) {
	uid, err := ContextUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var req CreateReceivedInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	invDate, err := parseInvoiceDate(req.InvoiceDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid invoiceDate"})
		return
	}
	supID, err := parseOptionalUUIDPtr(req.SupplierID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid supplierId"})
		return
	}
	inv := &domain.ReceivedInvoice{
		SupplierID:  supID,
		VendorName:  req.VendorName,
		Category:    req.Category,
		Amount:      req.Amount,
		InvoiceDate: invDate,
		Notes:       req.Notes,
	}
	out, err := h.svc.Create(c.Request.Context(), inv, uid)
	if err != nil {
		writeReceivedInvoiceServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, receivedInvoiceToResponse(out))
}

// GetReceivedInvoice GET /api/v1/received-invoices/:id
// @Summary     Obtener factura recibida
// @Tags        received-invoices
// @Security    BearerAuth
// @Produce     json
// @Param       id path string true "UUID"
// @Success     200 {object} ReceivedInvoiceResponse
// @Router      /api/v1/received-invoices/{id} [get]
func (h *ReceivedInvoiceHandler) GetReceivedInvoice(c *gin.Context) {
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
		writeReceivedInvoiceServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, receivedInvoiceToResponse(out))
}

// ListReceivedInvoices GET /api/v1/received-invoices
// @Summary     Listar facturas recibidas
// @Tags        received-invoices
// @Security    BearerAuth
// @Produce     json
// @Success     200 {object} map[string]interface{} "{items,total}"
// @Router      /api/v1/received-invoices [get]
func (h *ReceivedInvoiceHandler) ListReceivedInvoices(c *gin.Context) {
	uid, err := ContextUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	limit, offset := QueryLimitOffset(c, 50, 500)
	list, total, err := h.svc.List(c.Request.Context(), uid, limit, offset)
	if err != nil {
		writeReceivedInvoiceServiceError(c, err)
		return
	}
	items := make([]ReceivedInvoiceResponse, 0, len(list))
	for _, x := range list {
		items = append(items, receivedInvoiceToResponse(x))
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "total": total})
}

// UpdateReceivedInvoice PUT /api/v1/received-invoices/:id
// @Summary     Actualizar factura recibida
// @Tags        received-invoices
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       id path string true "UUID"
// @Param       body body UpdateReceivedInvoiceRequest true "Datos"
// @Success     200 {object} ReceivedInvoiceResponse
// @Router      /api/v1/received-invoices/{id} [put]
func (h *ReceivedInvoiceHandler) UpdateReceivedInvoice(c *gin.Context) {
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
	var req UpdateReceivedInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	invDate, err := parseInvoiceDate(req.InvoiceDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid invoiceDate"})
		return
	}
	supID, err := parseOptionalUUIDPtr(req.SupplierID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid supplierId"})
		return
	}
	inv := &domain.ReceivedInvoice{
		ID:          id,
		SupplierID:  supID,
		VendorName:  req.VendorName,
		Category:    req.Category,
		Amount:      req.Amount,
		InvoiceDate: invDate,
		Notes:       req.Notes,
	}
	out, err := h.svc.Update(c.Request.Context(), inv, uid)
	if err != nil {
		writeReceivedInvoiceServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, receivedInvoiceToResponse(out))
}

// DeleteReceivedInvoice DELETE /api/v1/received-invoices/:id
// @Summary     Eliminar factura recibida (baja lógica)
// @Tags        received-invoices
// @Security    BearerAuth
// @Param       id path string true "UUID"
// @Success     204
// @Router      /api/v1/received-invoices/{id} [delete]
func (h *ReceivedInvoiceHandler) DeleteReceivedInvoice(c *gin.Context) {
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
		writeReceivedInvoiceServiceError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

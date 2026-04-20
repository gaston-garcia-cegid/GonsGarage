package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestReceivedInvoice_Validate(t *testing.T) {
	t.Parallel()
	sid := uuid.New()
	r := &ReceivedInvoice{
		Amount:      10,
		InvoiceDate: time.Now().UTC(),
		Category:    "parts",
		SupplierID:  &sid,
	}
	assert.NoError(t, r.Validate())
	r2 := &ReceivedInvoice{Amount: 0, InvoiceDate: time.Now().UTC(), Category: "x"}
	assert.Error(t, r2.Validate())
	r3 := &ReceivedInvoice{Amount: 1, InvoiceDate: time.Time{}, Category: "x"}
	assert.Error(t, r3.Validate())
}

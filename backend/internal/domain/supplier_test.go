package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSupplier_Validate(t *testing.T) {
	t.Parallel()
	s := &Supplier{Name: "ACME Parts", ContactEmail: "a@acme.pt"}
	assert.NoError(t, s.Validate())
	assert.Error(t, (&Supplier{Name: ""}).Validate())
}

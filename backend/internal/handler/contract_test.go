package handler_test

import (
	"reflect"
	"testing"

	legacy "github.com/gaston-garcia-cegid/gonsgarage/internal/adapters/http/handlers"
	apihandler "github.com/gaston-garcia-cegid/gonsgarage/internal/handler"
	"github.com/stretchr/testify/require"
)

// Fase 1 (template): internal/handler es el borde HTTP público; las implementaciones
// siguen en adapters/http/handlers hasta migración completa.
func TestHandler_constructorsPointToLegacyImplementations(t *testing.T) {
	cases := []struct {
		name string
		fac  any
		leg  any
	}{
		{"NewAuthHandler", apihandler.NewAuthHandler, legacy.NewAuthHandler},
		{"NewEmployeeHandler", apihandler.NewEmployeeHandler, legacy.NewEmployeeHandler},
		{"NewCarHandler", apihandler.NewCarHandler, legacy.NewCarHandler},
		{"NewAppointmentHandler", apihandler.NewAppointmentHandler, legacy.NewAppointmentHandler},
		{"NewRepairHandler", apihandler.NewRepairHandler, legacy.NewRepairHandler},
		{"NewUserHandler", apihandler.NewUserHandler, legacy.NewUserHandler},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t,
				reflect.ValueOf(tc.leg).Pointer(),
				reflect.ValueOf(tc.fac).Pointer(),
				"facade must delegate to legacy constructor without duplicating logic",
			)
		})
	}
}

package service_test

import (
	"reflect"
	"testing"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/services/appointment"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/services/auth"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/services/car"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/services/employee"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/services/repair"
	appsvc "github.com/gaston-garcia-cegid/gonsgarage/internal/service"
	"github.com/stretchr/testify/require"
)

func TestService_constructorsPointToCore(t *testing.T) {
	cases := []struct {
		name string
		fac  any
		leg  any
	}{
		{"NewAuthService", appsvc.NewAuthService, auth.NewAuthService},
		{"NewEmployeeService", appsvc.NewEmployeeService, employee.NewEmployeeService},
		{"NewCarService", appsvc.NewCarService, car.NewCarService},
		{"NewAppointmentService", appsvc.NewAppointmentService, appointment.NewAppointmentService},
		{"NewRepairService", appsvc.NewRepairService, repair.NewRepairService},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t,
				reflect.ValueOf(tc.leg).Pointer(),
				reflect.ValueOf(tc.fac).Pointer(),
			)
		})
	}
}

package repository_test

import (
	"reflect"
	"testing"

	postgres "github.com/gaston-garcia-cegid/gonsgarage/internal/adapters/repository/postgres"
	apprepo "github.com/gaston-garcia-cegid/gonsgarage/internal/repository"
	"github.com/stretchr/testify/require"
)

func TestRepository_constructorsPointToPostgresAdapters(t *testing.T) {
	cases := []struct {
		name string
		fac  any
		leg  any
	}{
		{"NewPostgresUserRepository", apprepo.NewPostgresUserRepository, postgres.NewPostgresUserRepository},
		{"NewPostgresEmployeeRepository", apprepo.NewPostgresEmployeeRepository, postgres.NewPostgresEmployeeRepository},
		{"NewPostgresCarRepository", apprepo.NewPostgresCarRepository, postgres.NewPostgresCarRepository},
		{"NewPostgresRepairRepository", apprepo.NewPostgresRepairRepository, postgres.NewPostgresRepairRepository},
		{"NewPostgresAppointmentRepository", apprepo.NewPostgresAppointmentRepository, postgres.NewPostgresAppointmentRepository},
		{"NewWorkshopRepository", apprepo.NewWorkshopRepository, postgres.NewWorkshopRepository},
		{"NewAccountingRepository", apprepo.NewAccountingRepository, postgres.NewAccountingRepository},
		{"NewNotificationRepository", apprepo.NewNotificationRepository, postgres.NewNotificationRepository},
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

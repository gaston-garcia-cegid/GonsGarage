package domain_test

import (
	"testing"

	cd "github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
	appdomain "github.com/gaston-garcia-cegid/gonsgarage/internal/domain"
	"github.com/stretchr/testify/require"
)

// Composition root (cmd/api) needs these models for AutoMigrate; aliases must match core domain memory layout.
func TestDomain_modelAliasesMatchCore(t *testing.T) {
	var u appdomain.User
	var eu appdomain.Employee
	var c appdomain.Car
	var r appdomain.Repair
	var a appdomain.Appointment

	var cu cd.User = u
	var ce cd.Employee = eu
	var cc cd.Car = c
	var cr cd.Repair = r
	var ca cd.Appointment = a

	_ = cu
	_ = ce
	_ = cc
	_ = cr
	_ = ca

	require.Equal(t, cd.RoleClient, appdomain.RoleClient)
	require.Equal(t, cd.AppointmentStatusScheduled, appdomain.AppointmentStatusScheduled)
	require.Equal(t, cd.RepairStatusPending, appdomain.RepairStatusPending)
}

// Package service is the application / use-case layer entry (template §1).
// Phase 1: function aliases to internal/core/services/*.
package service

import (
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/services/appointment"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/services/auth"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/services/car"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/services/employee"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/services/repair"
)

var (
	NewAuthService        = auth.NewAuthService
	NewEmployeeService    = employee.NewEmployeeService
	NewCarService         = car.NewCarService
	NewAppointmentService = appointment.NewAppointmentService
	NewRepairService      = repair.NewRepairService
)

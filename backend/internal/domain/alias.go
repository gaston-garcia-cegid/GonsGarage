// Package domain re-exports entities used by the composition root (template §1).
// Phase 1: type aliases to internal/core/domain for GORM AutoMigrate and routing; services/handlers still import core/domain internally until a later phase.
package domain

import cd "github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"

type (
	User              = cd.User
	Employee          = cd.Employee
	Car               = cd.Car
	Repair            = cd.Repair
	Appointment       = cd.Appointment
	AppointmentStatus = cd.AppointmentStatus
	RepairStatus      = cd.RepairStatus
)

const (
	RoleAdmin    = cd.RoleAdmin
	RoleManager  = cd.RoleManager
	RoleEmployee = cd.RoleEmployee
	RoleClient   = cd.RoleClient

	AppointmentStatusScheduled = cd.AppointmentStatusScheduled
	AppointmentStatusConfirmed = cd.AppointmentStatusConfirmed
	AppointmentStatusCompleted = cd.AppointmentStatusCompleted
	AppointmentStatusCancelled = cd.AppointmentStatusCancelled

	RepairStatusPending    = cd.RepairStatusPending
	RepairStatusInProgress = cd.RepairStatusInProgress
	RepairStatusCompleted  = cd.RepairStatusCompleted
	RepairStatusCancelled  = cd.RepairStatusCancelled
)

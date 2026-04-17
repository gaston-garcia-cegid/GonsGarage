// Package handler is the HTTP transport entrypoint (template §1, gonsgarage-rules/01-architecture.md).
// Phase 1: aliases to internal/adapters/http/handlers — callers use this path; implementation stays in adapters until a vertical migration.
package handler

import legacy "github.com/gaston-garcia-cegid/gonsgarage/internal/adapters/http/handlers"

type (
	AuthHandler              = legacy.AuthHandler
	EmployeeHandler          = legacy.EmployeeHandler
	CarHandler               = legacy.CarHandler
	AppointmentHandler       = legacy.AppointmentHandler
	RepairHandler            = legacy.RepairHandler
	UserHandler              = legacy.UserHandler
	LoginRequest             = legacy.LoginRequest
	CreateCarRequest         = legacy.CreateCarRequest
	UpdateCarRequest         = legacy.UpdateCarRequest
	CarResponse              = legacy.CarResponse
	CreateAppointmentRequest = legacy.CreateAppointmentRequest
	UpdateAppointmentRequest = legacy.UpdateAppointmentRequest
	AppointmentResponse      = legacy.AppointmentResponse
	CreateRepairRequest      = legacy.CreateRepairRequest
	UpdateRepairRequest      = legacy.UpdateRepairRequest
	RepairResponse           = legacy.RepairResponse
	SwaggerLoginOK           = legacy.SwaggerLoginOK
	SwaggerRegisterUser      = legacy.SwaggerRegisterUser
	SwaggerRegisterOK        = legacy.SwaggerRegisterOK
	SwaggerMeUser            = legacy.SwaggerMeUser
	SwaggerMeOK              = legacy.SwaggerMeOK
	SwaggerMessage           = legacy.SwaggerMessage
)

var (
	NewAuthHandler        = legacy.NewAuthHandler
	NewEmployeeHandler    = legacy.NewEmployeeHandler
	NewCarHandler         = legacy.NewCarHandler
	NewAppointmentHandler = legacy.NewAppointmentHandler
	NewRepairHandler      = legacy.NewRepairHandler
	NewUserHandler        = legacy.NewUserHandler
)

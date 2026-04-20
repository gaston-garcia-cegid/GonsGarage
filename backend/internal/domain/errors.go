package domain

import "errors"

// Other error definitions...

var ErrUnauthorizedAccess = errors.New("unauthorized access")
var ErrCarNotFound = errors.New("car not found")
var ErrCarAlreadyExists = errors.New("car with the given license plate already exists")
var ErrInvalidCarData = errors.New("invalid car data")
var ErrRepairNotFound = errors.New("repair not found")
var ErrInvalidRepairData = errors.New("invalid repair data")
var ErrClientNotFound = errors.New("client not found")
var ErrInvalidClientData = errors.New("invalid client data")
var ErrAppointmentNotFound = errors.New("appointment not found")
var ErrAppointmentAlreadyExists = errors.New("appointment with the given ID already exists")
var ErrInvalidAppointmentData = errors.New("invalid appointment data")
var ErrAppointmentOutsideBusinessHours = errors.New("appointment outside business hours")
var ErrAppointmentDailyCapReached = errors.New("maximum appointments per day reached")
var ErrWorkshopNotFound = errors.New("workshop not found")
var ErrInvalidWorkshopData = errors.New("invalid workshop data")
var ErrAccountingEntryNotFound = errors.New("accounting entry not found")
var ErrInvalidAccountingEntryData = errors.New("invalid accounting entry data")
var ErrAccountingEntryAlreadyExists = errors.New("accounting entry with the given ID already exists")
var ErrInvoiceNotFound = errors.New("invoice not found")
var ErrSupplierNotFound = errors.New("supplier not found")
var ErrReceivedInvoiceNotFound = errors.New("received invoice not found")
var ErrBillingDocumentNotFound = errors.New("billing document not found")

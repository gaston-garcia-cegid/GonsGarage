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

package domain

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrDniAlreadyExist   = errors.New("dni already exists")
	ErrEmailAlreadyExist = errors.New("email already exists")
	ErrPhoneAlreadyExist = errors.New("phone already exists")
	ErrTenantNotFound    = errors.New("tenant not found in context")
	ErrInvalidTenantID   = errors.New("invalid tenant id")

	ErrDepartmentNotFound   = errors.New("department not found")
	ErrDepartmentCodeExists = errors.New("department code already exists")
	ErrDepartmentNameExists = errors.New("department name already exists")
	ErrPositionNotFound     = errors.New("position not found")
	ErrPositionCodeExists   = errors.New("position code already exists")
)

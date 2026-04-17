package domain

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrDniAlreadyExist   = errors.New("dni already exists")
	ErrEmailAlreadyExist = errors.New("email already exists")
	ErrPhoneAlreadyExist = errors.New("phone already exists")
	ErrTenantNotFound    = errors.New("tenant not found in context")
	ErrInvalidTenantID   = errors.New("invalid tenant id")
)

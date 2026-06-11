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

	ErrPositionNotFound      = errors.New("position not found")
	ErrPositionAlreadyExists = errors.New("position already exists")
	ErrPositionNameExists    = errors.New("position name already exists")

	ErrEmployeeNotFound       = errors.New("employee not found")
	ErrEmployeeAlreadyExists  = errors.New("employee already exists")
	ErrEmployeeNameExists     = errors.New("employee name already exists")

	ErrContractNotFound          = errors.New("contract not found")
	ErrContractTypeNotFound      = errors.New("contract type not found")
	ErrContractTypeAlreadyExists = errors.New("contract type already exists")

	ErrCountryParamNotFound      = errors.New("country params not found")
	ErrCountryParamAlreadyExists = errors.New("country params already exist for this country")

	ErrPayrollRecordNotFound = errors.New("payroll record not found")

	ErrInvalidCountryCode = errors.New("invalid country code")
	ErrInvalidCurrency    = errors.New("invalid currency")
	ErrInvalidPeriod      = errors.New("invalid period: start must be before end")
	ErrCalculatorNotFound = errors.New("calculator not found for country")
)

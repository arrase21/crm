package dto

import (
	"fmt"
	"strings"
	"time"

	"github.com/arrase21/crm/internal/domain"
)

type CreateUserRequest struct {
	FirstName string `json:"first_name" validate:"required,min=2,max=30,alpha"`
	LastName  string `json:"last_name" validate:"required,min=2,max=40,alpha"`
	Dni       string `json:"dni" validate:"required,len=8,numeric"`
	Gender    string `json:"gender" validate:"required,oneof=M F"`
	Phone     string `json:"phone" validate:"required,e164"`
	Email     string `json:"email" validate:"required,email,max=50"`
	BirthDay  string `json:"birth_day" validate:"required,datetime=2006-01-02"`
}

type UpdateUserRequest struct {
	ID        uint    `json:"id" validate:"required,min=1"`
	FirstName *string `json:"first_name,omitempty" validate:"omitempty,min=2,max=30,alpha"`
	LastName  *string `json:"last_name,omitempty" validate:"omitempty,min=2,max=40,alpha"`
	Dni       *string `json:"dni,omitempty" validate:"omitempty,len=8,numeric"`
	Gender    *string `json:"gender,omitempty" validate:"omitempty,oneof=M F"`
	Phone     *string `json:"phone,omitempty" validate:"omitempty,e164"`
	Email     *string `json:"email,omitempty" validate:"omitempty,email,max=50"`
	BirthDay  *string `json:"birth_day,omitempty" validate:"omitempty,datetime=2006-01-02"`
}

// type CreateRoleRequest struct {
// 	Name        string `json:"name" validate:"required,min=2,max=50,alphanum"`
// 	Description string `json:"description" validate:"max=255"`
// 	IsActive    *bool  `json:"is_active,omitempty"`
// }
//
// type UpdateRoleRequest struct {
// 	ID          uint    `json:"id" validate:"required,min=1"`
// 	Name        *string `json:"name,omitempty" validate:"omitempty,min=2,max=50,alphanum"`
// 	Description *string `json:"description,omitempty" validate:"omitempty,max=255"`
// 	IsActive    *bool   `json:"is_active,omitempty"`
// }
//
// type AssignRoleRequest struct {
// 	UserID uint `json:"user_id" validate:"required,min=1"`
// 	RoleID uint `json:"role_id" validate:"required,min=1"`
// }

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Tag     string `json:"tag"`
}

type ErrorResponse struct {
	Error       string            `json:"error"`
	Message     string            `json:"message"`
	Code        string            `json:"code"`
	Validations []ValidationError `json:"validations,omitempty"`
	Timestamp   time.Time         `json:"timestamp"`
}

func (r *CreateUserRequest) ToDomain(tenantID uint) (*domain.User, error) {
	birthDate, err := time.Parse("2006-01-02", r.BirthDay)
	if err != nil {
		return nil, fmt.Errorf("invalid birth_date format: %w", err)
	}

	return &domain.User{
		TenantID:  tenantID,
		FirstName: strings.TrimSpace(r.FirstName),
		LastName:  strings.TrimSpace(r.LastName),
		Dni:       strings.TrimSpace(r.Dni),
		Gender:    strings.ToUpper(r.Gender),
		Phone:     strings.TrimSpace(r.Phone),
		Email:     strings.ToLower(strings.TrimSpace(r.Email)),
		BirthDay:  birthDate,
	}, nil
}

func (r *CreateUserRequest) Sanitize() {
	r.FirstName = strings.TrimSpace(r.FirstName)
	r.LastName = strings.TrimSpace(r.LastName)
	r.Dni = strings.TrimSpace(r.Dni)
	r.Phone = strings.TrimSpace(r.Phone)
	r.Email = strings.ToLower(strings.TrimSpace(r.Email))
	r.Gender = strings.ToUpper(r.Gender)
}

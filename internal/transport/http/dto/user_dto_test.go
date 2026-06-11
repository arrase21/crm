package dto

import (
	"strings"
	"testing"
	"time"

	"github.com/arrase21/crm/internal/domain"
)

func TestCreateUserRequest_ToDomain(t *testing.T) {
	t.Run("valid request - success", func(t *testing.T) {
		req := &CreateUserRequest{
			FirstName: "  Juan  ",
			LastName:  "  Pérez  ",
			Dni:       " 12345678 ",
			Gender:    "m",
			Phone:     " +51999999999 ",
			Email:     "  JUAN@TEST.COM  ",
			BirthDay:  "1990-01-01",
		}
		user, err := req.ToDomain(1)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if user.TenantID != 1 {
			t.Errorf("expected tenant 1, got %d", user.TenantID)
		}
		if user.FirstName != "Juan" {
			t.Errorf("expected 'Juan', got '%s'", user.FirstName)
		}
		if user.LastName != "Pérez" {
			t.Errorf("expected 'Pérez', got '%s'", user.LastName)
		}
		if user.Dni != "12345678" {
			t.Errorf("expected '12345678', got '%s'", user.Dni)
		}
		if user.Gender != "M" {
			t.Errorf("expected 'M', got '%s'", user.Gender)
		}
		if user.Phone != "+51999999999" {
			t.Errorf("expected '+51999999999', got '%s'", user.Phone)
		}
		if user.Email != "juan@test.com" {
			t.Errorf("expected 'juan@test.com', got '%s'", user.Email)
		}
		if user.BirthDay.Year() != 1990 || user.BirthDay.Month() != 1 || user.BirthDay.Day() != 1 {
			t.Errorf("expected 1990-01-01, got %s", user.BirthDay.Format("2006-01-02"))
		}
	})

	t.Run("invalid birth_date format - should fail", func(t *testing.T) {
		req := &CreateUserRequest{
			FirstName: "Juan",
			LastName:  "Pérez",
			Dni:       "12345678",
			Gender:    "M",
			Phone:     "+51999999999",
			Email:     "juan@test.com",
			BirthDay:  "01-01-1990",
		}
		_, err := req.ToDomain(1)
		if err == nil {
			t.Fatal("expected error for invalid date format")
		}
	})
}

func TestCreateUserRequest_Sanitize(t *testing.T) {
	req := &CreateUserRequest{
		FirstName: "  Juan  ",
		LastName:  "  Pérez  ",
		Dni:       " 12345678 ",
		Gender:    "m",
		Phone:     " +51999999999 ",
		Email:     "  JUAN@TEST.COM  ",
		BirthDay:  "1990-01-01",
	}
	req.Sanitize()

	if req.FirstName != "Juan" {
		t.Errorf("expected 'Juan', got '%s'", req.FirstName)
	}
	if req.LastName != "Pérez" {
		t.Errorf("expected 'Pérez', got '%s'", req.LastName)
	}
	if req.Dni != "12345678" {
		t.Errorf("expected '12345678', got '%s'", req.Dni)
	}
	if req.Gender != "M" {
		t.Errorf("expected 'M', got '%s'", req.Gender)
	}
	if req.Phone != "+51999999999" {
		t.Errorf("expected '+51999999999', got '%s'", req.Phone)
	}
	if req.Email != "juan@test.com" {
		t.Errorf("expected 'juan@test.com', got '%s'", req.Email)
	}
}

func TestValidationError(t *testing.T) {
	v := ValidationError{
		Field:   "first_name",
		Message: "first name is required",
		Tag:     "required",
	}
	if v.Field != "first_name" {
		t.Errorf("expected 'first_name', got '%s'", v.Field)
	}
	if v.Message != "first name is required" {
		t.Errorf("expected 'first name is required', got '%s'", v.Message)
	}
	if v.Tag != "required" {
		t.Errorf("expected 'required', got '%s'", v.Tag)
	}
}

func TestErrorResponse(t *testing.T) {
	now := time.Now()
	errResp := ErrorResponse{
		Error:   "validation_error",
		Message: "Invalid input",
		Code:    "ERR_001",
		Validations: []ValidationError{
			{Field: "email", Message: "email is required", Tag: "required"},
		},
		Timestamp: now,
	}
	if errResp.Error != "validation_error" {
		t.Errorf("expected 'validation_error', got '%s'", errResp.Error)
	}
	if errResp.Code != "ERR_001" {
		t.Errorf("expected 'ERR_001', got '%s'", errResp.Code)
	}
	if len(errResp.Validations) != 1 {
		t.Errorf("expected 1 validation, got %d", len(errResp.Validations))
	}
}

func TestCreateUserRequest_ToDomain_Integration(t *testing.T) {
	req := &CreateUserRequest{
		FirstName: "Ana",
		LastName:  "López",
		Dni:       "87654321",
		Gender:    "F",
		Phone:     "+51988888888",
		Email:     "  ANA@TEST.COM  ",
		BirthDay:  "1995-05-15",
	}
	user, err := req.ToDomain(10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.TenantID != 10 {
		t.Errorf("expected tenant 10, got %d", user.TenantID)
	}
	birthStr := user.BirthDay.Format("2006-01-02")
	if birthStr != "1995-05-15" {
		t.Errorf("expected '1995-05-15', got '%s'", birthStr)
	}
	if strings.Contains(user.Email, " ") {
		t.Error("expected email to be trimmed")
	}
	if !strings.Contains(user.Email, "@test.com") {
		t.Errorf("expected email to contain @test.com, got '%s'", user.Email)
	}
	_ = domain.User{
		TenantID:  user.TenantID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Dni:       user.Dni,
		Gender:    user.Gender,
		Phone:     user.Phone,
		Email:     user.Email,
		BirthDay:  user.BirthDay,
	}
}

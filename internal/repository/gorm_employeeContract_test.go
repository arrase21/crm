package repository

import (
	"context"
	"testing"
	"time"

	"github.com/arrase21/crm/internal/domain"
	"gorm.io/gorm"
)

func setupTestDBWithEmployeeContract(t *testing.T) *gorm.DB {
	db := setupTestDB(t)
	db.AutoMigrate(&domain.EmployeeContract{}, &domain.Employee{}, &domain.User{}, &domain.ContractType{})
	return db
}

func TestGormEmployeeContractRepo_Create(t *testing.T) {
	db := setupTestDBWithEmployeeContract(t)
	repo := NewGormEmployeeContractRepository(db)
	ctx := contextWithTenant(context.Background(), 1)

	tests := []struct {
		name    string
		input   *domain.EmployeeContract
		wantErr bool
	}{
		{
			name: "create valid contract - success",
			input: &domain.EmployeeContract{
				EmployeeID:      1,
				CountryCode:     "MX",
				Currency:        "MXN",
				BaseSalary:      10000,
				StartDate:       time.Now(),
				WorkHoursPerDay: 8,
				WorkDaysPerWeek: 5,
			},
			wantErr: false,
		},
		{
			name:    "create nil - should fail",
			input:   nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(ctx, tt.input)
			if tt.wantErr && err == nil {
				t.Error("expected error but got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
		})
	}
}

func TestGormEmployeeContractRepo_GetByID(t *testing.T) {
	db := setupTestDBWithEmployeeContract(t)
	repo := NewGormEmployeeContractRepository(db)
	ctx := contextWithTenant(context.Background(), 1)

	ec := &domain.EmployeeContract{
		EmployeeID:      1,
		CountryCode:     "MX",
		Currency:        "MXN",
		BaseSalary:      10000,
		StartDate:       time.Now(),
		WorkHoursPerDay: 8,
		WorkDaysPerWeek: 5,
	}
	repo.Create(ctx, ec)

	tests := []struct {
		name    string
		inputID uint
		wantErr bool
	}{
		{name: "get existing - success", inputID: ec.ID, wantErr: false},
		{name: "get non-existing - should fail", inputID: 999, wantErr: true},
		{name: "get with id 0 - should fail", inputID: 0, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.GetByID(ctx, tt.inputID)
			if tt.wantErr && err == nil {
				t.Error("expected error but got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
			if !tt.wantErr && result == nil {
				t.Error("expected result but got nil")
			}
		})
	}
}

func TestGormEmployeeContractRepo_GetByEmployeeID(t *testing.T) {
	db := setupTestDBWithEmployeeContract(t)
	repo := NewGormEmployeeContractRepository(db)
	ctx := contextWithTenant(context.Background(), 1)

	for i := 1; i <= 3; i++ {
		repo.Create(ctx, &domain.EmployeeContract{
			EmployeeID:      1,
			CountryCode:     "MX",
			Currency:        "MXN",
			BaseSalary:      10000,
			StartDate:       time.Now().AddDate(0, 0, -i),
			WorkHoursPerDay: 8,
			WorkDaysPerWeek: 5,
		})
	}

	result, err := repo.GetByEmployeeID(ctx, 1)
	if err != nil {
		t.Errorf("expected no error but got: %v", err)
	}
	if len(result) != 3 {
		t.Errorf("expected 3 contracts but got %d", len(result))
	}
}

func TestGormEmployeeContractRepo_List(t *testing.T) {
	db := setupTestDBWithEmployeeContract(t)
	repo := NewGormEmployeeContractRepository(db)
	ctx := contextWithTenant(context.Background(), 1)

	for i := 1; i <= 5; i++ {
		repo.Create(ctx, &domain.EmployeeContract{
			EmployeeID:      uint(i),
			CountryCode:     "MX",
			Currency:        "MXN",
			BaseSalary:      10000,
			StartDate:       time.Now(),
			WorkHoursPerDay: 8,
			WorkDaysPerWeek: 5,
		})
	}

	result, total, err := repo.List(ctx, 1, 10)
	if err != nil {
		t.Errorf("expected no error but got: %v", err)
	}
	if total != 5 {
		t.Errorf("expected total 5 but got %d", total)
	}
	if len(result) != 5 {
		t.Errorf("expected 5 items but got %d", len(result))
	}
}

func TestGormEmployeeContractRepo_Update(t *testing.T) {
	db := setupTestDBWithEmployeeContract(t)
	repo := NewGormEmployeeContractRepository(db)
	ctx := contextWithTenant(context.Background(), 1)

	ec := &domain.EmployeeContract{
		EmployeeID:      1,
		CountryCode:     "MX",
		Currency:        "MXN",
		BaseSalary:      10000,
		StartDate:       time.Now(),
		WorkHoursPerDay: 8,
		WorkDaysPerWeek: 5,
	}
	repo.Create(ctx, ec)

	ec.BaseSalary = 20000
	err := repo.Update(ctx, ec)
	if err != nil {
		t.Errorf("expected no error but got: %v", err)
	}

	err = repo.Update(ctx, &domain.EmployeeContract{ID: 999})
	if err == nil {
		t.Error("expected error for non-existing")
	}
}

func TestGormEmployeeContractRepo_Delete(t *testing.T) {
	db := setupTestDBWithEmployeeContract(t)
	repo := NewGormEmployeeContractRepository(db)
	ctx := contextWithTenant(context.Background(), 1)

	ec := &domain.EmployeeContract{
		EmployeeID:      1,
		CountryCode:     "MX",
		Currency:        "MXN",
		BaseSalary:      10000,
		StartDate:       time.Now(),
		WorkHoursPerDay: 8,
		WorkDaysPerWeek: 5,
	}
	repo.Create(ctx, ec)

	tests := []struct {
		name    string
		inputID uint
		wantErr bool
	}{
		{name: "delete existing - success", inputID: ec.ID, wantErr: false},
		{name: "delete non-existing - should fail", inputID: 999, wantErr: true},
		{name: "delete with id 0 - should fail", inputID: 0, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Delete(ctx, tt.inputID)
			if tt.wantErr && err == nil {
				t.Error("expected error but got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
		})
	}
}

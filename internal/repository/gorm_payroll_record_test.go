package repository

import (
	"context"
	"testing"
	"time"

	"github.com/arrase21/crm/internal/domain"
	"gorm.io/gorm"
)

func setupTestDBWithPayroll(t *testing.T) *gorm.DB {
	db := setupTestDB(t)
	db.AutoMigrate(&domain.PayrollRecord{}, &domain.EmployeeContract{}, &domain.Employee{}, &domain.User{})
	return db
}

func TestGormPayrollRecordRepo_Create(t *testing.T) {
	db := setupTestDBWithPayroll(t)
	repo := NewGormPayrollRecordRepository(db)
	ctx := contextWithTenant(context.Background(), 1)

	tests := []struct {
		name    string
		input   *domain.PayrollRecord
		wantErr bool
	}{
		{
			name: "create valid payroll record - success",
			input: &domain.PayrollRecord{
				EmployeeContractID: 1,
				PeriodStart:        time.Now().AddDate(0, -1, 0),
				PeriodEnd:          time.Now(),
				BaseSalary:         10000,
				Currency:           "MXN",
				CountryCode:        "MX",
				GrossSalary:        10000,
				NetSalary:          8500,
				CalculatedAt:       time.Now(),
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

func TestGormPayrollRecordRepo_GetByID(t *testing.T) {
	db := setupTestDBWithPayroll(t)
	repo := NewGormPayrollRecordRepository(db)
	ctx := contextWithTenant(context.Background(), 1)

	pr := &domain.PayrollRecord{
		EmployeeContractID: 1,
		PeriodStart:        time.Now().AddDate(0, -1, 0),
		PeriodEnd:          time.Now(),
		BaseSalary:         10000,
		Currency:           "MXN",
		CountryCode:        "MX",
		GrossSalary:        10000,
		NetSalary:          8500,
		CalculatedAt:       time.Now(),
	}
	repo.Create(ctx, pr)

	tests := []struct {
		name    string
		inputID uint
		wantErr bool
	}{
		{name: "get existing - success", inputID: pr.ID, wantErr: false},
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

func TestGormPayrollRecordRepo_GetByContractID(t *testing.T) {
	db := setupTestDBWithPayroll(t)
	repo := NewGormPayrollRecordRepository(db)
	ctx := contextWithTenant(context.Background(), 1)

	for i := 1; i <= 3; i++ {
		repo.Create(ctx, &domain.PayrollRecord{
			EmployeeContractID: 1,
			PeriodStart:        time.Now().AddDate(0, 0, -i*30),
			PeriodEnd:          time.Now().AddDate(0, 0, -i*30+30),
			BaseSalary:         10000,
			Currency:           "MXN",
			CountryCode:        "MX",
			GrossSalary:        10000,
			NetSalary:          8500,
			CalculatedAt:       time.Now(),
		})
	}

	result, total, err := repo.GetByContractID(ctx, 1, 1, 10)
	if err != nil {
		t.Errorf("expected no error but got: %v", err)
	}
	if total != 3 {
		t.Errorf("expected total 3 but got %d", total)
	}
	if len(result) != 3 {
		t.Errorf("expected 3 items but got %d", len(result))
	}
}

func TestGormPayrollRecordRepo_List(t *testing.T) {
	db := setupTestDBWithPayroll(t)
	repo := NewGormPayrollRecordRepository(db)
	ctx := contextWithTenant(context.Background(), 1)

	for i := 1; i <= 5; i++ {
		repo.Create(ctx, &domain.PayrollRecord{
			EmployeeContractID: uint(i),
			PeriodStart:        time.Now().AddDate(0, -1, 0),
			PeriodEnd:          time.Now(),
			BaseSalary:         10000,
			Currency:           "MXN",
			CountryCode:        "MX",
			GrossSalary:        10000,
			NetSalary:          8500,
			CalculatedAt:       time.Now(),
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

package repository

import (
	"context"
	"testing"

	"github.com/arrase21/crm/internal/domain"
	"gorm.io/gorm"
)

func setupTestDBWithEmployee(t *testing.T) *gorm.DB {
	db := setupTestDB(t)
	db.AutoMigrate(&domain.Employee{}, &domain.User{}, &domain.Department{}, &domain.Position{}, &domain.EmployeeContract{}, &domain.ContractType{})
	return db
}

func TestGormEmployeeRepo_Create(t *testing.T) {
	db := setupTestDBWithEmployee(t)
	repo := NewGormEmployeeRepository(db)
	ctx := contextWithTenant(context.Background(), 1)

	tests := []struct {
		name    string
		input   *domain.Employee
		wantErr bool
	}{
		{
			name: "create valid employee - success",
			input: &domain.Employee{
				TenantID:     1,
				UserID:       1,
				DepartmentID: 1,
				PositionID:   1,
				IsActive:     true,
			},
			wantErr: false,
		},
		{
			name:    "create nil employee - should fail",
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

func TestGormEmployeeRepo_GetByID(t *testing.T) {
	db := setupTestDBWithEmployee(t)
	repo := NewGormEmployeeRepository(db)
	ctx := contextWithTenant(context.Background(), 1)

	emp := &domain.Employee{
		TenantID:     1,
		UserID:       1,
		DepartmentID: 1,
		PositionID:   1,
		IsActive:     true,
	}
	repo.Create(ctx, emp)

	tests := []struct {
		name       string
		inputID    uint
		wantErr    bool
		checkID    bool
		expectedID uint
	}{
		{
			name:       "get existing employee - success",
			inputID:    emp.ID,
			wantErr:    false,
			checkID:    true,
			expectedID: emp.ID,
		},
		{
			name:    "get non-existing employee - error",
			inputID: 999,
			wantErr: true,
		},
		{
			name:    "get with id 0 - error",
			inputID: 0,
			wantErr: true,
		},
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
			if tt.checkID && result != nil && result.ID != tt.expectedID {
				t.Errorf("expected ID %d but got %d", tt.expectedID, result.ID)
			}
		})
	}
}

func TestGormEmployeeRepo_GetByUserID(t *testing.T) {
	db := setupTestDBWithEmployee(t)
	repo := NewGormEmployeeRepository(db)
	ctx := contextWithTenant(context.Background(), 1)

	emp := &domain.Employee{
		TenantID:     1,
		UserID:       1,
		DepartmentID: 1,
		PositionID:   1,
		IsActive:     true,
	}
	repo.Create(ctx, emp)

	tests := []struct {
		name      string
		inputUserID uint
		wantErr   bool
	}{
		{
			name:        "get existing employee by user id - success",
			inputUserID: 1,
			wantErr:     false,
		},
		{
			name:        "get by user id - not found",
			inputUserID: 999,
			wantErr:     true,
		},
		{
			name:        "get by user id 0 - should fail",
			inputUserID: 0,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.GetByUserID(ctx, tt.inputUserID)
			if tt.wantErr && err == nil {
				t.Error("expected error but got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
			if !tt.wantErr && result == nil {
				t.Error("expected employee but got nil")
			}
		})
	}
}

func TestGormEmployeeRepo_List(t *testing.T) {
	db := setupTestDBWithEmployee(t)
	repo := NewGormEmployeeRepository(db)
	ctx := contextWithTenant(context.Background(), 1)

	for i := 1; i <= 5; i++ {
		repo.Create(ctx, &domain.Employee{
			TenantID:     1,
			UserID:       uint(i),
			DepartmentID: 1,
			PositionID:   1,
			IsActive:     true,
		})
	}

	tests := []struct {
		name      string
		page      int
		limit     int
		wantCount int
		wantTotal int64
	}{
		{
			name:      "list all employees",
			page:      1,
			limit:     10,
			wantCount: 5,
			wantTotal: 5,
		},
		{
			name:      "list with pagination - page 1, limit 2",
			page:      1,
			limit:     2,
			wantCount: 2,
			wantTotal: 5,
		},
		{
			name:      "list second page - page 2, limit 2",
			page:      2,
			limit:     2,
			wantCount: 2,
			wantTotal: 5,
		},
		{
			name:      "list page beyond available",
			page:      10,
			limit:     10,
			wantCount: 0,
			wantTotal: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			employees, total, err := repo.List(ctx, tt.page, tt.limit)
			if err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
			if total != tt.wantTotal {
				t.Errorf("expected total %d but got %d", tt.wantTotal, total)
			}
			if len(employees) != tt.wantCount {
				t.Errorf("expected %d employees but got %d", tt.wantCount, len(employees))
			}
		})
	}
}

func TestGormEmployeeRepo_ListActive(t *testing.T) {
	db := setupTestDBWithEmployee(t)
	repo := NewGormEmployeeRepository(db)
	ctx := contextWithTenant(context.Background(), 1)

	// Create employees (all active by default due to GORM zero-value handling)
	for i := 1; i <= 4; i++ {
		repo.Create(ctx, &domain.Employee{
			TenantID:     1,
			UserID:       uint(i),
			DepartmentID: 1,
			PositionID:   1,
			IsActive:     true,
		})
	}
	// Set some to inactive via Update
	repo.Update(ctx, &domain.Employee{ID: 1, TenantID: 1, UserID: 1, DepartmentID: 1, PositionID: 1, IsActive: false})
	repo.Update(ctx, &domain.Employee{ID: 3, TenantID: 1, UserID: 3, DepartmentID: 1, PositionID: 1, IsActive: false})

	employees, total, err := repo.ListActive(ctx, 1, 10)
	if err != nil {
		t.Errorf("expected no error but got: %v", err)
	}
	if total != 2 {
		t.Errorf("expected total 2 active employees but got %d", total)
	}
	if len(employees) != 2 {
		t.Errorf("expected 2 active employees but got %d", len(employees))
	}
}

func TestGormEmployeeRepo_Update(t *testing.T) {
	db := setupTestDBWithEmployee(t)
	repo := NewGormEmployeeRepository(db)
	ctx := contextWithTenant(context.Background(), 1)

	emp := &domain.Employee{
		TenantID:     1,
		UserID:       1,
		DepartmentID: 1,
		PositionID:   1,
		IsActive:     true,
	}
	repo.Create(ctx, emp)

	tests := []struct {
		name    string
		input   *domain.Employee
		wantErr bool
	}{
		{
			name: "update existing employee - success",
			input: &domain.Employee{
				ID:           emp.ID,
				TenantID:     1,
				UserID:       1,
				DepartmentID: 2,
				PositionID:   2,
				IsActive:     false,
			},
			wantErr: false,
		},
		{
			name:    "update with nil - should fail",
			input:   nil,
			wantErr: true,
		},
		{
			name: "update non-existing - should fail",
			input: &domain.Employee{
				ID: 999,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Update(ctx, tt.input)
			if tt.wantErr && err == nil {
				t.Error("expected error but got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
		})
	}
}

func TestGormEmployeeRepo_Delete(t *testing.T) {
	db := setupTestDBWithEmployee(t)
	repo := NewGormEmployeeRepository(db)
	ctx := contextWithTenant(context.Background(), 1)

	emp := &domain.Employee{
		TenantID:     1,
		UserID:       1,
		DepartmentID: 1,
		PositionID:   1,
		IsActive:     true,
	}
	repo.Create(ctx, emp)

	tests := []struct {
		name    string
		inputID uint
		wantErr bool
	}{
		{
			name:    "delete existing employee - success",
			inputID: emp.ID,
			wantErr: false,
		},
		{
			name:    "delete with id 0 - should fail",
			inputID: 0,
			wantErr: true,
		},
		{
			name:    "delete non-existing - should fail",
			inputID: 999,
			wantErr: true,
		},
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
			if !tt.wantErr {
				_, err := repo.GetByID(ctx, tt.inputID)
				if err == nil {
					t.Error("expected employee to be deleted but still exists")
				}
			}
		})
	}
}

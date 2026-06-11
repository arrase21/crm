package repository

import (
	"context"
	"testing"
	"time"

	"github.com/arrase21/crm/internal/domain"
)

func TestDepartmentRepo_Create(t *testing.T) {
	db := setupTestDB(t)
	db.AutoMigrate(&domain.Department{})
	repo := NewGormDepartmentRepository(db)
	ctx := contextWithTenant(context.Background(), 1)

	tests := []struct {
		name    string
		wantErr bool
		setup   func()
		input   *domain.Department
	}{
		{
			name:    "create valid department - success",
			wantErr: false,
			input: &domain.Department{
				TenantID: 1,
				Name:     "accounting",
				Code:     "acc",
				IsActive: true,
			},
		},
		{
			name:    "create nil department - should fail",
			wantErr: true,
			input:   nil,
		},
		{
			name:    "create duplicate code - should fail",
			wantErr: true,
			setup: func() {
				repo.Create(ctx, &domain.Department{
					TenantID: 1, Name: "technology", Code: "acc",
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				})
			},
			input: &domain.Department{
				TenantID: 1, Name: "accounting 2", Code: "acc",
			},
		},
		{
			name:    "create duplicate name - passes in repo (validation done in service)",
			wantErr: false,
			setup: func() {
				repo.Create(ctx, &domain.Department{
					TenantID: 1, Name: "accounting", Code: "acc2",
				})
			},
			input: &domain.Department{
				TenantID: 1, Name: "accounting", Code: "acc3",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
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

func TestDepartmentRepo_GetByID(t *testing.T) {
	db := setupTestDB(t)
	db.AutoMigrate(&domain.Department{})
	repo := NewGormDepartmentRepository(db)
	ctx := contextWithTenant(context.Background(), 1)

	// Setup: crear un department
	dept := &domain.Department{
		TenantID: 1,
		Name:     "accounting",
		Code:     "acc",
		IsActive: true,
	}
	repo.Create(ctx, dept)

	tests := []struct {
		name    string
		inputID uint
		wantErr bool
	}{
		{
			name:    "get existing department - success",
			inputID: dept.ID,
			wantErr: false,
		},
		{
			name:    "get non-existing department - should fail",
			inputID: 999,
			wantErr: true,
		},
		{
			name:    "get department with id 0 - should fail",
			inputID: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.GetByID(ctx, tt.inputID)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
				if result.ID != tt.inputID {
					t.Errorf("expected id %d but got %d", tt.inputID, result.ID)
				}
			}
		})
	}
}

func TestDepartmentRepo_GetByCode(t *testing.T) {
	db := setupTestDB(t)
	db.AutoMigrate(&domain.Department{})
	repo := NewGormDepartmentRepository(db)
	ctx := contextWithTenant(context.Background(), 1)

	// Setup: crear un department
	dept := &domain.Department{
		TenantID: 1,
		Name:     "accounting",
		Code:     "ACC",
		IsActive: true,
	}
	repo.Create(ctx, dept)

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "get by code - success",
			input:   "ACC",
			wantErr: false,
		},
		{
			name:    "get by code - not found",
			input:   "XXX",
			wantErr: true,
		},
		{
			name:    "get by code - empty string should fail",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.GetByCode(ctx, tt.input)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
				if result.Code != tt.input {
					t.Errorf("expected code %s but got %s", tt.input, result.Code)
				}
			}
		})
	}
}

func TestDepartmentRepo_GetByName(t *testing.T) {
	db := setupTestDB(t)
	db.AutoMigrate(&domain.Department{})
	repo := NewGormDepartmentRepository(db)
	ctx := contextWithTenant(context.Background(), 1)

	// Setup: crear un department
	dept := &domain.Department{
		TenantID: 1,
		Name:     "Accounting",
		Code:     "ACC",
		IsActive: true,
	}
	repo.Create(ctx, dept)

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "get by name - success",
			input:   "Accounting",
			wantErr: false,
		},
		{
			name:    "get by name - not found",
			input:   "NonExistent",
			wantErr: true,
		},
		{
			name:    "get by name - empty string should fail",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.GetByName(ctx, tt.input)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
				if result.Name != tt.input {
					t.Errorf("expected name %s but got %s", tt.input, result.Name)
				}
			}
		})
	}
}

func TestDepartmentRepo_List(t *testing.T) {
	db := setupTestDB(t)
	db.AutoMigrate(&domain.Department{})
	repo := NewGormDepartmentRepository(db)
	ctx := contextWithTenant(context.Background(), 1)

	// Setup: crear varios departments
	for i := 1; i <= 5; i++ {
		repo.Create(ctx, &domain.Department{
			TenantID: 1,
			Name:     "Dept" + string(rune('0'+i)),
			Code:     string(rune('A' + i)),
			IsActive: true,
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
			name:      "list all",
			page:      1,
			limit:     10,
			wantCount: 5,
			wantTotal: 5,
		},
		{
			name:      "list with pagination",
			page:      1,
			limit:     2,
			wantCount: 2,
			wantTotal: 5,
		},
		{
			name:      "list second page",
			page:      2,
			limit:     2,
			wantCount: 2,
			wantTotal: 5,
		},
		{
			name:      "list beyond page",
			page:      10,
			limit:     2,
			wantCount: 0,
			wantTotal: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, total, err := repo.List(ctx, tt.page, tt.limit)
			if err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
			if len(result) != tt.wantCount {
				t.Errorf("expected %d items but got %d", tt.wantCount, len(result))
			}
			if total != tt.wantTotal {
				t.Errorf("expected total %d but got %d", tt.wantTotal, total)
			}
		})
	}
}

func TestDepartmentRepo_Update(t *testing.T) {
	db := setupTestDB(t)
	db.AutoMigrate(&domain.Department{})
	repo := NewGormDepartmentRepository(db)
	ctx := contextWithTenant(context.Background(), 1)

	// Setup: crear un department
	original := &domain.Department{
		TenantID: 1,
		Name:     "Accounting",
		Code:     "ACC",
		IsActive: true,
	}
	repo.Create(ctx, original)

	tests := []struct {
		name    string
		input   *domain.Department
		wantErr bool
	}{
		{
			name: "update name - success",
			input: &domain.Department{
				ID:       original.ID,
				TenantID: 1,
				Name:     "Finance",
				Code:     "ACC",
				IsActive: true,
			},
			wantErr: false,
		},
		{
			name: "update non-existing - should fail",
			input: &domain.Department{
				ID:       999,
				TenantID: 1,
				Name:     "Finance",
				Code:     "FIN",
				IsActive: true,
			},
			wantErr: true,
		},
		{
			name:    "update with nil - should fail",
			input:   nil,
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

func TestDepartmentRepo_Delete(t *testing.T) {
	db := setupTestDB(t)
	db.AutoMigrate(&domain.Department{})
	repo := NewGormDepartmentRepository(db)
	ctx := contextWithTenant(context.Background(), 1)

	// Setup: crear un department
	dept := &domain.Department{
		TenantID: 1,
		Name:     "Accounting",
		Code:     "ACC",
		IsActive: true,
	}
	repo.Create(ctx, dept)

	tests := []struct {
		name    string
		inputID uint
		wantErr bool
	}{
		{
			name:    "delete existing - success",
			inputID: dept.ID,
			wantErr: false,
		},
		{
			name:    "delete non-existing - should fail",
			inputID: 999,
			wantErr: true,
		},
		{
			name:    "delete with id 0 - should fail",
			inputID: 0,
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
		})
	}
}

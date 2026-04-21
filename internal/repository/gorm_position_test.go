package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/arrase21/crm/internal/domain"
)

// ============================================
// TESTS - GormPositionRepo
// ============================================

func TestGormPositionRepo_Create(t *testing.T) {
	db := setupTestDB(t)
	db.AutoMigrate(&domain.Position{})
	repo := NewGormPositionRepository(db)
	ctx := contextWithTenant(context.Background(), 1)

	tests := []struct {
		name    string
		input   *domain.Position
		wantErr bool
	}{
		{
			name: "create valid position - success",
			input: &domain.Position{
				TenantID:    1,
				Name:        "Gerente",
				Description: "Gerente de área",
				IsActive:    true,
			},
			wantErr: false,
		},
		{
			name:    "create nil position - should fail",
			input:   nil,
			wantErr: true,
		},
		// Nota: El test de duplicate depende del constraint UNIQUE en la DB
		{
			name: "create duplicate name - should fail or succeed based on DB constraint",
			input: &domain.Position{
				TenantID:    1,
				Name:        "Gerente", // Same name
				Description: "Another description",
				IsActive:    true,
			},
			wantErr: false, // Puede passar si la DB no tiene constraint
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Si es el test de duplicate, primero creamos uno
			if tt.name == "create duplicate name - should fail or succeed based on DB constraint" {
				repo.Create(ctx, &domain.Position{
					TenantID:    1,
					Name:        "Gerente",
					Description: "First description",
					IsActive:    true,
				})
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

func TestGormPositionRepo_GetByID(t *testing.T) {
	db := setupTestDB(t)
	db.AutoMigrate(&domain.Position{})
	repo := NewGormPositionRepository(db)
	ctx := contextWithTenant(context.Background(), 1)

	position := &domain.Position{
		TenantID:    1,
		Name:        "Gerente",
		Description: "Gerente de área",
		IsActive:    true,
	}
	repo.Create(ctx, position)

	tests := []struct {
		name       string
		inputID    uint
		wantErr    bool
		checkName  bool
		expectedID uint
		expected   string
	}{
		{
			name:       "get existing position - success",
			inputID:    position.ID,
			wantErr:    false,
			checkName:  true,
			expectedID: position.ID,
			expected:   "Gerente",
		},
		{
			name:    "get non-existing position - error",
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
			if tt.checkName && result != nil && result.Name != tt.expected {
				t.Errorf("expected name %s but got %s", tt.expected, result.Name)
			}
		})
	}
}

func TestGormPositionRepo_GetByName(t *testing.T) {
	db := setupTestDB(t)
	db.AutoMigrate(&domain.Position{})
	repo := NewGormPositionRepository(db)
	ctx := contextWithTenant(context.Background(), 1)

	// Create position first
	position := &domain.Position{
		TenantID:    1,
		Name:        "Gerente",
		Description: "Gerente de área",
		IsActive:    true,
	}
	repo.Create(ctx, position)

	tests := []struct {
		name      string
		inputName string
		wantErr   bool
		checkDesc bool
		expected  string
	}{
		{
			name:      "get existing position by name - success",
			inputName: "Gerente",
			wantErr:   false,
			checkDesc: true,
			expected:  "Gerente de área",
		},
		{
			name:      "get by name - empty string - should fail",
			inputName: "",
			wantErr:   true,
		},
		{
			name:      "get by name - not found - should fail",
			inputName: "NoExiste",
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.GetByName(ctx, tt.inputName)
			if tt.wantErr && err == nil {
				t.Error("expected error but got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
			if tt.checkDesc && result != nil && result.Description != tt.expected {
				t.Errorf("expected description %s but got %s", tt.expected, result.Description)
			}
		})
	}
}

func TestGormPositionRepo_List(t *testing.T) {
	db := setupTestDB(t)
	db.AutoMigrate(&domain.Position{})
	repo := NewGormPositionRepository(db)
	ctx := contextWithTenant(context.Background(), 1)

	// Create multiple positions
	for i := 1; i <= 5; i++ {
		repo.Create(ctx, &domain.Position{
			TenantID:    1,
			Name:        fmt.Sprintf("Position%d", i),
			Description: fmt.Sprintf("Description %d", i),
			IsActive:    true,
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
			name:      "list all positions - no pagination",
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
			name:      "list page beyond available - should return empty",
			page:      10,
			limit:     10,
			wantCount: 0,
			wantTotal: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			positions, total, err := repo.List(ctx, tt.page, tt.limit)
			if err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
			if total != tt.wantTotal {
				t.Errorf("expected total %d but got %d", tt.wantTotal, total)
			}
			if len(positions) != tt.wantCount {
				t.Errorf("expected %d positions but got %d", tt.wantCount, len(positions))
			}
		})
	}
}

func TestGormPositionRepo_Update(t *testing.T) {
	db := setupTestDB(t)
	db.AutoMigrate(&domain.Position{})
	repo := NewGormPositionRepository(db)
	ctx := contextWithTenant(context.Background(), 1)

	// Create position first
	position := &domain.Position{
		TenantID:    1,
		Name:        "Gerente",
		Description: "Gerente de área",
		IsActive:    true,
	}
	repo.Create(ctx, position)

	tests := []struct {
		name    string
		input   *domain.Position
		wantErr bool
	}{
		{
			name: "update existing position - success",
			input: &domain.Position{
				ID:          position.ID,
				TenantID:    1,
				Name:        "Gerente",
				Description: "Gerente actualizado", // Changed
				IsActive:    false,
			},
			wantErr: false,
		},
		{
			name: "update with nil - should fail",
			input: &domain.Position{
				ID: 0,
			},
			wantErr: true,
		},
		{
			name: "update non-existing position - should fail",
			input: &domain.Position{
				ID:         999,
				TenantID:   1,
				Name:       "NoExiste",
				Description: "Description",
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

func TestGormPositionRepo_Delete(t *testing.T) {
	db := setupTestDB(t)
	db.AutoMigrate(&domain.Position{})
	repo := NewGormPositionRepository(db)
	ctx := contextWithTenant(context.Background(), 1)

	// Create position first
	position := &domain.Position{
		TenantID:    1,
		Name:        "Gerente",
		Description: "Gerente de área",
		IsActive:    true,
	}
	repo.Create(ctx, position)

	tests := []struct {
		name    string
		inputID uint
		wantErr bool
	}{
		{
			name:    "delete existing position - success",
			inputID: position.ID,
			wantErr: false,
		},
		{
			name:    "delete with id 0 - should fail",
			inputID: 0,
			wantErr: true,
		},
		{
			name:    "delete non-existing position - should fail",
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
		})
	}
}

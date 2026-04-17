package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/arrase21/crm/internal/domain"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ============================================
// HELPER - Setup DB en memoria
// ============================================
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to test db: %v", err)
	}
	db.AutoMigrate(&domain.User{})
	return db
}

// Helper para crear contexto con tenant
func contextWithTenant(ctx context.Context, tenantID uint) context.Context {
	return context.WithValue(ctx, domain.TenantIDKey, tenantID)
}

// ============================================
// TESTS - GormUserRepo
// ============================================
func TestGormUserRepo_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormUserRepository(db)
	ctx := contextWithTenant(context.Background(), 1)
	tests := []struct {
		name    string
		input   *domain.User
		wantErr bool
	}{
		{
			name: "create valid user - success",
			input: &domain.User{
				TenantID:  1,
				FirstName: "Juan",
				LastName:  "Pérez",
				Email:     "juan@test.com",
				Dni:       "12345678",
				Phone:     "1234567890",
				Gender:    "M",
				BirthDay:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name:    "create nil user - should fail",
			input:   nil,
			wantErr: true,
		},
		{
			name: "create duplicate DNI - should fail",
			input: &domain.User{
				TenantID:  1,
				FirstName: "Pedro",
				LastName:  "Gómez",
				Email:     "pedro@test.com",
				Dni:       "12345678", // Same DNI
				Phone:     "9999999999",
				Gender:    "M",
				BirthDay:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Si es el test de duplicate, primero creamos uno
			if tt.name == "create duplicate DNI - should fail" {
				repo.Create(ctx, &domain.User{
					TenantID: 1, FirstName: "Juan", LastName: "Pérez",
					Email: "juan@test.com", Dni: "12345678",
					Phone: "1234567890", Gender: "M",
					BirthDay: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
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
func TestGormUserRepo_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormUserRepository(db)
	ctx := contextWithTenant(context.Background(), 1)
	user := &domain.User{
		TenantID:  1,
		FirstName: "Juan",
		LastName:  "Pérez",
		Email:     "juan@test.com",
		Dni:       "12345678",
		Phone:     "1234567890",
		Gender:    "M",
		BirthDay:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	repo.Create(ctx, user)
	tests := []struct {
		name       string
		inputID    uint
		wantErr    bool
		checkID    bool
		expectedID uint
	}{
		{
			name:       "get existing user - success",
			inputID:    user.ID,
			wantErr:    false,
			checkID:    true,
			expectedID: user.ID,
		},
		{
			name:    "get non-existing user - error",
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
func TestGormUserRepo_GetByDNI(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormUserRepository(db)
	ctx := contextWithTenant(context.Background(), 1)

	// Create user first
	user := &domain.User{
		TenantID:  1,
		FirstName: "Juan",
		LastName:  "Pérez",
		Email:     "juan@test.com",
		Dni:       "12345678",
		Phone:     "1234567890",
		Gender:    "M",
		BirthDay:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	repo.Create(ctx, user)

	tests := []struct {
		name      string
		inputDNI  string
		wantErr   bool
		checkName bool
		expected  string
	}{
		{
			name:      "get existing user by DNI - success",
			inputDNI:  "12345678",
			wantErr:   false,
			checkName: true,
			expected:  "Juan",
		},
		{
			name:     "get by DNI - empty string - should fail",
			inputDNI: "",
			wantErr:  true,
		},
		{
			name:     "get by DNI - not found - should fail",
			inputDNI: "99999999",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.GetByDNI(ctx, tt.inputDNI)
			if tt.wantErr && err == nil {
				t.Error("expected error but got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
			if tt.checkName && result != nil && result.FirstName != tt.expected {
				t.Errorf("expected name %s but got %s", tt.expected, result.FirstName)
			}
		})
	}
}
func TestGormUserRepo_List(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormUserRepository(db)
	ctx := contextWithTenant(context.Background(), 1)

	// Create multiple users
	for i := 1; i <= 5; i++ {
		repo.Create(ctx, &domain.User{
			TenantID:  1,
			FirstName: "User",
			LastName:  fmt.Sprintf("Test%d", i),
			Email:     fmt.Sprintf("user%d@test.com", i),
			Dni:       fmt.Sprintf("1111111%d", i),
			Phone:     fmt.Sprintf("111111111%d", i),
			Gender:    "M",
			BirthDay:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
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
			name:      "list all users - no pagination",
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
			users, total, err := repo.List(ctx, tt.page, tt.limit)
			if err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
			if total != tt.wantTotal {
				t.Errorf("expected total %d but got %d", tt.wantTotal, total)
			}
			if len(users) != tt.wantCount {
				t.Errorf("expected %d users but got %d", tt.wantCount, len(users))
			}
		})
	}
}
func TestGormUserRepo_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormUserRepository(db)
	ctx := contextWithTenant(context.Background(), 1)

	// Create user first
	user := &domain.User{
		TenantID:  1,
		FirstName: "Juan",
		LastName:  "Pérez",
		Email:     "juan@test.com",
		Dni:       "12345678",
		Phone:     "1234567890",
		Gender:    "M",
		BirthDay:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	repo.Create(ctx, user)

	tests := []struct {
		name      string
		input     *domain.User
		wantErr   bool
		checkName string
	}{
		{
			name: "update existing user - success",
			input: &domain.User{
				ID:        user.ID,
				TenantID:  1,
				FirstName: "Carlos", // Changed
				LastName:  "Pérez",
				Email:     "carlos@test.com",
				Dni:       "12345678",
				Phone:     "1234567890",
				Gender:    "M",
				BirthDay:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			wantErr:   false,
			checkName: "Carlos",
		},
		{
			name: "update with nil - should fail",
			input: &domain.User{
				ID: 0,
			},
			wantErr: true,
		},
		{
			name: "update non-existing user - should fail",
			input: &domain.User{
				ID:        999,
				TenantID:  1,
				FirstName: "NotExists",
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
			// Verify update if successful
			if !tt.wantErr && tt.checkName != "" {
				updated, _ := repo.GetByID(ctx, tt.input.ID)
				if updated != nil && updated.FirstName != tt.checkName {
					t.Errorf("expected name %s but got %s", tt.checkName, updated.FirstName)
				}
			}
		})
	}
}
func TestGormUserRepo_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormUserRepository(db)
	ctx := contextWithTenant(context.Background(), 1)

	// Create user first
	user := &domain.User{
		TenantID:  1,
		FirstName: "Juan",
		LastName:  "Pérez",
		Email:     "juan@test.com",
		Dni:       "12345678",
		Phone:     "1234567890",
		Gender:    "M",
		BirthDay:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	repo.Create(ctx, user)

	tests := []struct {
		name    string
		inputID uint
		wantErr bool
	}{
		{
			name:    "delete existing user - success",
			inputID: user.ID,
			wantErr: false,
		},
		{
			name:    "delete with id 0 - should fail",
			inputID: 0,
			wantErr: true,
		},
		{
			name:    "delete non-existing user - should fail",
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
			// Verify deletion if successful
			if !tt.wantErr {
				_, err := repo.GetByID(ctx, tt.inputID)
				if err == nil {
					t.Error("expected user to be deleted but still exists")
				}
			}
		})
	}
}

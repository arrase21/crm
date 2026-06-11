package repository

import (
	"context"
	"testing"

	"github.com/arrase21/crm/internal/domain"
)

func TestContractTypeRepo_Create(t *testing.T) {
	db := setupTestDB(t)
	db.AutoMigrate(&domain.ContractType{})
	repo := NewGormContractTypeRepository(db)
	ctx := context.Background()

	tests := []struct {
		name    string
		input   *domain.ContractType
		wantErr bool
	}{
		{
			name: "create valid contract type - success",
			input: &domain.ContractType{
				Name:        "Full Time",
				Description: "Full time contract",
				CountryCode: "MX",
				IsActive:    true,
			},
			wantErr: false,
		},
		{
			name:    "create nil contract type - should fail",
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

func TestContractTypeRepo_GetByID(t *testing.T) {
	db := setupTestDB(t)
	db.AutoMigrate(&domain.ContractType{})
	repo := NewGormContractTypeRepository(db)
	ctx := context.Background()

	ct := &domain.ContractType{
		Name:        "Full Time",
		CountryCode: "MX",
		IsActive:    true,
	}
	repo.Create(ctx, ct)

	tests := []struct {
		name    string
		inputID uint
		wantErr bool
	}{
		{
			name:    "get existing - success",
			inputID: ct.ID,
			wantErr: false,
		},
		{
			name:    "get non-existing - should fail",
			inputID: 999,
			wantErr: true,
		},
		{
			name:    "get with id 0 - should fail",
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
			if !tt.wantErr && result == nil {
				t.Error("expected result but got nil")
			}
		})
	}
}

func TestContractTypeRepo_GetByCountry(t *testing.T) {
	db := setupTestDB(t)
	db.AutoMigrate(&domain.ContractType{})
	repo := NewGormContractTypeRepository(db)
	ctx := context.Background()

	repo.Create(ctx, &domain.ContractType{Name: "Full Time", CountryCode: "MX", IsActive: true})
	repo.Create(ctx, &domain.ContractType{Name: "Part Time", CountryCode: "MX", IsActive: true})
	repo.Create(ctx, &domain.ContractType{Name: "Full Time", CountryCode: "CO", IsActive: true})

	tests := []struct {
		name        string
		countryCode string
		wantErr     bool
		wantCount   int
	}{
		{
			name:        "get by MX - 2 types",
			countryCode: "MX",
			wantErr:     false,
			wantCount:   2,
		},
		{
			name:        "get by CO - 1 type",
			countryCode: "CO",
			wantErr:     false,
			wantCount:   1,
		},
		{
			name:        "get by non-existing - empty",
			countryCode: "XX",
			wantErr:     false,
			wantCount:   0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.GetByCountry(ctx, tt.countryCode)
			if tt.wantErr && err == nil {
				t.Error("expected error but got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
			if len(result) != tt.wantCount {
				t.Errorf("expected %d items but got %d", tt.wantCount, len(result))
			}
		})
	}
}

func TestContractTypeRepo_List(t *testing.T) {
	db := setupTestDB(t)
	db.AutoMigrate(&domain.ContractType{})
	repo := NewGormContractTypeRepository(db)
	ctx := context.Background()

	for i := 1; i <= 5; i++ {
		repo.Create(ctx, &domain.ContractType{
			Name:        "Type",
			CountryCode: "MX",
			IsActive:    true,
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

	// Test pagination
	result, total, err = repo.List(ctx, 1, 2)
	if err != nil {
		t.Errorf("expected no error but got: %v", err)
	}
	if total != 5 {
		t.Errorf("expected total 5 but got %d", total)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 items but got %d", len(result))
	}
}

func TestContractTypeRepo_Update(t *testing.T) {
	db := setupTestDB(t)
	db.AutoMigrate(&domain.ContractType{})
	repo := NewGormContractTypeRepository(db)
	ctx := context.Background()

	ct := &domain.ContractType{Name: "Full Time", CountryCode: "MX", IsActive: true}
	repo.Create(ctx, ct)

	tests := []struct {
		name    string
		input   *domain.ContractType
		wantErr bool
	}{
		{
			name: "update existing - success",
			input: &domain.ContractType{
				ID:          ct.ID,
				Name:        "Part Time",
				CountryCode: "MX",
				IsActive:    false,
			},
			wantErr: false,
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

func TestContractTypeRepo_Delete(t *testing.T) {
	db := setupTestDB(t)
	db.AutoMigrate(&domain.ContractType{})
	repo := NewGormContractTypeRepository(db)
	ctx := context.Background()

	ct := &domain.ContractType{Name: "Full Time", CountryCode: "MX", IsActive: true}
	repo.Create(ctx, ct)

	tests := []struct {
		name    string
		inputID uint
		wantErr bool
	}{
		{
			name:    "delete existing - success",
			inputID: ct.ID,
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

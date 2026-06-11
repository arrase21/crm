package repository

import (
	"context"
	"testing"

	"github.com/arrase21/crm/internal/domain"
)

func TestCountryParamRepo_Create(t *testing.T) {
	db := setupTestDB(t)
	db.AutoMigrate(&domain.CountryParam{})
	repo := NewGormCountryParamRepository(db)
	ctx := context.Background()

	tests := []struct {
		name    string
		input   *domain.CountryParam
		wantErr bool
	}{
		{
			name: "create valid - success",
			input: &domain.CountryParam{
				CountryCode: "MX",
				Name:        "Mexico",
				Currency:    "MXN",
				MinWage:     1000,
				HealthRate:  5,
				PensionRate: 6,
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

func TestCountryParamRepo_GetByID(t *testing.T) {
	db := setupTestDB(t)
	db.AutoMigrate(&domain.CountryParam{})
	repo := NewGormCountryParamRepository(db)
	ctx := context.Background()

	cp := &domain.CountryParam{
		CountryCode: "MX", Name: "Mexico", Currency: "MXN",
		MinWage: 1000, HealthRate: 5, PensionRate: 6,
	}
	repo.Create(ctx, cp)

	tests := []struct {
		name    string
		inputID uint
		wantErr bool
	}{
		{name: "get existing - success", inputID: cp.ID, wantErr: false},
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

func TestCountryParamRepo_GetByCountryCode(t *testing.T) {
	db := setupTestDB(t)
	db.AutoMigrate(&domain.CountryParam{})
	repo := NewGormCountryParamRepository(db)
	ctx := context.Background()

	repo.Create(ctx, &domain.CountryParam{
		CountryCode: "MX", Name: "Mexico", Currency: "MXN",
		MinWage: 1000, HealthRate: 5, PensionRate: 6,
	})

	tests := []struct {
		name        string
		countryCode string
		wantErr     bool
	}{
		{name: "get by MX - success", countryCode: "MX", wantErr: false},
		{name: "get by non-existing - should fail", countryCode: "XX", wantErr: true},
		{name: "get by empty - should fail", countryCode: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.GetByCountryCode(ctx, tt.countryCode)
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

func TestCountryParamRepo_List(t *testing.T) {
	db := setupTestDB(t)
	db.AutoMigrate(&domain.CountryParam{})
	repo := NewGormCountryParamRepository(db)
	ctx := context.Background()

	codes := []string{"MX", "CO", "BR"}
	for i, code := range codes {
		repo.Create(ctx, &domain.CountryParam{
			CountryCode: code, Name: "Country", Currency: "XXX",
			MinWage: 1000, HealthRate: 5, PensionRate: 6,
		})
		_ = i
	}

	result, total, err := repo.List(ctx, 1, 10)
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

func TestCountryParamRepo_Update(t *testing.T) {
	db := setupTestDB(t)
	db.AutoMigrate(&domain.CountryParam{})
	repo := NewGormCountryParamRepository(db)
	ctx := context.Background()

	cp := &domain.CountryParam{
		CountryCode: "MX", Name: "Mexico", Currency: "MXN",
		MinWage: 1000, HealthRate: 5, PensionRate: 6,
	}
	repo.Create(ctx, cp)

	cp.Name = "Mexico Updated"
	err := repo.Update(ctx, cp)
	if err != nil {
		t.Errorf("expected no error but got: %v", err)
	}

	updated, _ := repo.GetByID(ctx, cp.ID)
	if updated == nil || updated.Name != "Mexico Updated" {
		t.Error("expected updated name")
	}

	err = repo.Update(ctx, nil)
	if err == nil {
		t.Error("expected error for nil update")
	}
}

func TestCountryParamRepo_Delete(t *testing.T) {
	db := setupTestDB(t)
	db.AutoMigrate(&domain.CountryParam{})
	repo := NewGormCountryParamRepository(db)
	ctx := context.Background()

	cp := &domain.CountryParam{
		CountryCode: "MX", Name: "Mexico", Currency: "MXN",
		MinWage: 1000, HealthRate: 5, PensionRate: 6,
	}
	repo.Create(ctx, cp)

	tests := []struct {
		name    string
		inputID uint
		wantErr bool
	}{
		{name: "delete existing - success", inputID: cp.ID, wantErr: false},
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

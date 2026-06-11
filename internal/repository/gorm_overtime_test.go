package repository

import (
	"context"
	"testing"
	"time"

	"github.com/arrase21/crm/internal/domain"
	"gorm.io/gorm"
)

func setupTestDBWithOvertime(t *testing.T) *gorm.DB {
	db := setupTestDB(t)
	db.AutoMigrate(&domain.Overtime{}, &domain.Employee{}, &domain.User{})
	return db
}

func TestGormOvertimeRepo_Create(t *testing.T) {
	db := setupTestDBWithOvertime(t)
	repo := NewGormOvertimeRepository(db)
	ctx := contextWithTenant(context.Background(), 1)

	tests := []struct {
		name    string
		input   *domain.Overtime
		wantErr bool
	}{
		{
			name: "create valid overtime - success",
			input: &domain.Overtime{
				EmployeeID: 1,
				Date:       time.Now(),
				Hours:      2.5,
				Reason:     "Extra work",
			},
			wantErr: false,
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

func TestGormOvertimeRepo_GetByID(t *testing.T) {
	db := setupTestDBWithOvertime(t)
	repo := NewGormOvertimeRepository(db)
	ctx := contextWithTenant(context.Background(), 1)

	o := &domain.Overtime{EmployeeID: 1, Date: time.Now(), Hours: 2}
	repo.Create(ctx, o)

	tests := []struct {
		name    string
		inputID uint
		wantErr bool
	}{
		{name: "get existing - success", inputID: o.ID, wantErr: false},
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

func TestGormOvertimeRepo_ListByEmployee(t *testing.T) {
	db := setupTestDBWithOvertime(t)
	repo := NewGormOvertimeRepository(db)
	ctx := contextWithTenant(context.Background(), 1)

	for i := 1; i <= 3; i++ {
		repo.Create(ctx, &domain.Overtime{
			EmployeeID: 1,
			Date:       time.Now().AddDate(0, 0, -i),
			Hours:      2,
		})
	}
	repo.Create(ctx, &domain.Overtime{EmployeeID: 2, Date: time.Now(), Hours: 1})

	result, total, err := repo.ListByEmployee(ctx, 1, 1, 10)
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

func TestGormOvertimeRepo_List(t *testing.T) {
	db := setupTestDBWithOvertime(t)
	repo := NewGormOvertimeRepository(db)
	ctx := contextWithTenant(context.Background(), 1)

	for i := 1; i <= 5; i++ {
		repo.Create(ctx, &domain.Overtime{
			EmployeeID: uint(i),
			Date:       time.Now(),
			Hours:      1,
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

func TestGormOvertimeRepo_Update(t *testing.T) {
	db := setupTestDBWithOvertime(t)
	repo := NewGormOvertimeRepository(db)
	ctx := contextWithTenant(context.Background(), 1)

	o := &domain.Overtime{EmployeeID: 1, Date: time.Now(), Hours: 2}
	repo.Create(ctx, o)

	o.Hours = 4
	err := repo.Update(ctx, o)
	if err != nil {
		t.Errorf("expected no error but got: %v", err)
	}

	// Update non-existing (no error expected — GORM Updates doesn't error on 0 rows)
	err = repo.Update(ctx, &domain.Overtime{ID: 999})
	if err != nil {
		t.Errorf("expected no error but got: %v", err)
	}
}

func TestGormOvertimeRepo_Delete(t *testing.T) {
	db := setupTestDBWithOvertime(t)
	repo := NewGormOvertimeRepository(db)
	ctx := contextWithTenant(context.Background(), 1)

	o := &domain.Overtime{EmployeeID: 1, Date: time.Now(), Hours: 2}
	repo.Create(ctx, o)

	tests := []struct {
		name    string
		inputID uint
		wantErr bool
	}{
		{name: "delete existing - success", inputID: o.ID, wantErr: false},
		{name: "delete non-existing - no error", inputID: 999, wantErr: false},
		{name: "delete with id 0 - no error", inputID: 0, wantErr: false},
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

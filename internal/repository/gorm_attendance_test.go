package repository

import (
	"context"
	"testing"
	"time"

	"github.com/arrase21/crm/internal/domain"
	"gorm.io/gorm"
)

func setupTestDBWithAttendance(t *testing.T) *gorm.DB {
	db := setupTestDB(t)
	db.AutoMigrate(&domain.Attendance{}, &domain.Employee{}, &domain.User{})
	return db
}

func TestGormAttendanceRepo_Create(t *testing.T) {
	db := setupTestDBWithAttendance(t)
	repo := NewGormAttendanceRepository(db)
	ctx := contextWithTenant(context.Background(), 1)

	tests := []struct {
		name    string
		input   *domain.Attendance
		wantErr bool
	}{
		{
			name: "create valid attendance - success",
			input: &domain.Attendance{
				EmployeeID: 1,
				Date:       time.Now(),
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

func TestGormAttendanceRepo_GetByID(t *testing.T) {
	db := setupTestDBWithAttendance(t)
	repo := NewGormAttendanceRepository(db)
	ctx := contextWithTenant(context.Background(), 1)

	a := &domain.Attendance{EmployeeID: 1, Date: time.Now()}
	repo.Create(ctx, a)

	tests := []struct {
		name    string
		inputID uint
		wantErr bool
	}{
		{name: "get existing - success", inputID: a.ID, wantErr: false},
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

func TestGormAttendanceRepo_List(t *testing.T) {
	db := setupTestDBWithAttendance(t)
	repo := NewGormAttendanceRepository(db)
	ctx := contextWithTenant(context.Background(), 1)

	for i := 1; i <= 5; i++ {
		repo.Create(ctx, &domain.Attendance{EmployeeID: uint(i), Date: time.Now()})
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

func TestGormAttendanceRepo_ListByEmployee(t *testing.T) {
	db := setupTestDBWithAttendance(t)
	repo := NewGormAttendanceRepository(db)
	ctx := contextWithTenant(context.Background(), 1)

	for i := 1; i <= 3; i++ {
		repo.Create(ctx, &domain.Attendance{EmployeeID: 1, Date: time.Now().AddDate(0, 0, -i)})
	}
	repo.Create(ctx, &domain.Attendance{EmployeeID: 2, Date: time.Now()})

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

func TestGormAttendanceRepo_Update(t *testing.T) {
	db := setupTestDBWithAttendance(t)
	repo := NewGormAttendanceRepository(db)
	ctx := contextWithTenant(context.Background(), 1)

	a := &domain.Attendance{EmployeeID: 1, Date: time.Now()}
	repo.Create(ctx, a)

	a.Date = time.Now().AddDate(0, 0, 1)
	err := repo.Update(ctx, a)
	if err != nil {
		t.Errorf("expected no error but got: %v", err)
	}

	// Update non-existing (no error expected — GORM Updates doesn't error on 0 rows)
	err = repo.Update(ctx, &domain.Attendance{ID: 999})
	if err != nil {
		t.Errorf("expected no error but got: %v", err)
	}
}

func TestGormAttendanceRepo_Delete(t *testing.T) {
	db := setupTestDBWithAttendance(t)
	repo := NewGormAttendanceRepository(db)
	ctx := contextWithTenant(context.Background(), 1)

	a := &domain.Attendance{EmployeeID: 1, Date: time.Now()}
	repo.Create(ctx, a)

	tests := []struct {
		name    string
		inputID uint
		wantErr bool
	}{
		{name: "delete existing - success", inputID: a.ID, wantErr: false},
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

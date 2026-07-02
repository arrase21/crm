package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/arrase21/crm/internal/domain"
)

type MockAttendanceRepo struct {
	records       map[uint]*domain.Attendance
	employeeIdx   map[uint][]*domain.Attendance
	nextID        uint
	CreateErr     error
	GetByIDErr    error
	UpdateErr     error
	DeleteErr     error
	ListErr       error
}

func NewMockAttendanceRepo() *MockAttendanceRepo {
	return &MockAttendanceRepo{
		records:     make(map[uint]*domain.Attendance),
		employeeIdx: make(map[uint][]*domain.Attendance),
		nextID:      1,
	}
}

func (m *MockAttendanceRepo) Create(ctx context.Context, a *domain.Attendance) error {
	if m.CreateErr != nil {
		return m.CreateErr
	}
	a.ID = m.nextID
	m.nextID++
	m.records[a.ID] = a
	m.employeeIdx[a.EmployeeID] = append(m.employeeIdx[a.EmployeeID], a)
	return nil
}

func (m *MockAttendanceRepo) GetByID(ctx context.Context, id uint) (*domain.Attendance, error) {
	if m.GetByIDErr != nil {
		return nil, m.GetByIDErr
	}
	if id == 0 {
		return nil, errors.New("invalid id")
	}
	a, exists := m.records[id]
	if !exists {
		return nil, errors.New("attendance not found")
	}
	return a, nil
}

func (m *MockAttendanceRepo) ListByEmployee(ctx context.Context, employeeID uint, page, limit int) ([]domain.Attendance, int64, error) {
	records := m.employeeIdx[employeeID]
	result := make([]domain.Attendance, len(records))
	for i, a := range records {
		result[i] = *a
	}
	return result, int64(len(records)), nil
}

func (m *MockAttendanceRepo) ListByDepartment(ctx context.Context, departmentID uint, page, limit int) ([]domain.Attendance, int64, error) {
	return nil, 0, nil
}

func (m *MockAttendanceRepo) List(ctx context.Context, page, limit int) ([]domain.Attendance, int64, error) {
	if m.ListErr != nil {
		return nil, 0, m.ListErr
	}
	var all []domain.Attendance
	for _, a := range m.records {
		all = append(all, *a)
	}
	return all, int64(len(all)), nil
}

func (m *MockAttendanceRepo) Update(ctx context.Context, a *domain.Attendance) error {
	if m.UpdateErr != nil {
		return m.UpdateErr
	}
	if a.ID == 0 {
		return errors.New("invalid id")
	}
	if _, exists := m.records[a.ID]; !exists {
		return errors.New("attendance not found")
	}
	m.records[a.ID] = a
	return nil
}

func (m *MockAttendanceRepo) Delete(ctx context.Context, id uint) error {
	if m.DeleteErr != nil {
		return m.DeleteErr
	}
	if id == 0 {
		return errors.New("invalid id")
	}
	if _, exists := m.records[id]; !exists {
		return errors.New("attendance not found")
	}
	delete(m.records, id)
	return nil
}

func ctxWithClaims(userID uint, roles []string) context.Context {
	return context.WithValue(context.Background(), domain.ClaimsKey, &domain.Claims{
		UserID: userID,
		Roles:  roles,
	})
}

func TestAttendanceService_Create(t *testing.T) {
	t.Run("create valid - success", func(t *testing.T) {
		mockAtt := NewMockAttendanceRepo()
		svc := NewAttendanceService(mockAtt, NewMockEmployeeRepo())
		ctx := context.Background()

		err := svc.Create(ctx, &domain.Attendance{EmployeeID: 1, Date: time.Now()})
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
	})

	t.Run("create with employee id 0 - should fail", func(t *testing.T) {
		mockAtt := NewMockAttendanceRepo()
		svc := NewAttendanceService(mockAtt, NewMockEmployeeRepo())
		ctx := context.Background()

		err := svc.Create(ctx, &domain.Attendance{EmployeeID: 0})
		if err == nil {
			t.Error("Expected error for missing employee")
		}
	})
}

func TestAttendanceService_GetByID(t *testing.T) {
	mockAtt := NewMockAttendanceRepo()
	mockAtt.records[1] = &domain.Attendance{ID: 1, EmployeeID: 1}
	svc := NewAttendanceService(mockAtt, NewMockEmployeeRepo())
	ctx := context.Background()

	t.Run("get existing - success", func(t *testing.T) {
		result, err := svc.GetByID(ctx, 1)
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
		if result == nil || result.ID != 1 {
			t.Error("Expected attendance with ID 1")
		}
	})

	t.Run("get non-existing - error", func(t *testing.T) {
		_, err := svc.GetByID(ctx, 999)
		if err == nil {
			t.Error("Expected error for non-existing")
		}
	})
}

func TestAttendanceService_ListByEmployee(t *testing.T) {
	mockAtt := NewMockAttendanceRepo()
	svc := NewAttendanceService(mockAtt, NewMockEmployeeRepo())
	ctx := context.Background()

	records, total, err := svc.ListByEmployee(ctx, 1, 1, 10)
	if err != nil {
		t.Errorf("Expected no error but got: %v", err)
	}
	if total != 0 {
		t.Errorf("Expected total 0 but got %d", total)
	}
	if len(records) != 0 {
		t.Errorf("Expected 0 records but got %d", len(records))
	}
}

func TestAttendanceService_List(t *testing.T) {
	mockAtt := NewMockAttendanceRepo()
	for i := 1; i <= 3; i++ {
		mockAtt.records[uint(i)] = &domain.Attendance{ID: uint(i), EmployeeID: uint(i)}
	}
	svc := NewAttendanceService(mockAtt, NewMockEmployeeRepo())
	ctx := context.Background()

	records, total, err := svc.List(ctx, 1, 10)
	if err != nil {
		t.Errorf("Expected no error but got: %v", err)
	}
	if total != 3 {
		t.Errorf("Expected total 3 but got %d", total)
	}
	if len(records) != 3 {
		t.Errorf("Expected 3 records but got %d", len(records))
	}
}

func TestAttendanceService_Update(t *testing.T) {
	t.Run("update with id 0 - should fail", func(t *testing.T) {
		svc := NewAttendanceService(NewMockAttendanceRepo(), NewMockEmployeeRepo())
		ctx := ctxWithClaims(1, []string{"super_admin"})

		err := svc.Update(ctx, &domain.Attendance{ID: 0})
		if err == nil {
			t.Error("Expected error for id 0")
		}
	})

	t.Run("update without claims - should fail", func(t *testing.T) {
		svc := NewAttendanceService(NewMockAttendanceRepo(), NewMockEmployeeRepo())
		ctx := context.Background()

		err := svc.Update(ctx, &domain.Attendance{ID: 1})
		if err == nil {
			t.Error("Expected error for unauthorized")
		}
	})

	t.Run("update with wrong role - should fail", func(t *testing.T) {
		svc := NewAttendanceService(NewMockAttendanceRepo(), NewMockEmployeeRepo())
		ctx := ctxWithClaims(1, []string{"employee"})

		err := svc.Update(ctx, &domain.Attendance{ID: 1})
		if err == nil {
			t.Error("Expected error for wrong role")
		}
	})

	t.Run("update with super_admin - success", func(t *testing.T) {
		mockAtt := NewMockAttendanceRepo()
		mockAtt.records[1] = &domain.Attendance{ID: 1, EmployeeID: 1}
		svc := NewAttendanceService(mockAtt, NewMockEmployeeRepo())
		ctx := ctxWithClaims(1, []string{"super_admin"})

		err := svc.Update(ctx, &domain.Attendance{ID: 1, EmployeeID: 1})
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
	})

	t.Run("update with payroll_manager - success", func(t *testing.T) {
		mockAtt := NewMockAttendanceRepo()
		mockAtt.records[1] = &domain.Attendance{ID: 1, EmployeeID: 1}
		svc := NewAttendanceService(mockAtt, NewMockEmployeeRepo())
		ctx := ctxWithClaims(1, []string{"payroll_manager"})

		err := svc.Update(ctx, &domain.Attendance{ID: 1, EmployeeID: 1})
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
	})
}

func TestAttendanceService_Delete(t *testing.T) {
	t.Run("delete without claims - should fail", func(t *testing.T) {
		svc := NewAttendanceService(NewMockAttendanceRepo(), NewMockEmployeeRepo())
		ctx := context.Background()

		err := svc.Delete(ctx, 1)
		if err == nil {
			t.Error("Expected error for unauthorized")
		}
	})

	t.Run("delete with wrong role - should fail", func(t *testing.T) {
		svc := NewAttendanceService(NewMockAttendanceRepo(), NewMockEmployeeRepo())
		ctx := ctxWithClaims(1, []string{"payroll_manager"})

		err := svc.Delete(ctx, 1)
		if err == nil {
			t.Error("Expected error for wrong role")
		}
	})

	t.Run("delete with super_admin - success", func(t *testing.T) {
		mockAtt := NewMockAttendanceRepo()
		mockAtt.records[1] = &domain.Attendance{ID: 1, EmployeeID: 1}
		svc := NewAttendanceService(mockAtt, NewMockEmployeeRepo())
		ctx := ctxWithClaims(1, []string{"super_admin"})

		err := svc.Delete(ctx, 1)
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
	})
}

func hasAnyRole(roles []string, targets ...string) bool {
	for _, r := range roles {
		for _, t := range targets {
			if r == t {
				return true
			}
		}
	}
	return false
}

func TestHasAnyRole(t *testing.T) {
	tests := []struct {
		name    string
		roles   []string
		targets []string
		want    bool
	}{
		{name: "match found", roles: []string{"admin", "user"}, targets: []string{"admin"}, want: true},
		{name: "no match", roles: []string{"user"}, targets: []string{"admin"}, want: false},
		{name: "empty roles", roles: []string{}, targets: []string{"admin"}, want: false},
		{name: "empty targets", roles: []string{"admin"}, targets: []string{}, want: false},
		{name: "multiple targets - match", roles: []string{"user"}, targets: []string{"admin", "user"}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasAnyRole(tt.roles, tt.targets...)
			if got != tt.want {
				t.Errorf("hasAnyRole(%v, %v) = %v, want %v", tt.roles, tt.targets, got, tt.want)
			}
		})
	}
}

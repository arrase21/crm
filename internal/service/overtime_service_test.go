package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/arrase21/crm/internal/domain"
)

type MockOvertimeRepo struct {
	records       map[uint]*domain.Overtime
	employeeIdx   map[uint][]*domain.Overtime
	nextID        uint
	CreateErr     error
	GetByIDErr    error
	UpdateErr     error
	DeleteErr     error
	ListErr       error
}

func NewMockOvertimeRepo() *MockOvertimeRepo {
	return &MockOvertimeRepo{
		records:     make(map[uint]*domain.Overtime),
		employeeIdx: make(map[uint][]*domain.Overtime),
		nextID:      1,
	}
}

func (m *MockOvertimeRepo) Create(ctx context.Context, o *domain.Overtime) error {
	if m.CreateErr != nil {
		return m.CreateErr
	}
	o.ID = m.nextID
	m.nextID++
	m.records[o.ID] = o
	m.employeeIdx[o.EmployeeID] = append(m.employeeIdx[o.EmployeeID], o)
	return nil
}

func (m *MockOvertimeRepo) GetByID(ctx context.Context, id uint) (*domain.Overtime, error) {
	if m.GetByIDErr != nil {
		return nil, m.GetByIDErr
	}
	if id == 0 {
		return nil, errors.New("not found")
	}
	o, exists := m.records[id]
	if !exists {
		return nil, errors.New("not found")
	}
	return o, nil
}

func (m *MockOvertimeRepo) ListByEmployee(ctx context.Context, employeeID uint, page, limit int) ([]domain.Overtime, int64, error) {
	records := m.employeeIdx[employeeID]
	result := make([]domain.Overtime, len(records))
	for i, o := range records {
		result[i] = *o
	}
	return result, int64(len(records)), nil
}

func (m *MockOvertimeRepo) ListByDepartment(ctx context.Context, departmentID uint, page, limit int) ([]domain.Overtime, int64, error) {
	return nil, 0, nil
}

func (m *MockOvertimeRepo) List(ctx context.Context, page, limit int) ([]domain.Overtime, int64, error) {
	if m.ListErr != nil {
		return nil, 0, m.ListErr
	}
	var all []domain.Overtime
	for _, o := range m.records {
		all = append(all, *o)
	}
	return all, int64(len(all)), nil
}

func (m *MockOvertimeRepo) Update(ctx context.Context, o *domain.Overtime) error {
	if m.UpdateErr != nil {
		return m.UpdateErr
	}
	if o.ID == 0 {
		return errors.New("not found")
	}
	if _, exists := m.records[o.ID]; !exists {
		return errors.New("not found")
	}
	m.records[o.ID] = o
	return nil
}

func (m *MockOvertimeRepo) Delete(ctx context.Context, id uint) error {
	if m.DeleteErr != nil {
		return m.DeleteErr
	}
	if id == 0 {
		return errors.New("not found")
	}
	if _, exists := m.records[id]; !exists {
		return errors.New("not found")
	}
	delete(m.records, id)
	return nil
}

func TestOvertimeService_Create(t *testing.T) {
	t.Run("create valid - success", func(t *testing.T) {
		mockOT := NewMockOvertimeRepo()
		svc := NewOvertimeService(mockOT, NewMockEmployeeRepo())
		ctx := context.Background()

		err := svc.Create(ctx, &domain.Overtime{EmployeeID: 1, Date: time.Now(), Hours: 2})
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
	})

	t.Run("create with employee id 0 - should fail", func(t *testing.T) {
		mockOT := NewMockOvertimeRepo()
		svc := NewOvertimeService(mockOT, NewMockEmployeeRepo())
		ctx := context.Background()

		err := svc.Create(ctx, &domain.Overtime{EmployeeID: 0})
		if err == nil {
			t.Error("Expected error for missing employee")
		}
	})
}

func TestOvertimeService_GetByID(t *testing.T) {
	mockOT := NewMockOvertimeRepo()
	mockOT.records[1] = &domain.Overtime{ID: 1, EmployeeID: 1}
	svc := NewOvertimeService(mockOT, NewMockEmployeeRepo())
	ctx := context.Background()

	t.Run("get existing - success", func(t *testing.T) {
		result, err := svc.GetByID(ctx, 1)
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
		if result == nil || result.ID != 1 {
			t.Error("Expected overtime with ID 1")
		}
	})

	t.Run("get non-existing - error", func(t *testing.T) {
		_, err := svc.GetByID(ctx, 999)
		if err == nil {
			t.Error("Expected error for non-existing")
		}
	})
}

func TestOvertimeService_ListByEmployee(t *testing.T) {
	mockOT := NewMockOvertimeRepo()
	svc := NewOvertimeService(mockOT, NewMockEmployeeRepo())
	ctx := context.Background()

	_, total, err := svc.ListByEmployee(ctx, 1, 1, 10)
	if err != nil {
		t.Errorf("Expected no error but got: %v", err)
	}
	if total != 0 {
		t.Errorf("Expected total 0 but got %d", total)
	}
}

func TestOvertimeService_List(t *testing.T) {
	mockOT := NewMockOvertimeRepo()
	for i := 1; i <= 3; i++ {
		mockOT.records[uint(i)] = &domain.Overtime{ID: uint(i), EmployeeID: uint(i)}
	}
	svc := NewOvertimeService(mockOT, NewMockEmployeeRepo())
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

func TestOvertimeService_Update(t *testing.T) {
	t.Run("update with id 0 - should fail", func(t *testing.T) {
		svc := NewOvertimeService(NewMockOvertimeRepo(), NewMockEmployeeRepo())
		ctx := ctxWithClaims(1, []string{"super_admin"})

		err := svc.Update(ctx, &domain.Overtime{ID: 0})
		if err == nil {
			t.Error("Expected error for id 0")
		}
	})

	t.Run("update without claims - should fail", func(t *testing.T) {
		svc := NewOvertimeService(NewMockOvertimeRepo(), NewMockEmployeeRepo())
		ctx := context.Background()

		err := svc.Update(ctx, &domain.Overtime{ID: 1})
		if err == nil {
			t.Error("Expected error for unauthorized")
		}
	})

	t.Run("update with wrong role - should fail", func(t *testing.T) {
		svc := NewOvertimeService(NewMockOvertimeRepo(), NewMockEmployeeRepo())
		ctx := ctxWithClaims(1, []string{"employee"})

		err := svc.Update(ctx, &domain.Overtime{ID: 1})
		if err == nil {
			t.Error("Expected error for wrong role")
		}
	})

	t.Run("update with super_admin - success", func(t *testing.T) {
		mockOT := NewMockOvertimeRepo()
		mockOT.records[1] = &domain.Overtime{ID: 1, EmployeeID: 1}
		svc := NewOvertimeService(mockOT, NewMockEmployeeRepo())
		ctx := ctxWithClaims(1, []string{"super_admin"})

		err := svc.Update(ctx, &domain.Overtime{ID: 1, EmployeeID: 1})
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
	})

	t.Run("update with payroll_manager - success", func(t *testing.T) {
		mockOT := NewMockOvertimeRepo()
		mockOT.records[1] = &domain.Overtime{ID: 1, EmployeeID: 1}
		svc := NewOvertimeService(mockOT, NewMockEmployeeRepo())
		ctx := ctxWithClaims(1, []string{"payroll_manager"})

		err := svc.Update(ctx, &domain.Overtime{ID: 1, EmployeeID: 1})
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
	})
}

func TestOvertimeService_Delete(t *testing.T) {
	t.Run("delete without claims - should fail", func(t *testing.T) {
		svc := NewOvertimeService(NewMockOvertimeRepo(), NewMockEmployeeRepo())
		ctx := context.Background()

		err := svc.Delete(ctx, 1)
		if err == nil {
			t.Error("Expected error for unauthorized")
		}
	})

	t.Run("delete with wrong role - should fail", func(t *testing.T) {
		svc := NewOvertimeService(NewMockOvertimeRepo(), NewMockEmployeeRepo())
		ctx := ctxWithClaims(1, []string{"payroll_manager"})

		err := svc.Delete(ctx, 1)
		if err == nil {
			t.Error("Expected error for wrong role")
		}
	})

	t.Run("delete with super_admin - success", func(t *testing.T) {
		mockOT := NewMockOvertimeRepo()
		mockOT.records[1] = &domain.Overtime{ID: 1, EmployeeID: 1}
		svc := NewOvertimeService(mockOT, NewMockEmployeeRepo())
		ctx := ctxWithClaims(1, []string{"super_admin"})

		err := svc.Delete(ctx, 1)
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
	})
}

package service

import (
	"context"
	"errors"
	"testing"

	"github.com/arrase21/crm/internal/domain"
)

type MockEmployeeRepo struct {
	employees    map[uint]*domain.Employee
	userIndex    map[uint]*domain.Employee
	nextID       uint
	CreateErr    error
	GetByIDErr   error
	GetByUserIDErr error
	ListErr      error
	ListActiveErr error
	UpdateErr    error
	DeleteErr    error
}

func NewMockEmployeeRepo() *MockEmployeeRepo {
	return &MockEmployeeRepo{
		employees: make(map[uint]*domain.Employee),
		userIndex: make(map[uint]*domain.Employee),
		nextID:    1,
	}
}

func (m *MockEmployeeRepo) Create(ctx context.Context, emp *domain.Employee) error {
	if m.CreateErr != nil {
		return m.CreateErr
	}
	if emp == nil {
		return errors.New("employee cannot be nil")
	}
	emp.ID = m.nextID
	m.nextID++
	m.employees[emp.ID] = emp
	m.userIndex[emp.UserID] = emp
	return nil
}

func (m *MockEmployeeRepo) GetByID(ctx context.Context, id uint) (*domain.Employee, error) {
	if m.GetByIDErr != nil {
		return nil, m.GetByIDErr
	}
	if id == 0 {
		return nil, errors.New("invalid id")
	}
	emp, exists := m.employees[id]
	if !exists {
		return nil, domain.ErrEmployeeNotFound
	}
	return emp, nil
}

func (m *MockEmployeeRepo) GetByUserID(ctx context.Context, userID uint) (*domain.Employee, error) {
	if m.GetByUserIDErr != nil {
		return nil, m.GetByUserIDErr
	}
	if userID == 0 {
		return nil, errors.New("invalid user id")
	}
	emp, exists := m.userIndex[userID]
	if !exists {
		return nil, domain.ErrEmployeeNotFound
	}
	return emp, nil
}

func (m *MockEmployeeRepo) List(ctx context.Context, page, limit int) ([]domain.Employee, int64, error) {
	if m.ListErr != nil {
		return nil, 0, m.ListErr
	}
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	var all []domain.Employee
	for _, e := range m.employees {
		all = append(all, *e)
	}
	offset := (page - 1) * limit
	end := offset + limit
	if offset >= len(all) {
		return []domain.Employee{}, int64(len(all)), nil
	}
	if end > len(all) {
		end = len(all)
	}
	return all[offset:end], int64(len(all)), nil
}

func (m *MockEmployeeRepo) ListActive(ctx context.Context, page, limit int) ([]domain.Employee, int64, error) {
	if m.ListActiveErr != nil {
		return nil, 0, m.ListActiveErr
	}
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	var active []domain.Employee
	for _, e := range m.employees {
		if e.IsActive {
			active = append(active, *e)
		}
	}
	offset := (page - 1) * limit
	end := offset + limit
	if offset >= len(active) {
		return []domain.Employee{}, int64(len(active)), nil
	}
	if end > len(active) {
		end = len(active)
	}
	return active[offset:end], int64(len(active)), nil
}

func (m *MockEmployeeRepo) Update(ctx context.Context, emp *domain.Employee) error {
	if m.UpdateErr != nil {
		return m.UpdateErr
	}
	if emp == nil || emp.ID == 0 {
		return errors.New("employee cannot be nil or have 0")
	}
	if _, exists := m.employees[emp.ID]; !exists {
		return domain.ErrEmployeeNotFound
	}
	delete(m.userIndex, m.employees[emp.ID].UserID)
	m.employees[emp.ID] = emp
	m.userIndex[emp.UserID] = emp
	return nil
}

func (m *MockEmployeeRepo) Delete(ctx context.Context, id uint) error {
	if m.DeleteErr != nil {
		return m.DeleteErr
	}
	if id == 0 {
		return errors.New("invalid employee id")
	}
	emp, exists := m.employees[id]
	if !exists {
		return domain.ErrEmployeeNotFound
	}
	delete(m.userIndex, emp.UserID)
	delete(m.employees, id)
	return nil
}

func (m *MockEmployeeRepo) ListBySupervisor(ctx context.Context, supervisorEmployeeID uint, page, limit int) ([]domain.Employee, int64, error) {
	return nil, 0, nil
}

func validEmployee() *domain.Employee {
	return &domain.Employee{
		UserID:       1,
		DepartmentID: 1,
		PositionID:   1,
		IsActive:     true,
	}
}

func TestEmployeeService_Create(t *testing.T) {
	tests := []struct {
		name       string
		setupMocks func() (*MockEmployeeRepo, *MockUserRepo, *MockDepartmentRepo, *MockPositionRepo)
		input      *domain.Employee
		wantErr    bool
		errType    error
	}{
		{
			name: "create valid employee - success",
			setupMocks: func() (*MockEmployeeRepo, *MockUserRepo, *MockDepartmentRepo, *MockPositionRepo) {
				return NewMockEmployeeRepo(),
					&MockUserRepo{users: map[uint]*domain.User{1: {ID: 1}}, dniIndex: map[string]*domain.User{}, emailIndex: map[string]*domain.User{}, phoneIndex: map[string]*domain.User{}, nextID: 2},
					&MockDepartmentRepo{depts: map[uint]*domain.Department{1: {ID: 1}}, codeIndex: map[string]*domain.Department{}, nameIndex: map[string]*domain.Department{}, nextID: 2},
					&MockPositionRepo{positions: map[uint]*domain.Position{1: {ID: 1}}, nameIndex: map[string]*domain.Position{}, nextID: 2}
			},
			input:   validEmployee(),
			wantErr: false,
		},
		{
			name: "create employee with nil - should fail",
			setupMocks: func() (*MockEmployeeRepo, *MockUserRepo, *MockDepartmentRepo, *MockPositionRepo) {
				return NewMockEmployeeRepo(),
					NewMockUserRepo(),
					NewMockDepartmentRepo(),
					NewMockPositionRepo()
			},
			input:   nil,
			wantErr: true,
		},
		{
			name: "create employee with missing fields - should fail validation",
			setupMocks: func() (*MockEmployeeRepo, *MockUserRepo, *MockDepartmentRepo, *MockPositionRepo) {
				return NewMockEmployeeRepo(),
					NewMockUserRepo(),
					NewMockDepartmentRepo(),
					NewMockPositionRepo()
			},
			input:   &domain.Employee{UserID: 0},
			wantErr: true,
		},
		{
			name: "create employee - user not found",
			setupMocks: func() (*MockEmployeeRepo, *MockUserRepo, *MockDepartmentRepo, *MockPositionRepo) {
				return NewMockEmployeeRepo(),
					&MockUserRepo{GetByIDErr: domain.ErrUserNotFound, users: map[uint]*domain.User{}, dniIndex: map[string]*domain.User{}, emailIndex: map[string]*domain.User{}, phoneIndex: map[string]*domain.User{}},
					NewMockDepartmentRepo(),
					NewMockPositionRepo()
			},
			input:   validEmployee(),
			wantErr: true,
			errType: domain.ErrUserNotFound,
		},
		{
			name: "create employee - department not found",
			setupMocks: func() (*MockEmployeeRepo, *MockUserRepo, *MockDepartmentRepo, *MockPositionRepo) {
				return NewMockEmployeeRepo(),
					&MockUserRepo{users: map[uint]*domain.User{1: {ID: 1}}, dniIndex: map[string]*domain.User{}, emailIndex: map[string]*domain.User{}, phoneIndex: map[string]*domain.User{}, nextID: 2},
					&MockDepartmentRepo{GetByIDErr: domain.ErrDepartmentNotFound, depts: map[uint]*domain.Department{}, codeIndex: map[string]*domain.Department{}, nameIndex: map[string]*domain.Department{}},
					NewMockPositionRepo()
			},
			input:   validEmployee(),
			wantErr: true,
			errType: domain.ErrDepartmentNotFound,
		},
		{
			name: "create employee - position not found",
			setupMocks: func() (*MockEmployeeRepo, *MockUserRepo, *MockDepartmentRepo, *MockPositionRepo) {
				return NewMockEmployeeRepo(),
					&MockUserRepo{users: map[uint]*domain.User{1: {ID: 1}}, dniIndex: map[string]*domain.User{}, emailIndex: map[string]*domain.User{}, phoneIndex: map[string]*domain.User{}, nextID: 2},
					&MockDepartmentRepo{depts: map[uint]*domain.Department{1: {ID: 1}}, codeIndex: map[string]*domain.Department{}, nameIndex: map[string]*domain.Department{}, nextID: 2},
					&MockPositionRepo{GetByIDErr: domain.ErrPositionNotFound, positions: map[uint]*domain.Position{}, nameIndex: map[string]*domain.Position{}}
			},
			input:   validEmployee(),
			wantErr: true,
			errType: domain.ErrPositionNotFound,
		},
		{
			name: "create employee - already exists",
			setupMocks: func() (*MockEmployeeRepo, *MockUserRepo, *MockDepartmentRepo, *MockPositionRepo) {
				empRepo := NewMockEmployeeRepo()
				empRepo.userIndex[1] = &domain.Employee{UserID: 1}
				return empRepo,
					&MockUserRepo{users: map[uint]*domain.User{1: {ID: 1}}, dniIndex: map[string]*domain.User{}, emailIndex: map[string]*domain.User{}, phoneIndex: map[string]*domain.User{}, nextID: 2},
					&MockDepartmentRepo{depts: map[uint]*domain.Department{1: {ID: 1}}, codeIndex: map[string]*domain.Department{}, nameIndex: map[string]*domain.Department{}, nextID: 2},
					&MockPositionRepo{positions: map[uint]*domain.Position{1: {ID: 1}}, nameIndex: map[string]*domain.Position{}, nextID: 2}
			},
			input:   validEmployee(),
			wantErr: true,
			errType: domain.ErrEmployeeAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockEmp, mockUser, mockDept, mockPos := tt.setupMocks()
			svc := NewEmployeeService(mockEmp, mockUser, mockDept, mockPos)
			ctx := context.Background()

			err := svc.Create(ctx, tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got nil")
					return
				}
				if tt.errType != nil && !errors.Is(err, tt.errType) {
					t.Errorf("Expected error '%v' but got '%v'", tt.errType, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestEmployeeService_GetByID(t *testing.T) {
	mockEmp := NewMockEmployeeRepo()
	mockEmp.employees[1] = &domain.Employee{ID: 1, UserID: 1}
	mockEmp.userIndex[1] = mockEmp.employees[1]

	svc := NewEmployeeService(mockEmp, NewMockUserRepo(), NewMockDepartmentRepo(), NewMockPositionRepo())
	ctx := context.Background()

	t.Run("get existing employee - success", func(t *testing.T) {
		emp, err := svc.GetByID(ctx, 1)
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
		if emp == nil || emp.ID != 1 {
			t.Error("Expected employee with ID 1")
		}
	})

	t.Run("get non-existing employee - error", func(t *testing.T) {
		_, err := svc.GetByID(ctx, 999)
		if err == nil {
			t.Error("Expected error for non-existing employee but got nil")
		}
	})

	t.Run("get employee with id 0 - error", func(t *testing.T) {
		_, err := svc.GetByID(ctx, 0)
		if err == nil {
			t.Error("Expected error for id 0 but got nil")
		}
	})
}

func TestEmployeeService_GetByUserID(t *testing.T) {
	mockEmp := NewMockEmployeeRepo()
	mockEmp.userIndex[1] = &domain.Employee{ID: 1, UserID: 1}

	svc := NewEmployeeService(mockEmp, NewMockUserRepo(), NewMockDepartmentRepo(), NewMockPositionRepo())
	ctx := context.Background()

	t.Run("get by user id - success", func(t *testing.T) {
		emp, err := svc.GetByUserID(ctx, 1)
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
		if emp == nil || emp.UserID != 1 {
			t.Error("Expected employee with UserID 1")
		}
	})

	t.Run("get by user id - not found", func(t *testing.T) {
		_, err := svc.GetByUserID(ctx, 999)
		if err == nil {
			t.Error("Expected error for non-existing user id but got nil")
		}
	})

	t.Run("get by user id 0 - error", func(t *testing.T) {
		_, err := svc.GetByUserID(ctx, 0)
		if err == nil {
			t.Error("Expected error for user id 0 but got nil")
		}
	})
}

func TestEmployeeService_List(t *testing.T) {
	mockEmp := NewMockEmployeeRepo()
	for i := 1; i <= 5; i++ {
		e := &domain.Employee{ID: uint(i), UserID: uint(i), IsActive: i%2 == 0}
		mockEmp.employees[uint(i)] = e
		mockEmp.userIndex[uint(i)] = e
	}

	svc := NewEmployeeService(mockEmp, NewMockUserRepo(), NewMockDepartmentRepo(), NewMockPositionRepo())
	ctx := context.Background()

	t.Run("list all employees", func(t *testing.T) {
		emps, total, err := svc.List(ctx, 1, 10)
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
		if total != 5 {
			t.Errorf("Expected total 5 but got %d", total)
		}
		if len(emps) != 5 {
			t.Errorf("Expected 5 employees but got %d", len(emps))
		}
	})

	t.Run("list with pagination", func(t *testing.T) {
		emps, total, err := svc.List(ctx, 1, 2)
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
		if total != 5 {
			t.Errorf("Expected total 5 but got %d", total)
		}
		if len(emps) != 2 {
			t.Errorf("Expected 2 employees but got %d", len(emps))
		}
	})
}

func TestEmployeeService_ListActive(t *testing.T) {
	mockEmp := NewMockEmployeeRepo()
	for i := 1; i <= 5; i++ {
		e := &domain.Employee{ID: uint(i), UserID: uint(i), IsActive: i%2 == 0}
		mockEmp.employees[uint(i)] = e
		mockEmp.userIndex[uint(i)] = e
	}

	svc := NewEmployeeService(mockEmp, NewMockUserRepo(), NewMockDepartmentRepo(), NewMockPositionRepo())
	ctx := context.Background()

	emps, total, err := svc.ListActive(ctx, 1, 10)
	if err != nil {
		t.Errorf("Expected no error but got: %v", err)
	}
	if total != 2 {
		t.Errorf("Expected total 2 active employees but got %d", total)
	}
	if len(emps) != 2 {
		t.Errorf("Expected 2 active employees but got %d", len(emps))
	}
}

func TestEmployeeService_Update(t *testing.T) {
	tests := []struct {
		name       string
		setupMocks func() (*MockEmployeeRepo, *MockUserRepo, *MockDepartmentRepo, *MockPositionRepo)
		input      *domain.Employee
		wantErr    bool
		errType    error
	}{
		{
			name: "update valid employee - success",
			setupMocks: func() (*MockEmployeeRepo, *MockUserRepo, *MockDepartmentRepo, *MockPositionRepo) {
				empRepo := NewMockEmployeeRepo()
				empRepo.employees[1] = &domain.Employee{ID: 1, UserID: 1}
				empRepo.userIndex[1] = empRepo.employees[1]
				return empRepo,
					NewMockUserRepo(),
					&MockDepartmentRepo{depts: map[uint]*domain.Department{2: {ID: 2}}, codeIndex: map[string]*domain.Department{}, nameIndex: map[string]*domain.Department{}, nextID: 3},
					&MockPositionRepo{positions: map[uint]*domain.Position{2: {ID: 2}}, nameIndex: map[string]*domain.Position{}, nextID: 3}
			},
			input:   &domain.Employee{ID: 1, UserID: 1, DepartmentID: 2, PositionID: 2, IsActive: false},
			wantErr: false,
		},
		{
			name: "update with nil - should fail",
			setupMocks: func() (*MockEmployeeRepo, *MockUserRepo, *MockDepartmentRepo, *MockPositionRepo) {
				return NewMockEmployeeRepo(),
					NewMockUserRepo(),
					NewMockDepartmentRepo(),
					NewMockPositionRepo()
			},
			input:   nil,
			wantErr: true,
		},
		{
			name: "update with id 0 - should fail",
			setupMocks: func() (*MockEmployeeRepo, *MockUserRepo, *MockDepartmentRepo, *MockPositionRepo) {
				return NewMockEmployeeRepo(),
					NewMockUserRepo(),
					NewMockDepartmentRepo(),
					NewMockPositionRepo()
			},
			input:   &domain.Employee{ID: 0},
			wantErr: true,
		},
		{
			name: "update non-existing - should fail",
			setupMocks: func() (*MockEmployeeRepo, *MockUserRepo, *MockDepartmentRepo, *MockPositionRepo) {
				return NewMockEmployeeRepo(),
					NewMockUserRepo(),
					NewMockDepartmentRepo(),
					NewMockPositionRepo()
			},
			input:   &domain.Employee{ID: 999, UserID: 1},
			wantErr: true,
			errType: domain.ErrEmployeeNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockEmp, mockUser, mockDept, mockPos := tt.setupMocks()
			svc := NewEmployeeService(mockEmp, mockUser, mockDept, mockPos)
			ctx := context.Background()

			err := svc.Update(ctx, tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got nil")
					return
				}
				if tt.errType != nil && !errors.Is(err, tt.errType) {
					t.Errorf("Expected error '%v' but got '%v'", tt.errType, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestEmployeeService_Delete(t *testing.T) {
	mockEmp := NewMockEmployeeRepo()
	mockEmp.employees[1] = &domain.Employee{ID: 1, UserID: 1}
	mockEmp.userIndex[1] = mockEmp.employees[1]

	svc := NewEmployeeService(mockEmp, NewMockUserRepo(), NewMockDepartmentRepo(), NewMockPositionRepo())
	ctx := context.Background()

	t.Run("delete existing employee - success", func(t *testing.T) {
		err := svc.Delete(ctx, 1)
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
		if _, exists := mockEmp.employees[1]; exists {
			t.Error("Expected employee to be deleted but still exists")
		}
	})

	t.Run("delete non-existing employee - error", func(t *testing.T) {
		err := svc.Delete(ctx, 999)
		if err == nil {
			t.Error("Expected error for non-existing employee but got nil")
		}
	})

	t.Run("delete with id 0 - error", func(t *testing.T) {
		err := svc.Delete(ctx, 0)
		if err == nil {
			t.Error("Expected error for id 0 but got nil")
		}
	})
}

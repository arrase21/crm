package service

import (
	"context"
	"errors"
	"testing"

	"github.com/arrase21/crm/internal/domain"
)

// ============================================
// MOCK IMPLEMENTATION - Simula el Repository
// ============================================

type MockDepartmentRepo struct {
	depts     map[uint]*domain.Department
	nextID    uint
	codeIndex map[string]*domain.Department
	nameIndex map[string]*domain.Department
	// Control para simular errores
	CreateErr    error
	GetByIDErr   error
	GetByCodeErr error
	GetByNameErr error
	UpdateErr    error
	DeleteErr    error
	ListErr      error
}

func NewMockDepartmentRepo() *MockDepartmentRepo {
	return &MockDepartmentRepo{
		depts:     make(map[uint]*domain.Department),
		nextID:    1,
		codeIndex: make(map[string]*domain.Department),
		nameIndex: make(map[string]*domain.Department),
	}
}

func (m *MockDepartmentRepo) Create(ctx context.Context, dept *domain.Department) error {
	if m.CreateErr != nil {
		return m.CreateErr
	}
	if dept == nil {
		return errors.New("department cannot be nil")
	}
	// Check for duplicate code
	if _, exists := m.codeIndex[dept.Code]; exists {
		return domain.ErrDepartmentCodeExists
	}
	// Check for duplicate name
	if _, exists := m.nameIndex[dept.Name]; exists {
		return domain.ErrDepartmentNameExists
	}
	dept.ID = m.nextID
	m.nextID++
	m.depts[dept.ID] = dept
	m.codeIndex[dept.Code] = dept
	m.nameIndex[dept.Name] = dept
	return nil
}

func (m *MockDepartmentRepo) GetByID(ctx context.Context, id uint) (*domain.Department, error) {
	if m.GetByIDErr != nil {
		return nil, m.GetByIDErr
	}
	if id == 0 {
		return nil, errors.New("invalid department id")
	}
	dept, exists := m.depts[id]
	if !exists {
		return nil, domain.ErrDepartmentNotFound
	}
	return dept, nil
}

func (m *MockDepartmentRepo) GetByCode(ctx context.Context, code string) (*domain.Department, error) {
	if m.GetByCodeErr != nil {
		return nil, m.GetByCodeErr
	}
	if code == "" {
		return nil, errors.New("department code cannot be nil")
	}
	dept, exists := m.codeIndex[code]
	if !exists {
		return nil, domain.ErrDepartmentNotFound
	}
	return dept, nil
}

func (m *MockDepartmentRepo) GetByName(ctx context.Context, name string) (*domain.Department, error) {
	if m.GetByNameErr != nil {
		return nil, m.GetByNameErr
	}
	if name == "" {
		return nil, errors.New("department name cannot be nil")
	}
	dept, exists := m.nameIndex[name]
	if !exists {
		return nil, domain.ErrDepartmentNotFound
	}
	return dept, nil
}

func (m *MockDepartmentRepo) List(ctx context.Context, page, limit int) ([]domain.Department, int64, error) {
	if m.ListErr != nil {
		return nil, 0, m.ListErr
	}
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	var allDepts []domain.Department
	for _, d := range m.depts {
		allDepts = append(allDepts, *d)
	}
	offset := (page - 1) * limit
	end := offset + limit
	if offset >= len(allDepts) {
		return []domain.Department{}, int64(len(allDepts)), nil
	}
	if end > len(allDepts) {
		end = len(allDepts)
	}
	return allDepts[offset:end], int64(len(allDepts)), nil
}

func (m *MockDepartmentRepo) Update(ctx context.Context, dept *domain.Department) error {
	if m.UpdateErr != nil {
		return m.UpdateErr
	}
	if dept == nil || dept.ID == 0 {
		return errors.New("department cannot be nil or have zero id")
	}
	if _, exists := m.depts[dept.ID]; !exists {
		return domain.ErrDepartmentNotFound
	}
	// Check if code belongs to another department
	if existing, exists := m.codeIndex[dept.Code]; exists && existing.ID != dept.ID {
		return domain.ErrDepartmentCodeExists
	}
	// Remove from old indexes
	if old, exists := m.depts[dept.ID]; exists {
		delete(m.codeIndex, old.Code)
		delete(m.nameIndex, old.Name)
	}
	m.depts[dept.ID] = dept
	m.codeIndex[dept.Code] = dept
	m.nameIndex[dept.Name] = dept
	return nil
}

func (m *MockDepartmentRepo) Delete(ctx context.Context, id uint) error {
	if m.DeleteErr != nil {
		return m.DeleteErr
	}
	if id == 0 {
		return errors.New("invalid department id")
	}
	dept, exists := m.depts[id]
	if !exists {
		return domain.ErrDepartmentNotFound
	}
	delete(m.codeIndex, dept.Code)
	delete(m.nameIndex, dept.Name)
	delete(m.depts, id)
	return nil
}

// ============================================
// SERVICE TESTS - DepartmentService
// ============================================

func TestDepartmentService_Create(t *testing.T) {
	validDept := func() *domain.Department {
		return &domain.Department{
			Name:     "Accounting",
			Code:     "ACC",
			IsActive: true,
		}
	}

	tests := []struct {
		name      string
		setupMock func(*MockDepartmentRepo)
		input     *domain.Department
		wantErr   bool
		errType   error
	}{
		{
			name:      "create valid department - success",
			setupMock: func(m *MockDepartmentRepo) {},
			input:     validDept(),
			wantErr:   false,
		},
		{
			name:      "create department with nil - should fail",
			setupMock: func(m *MockDepartmentRepo) {},
			input:     nil,
			wantErr:   true,
			errType:   nil,
		},
		{
			name:      "create department with missing fields - should fail validation",
			setupMock: func(m *MockDepartmentRepo) {},
			input: &domain.Department{
				Name: "", // Falta required
			},
			wantErr: true,
		},
		{
			name: "create department with duplicate code - should fail",
			setupMock: func(m *MockDepartmentRepo) {
				m.codeIndex["ACC"] = &domain.Department{ID: 99, Code: "ACC"}
				m.depts[99] = m.codeIndex["ACC"]
			},
			input:   validDept(), // Same code
			wantErr: true,
			errType: domain.ErrDepartmentCodeExists,
		},
		{
			name: "create department with duplicate name - should fail",
			setupMock: func(m *MockDepartmentRepo) {
				m.nameIndex["Accounting"] = &domain.Department{ID: 99, Name: "Accounting"}
				m.depts[99] = m.nameIndex["Accounting"]
			},
			input:   validDept(), // Same name
			wantErr: true,
			errType: domain.ErrDepartmentNameExists,
		},
		{
			name: "create department - repository returns error",
			setupMock: func(m *MockDepartmentRepo) {
				m.CreateErr = errors.New("database connection failed")
			},
			input:   validDept(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := NewMockDepartmentRepo()
			tt.setupMock(mockRepo)
			service := NewDepartmentService(mockRepo)
			ctx := context.Background()

			// Act
			err := service.Create(ctx, tt.input)

			// Assert
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

func TestDepartmentService_GetByID(t *testing.T) {
	mockRepo := NewMockDepartmentRepo()
	existingDept := &domain.Department{
		ID:       1,
		Name:     "Accounting",
		Code:     "ACC",
		IsActive: true,
	}
	mockRepo.depts[1] = existingDept
	mockRepo.codeIndex["ACC"] = existingDept
	mockRepo.nameIndex["Accounting"] = existingDept
	service := NewDepartmentService(mockRepo)
	ctx := context.Background()

	t.Run("get existing department - success", func(t *testing.T) {
		dept, err := service.GetByID(ctx, 1)
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
		if dept == nil {
			t.Error("Expected department but got nil")
			return
		}
		if dept.ID != 1 {
			t.Errorf("Expected department ID 1 but got %d", dept.ID)
		}
	})

	t.Run("get non-existing department - error", func(t *testing.T) {
		_, err := service.GetByID(ctx, 999)
		if err == nil {
			t.Error("Expected error for non-existing department but got nil")
		}
	})

	t.Run("get department with id 0 - error", func(t *testing.T) {
		_, err := service.GetByID(ctx, 0)
		if err == nil {
			t.Error("Expected error for id 0 but got nil")
		}
	})
}

func TestDepartmentService_GetByCode(t *testing.T) {
	mockRepo := NewMockDepartmentRepo()
	existingDept := &domain.Department{
		ID:       1,
		Name:     "Accounting",
		Code:     "ACC",
		IsActive: true,
	}
	mockRepo.depts[1] = existingDept
	mockRepo.codeIndex["ACC"] = existingDept
	mockRepo.nameIndex["Accounting"] = existingDept
	service := NewDepartmentService(mockRepo)
	ctx := context.Background()

	t.Run("get by code - success", func(t *testing.T) {
		dept, err := service.GetByCode(ctx, "ACC")
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
		if dept == nil || dept.ID != 1 {
			t.Errorf("Expected department with ID 1")
		}
	})

	t.Run("get by code - empty string error", func(t *testing.T) {
		_, err := service.GetByCode(ctx, "")
		if err == nil {
			t.Error("Expected error for empty code but got nil")
		}
	})

	t.Run("get by code - not found", func(t *testing.T) {
		_, err := service.GetByCode(ctx, "XXX")
		if err == nil {
			t.Error("Expected error for non-existing code but got nil")
		}
	})
}

func TestDepartmentService_GetByName(t *testing.T) {
	mockRepo := NewMockDepartmentRepo()
	existingDept := &domain.Department{
		ID:       1,
		Name:     "Accounting",
		Code:     "ACC",
		IsActive: true,
	}
	mockRepo.depts[1] = existingDept
	mockRepo.codeIndex["ACC"] = existingDept
	mockRepo.nameIndex["Accounting"] = existingDept
	service := NewDepartmentService(mockRepo)
	ctx := context.Background()

	t.Run("get by name - success", func(t *testing.T) {
		dept, err := service.GetByName(ctx, "Accounting")
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
		if dept == nil || dept.ID != 1 {
			t.Errorf("Expected department with ID 1")
		}
	})

	t.Run("get by name - empty string error", func(t *testing.T) {
		_, err := service.GetByName(ctx, "")
		if err == nil {
			t.Error("Expected error for empty name but got nil")
		}
	})

	t.Run("get by name - not found", func(t *testing.T) {
		_, err := service.GetByName(ctx, "NonExistent")
		if err == nil {
			t.Error("Expected error for non-existing name but got nil")
		}
	})
}

func TestDepartmentService_Update(t *testing.T) {
	validDept := func() *domain.Department {
		return &domain.Department{
			ID:       1,
			Name:     "Accounting",
			Code:     "ACC",
			IsActive: true,
		}
	}

	tests := []struct {
		name      string
		setupMock func(*MockDepartmentRepo)
		input     *domain.Department
		wantErr   bool
		errType   error
	}{
		{
			name: "update valid department - success",
			setupMock: func(m *MockDepartmentRepo) {
				d := validDept()
				m.depts[1] = d
				m.codeIndex["ACC"] = d
				m.nameIndex["Accounting"] = d
			},
			input: func() *domain.Department {
				d := validDept()
				d.Name = "Finance" // Cambiamos el nombre
				return d
			}(),
			wantErr: false,
		},
		{
			name:      "update with nil department - should fail",
			setupMock: func(m *MockDepartmentRepo) {},
			input:     nil,
			wantErr:   true,
		},
		{
			name:      "update with id 0 - should fail",
			setupMock: func(m *MockDepartmentRepo) {},
			input:     &domain.Department{ID: 0},
			wantErr:   true,
		},
		{
			name:      "update non-existing department - should fail",
			setupMock: func(m *MockDepartmentRepo) {},
			input:     validDept(),
			wantErr:   true,
			errType:   domain.ErrDepartmentNotFound,
		},
		{
			name: "update department to use existing code - should fail",
			setupMock: func(m *MockDepartmentRepo) {
				// Dept 1 con código "AAA"
				m.depts[1] = &domain.Department{ID: 1, Code: "AAA", Name: "Dept One"}
				m.codeIndex["AAA"] = m.depts[1]
				m.nameIndex["Dept One"] = m.depts[1]
				// Dept 2 con código "BBB"
				m.depts[2] = &domain.Department{ID: 2, Code: "BBB", Name: "Dept Two"}
				m.codeIndex["BBB"] = m.depts[2]
				m.nameIndex["Dept Two"] = m.depts[2]
			},
			input: &domain.Department{
				ID:       2,
				Name:     "Updated",
				Code:     "AAA", // Ya existe en dept 1
				IsActive: true,
			},
			wantErr: true,
			errType: domain.ErrDepartmentCodeExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockDepartmentRepo()
			tt.setupMock(mockRepo)
			service := NewDepartmentService(mockRepo)
			ctx := context.Background()
			err := service.Update(ctx, tt.input)
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

func TestDepartmentService_Delete(t *testing.T) {
	t.Run("delete existing department - success", func(t *testing.T) {
		mockRepo := NewMockDepartmentRepo()
		dept := &domain.Department{ID: 1, Code: "ACC", Name: "Accounting"}
		mockRepo.depts[1] = dept
		mockRepo.codeIndex["ACC"] = dept
		mockRepo.nameIndex["Accounting"] = dept
		service := NewDepartmentService(mockRepo)
		ctx := context.Background()
		err := service.Delete(ctx, 1)
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
		// Verify department was deleted
		if _, exists := mockRepo.depts[1]; exists {
			t.Error("Expected department to be deleted but still exists")
		}
	})

	t.Run("delete non-existing department - error", func(t *testing.T) {
		mockRepo := NewMockDepartmentRepo()
		service := NewDepartmentService(mockRepo)
		ctx := context.Background()
		err := service.Delete(ctx, 999)
		if err == nil {
			t.Error("Expected error for non-existing department but got nil")
		}
	})

	t.Run("delete with id 0 - error", func(t *testing.T) {
		mockRepo := NewMockDepartmentRepo()
		service := NewDepartmentService(mockRepo)
		ctx := context.Background()
		err := service.Delete(ctx, 0)
		if err == nil {
			t.Error("Expected error for id 0 but got nil")
		}
	})
}

func TestDepartmentService_List(t *testing.T) {
	mockRepo := NewMockDepartmentRepo()
	// Agregar varios departments
	for i := 1; i <= 5; i++ {
		d := &domain.Department{ID: uint(i), Code: string(rune('A' + i)), Name: "Dept"}
		mockRepo.depts[uint(i)] = d
		mockRepo.codeIndex[d.Code] = d
		mockRepo.nameIndex[d.Name] = d
	}
	service := NewDepartmentService(mockRepo)
	ctx := context.Background()

	t.Run("list all departments", func(t *testing.T) {
		depts, total, err := service.List(ctx, 1, 10)
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
		if total != 5 {
			t.Errorf("Expected total 5 but got %d", total)
		}
		if len(depts) != 5 {
			t.Errorf("Expected 5 departments but got %d", len(depts))
		}
	})

	t.Run("list with pagination", func(t *testing.T) {
		depts, total, err := service.List(ctx, 1, 2)
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
		if total != 5 {
			t.Errorf("Expected total 5 but got %d", total)
		}
		if len(depts) != 2 {
			t.Errorf("Expected 2 departments per page but got %d", len(depts))
		}
	})

	t.Run("list second page", func(t *testing.T) {
		depts, total, err := service.List(ctx, 2, 2)
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
		if total != 5 {
			t.Errorf("Expected total 5 but got %d", total)
		}
		if len(depts) != 2 {
			t.Errorf("Expected 2 departments on page 2 but got %d", len(depts))
		}
	})
}

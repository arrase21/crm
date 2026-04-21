package service

import (
	"context"
	"errors"
	"github.com/arrase21/crm/internal/domain"
	"testing"
	"time"
)

// ============================================
// MOCK IMPLEMENTATION - Simula el Repository
// ============================================
type MockUserRepo struct {
	// Storage simula la DB
	users      map[uint]*domain.User
	nextID     uint
	dniIndex   map[string]*domain.User
	emailIndex map[string]*domain.User
	// Control para simular errores
	CreateErr   error
	GetByIDErr  error
	GetByDNIErr error
	UpdateErr   error
	DeleteErr   error
	ListErr     error
}

// Helper para crear un mock con datos iniciales
func NewMockUserRepo() *MockUserRepo {
	return &MockUserRepo{
		users:      make(map[uint]*domain.User),
		nextID:     1,
		dniIndex:   make(map[string]*domain.User),
		emailIndex: make(map[string]*domain.User),
	}
}
func (m *MockUserRepo) Create(ctx context.Context, usr *domain.User) error {
	if m.CreateErr != nil {
		return m.CreateErr
	}
	if usr == nil {
		return errors.New("user cannot be nil")
	}
	// Check for duplicate DNI
	if _, exists := m.dniIndex[usr.Dni]; exists {
		return domain.ErrDniAlreadyExist
	}
	usr.ID = m.nextID
	m.nextID++
	m.users[usr.ID] = usr
	m.dniIndex[usr.Dni] = usr
	m.emailIndex[usr.Email] = usr
	return nil
}
func (m *MockUserRepo) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	if m.GetByIDErr != nil {
		return nil, m.GetByIDErr
	}
	if id == 0 {
		return nil, errors.New("invalid user id")
	}
	user, exists := m.users[id]
	if !exists {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}
func (m *MockUserRepo) GetByDNI(ctx context.Context, dni string) (*domain.User, error) {
	if m.GetByDNIErr != nil {
		return nil, m.GetByDNIErr
	}
	if dni == "" {
		return nil, errors.New("dni cannot be empty")
	}
	user, exists := m.dniIndex[dni]
	if !exists {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}
func (m *MockUserRepo) List(ctx context.Context, page, limit int) ([]domain.User, int64, error) {
	if m.ListErr != nil {
		return nil, 0, m.ListErr
	}
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	var allUsers []domain.User
	for _, u := range m.users {
		allUsers = append(allUsers, *u)
	}
	offset := (page - 1) * limit
	end := offset + limit
	if offset >= len(allUsers) {
		return []domain.User{}, int64(len(allUsers)), nil
	}
	if end > len(allUsers) {
		end = len(allUsers)
	}
	return allUsers[offset:end], int64(len(allUsers)), nil
}
func (m *MockUserRepo) Update(ctx context.Context, usr *domain.User) error {
	if m.UpdateErr != nil {
		return m.UpdateErr
	}
	if usr == nil || usr.ID == 0 {
		return errors.New("user cannot be nil or have zero id")
	}
	if _, exists := m.users[usr.ID]; !exists {
		return domain.ErrUserNotFound
	}
	// Check if DNI belongs to another user
	if existing, exists := m.dniIndex[usr.Dni]; exists && existing.ID != usr.ID {
		return domain.ErrDniAlreadyExist
	}
	// Remove from old indexes
	if old, exists := m.users[usr.ID]; exists {
		delete(m.dniIndex, old.Dni)
		delete(m.emailIndex, old.Email)
	}
	m.users[usr.ID] = usr
	m.dniIndex[usr.Dni] = usr
	m.emailIndex[usr.Email] = usr
	return nil
}
func (m *MockUserRepo) Delete(ctx context.Context, id uint) error {
	if m.DeleteErr != nil {
		return m.DeleteErr
	}
	if id == 0 {
		return errors.New("invalid user id")
	}
	user, exists := m.users[id]
	if !exists {
		return domain.ErrUserNotFound
	}
	delete(m.dniIndex, user.Dni)
	delete(m.emailIndex, user.Email)
	delete(m.users, id)
	return nil
}

// ============================================
// SERVICE TESTS - UserService
// ============================================
func TestUserService_Create(t *testing.T) {
	// Helper para crear usuario válido
	validUser := func() *domain.User {
		return &domain.User{
			FirstName: "Juan",
			LastName:  "Pérez",
			Email:     "juan@example.com",
			Dni:       "12345678",
			Phone:     "1234567890",
			Gender:    "M",
			BirthDay:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		}
	}
	tests := []struct {
		name      string
		setupMock func(*MockUserRepo)
		input     *domain.User
		wantErr   bool
		errType   error // El error específico que esperamos
	}{
		{
			name: "create valid user - success",
			setupMock: func(m *MockUserRepo) {
				// No hay setup especial, repo vacío
			},
			input:   validUser(),
			wantErr: false,
		},
		{
			name:      "create user with nil - should fail",
			setupMock: func(m *MockUserRepo) {},
			input:     nil,
			wantErr:   true,
			errType:   errors.New("user cannot be nil"),
		},
		{
			name:      "create user with missing fields - should fail validation",
			setupMock: func(m *MockUserRepo) {},
			input: &domain.User{
				FirstName: "", // Falta required
			},
			wantErr: true,
		},
		{
			name:      "create user with invalid gender - should fail",
			setupMock: func(m *MockUserRepo) {},
			input: func() *domain.User {
				u := validUser()
				u.Gender = "X"
				return u
			}(),
			wantErr: true,
		},
		{
			name:      "create user that is minor - should fail",
			setupMock: func(m *MockUserRepo) {},
			input: func() *domain.User {
				u := validUser()
				u.BirthDay = time.Now().AddDate(-17, 0, 0)
				return u
			}(),
			wantErr: true,
		},
		{
			name: "create user with duplicate DNI - should fail",
			setupMock: func(m *MockUserRepo) {
				// Pre-populate with existing user
				m.dniIndex["12345678"] = &domain.User{ID: 99, Dni: "12345678"}
				m.users[99] = &domain.User{ID: 99, Dni: "12345678"}
			},
			input:   validUser(), // Same DNI
			wantErr: true,
			errType: domain.ErrDniAlreadyExist,
		},
		{
			name: "create user - repository returns error",
			setupMock: func(m *MockUserRepo) {
				m.CreateErr = errors.New("database connection failed")
			},
			input:   validUser(),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := NewMockUserRepo()
			tt.setupMock(mockRepo)
			service := NewUserService(mockRepo)
			ctx := context.Background()
			// Act
			err := service.Create(ctx, tt.input)
			// Assert
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got nil")
					return
				}
				// Si especificamos errType, verificamos que sea ese error
				if tt.errType != nil && err.Error() != tt.errType.Error() {
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
func TestUserService_GetByID(t *testing.T) {
	mockRepo := NewMockUserRepo()
	// Setup: agregar un usuario al mock
	existingUser := &domain.User{
		ID:        1,
		FirstName: "Juan",
		LastName:  "Pérez",
		Email:     "juan@example.com",
		Dni:       "12345678",
		Phone:     "1234567890",
		Gender:    "M",
		BirthDay:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	mockRepo.users[1] = existingUser
	mockRepo.dniIndex["12345678"] = existingUser
	mockRepo.emailIndex["juan@example.com"] = existingUser
	service := NewUserService(mockRepo)
	ctx := context.Background()
	t.Run("get existing user - success", func(t *testing.T) {
		user, err := service.GetByID(ctx, 1)
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
		if user == nil {
			t.Error("Expected user but got nil")
			return
		}
		if user.ID != 1 {
			t.Errorf("Expected user ID 1 but got %d", user.ID)
		}
	})
	t.Run("get non-existing user - error", func(t *testing.T) {
		_, err := service.GetByID(ctx, 999)
		if err == nil {
			t.Error("Expected error for non-existing user but got nil")
		}
	})
	t.Run("get user with id 0 - error", func(t *testing.T) {
		_, err := service.GetByID(ctx, 0)
		if err == nil {
			t.Error("Expected error for id 0 but got nil")
		}
	})
}
func TestUserService_Update(t *testing.T) {
	validUser := func() *domain.User {
		return &domain.User{
			ID:        1,
			FirstName: "Juan",
			LastName:  "Pérez",
			Email:     "juan@example.com",
			Dni:       "12345678",
			Phone:     "1234567890",
			Gender:    "M",
			BirthDay:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		}
	}
	tests := []struct {
		name      string
		setupMock func(*MockUserRepo)
		input     *domain.User
		wantErr   bool
		errType   error
	}{
		{
			name: "update valid user - success",
			setupMock: func(m *MockUserRepo) {
				u := validUser()
				m.users[1] = u
				m.dniIndex["12345678"] = u
			},
			input: func() *domain.User {
				u := validUser()
				u.FirstName = "Carlos" // Cambiamos el nombre
				return u
			}(),
			wantErr: false,
		},
		{
			name:      "update with nil user - should fail",
			setupMock: func(m *MockUserRepo) {},
			input:     nil,
			wantErr:   true,
			errType:   nil,
		},
		{
			name:      "update with id 0 - should fail",
			setupMock: func(m *MockUserRepo) {},
			input:     &domain.User{ID: 0},
			wantErr:   true,
			errType:   nil,
		},
		{
			name: "update non-existing user - should fail",
			setupMock: func(m *MockUserRepo) {
			},
			input:   validUser(),
			wantErr: true,
			errType: domain.ErrUserNotFound,
		},
		{
			name: "update user to use existing DNI - should fail",
			setupMock: func(m *MockUserRepo) {
				// User 1 con DNI "11111111"
				m.users[1] = &domain.User{ID: 1, Dni: "11111111", FirstName: "User", LastName: "One", Email: "user1@test.com"}
				// m.users[1] = &domain.User{ID: 1, Dni: "11111111"}
				m.dniIndex["11111111"] = m.users[1]
				// User 2 con DNI "22222222"
				// m.users[2] = &domain.User{ID: 2, Dni: "22222222"}
				m.users[2] = &domain.User{ID: 2, Dni: "22222222", FirstName: "User", LastName: "Two", Email: "user2@test.com"}
				m.dniIndex["22222222"] = m.users[2]
			},
			input: &domain.User{
				ID:        2,
				FirstName: "Updated",
				LastName:  "Name",
				Email:     "updated@test.com",
				Dni:       "11111111", // Ya existe en user 1
				Phone:     "123456",
				Gender:    "M",
				BirthDay:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			wantErr: true,
			errType: domain.ErrDniAlreadyExist,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockUserRepo()
			tt.setupMock(mockRepo)
			service := NewUserService(mockRepo)
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
func TestUserService_Delete(t *testing.T) {
	t.Run("delete existing user - success", func(t *testing.T) {
		mockRepo := NewMockUserRepo()
		mockRepo.users[1] = &domain.User{ID: 1, Dni: "12345678"}
		mockRepo.dniIndex["12345678"] = mockRepo.users[1]
		service := NewUserService(mockRepo)
		ctx := context.Background()
		err := service.Delete(ctx, 1)
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
		// Verify user was deleted
		if _, exists := mockRepo.users[1]; exists {
			t.Error("Expected user to be deleted but still exists")
		}
	})
	t.Run("delete non-existing user - error", func(t *testing.T) {
		mockRepo := NewMockUserRepo()
		service := NewUserService(mockRepo)
		ctx := context.Background()
		err := service.Delete(ctx, 999)
		if err == nil {
			t.Error("Expected error for non-existing user but got nil")
		}
	})
	t.Run("delete with id 0 - error", func(t *testing.T) {
		mockRepo := NewMockUserRepo()
		service := NewUserService(mockRepo)
		ctx := context.Background()
		err := service.Delete(ctx, 0)
		if err == nil {
			t.Error("Expected error for id 0 but got nil")
		}
	})
}
func TestUserService_GetByDNI(t *testing.T) {
	mockRepo := NewMockUserRepo()
	mockRepo.users[1] = &domain.User{ID: 1, Dni: "12345678"}
	mockRepo.dniIndex["12345678"] = mockRepo.users[1]
	service := NewUserService(mockRepo)
	ctx := context.Background()
	t.Run("get by DNI - success", func(t *testing.T) {
		user, err := service.GetByDNI(ctx, "12345678")
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
		if user == nil || user.ID != 1 {
			t.Errorf("Expected user with ID 1")
		}
	})
	t.Run("get by DNI - empty string error", func(t *testing.T) {
		_, err := service.GetByDNI(ctx, "")
		if err == nil {
			t.Error("Expected error for empty DNI but got nil")
		}
	})
	t.Run("get by DNI - not found", func(t *testing.T) {
		_, err := service.GetByDNI(ctx, "99999999")
		if err == nil {
			t.Error("Expected error for non-existing DNI but got nil")
		}
	})
}
func TestUserService_List(t *testing.T) {
	mockRepo := NewMockUserRepo()
	// Agregar varios usuarios
	for i := 1; i <= 5; i++ {
		u := &domain.User{ID: uint(i), Dni: string(rune('0' + i))}
		mockRepo.users[uint(i)] = u
	}
	service := NewUserService(mockRepo)
	ctx := context.Background()
	t.Run("list all users", func(t *testing.T) {
		users, total, err := service.List(ctx, 1, 10)
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
		if total != 5 {
			t.Errorf("Expected total 5 but got %d", total)
		}
		if len(users) != 5 {
			t.Errorf("Expected 5 users but got %d", len(users))
		}
	})
	t.Run("list with pagination", func(t *testing.T) {
		users, total, err := service.List(ctx, 1, 2)
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
		if total != 5 {
			t.Errorf("Expected total 5 but got %d", total)
		}
		if len(users) != 2 {
			t.Errorf("Expected 2 users per page but got %d", len(users))
		}
	})
	t.Run("list second page", func(t *testing.T) {
		users, total, err := service.List(ctx, 2, 2)
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
		if total != 5 {
			t.Errorf("Expected total 5 but got %d", total)
		}
		if len(users) != 2 {
			t.Errorf("Expected 2 users on page 2 but got %d", len(users))
		}
	})
}

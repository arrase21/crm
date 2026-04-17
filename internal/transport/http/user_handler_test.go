package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/arrase21/crm/internal/domain"
	"github.com/arrase21/crm/internal/service"
	"github.com/gin-gonic/gin"
)

// ============================================
// MOCK SERVICE - Simula el Service para tests HTTP
// ============================================
type MockUserService struct {
	users       map[uint]*domain.User
	nextID      uint
	dniIndex    map[string]*domain.User
	createErr   error
	getByIDErr  error
	getByDNIErr error
	updateErr   error
	deleteErr   error
	listErr     error
}

func NewMockUserService() *MockUserService {
	return &MockUserService{
		users:    make(map[uint]*domain.User),
		nextID:   1,
		dniIndex: make(map[string]*domain.User),
	}
}

func (m *MockUserService) Create(ctx context.Context, usr *domain.User) error {
	if m.createErr != nil {
		return m.createErr
	}
	if usr == nil {
		return domain.ErrUserNotFound
	}
	if _, exists := m.dniIndex[usr.Dni]; exists {
		return domain.ErrDniAlreadyExist
	}
	usr.ID = m.nextID
	m.nextID++
	m.users[usr.ID] = usr
	m.dniIndex[usr.Dni] = usr
	return nil
}

func (m *MockUserService) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}
	if id == 0 {
		return nil, domain.ErrUserNotFound
	}
	user, exists := m.users[id]
	if !exists {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}

func (m *MockUserService) GetByDNI(ctx context.Context, dni string) (*domain.User, error) {
	if m.getByDNIErr != nil {
		return nil, m.getByDNIErr
	}
	if dni == "" {
		return nil, domain.ErrUserNotFound
	}
	user, exists := m.dniIndex[dni]
	if !exists {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}

func (m *MockUserService) List(ctx context.Context, page, limit int) ([]domain.User, int64, error) {
	if m.listErr != nil {
		return nil, 0, m.listErr
	}
	var allUsers []domain.User
	for _, u := range m.users {
		allUsers = append(allUsers, *u)
	}
	return allUsers, int64(len(allUsers)), nil
}

func (m *MockUserService) Update(ctx context.Context, usr *domain.User) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	if usr == nil || usr.ID == 0 {
		return domain.ErrUserNotFound
	}
	if _, exists := m.users[usr.ID]; !exists {
		return domain.ErrUserNotFound
	}
	if existing, exists := m.dniIndex[usr.Dni]; exists && existing.ID != usr.ID {
		return domain.ErrDniAlreadyExist
	}
	m.users[usr.ID] = usr
	delete(m.dniIndex, usr.Dni)
	m.dniIndex[usr.Dni] = usr
	return nil
}

func (m *MockUserService) Delete(ctx context.Context, id uint) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	if id == 0 {
		return domain.ErrUserNotFound
	}
	user, exists := m.users[id]
	if !exists {
		return domain.ErrUserNotFound
	}
	delete(m.dniIndex, user.Dni)
	delete(m.users, id)
	return nil
}

// Helper para setup del router con mock
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	return r
}

// ============================================
// HTTP TESTS - User Handler
// ============================================

// --- CREATE TESTS ---

func TestUserHandler_Create_Success(t *testing.T) {
	mockSvc := NewMockUserService()
	handler := NewUserHandler(service.NewUserService(mockSvc))
	router := setupTestRouter()

	router.POST("/users", handler.Create)

	body := map[string]interface{}{
		"first_name": "Juan",
		"last_name":  "Pérez",
		"email":      "juan@test.com",
		"dni":        "12345678",
		"phone":      "1234567890",
		"gender":     "M",
		"birth_day":  "1990-01-01",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201 but got %d", w.Code)
	}
}

func TestUserHandler_Create_InvalidJSON(t *testing.T) {
	mockSvc := NewMockUserService()
	handler := NewUserHandler(service.NewUserService(mockSvc))
	router := setupTestRouter()

	router.POST("/users", handler.Create)

	req, _ := http.NewRequest("POST", "/users", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 but got %d", w.Code)
	}
}

func TestUserHandler_Create_MissingFields(t *testing.T) {
	mockSvc := NewMockUserService()
	handler := NewUserHandler(service.NewUserService(mockSvc))
	router := setupTestRouter()

	router.POST("/users", handler.Create)

	// Missing required fields
	body := map[string]interface{}{
		"first_name": "Juan",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 but got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestUserHandler_Create_InvalidGender(t *testing.T) {
	mockSvc := NewMockUserService()
	handler := NewUserHandler(service.NewUserService(mockSvc))
	router := setupTestRouter()

	router.POST("/users", handler.Create)

	body := map[string]interface{}{
		"first_name": "Juan",
		"last_name":  "Pérez",
		"email":      "juan@test.com",
		"dni":        "12345678",
		"phone":      "1234567890",
		"gender":     "X", // Invalid gender
		"birth_day":  "1990-01-01",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 but got %d", w.Code)
	}
}

func TestUserHandler_Create_DuplicateDNI(t *testing.T) {
	mockSvc := NewMockUserService()
	// Pre-create user with same DNI
	mockSvc.dniIndex["12345678"] = &domain.User{ID: 1, Dni: "12345678"}
	mockSvc.users[1] = &domain.User{ID: 1, Dni: "12345678"}

	handler := NewUserHandler(service.NewUserService(mockSvc))
	router := setupTestRouter()

	router.POST("/users", handler.Create)

	body := map[string]interface{}{
		"first_name": "Pedro",
		"last_name":  "Gómez",
		"email":      "pedro@test.com",
		"dni":        "12345678", // Duplicate
		"phone":      "9999999999",
		"gender":     "M",
		"birth_day":  "1990-01-01",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("Expected status 409 but got %d", w.Code)
	}
}

// --- GET BY ID TESTS ---

func TestUserHandler_GetByID_Success(t *testing.T) {
	mockSvc := NewMockUserService()
	// Create user first
	mockSvc.users[1] = &domain.User{
		ID:        1,
		FirstName: "Juan",
		LastName:  "Pérez",
		Email:     "juan@test.com",
		Dni:       "12345678",
		Phone:     "1234567890",
		Gender:    "M",
		BirthDay:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	mockSvc.dniIndex["12345678"] = mockSvc.users[1]

	handler := NewUserHandler(service.NewUserService(mockSvc))
	router := setupTestRouter()

	router.GET("/users/:id", handler.GetByID)

	req, _ := http.NewRequest("GET", "/users/1", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 but got %d", w.Code)
	}
}

func TestUserHandler_GetByID_InvalidID(t *testing.T) {
	mockSvc := NewMockUserService()
	handler := NewUserHandler(service.NewUserService(mockSvc))
	router := setupTestRouter()

	router.GET("/users/:id", handler.GetByID)

	req, _ := http.NewRequest("GET", "/users/abc", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 but got %d", w.Code)
	}
}

func TestUserHandler_GetByID_NotFound(t *testing.T) {
	mockSvc := NewMockUserService()
	handler := NewUserHandler(service.NewUserService(mockSvc))
	router := setupTestRouter()

	router.GET("/users/:id", handler.GetByID)

	req, _ := http.NewRequest("GET", "/users/999", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 but got %d", w.Code)
	}
}

// --- GET BY DNI TESTS ---

func TestUserHandler_GetByDni_Success(t *testing.T) {
	mockSvc := NewMockUserService()
	mockSvc.users[1] = &domain.User{
		ID:        1,
		FirstName: "Juan",
		Dni:       "12345678",
	}
	mockSvc.dniIndex["12345678"] = mockSvc.users[1]

	handler := NewUserHandler(service.NewUserService(mockSvc))
	router := setupTestRouter()

	router.GET("/users/dni", handler.GetByDni)

	req, _ := http.NewRequest("GET", "/users/dni?dni=12345678", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 but got %d", w.Code)
	}
}

func TestUserHandler_GetByDni_MissingParam(t *testing.T) {
	mockSvc := NewMockUserService()
	handler := NewUserHandler(service.NewUserService(mockSvc))
	router := setupTestRouter()

	router.GET("/users/dni", handler.GetByDni)

	req, _ := http.NewRequest("GET", "/users/dni", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 but got %d", w.Code)
	}
}

// --- LIST TESTS ---

func TestUserHandler_List_Success(t *testing.T) {
	mockSvc := NewMockUserService()
	// Add some users
	for i := 1; i <= 3; i++ {
		mockSvc.users[uint(i)] = &domain.User{ID: uint(i), FirstName: "User"}
	}

	handler := NewUserHandler(service.NewUserService(mockSvc))
	router := setupTestRouter()

	router.GET("/users", handler.List)

	req, _ := http.NewRequest("GET", "/users?page=1&limit=10", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 but got %d", w.Code)
	}
}

// --- UPDATE TESTS ---

func TestUserHandler_Update_Success(t *testing.T) {
	mockSvc := NewMockUserService()
	mockSvc.users[1] = &domain.User{
		ID:        1,
		FirstName: "Juan",
		LastName:  "Pérez",
		Email:     "juan@test.com",
		Dni:       "12345678",
		Phone:     "1234567890",
		Gender:    "M",
		BirthDay:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	mockSvc.dniIndex["12345678"] = mockSvc.users[1]

	handler := NewUserHandler(service.NewUserService(mockSvc))
	router := setupTestRouter()

	router.PUT("/users/:id", handler.Update)

	body := map[string]interface{}{
		"first_name": "Carlos",
		"last_name":  "Pérez",
		"email":      "juan@test.com",
		"dni":        "12345678",
		"phone":      "1234567890",
		"gender":     "M",
		"birth_day":  "1990-01-01",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("PUT", "/users/1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 but got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestUserHandler_Update_NotFound(t *testing.T) {
	mockSvc := NewMockUserService()
	handler := NewUserHandler(service.NewUserService(mockSvc))
	router := setupTestRouter()

	router.PUT("/users/:id", handler.Update)

	body := map[string]interface{}{
		"first_name": "Carlos",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("PUT", "/users/999", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 but got %d", w.Code)
	}
}

func TestUserHandler_Update_InvalidID(t *testing.T) {
	mockSvc := NewMockUserService()
	handler := NewUserHandler(service.NewUserService(mockSvc))
	router := setupTestRouter()

	router.PUT("/users/:id", handler.Update)

	body := map[string]interface{}{
		"first_name": "Carlos",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("PUT", "/users/abc", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 but got %d", w.Code)
	}
}

// --- DELETE TESTS ---

func TestUserHandler_Delete_Success(t *testing.T) {
	mockSvc := NewMockUserService()
	mockSvc.users[1] = &domain.User{ID: 1, Dni: "12345678"}
	mockSvc.dniIndex["12345678"] = mockSvc.users[1]

	handler := NewUserHandler(service.NewUserService(mockSvc))
	router := setupTestRouter()

	router.DELETE("/users/:id", handler.Delete)

	req, _ := http.NewRequest("DELETE", "/users/1", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204 but got %d", w.Code)
	}
}

func TestUserHandler_Delete_NotFound(t *testing.T) {
	mockSvc := NewMockUserService()
	handler := NewUserHandler(service.NewUserService(mockSvc))
	router := setupTestRouter()

	router.DELETE("/users/:id", handler.Delete)

	req, _ := http.NewRequest("DELETE", "/users/999", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 but got %d", w.Code)
	}
}

func TestUserHandler_Delete_InvalidID(t *testing.T) {
	mockSvc := NewMockUserService()
	handler := NewUserHandler(service.NewUserService(mockSvc))
	router := setupTestRouter()

	router.DELETE("/users/:id", handler.Delete)

	req, _ := http.NewRequest("DELETE", "/users/abc", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 but got %d", w.Code)
	}
}

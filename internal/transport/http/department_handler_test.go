package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arrase21/crm/internal/domain"
	"github.com/arrase21/crm/internal/service"
	"github.com/gin-gonic/gin"
)

type MockDepartmentService struct {
	departments  map[uint]*domain.Department
	nextID       uint
	codeIndex    map[string]*domain.Department
	nameIndex    map[string]*domain.Department
	createErr    error
	getByIDErr   error
	getByCodeErr error
	getByNameErr error
	updateErr    error
	deleteErr    error
	listErr      error
}

func NewMockDepartmentService() *MockDepartmentService {
	return &MockDepartmentService{
		departments: make(map[uint]*domain.Department),
		nextID:      1,
		codeIndex:   make(map[string]*domain.Department),
		nameIndex:   make(map[string]*domain.Department),
	}
}

func (m *MockDepartmentService) Create(ctx context.Context, dept *domain.Department) error {
	if m.createErr != nil {
		return m.createErr
	}
	if dept == nil {
		return domain.ErrDepartmentNotFound
	}
	if _, exists := m.codeIndex[dept.Code]; exists {
		return domain.ErrDepartmentCodeExists
	}
	dept.ID = m.nextID
	m.nextID++
	m.departments[dept.ID] = dept
	m.codeIndex[dept.Code] = dept
	return nil
}

func (m *MockDepartmentService) GetByID(ctx context.Context, id uint) (*domain.Department, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}
	if id == 0 {
		return nil, domain.ErrUserNotFound
	}
	department, exists := m.departments[id]
	if !exists {
		return nil, domain.ErrDepartmentNotFound
	}
	return department, nil
}
func (m *MockDepartmentService) GetByCode(ctx context.Context, code string) (*domain.Department, error) {
	if m.getByCodeErr != nil {
		return nil, m.getByCodeErr
	}
	if code == "" {
		return nil, domain.ErrDepartmentNotFound
	}
	dept, exists := m.codeIndex[code]
	if !exists {
		return nil, domain.ErrDepartmentNotFound
	}
	return dept, nil
}

func (m *MockDepartmentService) GetByName(ctx context.Context, name string) (*domain.Department, error) {
	if m.getByNameErr != nil {
		return nil, m.getByNameErr
	}
	if name == "" {
		return nil, domain.ErrDepartmentNotFound
	}
	dept, exits := m.nameIndex[name]
	if !exits {
		return nil, domain.ErrDepartmentNotFound
	}
	return dept, nil
}

func (m *MockDepartmentService) List(ctx context.Context, page, limit int) ([]domain.Department, int64, error) {
	if m.listErr != nil {
		return nil, 0, m.listErr
	}
	var allDepartments []domain.Department
	for _, u := range m.departments {
		allDepartments = append(allDepartments, *u)
	}
	return allDepartments, int64(len(allDepartments)), nil
}

func (m *MockDepartmentService) Update(ctx context.Context, dept *domain.Department) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	if dept == nil || dept.ID == 0 {
		return domain.ErrDepartmentNotFound
	}
	if _, exists := m.departments[dept.ID]; !exists {
		return domain.ErrDepartmentNotFound
	}
	if existing, exists := m.codeIndex[dept.Code]; exists && existing.ID != dept.ID {
		return domain.ErrDniAlreadyExist
	}
	m.departments[dept.ID] = dept
	delete(m.codeIndex, dept.Code)
	m.codeIndex[dept.Code] = dept
	return nil
}

func (m *MockDepartmentService) Delete(ctx context.Context, id uint) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	if id == 0 {
		return domain.ErrDepartmentNotFound
	}
	dept, exists := m.departments[id]
	if !exists {
		return domain.ErrDepartmentNotFound
	}
	delete(m.codeIndex, dept.Code)
	delete(m.departments, id)
	return nil
}

// Helper test
func setupTestRouterDept() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	return r
}

func TestDepartmentHandler_Create_Success(t *testing.T) {
	mockSvc := NewMockDepartmentService()
	handler := NewDepartmentHandler(service.NewDepartmentService(mockSvc))
	router := setupTestRouterDept()
	router.POST("/departments", handler.Create)
	body := map[string]interface{}{
		"name":      "accounting",
		"code":      "acc",
		"is_active": true,
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/departments", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201 but got %d", w.Code)
	}
}

func TestDepartmentHandler_Create_InvalidJSON(t *testing.T) {
	mockSvc := NewMockDepartmentService()
	handler := NewDepartmentHandler(service.NewDepartmentService(mockSvc))
	router := setupTestRouter()

	router.POST("/departments", handler.Create)

	req, _ := http.NewRequest("POST", "/departments", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 but got %d", w.Code)
	}
}

func TestDepartmentHandler_Delete_Success(t *testing.T) {
	mockSvc := NewMockDepartmentService()
	mockSvc.departments[1] = &domain.Department{ID: 1, Code: "acc"}
	mockSvc.codeIndex["acc"] = mockSvc.departments[1]

	handler := NewDepartmentHandler(service.NewDepartmentService(mockSvc))
	router := setupTestRouter()

	router.DELETE("/departments/:id", handler.Delete)

	req, _ := http.NewRequest("DELETE", "/departments/1", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204 but got %d", w.Code)
	}
}

func TestDepartmentHandler_Delete_InvalidID(t *testing.T) {
	mockSvc := NewMockDepartmentService()
	handler := NewDepartmentHandler(service.NewDepartmentService(mockSvc))
	router := setupTestRouter()

	router.DELETE("/departments/:id", handler.Delete)

	req, _ := http.NewRequest("DELETE", "/departments/abc", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 but got %d", w.Code)
	}
}

// ============================================
// GET BY ID TESTS
// ============================================

func TestDepartmentHandler_GetByID_Success(t *testing.T) {
	mockSvc := NewMockDepartmentService()
	mockSvc.departments[1] = &domain.Department{
		ID:   1,
		Name: "Accounting",
		Code: "ACC",
	}
	mockSvc.codeIndex["ACC"] = mockSvc.departments[1]

	handler := NewDepartmentHandler(service.NewDepartmentService(mockSvc))
	router := setupTestRouter()

	router.GET("/departments/:id", handler.GetByID)

	req, _ := http.NewRequest("GET", "/departments/1", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 but got %d", w.Code)
	}
}

func TestDepartmentHandler_GetByID_InvalidID(t *testing.T) {
	mockSvc := NewMockDepartmentService()
	handler := NewDepartmentHandler(service.NewDepartmentService(mockSvc))
	router := setupTestRouter()

	router.GET("/departments/:id", handler.GetByID)

	req, _ := http.NewRequest("GET", "/departments/abc", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 but got %d", w.Code)
	}
}

func TestDepartmentHandler_GetByID_NotFound(t *testing.T) {
	mockSvc := NewMockDepartmentService()
	handler := NewDepartmentHandler(service.NewDepartmentService(mockSvc))
	router := setupTestRouter()

	router.GET("/departments/:id", handler.GetByID)

	req, _ := http.NewRequest("GET", "/departments/999", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 but got %d", w.Code)
	}
}

// ============================================
// GET BY CODE TESTS
// ============================================

func TestDepartmentHandler_GetByCode_Success(t *testing.T) {
	mockSvc := NewMockDepartmentService()
	mockSvc.departments[1] = &domain.Department{
		ID:   1,
		Name: "Accounting",
		Code: "ACC",
	}
	mockSvc.codeIndex["ACC"] = mockSvc.departments[1]

	handler := NewDepartmentHandler(service.NewDepartmentService(mockSvc))
	router := setupTestRouter()

	router.GET("/departments", handler.GetByCode)

	req, _ := http.NewRequest("GET", "/departments?code=ACC", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 but got %d", w.Code)
	}
}

func TestDepartmentHandler_GetByCode_MissingQueryParam(t *testing.T) {
	mockSvc := NewMockDepartmentService()
	handler := NewDepartmentHandler(service.NewDepartmentService(mockSvc))
	router := setupTestRouter()

	router.GET("/departments", handler.GetByCode)

	req, _ := http.NewRequest("GET", "/departments", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 but got %d", w.Code)
	}
}

func TestDepartmentHandler_GetByCode_NotFound(t *testing.T) {
	mockSvc := NewMockDepartmentService()
	handler := NewDepartmentHandler(service.NewDepartmentService(mockSvc))
	router := setupTestRouter()

	router.GET("/departments", handler.GetByCode)

	req, _ := http.NewRequest("GET", "/departments?code=NONEXISTENT", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 but got %d", w.Code)
	}
}

// ============================================
// GET BY NAME TESTS
// ============================================

func TestDepartmentHandler_GetByName_Success(t *testing.T) {
	mockSvc := NewMockDepartmentService()
	mockSvc.departments[1] = &domain.Department{
		ID:   1,
		Name: "Accounting",
		Code: "ACC",
	}
	mockSvc.nameIndex["Accounting"] = mockSvc.departments[1]

	handler := NewDepartmentHandler(service.NewDepartmentService(mockSvc))
	router := setupTestRouter()

	router.GET("/departments", handler.GetByName)

	req, _ := http.NewRequest("GET", "/departments?name=Accounting", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 but got %d", w.Code)
	}
}

func TestDepartmentHandler_GetByName_MissingQueryParam(t *testing.T) {
	mockSvc := NewMockDepartmentService()
	handler := NewDepartmentHandler(service.NewDepartmentService(mockSvc))
	router := setupTestRouter()

	router.GET("/departments", handler.GetByName)

	req, _ := http.NewRequest("GET", "/departments", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 but got %d", w.Code)
	}
}

// ============================================
// LIST TESTS
// ============================================

func TestDepartmentHandler_List_Success(t *testing.T) {
	mockSvc := NewMockDepartmentService()
	// Agregar departments
	for i := 1; i <= 3; i++ {
		d := &domain.Department{
			ID:   uint(i),
			Name: "Dept" + string(rune('0'+i)),
			Code: string(rune('A' + i)),
		}
		mockSvc.departments[uint(i)] = d
		mockSvc.codeIndex[d.Code] = d
	}

	handler := NewDepartmentHandler(service.NewDepartmentService(mockSvc))
	router := setupTestRouter()

	router.GET("/departments", handler.List)

	req, _ := http.NewRequest("GET", "/departments?page=1&limit=10", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 but got %d", w.Code)
	}
}

func TestDepartmentHandler_List_DefaultPagination(t *testing.T) {
	mockSvc := NewMockDepartmentService()
	handler := NewDepartmentHandler(service.NewDepartmentService(mockSvc))
	router := setupTestRouter()

	router.GET("/departments", handler.List)

	req, _ := http.NewRequest("GET", "/departments", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 but got %d", w.Code)
	}
}

// ============================================
// UPDATE TESTS
// ============================================

func TestDepartmentHandler_Update_Success(t *testing.T) {
	mockSvc := NewMockDepartmentService()
	mockSvc.departments[1] = &domain.Department{
		ID:   1,
		Name: "Accounting",
		Code: "ACC",
	}
	mockSvc.codeIndex["ACC"] = mockSvc.departments[1]

	handler := NewDepartmentHandler(service.NewDepartmentService(mockSvc))
	router := setupTestRouter()

	router.PUT("/departments/:id", handler.Update)
	body := map[string]interface{}{
		"name": "Finance",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("PUT", "/departments/1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 but got %d", w.Code)
	}
}

func TestDepartmentHandler_Update_InvalidID(t *testing.T) {
	mockSvc := NewMockDepartmentService()
	handler := NewDepartmentHandler(service.NewDepartmentService(mockSvc))
	router := setupTestRouter()

	router.PUT("/departments/:id", handler.Update)
	body := map[string]interface{}{
		"name": "Finance",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("PUT", "/departments/abc", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 but got %d", w.Code)
	}
}

func TestDepartmentHandler_Update_NotFound(t *testing.T) {
	mockSvc := NewMockDepartmentService()
	handler := NewDepartmentHandler(service.NewDepartmentService(mockSvc))
	router := setupTestRouter()

	router.PUT("/departments/:id", handler.Update)
	body := map[string]interface{}{
		"name": "Finance",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("PUT", "/departments/999", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 but got %d", w.Code)
	}
}

// ============================================
// CREATE ERROR CASES
// ============================================

func TestDepartmentHandler_Create_DuplicateCode(t *testing.T) {
	mockSvc := NewMockDepartmentService()
	mockSvc.codeIndex["ACC"] = &domain.Department{ID: 1, Code: "ACC"}

	handler := NewDepartmentHandler(service.NewDepartmentService(mockSvc))
	router := setupTestRouter()

	router.POST("/departments", handler.Create)
	body := map[string]interface{}{
		"name": "Accounting",
		"code": "ACC",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/departments", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("Expected status 409 but got %d", w.Code)
	}
}

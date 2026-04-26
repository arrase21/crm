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

// ============================================
// MOCK SERVICE - Position
// ============================================

type MockPositionService struct {
	positions    map[uint]*domain.Position
	nextID       uint
	nameIndex    map[string]*domain.Position
	createErr    error
	getByIDErr   error
	getByNameErr error
	updateErr    error
	deleteErr    error
	listErr      error
}

// Add new methods to implement domain.PositionRepo interface
func (m *MockPositionService) GetByIDWithDepartment(ctx context.Context, id uint) (*domain.Position, error) {
	return m.GetByID(ctx, id)
}

func (m *MockPositionService) ListByDepartment(ctx context.Context, deptID uint) ([]domain.Position, int64, error) {
	var result []domain.Position
	for _, p := range m.positions {
		if p.DepartmentID == deptID {
			result = append(result, *p)
		}
	}
	return result, int64(len(result)), nil
}

func (m *MockPositionService) CountByDepartment(ctx context.Context, deptID uint) (int64, error) {
	var count int64
	for _, p := range m.positions {
		if p.DepartmentID == deptID {
			count++
		}
	}
	return count, nil
}

func NewMockPositionService() *MockPositionService {
	return &MockPositionService{
		positions: make(map[uint]*domain.Position),
		nextID:    1,
		nameIndex: make(map[string]*domain.Position),
	}
}

func (m *MockPositionService) Create(ctx context.Context, pos *domain.Position) error {
	if m.createErr != nil {
		return m.createErr
	}
	if pos == nil {
		return domain.ErrPositionNotFound
	}
	if _, exists := m.nameIndex[pos.Name]; exists {
		return domain.ErrPositionNameExists
	}
	pos.ID = m.nextID
	m.nextID++
	m.positions[pos.ID] = pos
	m.nameIndex[pos.Name] = pos
	return nil
}

func (m *MockPositionService) GetByID(ctx context.Context, id uint) (*domain.Position, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}
	if id == 0 {
		return nil, domain.ErrPositionNotFound
	}
	position, exists := m.positions[id]
	if !exists {
		return nil, domain.ErrPositionNotFound
	}
	return position, nil
}

func (m *MockPositionService) GetByName(ctx context.Context, name string) (*domain.Position, error) {
	if m.getByNameErr != nil {
		return nil, m.getByNameErr
	}
	if name == "" {
		return nil, domain.ErrPositionNotFound
	}
	position, exists := m.nameIndex[name]
	if !exists {
		return nil, domain.ErrPositionNotFound
	}
	return position, nil
}

func (m *MockPositionService) List(ctx context.Context, page, limit int) ([]domain.Position, int64, error) {
	if m.listErr != nil {
		return nil, 0, m.listErr
	}
	var allPositions []domain.Position
	for _, p := range m.positions {
		allPositions = append(allPositions, *p)
	}
	return allPositions, int64(len(allPositions)), nil
}

func (m *MockPositionService) Update(ctx context.Context, pos *domain.Position) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	if pos == nil || pos.ID == 0 {
		return domain.ErrPositionNotFound
	}
	if _, exists := m.positions[pos.ID]; !exists {
		return domain.ErrPositionNotFound
	}
	if existing, exists := m.nameIndex[pos.Name]; exists && existing.ID != pos.ID {
		return domain.ErrPositionNameExists
	}
	m.positions[pos.ID] = pos
	delete(m.nameIndex, pos.Name)
	m.nameIndex[pos.Name] = pos
	return nil
}

func (m *MockPositionService) Delete(ctx context.Context, id uint) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	if id == 0 {
		return domain.ErrPositionNotFound
	}
	position, exists := m.positions[id]
	if !exists {
		return domain.ErrPositionNotFound
	}
	delete(m.nameIndex, position.Name)
	delete(m.positions, id)
	return nil
}

// Helper para setup del router
func setupTestRouterPosition() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	return r
}

// ============================================
// HTTP TESTS - Position Handler
// ============================================

// --- CREATE TESTS ---

func TestPositionHandler_Create_Success(t *testing.T) {
	mockSvc := NewMockPositionService()
	handler := NewPositionHandler(service.NewPositionService(mockSvc))
	router := setupTestRouterPosition()

	router.POST("/positions", handler.Create)

	body := map[string]interface{}{
		"name":         "Software Engineer",
		"description":  "Develops software",
		"department_id": uint(1),
		"is_active":    true,
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/positions", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201 but got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestPositionHandler_Create_InvalidJSON(t *testing.T) {
	mockSvc := NewMockPositionService()
	handler := NewPositionHandler(service.NewPositionService(mockSvc))
	router := setupTestRouterPosition()

	router.POST("/positions", handler.Create)

	req, _ := http.NewRequest("POST", "/positions", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 but got %d", w.Code)
	}
}

func TestPositionHandler_Create_MissingFields(t *testing.T) {
	mockSvc := NewMockPositionService()
	handler := NewPositionHandler(service.NewPositionService(mockSvc))
	router := setupTestRouterPosition()

	router.POST("/positions", handler.Create)

	// Missing required fields (name and description)
	body := map[string]interface{}{
		"department_id": uint(1),
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/positions", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 but got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestPositionHandler_Create_DuplicateName(t *testing.T) {
	mockSvc := NewMockPositionService()
	// Pre-create position with same name
	mockSvc.nameIndex["Software Engineer"] = &domain.Position{ID: 1, Name: "Software Engineer"}
	mockSvc.positions[1] = &domain.Position{ID: 1, Name: "Software Engineer"}

	handler := NewPositionHandler(service.NewPositionService(mockSvc))
	router := setupTestRouterPosition()

	router.POST("/positions", handler.Create)

	body := map[string]interface{}{
		"name":        "Software Engineer", // Duplicate
		"description": "Another position",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/positions", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("Expected status 409 but got %d", w.Code)
	}
}

// --- GET BY ID TESTS ---

func TestPositionHandler_GetByID_Success(t *testing.T) {
	mockSvc := NewMockPositionService()
	mockSvc.positions[1] = &domain.Position{
		ID:          1,
		Name:        "Software Engineer",
		Description: "Develops software",
		IsActive:    true,
	}
	mockSvc.nameIndex["Software Engineer"] = mockSvc.positions[1]

	handler := NewPositionHandler(service.NewPositionService(mockSvc))
	router := setupTestRouterPosition()

	router.GET("/positions/:id", handler.GetByID)

	req, _ := http.NewRequest("GET", "/positions/1", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 but got %d", w.Code)
	}
}

func TestPositionHandler_GetByID_InvalidID(t *testing.T) {
	mockSvc := NewMockPositionService()
	handler := NewPositionHandler(service.NewPositionService(mockSvc))
	router := setupTestRouterPosition()

	router.GET("/positions/:id", handler.GetByID)

	req, _ := http.NewRequest("GET", "/positions/abc", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 but got %d", w.Code)
	}
}

func TestPositionHandler_GetByID_NotFound(t *testing.T) {
	mockSvc := NewMockPositionService()
	handler := NewPositionHandler(service.NewPositionService(mockSvc))
	router := setupTestRouterPosition()

	router.GET("/positions/:id", handler.GetByID)

	req, _ := http.NewRequest("GET", "/positions/999", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 but got %d", w.Code)
	}
}

// --- GET BY NAME TESTS ---

func TestPositionHandler_GetByName_Success(t *testing.T) {
	mockSvc := NewMockPositionService()
	mockSvc.positions[1] = &domain.Position{
		ID:   1,
		Name: "Software Engineer",
	}
	mockSvc.nameIndex["Software Engineer"] = mockSvc.positions[1]

	handler := NewPositionHandler(service.NewPositionService(mockSvc))
	router := setupTestRouterPosition()

	router.GET("/positions", handler.GetByName)

	req, _ := http.NewRequest("GET", "/positions?name=Software Engineer", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 but got %d", w.Code)
	}
}

func TestPositionHandler_GetByName_MissingParam(t *testing.T) {
	mockSvc := NewMockPositionService()
	handler := NewPositionHandler(service.NewPositionService(mockSvc))
	router := setupTestRouterPosition()

	router.GET("/positions", handler.GetByName)

	req, _ := http.NewRequest("GET", "/positions", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 but got %d", w.Code)
	}
}

func TestPositionHandler_GetByName_NotFound(t *testing.T) {
	mockSvc := NewMockPositionService()
	handler := NewPositionHandler(service.NewPositionService(mockSvc))
	router := setupTestRouterPosition()

	router.GET("/positions", handler.GetByName)

	req, _ := http.NewRequest("GET", "/positions?name=NonExistent", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 but got %d", w.Code)
	}
}

// --- LIST TESTS ---

func TestPositionHandler_List_Success(t *testing.T) {
	mockSvc := NewMockPositionService()
	// Add some positions
	for i := 1; i <= 3; i++ {
		p := &domain.Position{
			ID:   uint(i),
			Name: "Position" + string(rune('0'+i)),
		}
		mockSvc.positions[uint(i)] = p
		mockSvc.nameIndex[p.Name] = p
	}

	handler := NewPositionHandler(service.NewPositionService(mockSvc))
	router := setupTestRouterPosition()

	router.GET("/positions", handler.List)

	req, _ := http.NewRequest("GET", "/positions?page=1&limit=10", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 but got %d", w.Code)
	}
}

func TestPositionHandler_List_DefaultPagination(t *testing.T) {
	mockSvc := NewMockPositionService()
	handler := NewPositionHandler(service.NewPositionService(mockSvc))
	router := setupTestRouterPosition()

	router.GET("/positions", handler.List)

	req, _ := http.NewRequest("GET", "/positions", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 but got %d", w.Code)
	}
}

// --- UPDATE TESTS ---

func TestPositionHandler_Update_Success(t *testing.T) {
	mockSvc := NewMockPositionService()
	mockSvc.positions[1] = &domain.Position{
		ID:          1,
		Name:        "Software Engineer",
		Description: "Develops software",
		IsActive:    true,
	}
	mockSvc.nameIndex["Software Engineer"] = mockSvc.positions[1]

	handler := NewPositionHandler(service.NewPositionService(mockSvc))
	router := setupTestRouterPosition()

	router.PUT("/positions/:id", handler.Update)

	body := map[string]interface{}{
		"name":        "Senior Software Engineer",
		"description": "Develops complex software",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("PUT", "/positions/1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 but got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestPositionHandler_Update_InvalidID(t *testing.T) {
	mockSvc := NewMockPositionService()
	handler := NewPositionHandler(service.NewPositionService(mockSvc))
	router := setupTestRouterPosition()

	router.PUT("/positions/:id", handler.Update)

	body := map[string]interface{}{
		"name": "Senior Engineer",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("PUT", "/positions/abc", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 but got %d", w.Code)
	}
}

func TestPositionHandler_Update_NotFound(t *testing.T) {
	mockSvc := NewMockPositionService()
	handler := NewPositionHandler(service.NewPositionService(mockSvc))
	router := setupTestRouterPosition()

	router.PUT("/positions/:id", handler.Update)

	body := map[string]interface{}{
		"name": "Senior Engineer",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("PUT", "/positions/999", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 but got %d", w.Code)
	}
}

func TestPositionHandler_Update_DuplicateName(t *testing.T) {
	mockSvc := NewMockPositionService()
	// Existing position
	mockSvc.positions[1] = &domain.Position{
		ID:   1,
		Name: "Software Engineer",
	}
	mockSvc.nameIndex["Software Engineer"] = mockSvc.positions[1]
	// Another position with different name
	mockSvc.positions[2] = &domain.Position{
		ID:   2,
		Name: "QA Engineer",
	}
	mockSvc.nameIndex["QA Engineer"] = mockSvc.positions[2]

	handler := NewPositionHandler(service.NewPositionService(mockSvc))
	router := setupTestRouterPosition()

	router.PUT("/positions/:id", handler.Update)

	// Try to update position 2 with same name as position 1
	body := map[string]interface{}{
		"name": "Software Engineer",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("PUT", "/positions/2", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("Expected status 409 but got %d", w.Code)
	}
}

// --- DELETE TESTS ---

func TestPositionHandler_Delete_Success(t *testing.T) {
	mockSvc := NewMockPositionService()
	mockSvc.positions[1] = &domain.Position{ID: 1, Name: "Software Engineer"}
	mockSvc.nameIndex["Software Engineer"] = mockSvc.positions[1]

	handler := NewPositionHandler(service.NewPositionService(mockSvc))
	router := setupTestRouterPosition()

	router.DELETE("/positions/:id", handler.Delete)

	req, _ := http.NewRequest("DELETE", "/positions/1", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204 but got %d", w.Code)
	}
}

func TestPositionHandler_Delete_NotFound(t *testing.T) {
	mockSvc := NewMockPositionService()
	handler := NewPositionHandler(service.NewPositionService(mockSvc))
	router := setupTestRouterPosition()

	router.DELETE("/positions/:id", handler.Delete)

	req, _ := http.NewRequest("DELETE", "/positions/999", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 but got %d", w.Code)
	}
}

func TestPositionHandler_Delete_InvalidID(t *testing.T) {
	mockSvc := NewMockPositionService()
	handler := NewPositionHandler(service.NewPositionService(mockSvc))
	router := setupTestRouterPosition()

	router.DELETE("/positions/:id", handler.Delete)

	req, _ := http.NewRequest("DELETE", "/positions/abc", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 but got %d", w.Code)
	}
}

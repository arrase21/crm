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
)

type mockEmpRepo struct {
	employees    map[uint]*domain.Employee
	userIdx      map[uint]*domain.Employee
	nextID       uint
	createErr    error
	getByIDErr   error
	getByUserIDErr error
	listErr      error
	listActiveErr error
	updateErr    error
	deleteErr    error
}

func newMockEmpRepo() *mockEmpRepo {
	return &mockEmpRepo{
		employees: make(map[uint]*domain.Employee),
		userIdx:   make(map[uint]*domain.Employee),
		nextID:    1,
	}
}

func (m *mockEmpRepo) Create(ctx context.Context, emp *domain.Employee) error {
	if m.createErr != nil {
		return m.createErr
	}
	if emp == nil {
		return domain.ErrEmployeeNotFound
	}
	emp.ID = m.nextID
	m.nextID++
	m.employees[emp.ID] = emp
	m.userIdx[emp.UserID] = emp
	return nil
}

func (m *mockEmpRepo) GetByID(ctx context.Context, id uint) (*domain.Employee, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}
	if id == 0 {
		return nil, domain.ErrEmployeeNotFound
	}
	e, ok := m.employees[id]
	if !ok {
		return nil, domain.ErrEmployeeNotFound
	}
	return e, nil
}

func (m *mockEmpRepo) GetByUserID(ctx context.Context, uid uint) (*domain.Employee, error) {
	if m.getByUserIDErr != nil {
		return nil, m.getByUserIDErr
	}
	if uid == 0 {
		return nil, domain.ErrEmployeeNotFound
	}
	e, ok := m.userIdx[uid]
	if !ok {
		return nil, domain.ErrEmployeeNotFound
	}
	return e, nil
}

func (m *mockEmpRepo) List(ctx context.Context, page, limit int) ([]domain.Employee, int64, error) {
	if m.listErr != nil {
		return nil, 0, m.listErr
	}
	var all []domain.Employee
	for _, e := range m.employees {
		all = append(all, *e)
	}
	return all, int64(len(all)), nil
}

func (m *mockEmpRepo) ListActive(ctx context.Context, page, limit int) ([]domain.Employee, int64, error) {
	if m.listActiveErr != nil {
		return nil, 0, m.listActiveErr
	}
	var act []domain.Employee
	for _, e := range m.employees {
		if e.IsActive {
			act = append(act, *e)
		}
	}
	return act, int64(len(act)), nil
}

func (m *mockEmpRepo) Update(ctx context.Context, emp *domain.Employee) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	if emp == nil || emp.ID == 0 {
		return domain.ErrEmployeeNotFound
	}
	if _, ok := m.employees[emp.ID]; !ok {
		return domain.ErrEmployeeNotFound
	}
	m.employees[emp.ID] = emp
	return nil
}

func (m *mockEmpRepo) Delete(ctx context.Context, id uint) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	if id == 0 {
		return domain.ErrEmployeeNotFound
	}
	if _, ok := m.employees[id]; !ok {
		return domain.ErrEmployeeNotFound
	}
	delete(m.employees, id)
	return nil
}

type mockUserRepoForEmp struct {
	getByIDErr error
	user       *domain.User
}

func (m *mockUserRepoForEmp) Create(ctx context.Context, _ *domain.User) error  { return nil }
func (m *mockUserRepoForEmp) GetByID(ctx context.Context, _ uint) (*domain.User, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}
	return m.user, nil
}
func (m *mockUserRepoForEmp) GetByDNI(_ context.Context, _ string) (*domain.User, error)    { return nil, domain.ErrUserNotFound }
func (m *mockUserRepoForEmp) GetByEmail(_ context.Context, _ string) (*domain.User, error)   { return nil, domain.ErrUserNotFound }
func (m *mockUserRepoForEmp) GetByPhone(_ context.Context, _ string) (*domain.User, error)   { return nil, domain.ErrUserNotFound }
func (m *mockUserRepoForEmp) List(_ context.Context, _, _ int) ([]domain.User, int64, error) { return nil, 0, nil }
func (m *mockUserRepoForEmp) Update(_ context.Context, _ *domain.User) error                 { return nil }
func (m *mockUserRepoForEmp) Delete(_ context.Context, _ uint) error                         { return nil }

type mockDeptRepoForEmp struct {
	getByIDErr error
	dept       *domain.Department
}

func (m *mockDeptRepoForEmp) Create(_ context.Context, _ *domain.Department) error  { return nil }
func (m *mockDeptRepoForEmp) GetByID(_ context.Context, _ uint) (*domain.Department, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}
	return m.dept, nil
}
func (m *mockDeptRepoForEmp) GetByCode(_ context.Context, _ string) (*domain.Department, error)     { return nil, domain.ErrDepartmentNotFound }
func (m *mockDeptRepoForEmp) GetByName(_ context.Context, _ string) (*domain.Department, error)     { return nil, domain.ErrDepartmentNotFound }
func (m *mockDeptRepoForEmp) List(_ context.Context, _, _ int) ([]domain.Department, int64, error)  { return nil, 0, nil }
func (m *mockDeptRepoForEmp) Update(_ context.Context, _ *domain.Department) error                  { return nil }
func (m *mockDeptRepoForEmp) Delete(_ context.Context, _ uint) error                                { return nil }

type mockPosRepoForEmp struct {
	getByIDErr error
	pos        *domain.Position
}

func (m *mockPosRepoForEmp) Create(_ context.Context, _ *domain.Position) error               { return nil }
func (m *mockPosRepoForEmp) GetByID(_ context.Context, _ uint) (*domain.Position, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}
	return m.pos, nil
}
func (m *mockPosRepoForEmp) GetByIDWithDepartment(_ context.Context, _ uint) (*domain.Position, error) { return nil, domain.ErrPositionNotFound }
func (m *mockPosRepoForEmp) GetByName(_ context.Context, _ string) (*domain.Position, error)           { return nil, domain.ErrPositionNotFound }
func (m *mockPosRepoForEmp) List(_ context.Context, _, _ int) ([]domain.Position, int64, error)         { return nil, 0, nil }
func (m *mockPosRepoForEmp) ListByDepartment(_ context.Context, _ uint) ([]domain.Position, int64, error) { return nil, 0, nil }
func (m *mockPosRepoForEmp) CountByDepartment(_ context.Context, _ uint) (int64, error)                { return 0, nil }
func (m *mockPosRepoForEmp) Update(_ context.Context, _ *domain.Position) error                        { return nil }
func (m *mockPosRepoForEmp) Delete(_ context.Context, _ uint) error                                    { return nil }

func TestEmployeeHandler_Create_Success(t *testing.T) {
	empRepo := newMockEmpRepo()
	userRepo := &mockUserRepoForEmp{user: &domain.User{ID: 1}}
	deptRepo := &mockDeptRepoForEmp{dept: &domain.Department{ID: 1}}
	posRepo := &mockPosRepoForEmp{pos: &domain.Position{ID: 1}}

	svc := service.NewEmployeeService(empRepo, userRepo, deptRepo, posRepo)
	handler := NewEmployeeHandler(svc)
	router := setupTestRouter()

	router.POST("/employees", handler.Create)

	body := map[string]interface{}{
		"user_id":       1,
		"department_id": 1,
		"position_id":   1,
		"is_active":     true,
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/employees", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201 but got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestEmployeeHandler_Create_InvalidJSON(t *testing.T) {
	empRepo := newMockEmpRepo()
	svc := service.NewEmployeeService(empRepo, &mockUserRepoForEmp{}, &mockDeptRepoForEmp{}, &mockPosRepoForEmp{})
	handler := NewEmployeeHandler(svc)
	router := setupTestRouter()

	router.POST("/employees", handler.Create)

	req, _ := http.NewRequest("POST", "/employees", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 but got %d", w.Code)
	}
}

func TestEmployeeHandler_Create_MissingFields(t *testing.T) {
	empRepo := newMockEmpRepo()
	svc := service.NewEmployeeService(empRepo, &mockUserRepoForEmp{}, &mockDeptRepoForEmp{}, &mockPosRepoForEmp{})
	handler := NewEmployeeHandler(svc)
	router := setupTestRouter()

	router.POST("/employees", handler.Create)

	body := map[string]interface{}{"user_id": 1}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/employees", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 but got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestEmployeeHandler_Create_AlreadyExists(t *testing.T) {
	empRepo := newMockEmpRepo()
	empRepo.userIdx[1] = &domain.Employee{UserID: 1}

	svc := service.NewEmployeeService(empRepo, &mockUserRepoForEmp{}, &mockDeptRepoForEmp{}, &mockPosRepoForEmp{})
	handler := NewEmployeeHandler(svc)
	router := setupTestRouter()

	router.POST("/employees", handler.Create)

	body := map[string]interface{}{
		"user_id":       1,
		"department_id": 1,
		"position_id":   1,
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/employees", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("Expected status 409 but got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestEmployeeHandler_GetByID_Success(t *testing.T) {
	empRepo := newMockEmpRepo()
	empRepo.employees[1] = &domain.Employee{ID: 1, UserID: 1, DepartmentID: 1, PositionID: 1, IsActive: true}
	empRepo.userIdx[1] = empRepo.employees[1]

	svc := service.NewEmployeeService(empRepo, &mockUserRepoForEmp{}, &mockDeptRepoForEmp{}, &mockPosRepoForEmp{})
	handler := NewEmployeeHandler(svc)
	router := setupTestRouter()

	router.GET("/employees/:id", handler.GetByID)

	req, _ := http.NewRequest("GET", "/employees/1", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 but got %d", w.Code)
	}
}

func TestEmployeeHandler_GetByID_InvalidID(t *testing.T) {
	empRepo := newMockEmpRepo()
	svc := service.NewEmployeeService(empRepo, &mockUserRepoForEmp{}, &mockDeptRepoForEmp{}, &mockPosRepoForEmp{})
	handler := NewEmployeeHandler(svc)
	router := setupTestRouter()

	router.GET("/employees/:id", handler.GetByID)

	req, _ := http.NewRequest("GET", "/employees/abc", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 but got %d", w.Code)
	}
}

func TestEmployeeHandler_GetByID_NotFound(t *testing.T) {
	empRepo := newMockEmpRepo()
	svc := service.NewEmployeeService(empRepo, &mockUserRepoForEmp{}, &mockDeptRepoForEmp{}, &mockPosRepoForEmp{})
	handler := NewEmployeeHandler(svc)
	router := setupTestRouter()

	router.GET("/employees/:id", handler.GetByID)

	req, _ := http.NewRequest("GET", "/employees/999", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 but got %d", w.Code)
	}
}

func TestEmployeeHandler_GetByUserID_Success(t *testing.T) {
	empRepo := newMockEmpRepo()
	empRepo.userIdx[1] = &domain.Employee{ID: 1, UserID: 1}

	svc := service.NewEmployeeService(empRepo, &mockUserRepoForEmp{}, &mockDeptRepoForEmp{}, &mockPosRepoForEmp{})
	handler := NewEmployeeHandler(svc)
	router := setupTestRouter()

	router.GET("/employees", handler.GetByUserID)

	req, _ := http.NewRequest("GET", "/employees?user_id=1", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 but got %d", w.Code)
	}
}

func TestEmployeeHandler_GetByUserID_MissingParam(t *testing.T) {
	empRepo := newMockEmpRepo()
	svc := service.NewEmployeeService(empRepo, &mockUserRepoForEmp{}, &mockDeptRepoForEmp{}, &mockPosRepoForEmp{})
	handler := NewEmployeeHandler(svc)
	router := setupTestRouter()

	router.GET("/employees", handler.GetByUserID)

	req, _ := http.NewRequest("GET", "/employees", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 but got %d", w.Code)
	}
}

func TestEmployeeHandler_GetByUserID_InvalidParam(t *testing.T) {
	empRepo := newMockEmpRepo()
	svc := service.NewEmployeeService(empRepo, &mockUserRepoForEmp{}, &mockDeptRepoForEmp{}, &mockPosRepoForEmp{})
	handler := NewEmployeeHandler(svc)
	router := setupTestRouter()

	router.GET("/employees", handler.GetByUserID)

	req, _ := http.NewRequest("GET", "/employees?user_id=abc", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 but got %d", w.Code)
	}
}

func TestEmployeeHandler_List_Success(t *testing.T) {
	empRepo := newMockEmpRepo()
	for i := 1; i <= 3; i++ {
		empRepo.employees[uint(i)] = &domain.Employee{ID: uint(i), UserID: uint(i), IsActive: true}
	}

	svc := service.NewEmployeeService(empRepo, &mockUserRepoForEmp{}, &mockDeptRepoForEmp{}, &mockPosRepoForEmp{})
	handler := NewEmployeeHandler(svc)
	router := setupTestRouter()

	router.GET("/employees", handler.List)

	req, _ := http.NewRequest("GET", "/employees?page=1&limit=10", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 but got %d", w.Code)
	}
}

func TestEmployeeHandler_List_ActiveOnly(t *testing.T) {
	empRepo := newMockEmpRepo()
	empRepo.employees[1] = &domain.Employee{ID: 1, UserID: 1, IsActive: true}
	empRepo.employees[2] = &domain.Employee{ID: 2, UserID: 2, IsActive: false}

	svc := service.NewEmployeeService(empRepo, &mockUserRepoForEmp{}, &mockDeptRepoForEmp{}, &mockPosRepoForEmp{})
	handler := NewEmployeeHandler(svc)
	router := setupTestRouter()

	router.GET("/employees", handler.List)

	req, _ := http.NewRequest("GET", "/employees?active=true", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 but got %d", w.Code)
	}
}

func TestEmployeeHandler_Update_Success(t *testing.T) {
	empRepo := newMockEmpRepo()
	empRepo.employees[1] = &domain.Employee{ID: 1, UserID: 1, DepartmentID: 1, PositionID: 1, IsActive: true}
	empRepo.userIdx[1] = empRepo.employees[1]

	svc := service.NewEmployeeService(empRepo, &mockUserRepoForEmp{}, &mockDeptRepoForEmp{}, &mockPosRepoForEmp{})
	handler := NewEmployeeHandler(svc)
	router := setupTestRouter()

	router.PUT("/employees/:id", handler.Update)

	body := map[string]interface{}{"department_id": 2, "position_id": 2}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("PUT", "/employees/1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 but got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestEmployeeHandler_Update_InvalidID(t *testing.T) {
	empRepo := newMockEmpRepo()
	svc := service.NewEmployeeService(empRepo, &mockUserRepoForEmp{}, &mockDeptRepoForEmp{}, &mockPosRepoForEmp{})
	handler := NewEmployeeHandler(svc)
	router := setupTestRouter()

	router.PUT("/employees/:id", handler.Update)

	body := map[string]interface{}{"department_id": 2}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("PUT", "/employees/abc", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 but got %d", w.Code)
	}
}

func TestEmployeeHandler_Update_NotFound(t *testing.T) {
	empRepo := newMockEmpRepo()
	svc := service.NewEmployeeService(empRepo, &mockUserRepoForEmp{}, &mockDeptRepoForEmp{}, &mockPosRepoForEmp{})
	handler := NewEmployeeHandler(svc)
	router := setupTestRouter()

	router.PUT("/employees/:id", handler.Update)

	body := map[string]interface{}{"department_id": 2}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("PUT", "/employees/999", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 but got %d", w.Code)
	}
}

func TestEmployeeHandler_Delete_Success(t *testing.T) {
	empRepo := newMockEmpRepo()
	empRepo.employees[1] = &domain.Employee{ID: 1, UserID: 1}

	svc := service.NewEmployeeService(empRepo, &mockUserRepoForEmp{}, &mockDeptRepoForEmp{}, &mockPosRepoForEmp{})
	handler := NewEmployeeHandler(svc)
	router := setupTestRouter()

	router.DELETE("/employees/:id", handler.Delete)

	req, _ := http.NewRequest("DELETE", "/employees/1", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204 but got %d", w.Code)
	}
}

func TestEmployeeHandler_Delete_NotFound(t *testing.T) {
	empRepo := newMockEmpRepo()
	svc := service.NewEmployeeService(empRepo, &mockUserRepoForEmp{}, &mockDeptRepoForEmp{}, &mockPosRepoForEmp{})
	handler := NewEmployeeHandler(svc)
	router := setupTestRouter()

	router.DELETE("/employees/:id", handler.Delete)

	req, _ := http.NewRequest("DELETE", "/employees/999", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 but got %d", w.Code)
	}
}

func TestEmployeeHandler_Delete_InvalidID(t *testing.T) {
	empRepo := newMockEmpRepo()
	svc := service.NewEmployeeService(empRepo, &mockUserRepoForEmp{}, &mockDeptRepoForEmp{}, &mockPosRepoForEmp{})
	handler := NewEmployeeHandler(svc)
	router := setupTestRouter()

	router.DELETE("/employees/:id", handler.Delete)

	req, _ := http.NewRequest("DELETE", "/employees/abc", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 but got %d", w.Code)
	}
}

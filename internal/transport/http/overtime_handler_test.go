package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arrase21/crm/internal/domain"
	"github.com/arrase21/crm/internal/service"
)

var errOvertimeNotFound = errors.New("overtime not found")

type mockOvertimeRepo struct {
	records map[uint]*domain.Overtime
	nextID  uint

	createErr           error
	getByIDErr          error
	listByEmployeeErr   error
	listByDepartmentErr error
	listErr             error
	updateErr           error
	deleteErr           error
}

func newMockOvertimeRepo() *mockOvertimeRepo {
	return &mockOvertimeRepo{
		records: make(map[uint]*domain.Overtime),
		nextID:  1,
	}
}

func (m *mockOvertimeRepo) Create(_ context.Context, o *domain.Overtime) error {
	if m.createErr != nil {
		return m.createErr
	}
	o.ID = m.nextID
	m.nextID++
	m.records[o.ID] = o
	return nil
}

func (m *mockOvertimeRepo) GetByID(_ context.Context, id uint) (*domain.Overtime, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}
	o, ok := m.records[id]
	if !ok {
		return nil, errOvertimeNotFound
	}
	return o, nil
}

func (m *mockOvertimeRepo) ListByEmployee(_ context.Context, employeeID uint, _, _ int) ([]domain.Overtime, int64, error) {
	if m.listByEmployeeErr != nil {
		return nil, 0, m.listByEmployeeErr
	}
	var res []domain.Overtime
	for _, o := range m.records {
		if o.EmployeeID == employeeID {
			res = append(res, *o)
		}
	}
	return res, int64(len(res)), nil
}

func (m *mockOvertimeRepo) ListByDepartment(_ context.Context, _ uint, _, _ int) ([]domain.Overtime, int64, error) {
	if m.listByDepartmentErr != nil {
		return nil, 0, m.listByDepartmentErr
	}
	var all []domain.Overtime
	for _, o := range m.records {
		all = append(all, *o)
	}
	return all, int64(len(all)), nil
}

func (m *mockOvertimeRepo) List(_ context.Context, _, _ int) ([]domain.Overtime, int64, error) {
	if m.listErr != nil {
		return nil, 0, m.listErr
	}
	var all []domain.Overtime
	for _, o := range m.records {
		all = append(all, *o)
	}
	return all, int64(len(all)), nil
}

func (m *mockOvertimeRepo) Update(_ context.Context, o *domain.Overtime) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	if _, ok := m.records[o.ID]; !ok {
		return errOvertimeNotFound
	}
	m.records[o.ID] = o
	return nil
}

func (m *mockOvertimeRepo) Delete(_ context.Context, id uint) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	if _, ok := m.records[id]; !ok {
		return errOvertimeNotFound
	}
	delete(m.records, id)
	return nil
}

type mockEmpRepoForOvertime struct {
	employees      map[uint]*domain.Employee
	userIdx        map[uint]*domain.Employee
	getByIDErr     error
	getByUserIDErr error
}

func newMockEmpRepoForOvertime() *mockEmpRepoForOvertime {
	return &mockEmpRepoForOvertime{
		employees: make(map[uint]*domain.Employee),
		userIdx:   make(map[uint]*domain.Employee),
	}
}

func (m *mockEmpRepoForOvertime) Create(_ context.Context, _ *domain.Employee) error { return nil }

func (m *mockEmpRepoForOvertime) GetByID(_ context.Context, id uint) (*domain.Employee, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}
	e, ok := m.employees[id]
	if !ok {
		return nil, domain.ErrEmployeeNotFound
	}
	return e, nil
}

func (m *mockEmpRepoForOvertime) GetByUserID(_ context.Context, uid uint) (*domain.Employee, error) {
	if m.getByUserIDErr != nil {
		return nil, m.getByUserIDErr
	}
	e, ok := m.userIdx[uid]
	if !ok {
		return nil, domain.ErrEmployeeNotFound
	}
	return e, nil
}

func (m *mockEmpRepoForOvertime) List(_ context.Context, _, _ int) ([]domain.Employee, int64, error) {
	return nil, 0, nil
}

func (m *mockEmpRepoForOvertime) ListActive(_ context.Context, _, _ int) ([]domain.Employee, int64, error) {
	return nil, 0, nil
}

func (m *mockEmpRepoForOvertime) Update(_ context.Context, _ *domain.Employee) error { return nil }
func (m *mockEmpRepoForOvertime) Delete(_ context.Context, _ uint) error             { return nil }

func TestOvertimeHandler_Create_Success(t *testing.T) {
	otRepo := newMockOvertimeRepo()
	empRepo := newMockEmpRepoForOvertime()
	svc := service.NewOvertimeService(otRepo, empRepo)
	handler := NewOvertimeHandler(svc)
	router := setupTestRouter()

	router.POST("/overtime", handler.Create)

	body := map[string]interface{}{
		"employee_id": 1,
		"date":        "2024-01-15",
		"hours":       2.5,
		"reason":      "Extra work",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/overtime", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201 but got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestOvertimeHandler_Create_InvalidJSON(t *testing.T) {
	otRepo := newMockOvertimeRepo()
	empRepo := newMockEmpRepoForOvertime()
	svc := service.NewOvertimeService(otRepo, empRepo)
	handler := NewOvertimeHandler(svc)
	router := setupTestRouter()

	router.POST("/overtime", handler.Create)

	req, _ := http.NewRequest("POST", "/overtime", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 but got %d", w.Code)
	}
}

func TestOvertimeHandler_GetByID_Success(t *testing.T) {
	otRepo := newMockOvertimeRepo()
	otRepo.records[1] = &domain.Overtime{ID: 1, EmployeeID: 1}
	empRepo := newMockEmpRepoForOvertime()
	svc := service.NewOvertimeService(otRepo, empRepo)
	handler := NewOvertimeHandler(svc)
	router := setupTestRouter()

	router.GET("/overtime/:id", handler.GetByID)

	req, _ := http.NewRequest("GET", "/overtime/1", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 but got %d", w.Code)
	}
}

func TestOvertimeHandler_GetByID_NotFound(t *testing.T) {
	otRepo := newMockOvertimeRepo()
	empRepo := newMockEmpRepoForOvertime()
	svc := service.NewOvertimeService(otRepo, empRepo)
	handler := NewOvertimeHandler(svc)
	router := setupTestRouter()

	router.GET("/overtime/:id", handler.GetByID)

	req, _ := http.NewRequest("GET", "/overtime/999", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 but got %d", w.Code)
	}
}

func TestOvertimeHandler_List_Success(t *testing.T) {
	otRepo := newMockOvertimeRepo()
	for i := 1; i <= 3; i++ {
		otRepo.records[uint(i)] = &domain.Overtime{ID: uint(i), EmployeeID: uint(i)}
	}
	empRepo := newMockEmpRepoForOvertime()
	svc := service.NewOvertimeService(otRepo, empRepo)
	handler := NewOvertimeHandler(svc)
	router := setupTestRouter()

	router.GET("/overtime", handler.List)

	req, _ := http.NewRequest("GET", "/overtime?page=1&limit=10", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 but got %d", w.Code)
	}
}

func TestOvertimeHandler_Update_Success(t *testing.T) {
	otRepo := newMockOvertimeRepo()
	otRepo.records[1] = &domain.Overtime{ID: 1, EmployeeID: 1}
	empRepo := newMockEmpRepoForOvertime()
	svc := service.NewOvertimeService(otRepo, empRepo)
	handler := NewOvertimeHandler(svc)
	router := setupTestRouter()

	router.PUT("/overtime/:id", handler.Update)

	body := map[string]interface{}{
		"hours": 4.0,
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("PUT", "/overtime/1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), domain.ClaimsKey, &domain.Claims{
		UserID:   1,
		TenantID: 1,
		Roles:    []string{"super_admin"},
	})
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 but got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestOvertimeHandler_Update_NotFound(t *testing.T) {
	otRepo := newMockOvertimeRepo()
	empRepo := newMockEmpRepoForOvertime()
	svc := service.NewOvertimeService(otRepo, empRepo)
	handler := NewOvertimeHandler(svc)
	router := setupTestRouter()

	router.PUT("/overtime/:id", handler.Update)

	body := map[string]interface{}{
		"hours": 4.0,
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("PUT", "/overtime/999", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 but got %d", w.Code)
	}
}

func TestOvertimeHandler_Delete_Success(t *testing.T) {
	otRepo := newMockOvertimeRepo()
	otRepo.records[1] = &domain.Overtime{ID: 1, EmployeeID: 1}
	empRepo := newMockEmpRepoForOvertime()
	svc := service.NewOvertimeService(otRepo, empRepo)
	handler := NewOvertimeHandler(svc)
	router := setupTestRouter()

	router.DELETE("/overtime/:id", handler.Delete)

	req, _ := http.NewRequest("DELETE", "/overtime/1", nil)
	ctx := context.WithValue(req.Context(), domain.ClaimsKey, &domain.Claims{
		UserID:   1,
		TenantID: 1,
		Roles:    []string{"super_admin"},
	})
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204 but got %d", w.Code)
	}
}

func TestOvertimeHandler_Delete_NotFound(t *testing.T) {
	otRepo := newMockOvertimeRepo()
	empRepo := newMockEmpRepoForOvertime()
	svc := service.NewOvertimeService(otRepo, empRepo)
	handler := NewOvertimeHandler(svc)
	router := setupTestRouter()

	router.DELETE("/overtime/:id", handler.Delete)

	req, _ := http.NewRequest("DELETE", "/overtime/999", nil)
	ctx := context.WithValue(req.Context(), domain.ClaimsKey, &domain.Claims{
		UserID:   1,
		TenantID: 1,
		Roles:    []string{"super_admin"},
	})
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403 but got %d", w.Code)
	}
}

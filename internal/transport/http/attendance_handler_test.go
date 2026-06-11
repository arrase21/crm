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

var errAttendanceNotFound = errors.New("attendance not found")

type mockAttendanceRepo struct {
	records map[uint]*domain.Attendance
	nextID  uint

	createErr           error
	getByIDErr          error
	listByEmployeeErr   error
	listByDepartmentErr error
	listErr             error
	updateErr           error
	deleteErr           error
}

func newMockAttendanceRepo() *mockAttendanceRepo {
	return &mockAttendanceRepo{
		records: make(map[uint]*domain.Attendance),
		nextID:  1,
	}
}

func (m *mockAttendanceRepo) Create(_ context.Context, a *domain.Attendance) error {
	if m.createErr != nil {
		return m.createErr
	}
	a.ID = m.nextID
	m.nextID++
	m.records[a.ID] = a
	return nil
}

func (m *mockAttendanceRepo) GetByID(_ context.Context, id uint) (*domain.Attendance, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}
	a, ok := m.records[id]
	if !ok {
		return nil, errAttendanceNotFound
	}
	return a, nil
}

func (m *mockAttendanceRepo) ListByEmployee(_ context.Context, employeeID uint, _, _ int) ([]domain.Attendance, int64, error) {
	if m.listByEmployeeErr != nil {
		return nil, 0, m.listByEmployeeErr
	}
	var res []domain.Attendance
	for _, a := range m.records {
		if a.EmployeeID == employeeID {
			res = append(res, *a)
		}
	}
	return res, int64(len(res)), nil
}

func (m *mockAttendanceRepo) ListByDepartment(_ context.Context, _ uint, _, _ int) ([]domain.Attendance, int64, error) {
	if m.listByDepartmentErr != nil {
		return nil, 0, m.listByDepartmentErr
	}
	var all []domain.Attendance
	for _, a := range m.records {
		all = append(all, *a)
	}
	return all, int64(len(all)), nil
}

func (m *mockAttendanceRepo) List(_ context.Context, _, _ int) ([]domain.Attendance, int64, error) {
	if m.listErr != nil {
		return nil, 0, m.listErr
	}
	var all []domain.Attendance
	for _, a := range m.records {
		all = append(all, *a)
	}
	return all, int64(len(all)), nil
}

func (m *mockAttendanceRepo) Update(_ context.Context, a *domain.Attendance) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	if _, ok := m.records[a.ID]; !ok {
		return errAttendanceNotFound
	}
	m.records[a.ID] = a
	return nil
}

func (m *mockAttendanceRepo) Delete(_ context.Context, id uint) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	if _, ok := m.records[id]; !ok {
		return errAttendanceNotFound
	}
	delete(m.records, id)
	return nil
}

type mockEmpRepoForAtt struct {
	employees      map[uint]*domain.Employee
	userIdx        map[uint]*domain.Employee
	getByIDErr     error
	getByUserIDErr error
}

func newMockEmpRepoForAtt() *mockEmpRepoForAtt {
	return &mockEmpRepoForAtt{
		employees: make(map[uint]*domain.Employee),
		userIdx:   make(map[uint]*domain.Employee),
	}
}

func (m *mockEmpRepoForAtt) Create(_ context.Context, _ *domain.Employee) error { return nil }

func (m *mockEmpRepoForAtt) GetByID(_ context.Context, id uint) (*domain.Employee, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}
	e, ok := m.employees[id]
	if !ok {
		return nil, domain.ErrEmployeeNotFound
	}
	return e, nil
}

func (m *mockEmpRepoForAtt) GetByUserID(_ context.Context, uid uint) (*domain.Employee, error) {
	if m.getByUserIDErr != nil {
		return nil, m.getByUserIDErr
	}
	e, ok := m.userIdx[uid]
	if !ok {
		return nil, domain.ErrEmployeeNotFound
	}
	return e, nil
}

func (m *mockEmpRepoForAtt) List(_ context.Context, _, _ int) ([]domain.Employee, int64, error) {
	return nil, 0, nil
}

func (m *mockEmpRepoForAtt) ListActive(_ context.Context, _, _ int) ([]domain.Employee, int64, error) {
	return nil, 0, nil
}

func (m *mockEmpRepoForAtt) Update(_ context.Context, _ *domain.Employee) error { return nil }
func (m *mockEmpRepoForAtt) Delete(_ context.Context, _ uint) error             { return nil }

func TestAttendanceHandler_Create_Success(t *testing.T) {
	attRepo := newMockAttendanceRepo()
	empRepo := newMockEmpRepoForAtt()
	svc := service.NewAttendanceService(attRepo, empRepo)
	handler := NewAttendanceHandler(svc)
	router := setupTestRouter()

	router.POST("/attendance", handler.Create)

	body := map[string]interface{}{
		"employee_id": 1,
		"date":        "2024-01-15",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/attendance", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201 but got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestAttendanceHandler_Create_InvalidJSON(t *testing.T) {
	attRepo := newMockAttendanceRepo()
	empRepo := newMockEmpRepoForAtt()
	svc := service.NewAttendanceService(attRepo, empRepo)
	handler := NewAttendanceHandler(svc)
	router := setupTestRouter()

	router.POST("/attendance", handler.Create)

	req, _ := http.NewRequest("POST", "/attendance", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 but got %d", w.Code)
	}
}

func TestAttendanceHandler_GetByID_Success(t *testing.T) {
	attRepo := newMockAttendanceRepo()
	attRepo.records[1] = &domain.Attendance{ID: 1, EmployeeID: 1}
	empRepo := newMockEmpRepoForAtt()
	svc := service.NewAttendanceService(attRepo, empRepo)
	handler := NewAttendanceHandler(svc)
	router := setupTestRouter()

	router.GET("/attendance/:id", handler.GetByID)

	req, _ := http.NewRequest("GET", "/attendance/1", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 but got %d", w.Code)
	}
}

func TestAttendanceHandler_GetByID_NotFound(t *testing.T) {
	attRepo := newMockAttendanceRepo()
	empRepo := newMockEmpRepoForAtt()
	svc := service.NewAttendanceService(attRepo, empRepo)
	handler := NewAttendanceHandler(svc)
	router := setupTestRouter()

	router.GET("/attendance/:id", handler.GetByID)

	req, _ := http.NewRequest("GET", "/attendance/999", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 but got %d", w.Code)
	}
}

func TestAttendanceHandler_List_Success(t *testing.T) {
	attRepo := newMockAttendanceRepo()
	for i := 1; i <= 3; i++ {
		attRepo.records[uint(i)] = &domain.Attendance{ID: uint(i), EmployeeID: uint(i)}
	}
	empRepo := newMockEmpRepoForAtt()
	svc := service.NewAttendanceService(attRepo, empRepo)
	handler := NewAttendanceHandler(svc)
	router := setupTestRouter()

	router.GET("/attendance", handler.List)

	req, _ := http.NewRequest("GET", "/attendance?page=1&limit=10", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 but got %d", w.Code)
	}
}

func TestAttendanceHandler_Update_Success(t *testing.T) {
	attRepo := newMockAttendanceRepo()
	attRepo.records[1] = &domain.Attendance{ID: 1, EmployeeID: 1}
	empRepo := newMockEmpRepoForAtt()
	svc := service.NewAttendanceService(attRepo, empRepo)
	handler := NewAttendanceHandler(svc)
	router := setupTestRouter()

	router.PUT("/attendance/:id", handler.Update)

	body := map[string]interface{}{
		"check_in": "09:00",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("PUT", "/attendance/1", bytes.NewBuffer(jsonBody))
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

func TestAttendanceHandler_Update_NotFound(t *testing.T) {
	attRepo := newMockAttendanceRepo()
	empRepo := newMockEmpRepoForAtt()
	svc := service.NewAttendanceService(attRepo, empRepo)
	handler := NewAttendanceHandler(svc)
	router := setupTestRouter()

	router.PUT("/attendance/:id", handler.Update)

	body := map[string]interface{}{
		"check_in": "09:00",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("PUT", "/attendance/999", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 but got %d", w.Code)
	}
}

func TestAttendanceHandler_Delete_Success(t *testing.T) {
	attRepo := newMockAttendanceRepo()
	attRepo.records[1] = &domain.Attendance{ID: 1, EmployeeID: 1}
	empRepo := newMockEmpRepoForAtt()
	svc := service.NewAttendanceService(attRepo, empRepo)
	handler := NewAttendanceHandler(svc)
	router := setupTestRouter()

	router.DELETE("/attendance/:id", handler.Delete)

	req, _ := http.NewRequest("DELETE", "/attendance/1", nil)
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

func TestAttendanceHandler_Delete_NotFound(t *testing.T) {
	attRepo := newMockAttendanceRepo()
	empRepo := newMockEmpRepoForAtt()
	svc := service.NewAttendanceService(attRepo, empRepo)
	handler := NewAttendanceHandler(svc)
	router := setupTestRouter()

	router.DELETE("/attendance/:id", handler.Delete)

	req, _ := http.NewRequest("DELETE", "/attendance/999", nil)
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

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
)

type mockEmployeeContractRepo struct {
	contracts   map[uint]*domain.EmployeeContract
	empIdx      map[uint][]domain.EmployeeContract
	nextID      uint
	createErr   error
	getByIDErr  error
	getByEmpErr error
	listErr     error
	updateErr   error
	deleteErr   error
}

func newMockEmployeeContractRepo() *mockEmployeeContractRepo {
	return &mockEmployeeContractRepo{
		contracts: make(map[uint]*domain.EmployeeContract),
		empIdx:    make(map[uint][]domain.EmployeeContract),
		nextID:    1,
	}
}

func (m *mockEmployeeContractRepo) Create(_ context.Context, ec *domain.EmployeeContract) error {
	if m.createErr != nil {
		return m.createErr
	}
	if ec == nil {
		return domain.ErrContractNotFound
	}
	ec.ID = m.nextID
	m.nextID++
	m.contracts[ec.ID] = ec
	m.empIdx[ec.EmployeeID] = append(m.empIdx[ec.EmployeeID], *ec)
	return nil
}

func (m *mockEmployeeContractRepo) GetByID(_ context.Context, id uint) (*domain.EmployeeContract, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}
	if id == 0 {
		return nil, domain.ErrContractNotFound
	}
	c, ok := m.contracts[id]
	if !ok {
		return nil, domain.ErrContractNotFound
	}
	return c, nil
}

func (m *mockEmployeeContractRepo) GetByEmployeeID(_ context.Context, employeeID uint) ([]domain.EmployeeContract, error) {
	if m.getByEmpErr != nil {
		return nil, m.getByEmpErr
	}
	return m.empIdx[employeeID], nil
}

func (m *mockEmployeeContractRepo) List(_ context.Context, page, limit int) ([]domain.EmployeeContract, int64, error) {
	if m.listErr != nil {
		return nil, 0, m.listErr
	}
	var all []domain.EmployeeContract
	for _, c := range m.contracts {
		all = append(all, *c)
	}
	return all, int64(len(all)), nil
}

func (m *mockEmployeeContractRepo) Update(_ context.Context, ec *domain.EmployeeContract) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	if ec == nil || ec.ID == 0 {
		return domain.ErrContractNotFound
	}
	if _, ok := m.contracts[ec.ID]; !ok {
		return domain.ErrContractNotFound
	}
	m.contracts[ec.ID] = ec
	return nil
}

func (m *mockEmployeeContractRepo) Delete(_ context.Context, id uint) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	if id == 0 {
		return domain.ErrContractNotFound
	}
	if _, ok := m.contracts[id]; !ok {
		return domain.ErrContractNotFound
	}
	delete(m.contracts, id)
	return nil
}

type mockEmpRepoForContract struct {
	employees  map[uint]*domain.Employee
	nextID     uint
	getByIDErr error
}

func newMockEmpRepoForContract() *mockEmpRepoForContract {
	return &mockEmpRepoForContract{
		employees: make(map[uint]*domain.Employee),
		nextID:    1,
	}
}

func (m *mockEmpRepoForContract) Create(_ context.Context, emp *domain.Employee) error {
	if emp == nil {
		return domain.ErrEmployeeNotFound
	}
	emp.ID = m.nextID
	m.nextID++
	m.employees[emp.ID] = emp
	return nil
}

func (m *mockEmpRepoForContract) GetByID(_ context.Context, id uint) (*domain.Employee, error) {
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

func (m *mockEmpRepoForContract) GetByUserID(_ context.Context, _ uint) (*domain.Employee, error) {
	return nil, domain.ErrEmployeeNotFound
}

func (m *mockEmpRepoForContract) List(_ context.Context, _, _ int) ([]domain.Employee, int64, error) {
	return nil, 0, nil
}

func (m *mockEmpRepoForContract) ListActive(_ context.Context, _, _ int) ([]domain.Employee, int64, error) {
	return nil, 0, nil
}

func (m *mockEmpRepoForContract) Update(_ context.Context, _ *domain.Employee) error {
	return nil
}

func (m *mockEmpRepoForContract) Delete(_ context.Context, _ uint) error {
	return nil
}

func (m *mockEmpRepoForContract) ListBySupervisor(_ context.Context, _ uint, _, _ int) ([]domain.Employee, int64, error) {
	return nil, 0, nil
}

type mockContractTypeRepoForContract struct {
	getByIDErr error
}

func (m *mockContractTypeRepoForContract) Create(_ context.Context, _ *domain.ContractType) error {
	return nil
}

func (m *mockContractTypeRepoForContract) GetByID(_ context.Context, id uint) (*domain.ContractType, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}
	return &domain.ContractType{ID: id}, nil
}

func (m *mockContractTypeRepoForContract) GetByCountry(_ context.Context, _ string) ([]domain.ContractType, error) {
	return nil, nil
}

func (m *mockContractTypeRepoForContract) List(_ context.Context, _, _ int) ([]domain.ContractType, int64, error) {
	return nil, 0, nil
}

func (m *mockContractTypeRepoForContract) Update(_ context.Context, _ *domain.ContractType) error {
	return nil
}

func (m *mockContractTypeRepoForContract) Delete(_ context.Context, _ uint) error {
	return nil
}

type mockCountryParamRepoForContract struct{}

func (m *mockCountryParamRepoForContract) Create(_ context.Context, _ *domain.CountryParam) error {
	return nil
}

func (m *mockCountryParamRepoForContract) GetByID(_ context.Context, _ uint) (*domain.CountryParam, error) {
	return nil, domain.ErrCountryParamNotFound
}

func (m *mockCountryParamRepoForContract) GetByCountryCode(_ context.Context, _ string) (*domain.CountryParam, error) {
	return nil, domain.ErrCountryParamNotFound
}

func (m *mockCountryParamRepoForContract) List(_ context.Context, _, _ int) ([]domain.CountryParam, int64, error) {
	return nil, 0, nil
}

func (m *mockCountryParamRepoForContract) Update(_ context.Context, _ *domain.CountryParam) error {
	return nil
}

func (m *mockCountryParamRepoForContract) Delete(_ context.Context, _ uint) error {
	return nil
}

type mockOvertimeRepoForContract struct{}

func (m *mockOvertimeRepoForContract) Create(_ context.Context, _ *domain.Overtime) error {
	return nil
}
func (m *mockOvertimeRepoForContract) GetByID(_ context.Context, _ uint) (*domain.Overtime, error) {
	return nil, nil
}
func (m *mockOvertimeRepoForContract) ListByEmployee(_ context.Context, _ uint, _, _ int) ([]domain.Overtime, int64, error) {
	return nil, 0, nil
}
func (m *mockOvertimeRepoForContract) ListByEmployeeAndPeriod(_ context.Context, _ uint, _, _ time.Time) ([]domain.Overtime, error) {
	return nil, nil
}
func (m *mockOvertimeRepoForContract) ListByDepartment(_ context.Context, _ uint, _, _ int) ([]domain.Overtime, int64, error) {
	return nil, 0, nil
}
func (m *mockOvertimeRepoForContract) List(_ context.Context, _, _ int) ([]domain.Overtime, int64, error) {
	return nil, 0, nil
}
func (m *mockOvertimeRepoForContract) Update(_ context.Context, _ *domain.Overtime) error {
	return nil
}
func (m *mockOvertimeRepoForContract) Delete(_ context.Context, _ uint) error {
	return nil
}

type mockPayrollRecordRepoForContract struct{}

func (m *mockPayrollRecordRepoForContract) Create(_ context.Context, _ *domain.PayrollRecord) error {
	return nil
}

func (m *mockPayrollRecordRepoForContract) GetByID(_ context.Context, _ uint) (*domain.PayrollRecord, error) {
	return nil, domain.ErrPayrollRecordNotFound
}

func (m *mockPayrollRecordRepoForContract) GetByContractID(_ context.Context, _ uint, _, _ int) ([]domain.PayrollRecord, int64, error) {
	return nil, 0, nil
}

func (m *mockPayrollRecordRepoForContract) List(_ context.Context, _, _ int) ([]domain.PayrollRecord, int64, error) {
	return nil, 0, nil
}

func TestEmployeeContractHandler_Create_Success(t *testing.T) {
	contractRepo := newMockEmployeeContractRepo()
	empRepo := newMockEmpRepoForContract()
	empRepo.employees[1] = &domain.Employee{ID: 1}

	svc := service.NewEmployeeContractService(contractRepo, empRepo, &mockContractTypeRepoForContract{}, &mockCountryParamRepoForContract{}, &mockPayrollRecordRepoForContract{}, &mockOvertimeRepoForContract{})
	handler := NewEmployeeContractHandler(svc)
	router := setupTestRouter()

	router.POST("/contracts", handler.Create)

	body := map[string]interface{}{
		"employee_id":       1,
		"base_salary":       5000,
		"currency":          "USD",
		"country_code":      "CO",
		"start_date":        "2024-01-01",
		"work_hours_per_day": 8,
		"work_days_per_week": 5,
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/contracts", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201 but got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestEmployeeContractHandler_Create_InvalidJSON(t *testing.T) {
	svc := service.NewEmployeeContractService(newMockEmployeeContractRepo(), newMockEmpRepoForContract(), &mockContractTypeRepoForContract{}, &mockCountryParamRepoForContract{}, &mockPayrollRecordRepoForContract{}, &mockOvertimeRepoForContract{})
	handler := NewEmployeeContractHandler(svc)
	router := setupTestRouter()

	router.POST("/contracts", handler.Create)

	req, _ := http.NewRequest("POST", "/contracts", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 but got %d", w.Code)
	}
}

func TestEmployeeContractHandler_GetByID_Success(t *testing.T) {
	contractRepo := newMockEmployeeContractRepo()
	contractRepo.contracts[1] = &domain.EmployeeContract{ID: 1, EmployeeID: 1, BaseSalary: 5000}

	svc := service.NewEmployeeContractService(contractRepo, newMockEmpRepoForContract(), &mockContractTypeRepoForContract{}, &mockCountryParamRepoForContract{}, &mockPayrollRecordRepoForContract{}, &mockOvertimeRepoForContract{})
	handler := NewEmployeeContractHandler(svc)
	router := setupTestRouter()

	router.GET("/contracts/:id", handler.GetByID)

	req, _ := http.NewRequest("GET", "/contracts/1", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 but got %d", w.Code)
	}
}

func TestEmployeeContractHandler_GetByID_NotFound(t *testing.T) {
	contractRepo := newMockEmployeeContractRepo()

	svc := service.NewEmployeeContractService(contractRepo, newMockEmpRepoForContract(), &mockContractTypeRepoForContract{}, &mockCountryParamRepoForContract{}, &mockPayrollRecordRepoForContract{}, &mockOvertimeRepoForContract{})
	handler := NewEmployeeContractHandler(svc)
	router := setupTestRouter()

	router.GET("/contracts/:id", handler.GetByID)

	req, _ := http.NewRequest("GET", "/contracts/999", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 but got %d", w.Code)
	}
}

func TestEmployeeContractHandler_GetByEmployeeID_Success(t *testing.T) {
	contractRepo := newMockEmployeeContractRepo()
	contractRepo.empIdx[1] = []domain.EmployeeContract{{ID: 1, EmployeeID: 1, BaseSalary: 5000}}

	svc := service.NewEmployeeContractService(contractRepo, newMockEmpRepoForContract(), &mockContractTypeRepoForContract{}, &mockCountryParamRepoForContract{}, &mockPayrollRecordRepoForContract{}, &mockOvertimeRepoForContract{})
	handler := NewEmployeeContractHandler(svc)
	router := setupTestRouter()

	router.GET("/contracts", handler.GetByEmployeeID)

	req, _ := http.NewRequest("GET", "/contracts?employee_id=1", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 but got %d", w.Code)
	}
}

func TestEmployeeContractHandler_List_Success(t *testing.T) {
	contractRepo := newMockEmployeeContractRepo()
	for i := 1; i <= 3; i++ {
		contractRepo.contracts[uint(i)] = &domain.EmployeeContract{ID: uint(i), EmployeeID: uint(i), BaseSalary: 5000}
	}

	svc := service.NewEmployeeContractService(contractRepo, newMockEmpRepoForContract(), &mockContractTypeRepoForContract{}, &mockCountryParamRepoForContract{}, &mockPayrollRecordRepoForContract{}, &mockOvertimeRepoForContract{})
	handler := NewEmployeeContractHandler(svc)
	router := setupTestRouter()

	router.GET("/contracts", handler.List)

	req, _ := http.NewRequest("GET", "/contracts?page=1&limit=10", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 but got %d", w.Code)
	}
}

func TestEmployeeContractHandler_Update_Success(t *testing.T) {
	contractRepo := newMockEmployeeContractRepo()
	contractRepo.contracts[1] = &domain.EmployeeContract{ID: 1, EmployeeID: 1, BaseSalary: 5000}

	svc := service.NewEmployeeContractService(contractRepo, newMockEmpRepoForContract(), &mockContractTypeRepoForContract{}, &mockCountryParamRepoForContract{}, &mockPayrollRecordRepoForContract{}, &mockOvertimeRepoForContract{})
	handler := NewEmployeeContractHandler(svc)
	router := setupTestRouter()

	router.PUT("/contracts/:id", handler.Update)

	body := map[string]interface{}{"base_salary": 6000}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("PUT", "/contracts/1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 but got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestEmployeeContractHandler_Update_NotFound(t *testing.T) {
	contractRepo := newMockEmployeeContractRepo()

	svc := service.NewEmployeeContractService(contractRepo, newMockEmpRepoForContract(), &mockContractTypeRepoForContract{}, &mockCountryParamRepoForContract{}, &mockPayrollRecordRepoForContract{}, &mockOvertimeRepoForContract{})
	handler := NewEmployeeContractHandler(svc)
	router := setupTestRouter()

	router.PUT("/contracts/:id", handler.Update)

	body := map[string]interface{}{"base_salary": 6000}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("PUT", "/contracts/999", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 but got %d", w.Code)
	}
}

func TestEmployeeContractHandler_Delete_Success(t *testing.T) {
	contractRepo := newMockEmployeeContractRepo()
	contractRepo.contracts[1] = &domain.EmployeeContract{ID: 1}

	svc := service.NewEmployeeContractService(contractRepo, newMockEmpRepoForContract(), &mockContractTypeRepoForContract{}, &mockCountryParamRepoForContract{}, &mockPayrollRecordRepoForContract{}, &mockOvertimeRepoForContract{})
	handler := NewEmployeeContractHandler(svc)
	router := setupTestRouter()

	router.DELETE("/contracts/:id", handler.Delete)

	req, _ := http.NewRequest("DELETE", "/contracts/1", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204 but got %d", w.Code)
	}
}

func TestEmployeeContractHandler_Delete_NotFound(t *testing.T) {
	contractRepo := newMockEmployeeContractRepo()

	svc := service.NewEmployeeContractService(contractRepo, newMockEmpRepoForContract(), &mockContractTypeRepoForContract{}, &mockCountryParamRepoForContract{}, &mockPayrollRecordRepoForContract{}, &mockOvertimeRepoForContract{})
	handler := NewEmployeeContractHandler(svc)
	router := setupTestRouter()

	router.DELETE("/contracts/:id", handler.Delete)

	req, _ := http.NewRequest("DELETE", "/contracts/999", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 but got %d", w.Code)
	}
}

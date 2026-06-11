package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/arrase21/crm/internal/domain"
)

type MockEmployeeContractRepo struct {
	contracts   map[uint]*domain.EmployeeContract
	employeeIdx map[uint][]*domain.EmployeeContract
	nextID      uint
	CreateErr   error
	GetByIDErr  error
	ListErr     error
	UpdateErr   error
	DeleteErr   error
}

func NewMockEmployeeContractRepo() *MockEmployeeContractRepo {
	return &MockEmployeeContractRepo{
		contracts:   make(map[uint]*domain.EmployeeContract),
		employeeIdx: make(map[uint][]*domain.EmployeeContract),
		nextID:      1,
	}
}

func (m *MockEmployeeContractRepo) Create(ctx context.Context, ec *domain.EmployeeContract) error {
	if m.CreateErr != nil {
		return m.CreateErr
	}
	if ec == nil {
		return errors.New("contract cannot be nil")
	}
	ec.ID = m.nextID
	m.nextID++
	m.contracts[ec.ID] = ec
	m.employeeIdx[ec.EmployeeID] = append(m.employeeIdx[ec.EmployeeID], ec)
	return nil
}

func (m *MockEmployeeContractRepo) GetByID(ctx context.Context, id uint) (*domain.EmployeeContract, error) {
	if m.GetByIDErr != nil {
		return nil, m.GetByIDErr
	}
	if id == 0 {
		return nil, errors.New("invalid contract id")
	}
	ec, exists := m.contracts[id]
	if !exists {
		return nil, domain.ErrContractNotFound
	}
	return ec, nil
}

func (m *MockEmployeeContractRepo) GetByEmployeeID(ctx context.Context, employeeID uint) ([]domain.EmployeeContract, error) {
	if employeeID == 0 {
		return nil, errors.New("invalid employee id")
	}
	records := m.employeeIdx[employeeID]
	result := make([]domain.EmployeeContract, len(records))
	for i, ec := range records {
		result[i] = *ec
	}
	return result, nil
}

func (m *MockEmployeeContractRepo) List(ctx context.Context, page, limit int) ([]domain.EmployeeContract, int64, error) {
	if m.ListErr != nil {
		return nil, 0, m.ListErr
	}
	var all []domain.EmployeeContract
	for _, ec := range m.contracts {
		all = append(all, *ec)
	}
	return all, int64(len(all)), nil
}

func (m *MockEmployeeContractRepo) Update(ctx context.Context, ec *domain.EmployeeContract) error {
	if m.UpdateErr != nil {
		return m.UpdateErr
	}
	if ec == nil || ec.ID == 0 {
		return errors.New("contract cannot be nil or have zero id")
	}
	if _, exists := m.contracts[ec.ID]; !exists {
		return domain.ErrContractNotFound
	}
	m.contracts[ec.ID] = ec
	return nil
}

func (m *MockEmployeeContractRepo) Delete(ctx context.Context, id uint) error {
	if m.DeleteErr != nil {
		return m.DeleteErr
	}
	if id == 0 {
		return errors.New("invalid contract id")
	}
	if _, exists := m.contracts[id]; !exists {
		return domain.ErrContractNotFound
	}
	delete(m.contracts, id)
	return nil
}

type MockContractTypeRepo struct {
	ctypes     map[uint]*domain.ContractType
	nextID     uint
	CreateErr  error
	GetByIDErr error
}

func NewMockContractTypeRepo() *MockContractTypeRepo {
	return &MockContractTypeRepo{
		ctypes: make(map[uint]*domain.ContractType),
		nextID: 1,
	}
}

func (m *MockContractTypeRepo) Create(ctx context.Context, ct *domain.ContractType) error {
	if m.CreateErr != nil {
		return m.CreateErr
	}
	ct.ID = m.nextID
	m.nextID++
	m.ctypes[ct.ID] = ct
	return nil
}

func (m *MockContractTypeRepo) GetByID(ctx context.Context, id uint) (*domain.ContractType, error) {
	if m.GetByIDErr != nil {
		return nil, m.GetByIDErr
	}
	if id == 0 {
		return nil, errors.New("invalid id")
	}
	ct, exists := m.ctypes[id]
	if !exists {
		return nil, domain.ErrContractTypeNotFound
	}
	return ct, nil
}

func (m *MockContractTypeRepo) GetByCountry(ctx context.Context, countryCode string) ([]domain.ContractType, error) {
	return nil, nil
}

func (m *MockContractTypeRepo) List(ctx context.Context, page, limit int) ([]domain.ContractType, int64, error) {
	return nil, 0, nil
}

func (m *MockContractTypeRepo) Update(ctx context.Context, ct *domain.ContractType) error {
	return nil
}

func (m *MockContractTypeRepo) Delete(ctx context.Context, id uint) error {
	return nil
}

type MockCountryParamRepo struct {
	params     map[string]*domain.CountryParam
	CreateErr  error
	GetByCountryCodeErr error
}

func NewMockCountryParamRepo() *MockCountryParamRepo {
	return &MockCountryParamRepo{
		params: make(map[string]*domain.CountryParam),
	}
}

func (m *MockCountryParamRepo) Create(ctx context.Context, cp *domain.CountryParam) error {
	return nil
}

func (m *MockCountryParamRepo) GetByID(ctx context.Context, id uint) (*domain.CountryParam, error) {
	return nil, nil
}

func (m *MockCountryParamRepo) GetByCountryCode(ctx context.Context, countryCode string) (*domain.CountryParam, error) {
	if m.GetByCountryCodeErr != nil {
		return nil, m.GetByCountryCodeErr
	}
	cp, exists := m.params[countryCode]
	if !exists {
		return nil, domain.ErrCountryParamNotFound
	}
	return cp, nil
}

func (m *MockCountryParamRepo) List(ctx context.Context, page, limit int) ([]domain.CountryParam, int64, error) {
	return nil, 0, nil
}

func (m *MockCountryParamRepo) Update(ctx context.Context, cp *domain.CountryParam) error {
	return nil
}

func (m *MockCountryParamRepo) Delete(ctx context.Context, id uint) error {
	return nil
}

type MockPayrollRecordRepo struct {
	records    map[uint]*domain.PayrollRecord
	nextID     uint
	CreateErr  error
	GetByIDErr error
}

func NewMockPayrollRecordRepo() *MockPayrollRecordRepo {
	return &MockPayrollRecordRepo{
		records: make(map[uint]*domain.PayrollRecord),
		nextID:  1,
	}
}

func (m *MockPayrollRecordRepo) Create(ctx context.Context, pr *domain.PayrollRecord) error {
	if m.CreateErr != nil {
		return m.CreateErr
	}
	pr.ID = m.nextID
	m.nextID++
	m.records[pr.ID] = pr
	return nil
}

func (m *MockPayrollRecordRepo) GetByID(ctx context.Context, id uint) (*domain.PayrollRecord, error) {
	return nil, nil
}

func (m *MockPayrollRecordRepo) GetByContractID(ctx context.Context, contractID uint, page, limit int) ([]domain.PayrollRecord, int64, error) {
	return nil, 0, nil
}

func (m *MockPayrollRecordRepo) List(ctx context.Context, page, limit int) ([]domain.PayrollRecord, int64, error) {
	return nil, 0, nil
}

func validEmployeeContract() *domain.EmployeeContract {
	return &domain.EmployeeContract{
		EmployeeID:      1,
		CountryCode:     "MX",
		Currency:        "MXN",
		BaseSalary:      10000,
		StartDate:       time.Now().AddDate(0, -1, 0),
		WorkHoursPerDay: 8,
		WorkDaysPerWeek: 5,
	}
}

func TestEmployeeContractService_Create(t *testing.T) {
	tests := []struct {
		name      string
		setupMocks func() (*MockEmployeeContractRepo, *MockEmployeeRepo, *MockContractTypeRepo, *MockCountryParamRepo, *MockPayrollRecordRepo)
		input     *domain.EmployeeContract
		wantErr   bool
		errType   error
	}{
		{
			name: "create valid contract - success",
			setupMocks: func() (*MockEmployeeContractRepo, *MockEmployeeRepo, *MockContractTypeRepo, *MockCountryParamRepo, *MockPayrollRecordRepo) {
				return NewMockEmployeeContractRepo(),
					&MockEmployeeRepo{employees: map[uint]*domain.Employee{1: {ID: 1}}, userIndex: map[uint]*domain.Employee{}, nextID: 2},
					NewMockContractTypeRepo(),
					NewMockCountryParamRepo(),
					NewMockPayrollRecordRepo()
			},
			input:   validEmployeeContract(),
			wantErr: false,
		},
		{
			name: "create with nil - should fail",
			setupMocks: func() (*MockEmployeeContractRepo, *MockEmployeeRepo, *MockContractTypeRepo, *MockCountryParamRepo, *MockPayrollRecordRepo) {
				return NewMockEmployeeContractRepo(),
					NewMockEmployeeRepo(),
					NewMockContractTypeRepo(),
					NewMockCountryParamRepo(),
					NewMockPayrollRecordRepo()
			},
			input:   nil,
			wantErr: true,
		},
		{
			name: "create with missing fields - should fail validation",
			setupMocks: func() (*MockEmployeeContractRepo, *MockEmployeeRepo, *MockContractTypeRepo, *MockCountryParamRepo, *MockPayrollRecordRepo) {
				return NewMockEmployeeContractRepo(),
					NewMockEmployeeRepo(),
					NewMockContractTypeRepo(),
					NewMockCountryParamRepo(),
					NewMockPayrollRecordRepo()
			},
			input:   &domain.EmployeeContract{EmployeeID: 0},
			wantErr: true,
		},
		{
			name: "create - employee not found",
			setupMocks: func() (*MockEmployeeContractRepo, *MockEmployeeRepo, *MockContractTypeRepo, *MockCountryParamRepo, *MockPayrollRecordRepo) {
				empRepo := NewMockEmployeeRepo()
				empRepo.GetByIDErr = domain.ErrEmployeeNotFound
				return NewMockEmployeeContractRepo(),
					empRepo,
					NewMockContractTypeRepo(),
					NewMockCountryParamRepo(),
					NewMockPayrollRecordRepo()
			},
			input:   validEmployeeContract(),
			wantErr: true,
			errType: domain.ErrEmployeeNotFound,
		},
		{
			name: "create with contract type - success",
			setupMocks: func() (*MockEmployeeContractRepo, *MockEmployeeRepo, *MockContractTypeRepo, *MockCountryParamRepo, *MockPayrollRecordRepo) {
				ctRepo := NewMockContractTypeRepo()
				ctRepo.ctypes[1] = &domain.ContractType{ID: 1, Name: "Full Time"}
				ctRepo.nextID = 2
				return NewMockEmployeeContractRepo(),
					&MockEmployeeRepo{employees: map[uint]*domain.Employee{1: {ID: 1}}, userIndex: map[uint]*domain.Employee{}, nextID: 2},
					ctRepo,
					NewMockCountryParamRepo(),
					NewMockPayrollRecordRepo()
			},
			input: func() *domain.EmployeeContract {
				ec := validEmployeeContract()
				ec.ContractTypeID = 1
				return ec
			}(),
			wantErr: false,
		},
		{
			name: "create - contract type not found",
			setupMocks: func() (*MockEmployeeContractRepo, *MockEmployeeRepo, *MockContractTypeRepo, *MockCountryParamRepo, *MockPayrollRecordRepo) {
				ctRepo := NewMockContractTypeRepo()
				ctRepo.GetByIDErr = domain.ErrContractTypeNotFound
				return NewMockEmployeeContractRepo(),
					&MockEmployeeRepo{employees: map[uint]*domain.Employee{1: {ID: 1}}, userIndex: map[uint]*domain.Employee{}, nextID: 2},
					ctRepo,
					NewMockCountryParamRepo(),
					NewMockPayrollRecordRepo()
			},
			input: func() *domain.EmployeeContract {
				ec := validEmployeeContract()
				ec.ContractTypeID = 99
				return ec
			}(),
			wantErr: true,
			errType: domain.ErrContractTypeNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContract, mockEmp, mockCT, mockCP, mockPR := tt.setupMocks()
			svc := NewEmployeeContractService(mockContract, mockEmp, mockCT, mockCP, mockPR)
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

func TestEmployeeContractService_GetByID(t *testing.T) {
	mockContract := NewMockEmployeeContractRepo()
	mockContract.contracts[1] = &domain.EmployeeContract{ID: 1, EmployeeID: 1}

	svc := NewEmployeeContractService(mockContract, NewMockEmployeeRepo(), NewMockContractTypeRepo(), NewMockCountryParamRepo(), NewMockPayrollRecordRepo())
	ctx := context.Background()

	t.Run("get existing - success", func(t *testing.T) {
		result, err := svc.GetByID(ctx, 1)
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
		if result == nil || result.ID != 1 {
			t.Error("Expected contract with ID 1")
		}
	})

	t.Run("get non-existing - error", func(t *testing.T) {
		_, err := svc.GetByID(ctx, 999)
		if err == nil {
			t.Error("Expected error for non-existing")
		}
	})

	t.Run("get with id 0 - error", func(t *testing.T) {
		_, err := svc.GetByID(ctx, 0)
		if err == nil {
			t.Error("Expected error for id 0")
		}
	})
}

func TestEmployeeContractService_GetByEmployeeID(t *testing.T) {
	mockContract := NewMockEmployeeContractRepo()
	svc := NewEmployeeContractService(mockContract, NewMockEmployeeRepo(), NewMockContractTypeRepo(), NewMockCountryParamRepo(), NewMockPayrollRecordRepo())
	ctx := context.Background()

	t.Run("get by employee id 0 - error", func(t *testing.T) {
		_, err := svc.GetByEmployeeID(ctx, 0)
		if err == nil {
			t.Error("Expected error for invalid employee id")
		}
	})
}

func TestEmployeeContractService_List(t *testing.T) {
	mockContract := NewMockEmployeeContractRepo()
	for i := 1; i <= 5; i++ {
		mockContract.contracts[uint(i)] = &domain.EmployeeContract{ID: uint(i), EmployeeID: uint(i)}
	}

	svc := NewEmployeeContractService(mockContract, NewMockEmployeeRepo(), NewMockContractTypeRepo(), NewMockCountryParamRepo(), NewMockPayrollRecordRepo())
	ctx := context.Background()

	contracts, total, err := svc.List(ctx, 1, 10)
	if err != nil {
		t.Errorf("Expected no error but got: %v", err)
	}
	if total != 5 {
		t.Errorf("Expected total 5 but got %d", total)
	}
	if len(contracts) != 5 {
		t.Errorf("Expected 5 contracts but got %d", len(contracts))
	}
}

func TestEmployeeContractService_Update(t *testing.T) {
	mockContract := NewMockEmployeeContractRepo()
	mockContract.contracts[1] = &domain.EmployeeContract{ID: 1, EmployeeID: 1}

	svc := NewEmployeeContractService(mockContract, NewMockEmployeeRepo(), NewMockContractTypeRepo(), NewMockCountryParamRepo(), NewMockPayrollRecordRepo())
	ctx := context.Background()

	t.Run("update with nil - should fail", func(t *testing.T) {
		err := svc.Update(ctx, nil)
		if err == nil {
			t.Error("Expected error for nil")
		}
	})

	t.Run("update with id 0 - should fail", func(t *testing.T) {
		err := svc.Update(ctx, &domain.EmployeeContract{ID: 0})
		if err == nil {
			t.Error("Expected error for id 0")
		}
	})

	t.Run("update non-existing - should fail", func(t *testing.T) {
		err := svc.Update(ctx, &domain.EmployeeContract{ID: 999})
		if err == nil {
			t.Error("Expected error for non-existing")
		}
	})

	t.Run("update valid - success", func(t *testing.T) {
		err := svc.Update(ctx, &domain.EmployeeContract{ID: 1, EmployeeID: 1})
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
	})
}

func TestEmployeeContractService_Delete(t *testing.T) {
	mockContract := NewMockEmployeeContractRepo()
	mockContract.contracts[1] = &domain.EmployeeContract{ID: 1, EmployeeID: 1}

	svc := NewEmployeeContractService(mockContract, NewMockEmployeeRepo(), NewMockContractTypeRepo(), NewMockCountryParamRepo(), NewMockPayrollRecordRepo())
	ctx := context.Background()

	t.Run("delete existing - success", func(t *testing.T) {
		err := svc.Delete(ctx, 1)
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
		if _, exists := mockContract.contracts[1]; exists {
			t.Error("Expected contract to be deleted")
		}
	})

	t.Run("delete non-existing - error", func(t *testing.T) {
		err := svc.Delete(ctx, 999)
		if err == nil {
			t.Error("Expected error for non-existing")
		}
	})

	t.Run("delete with id 0 - error", func(t *testing.T) {
		err := svc.Delete(ctx, 0)
		if err == nil {
			t.Error("Expected error for id 0")
		}
	})
}

func TestEmployeeContractService_GetPayrollRecords(t *testing.T) {
	mockPR := NewMockPayrollRecordRepo()
	svc := NewEmployeeContractService(NewMockEmployeeContractRepo(), NewMockEmployeeRepo(), NewMockContractTypeRepo(), NewMockCountryParamRepo(), mockPR)
	ctx := context.Background()

	_, _, err := svc.GetPayrollRecords(ctx, 1, 1, 10)
	if err != nil {
		t.Errorf("Expected no error but got: %v", err)
	}
}

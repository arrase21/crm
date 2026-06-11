package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/arrase21/crm/internal/calculator"
	"github.com/arrase21/crm/internal/domain"
)

type EmployeeContractService struct {
	contractRepo     domain.EmployeeContractRepo
	employeeRepo     domain.EmployeeRepo
	contractTypeRepo domain.ContractTypeRepo
	countryParamRepo domain.CountryParamRepo
	payrollRepo      domain.PayrollRecordRepo
	overtimeRepo     domain.OvertimeRepo
}

func NewEmployeeContractService(
	contractRepo domain.EmployeeContractRepo,
	employeeRepo domain.EmployeeRepo,
	contractTypeRepo domain.ContractTypeRepo,
	countryParamRepo domain.CountryParamRepo,
	payrollRepo domain.PayrollRecordRepo,
	overtimeRepo domain.OvertimeRepo,
) *EmployeeContractService {
	return &EmployeeContractService{
		contractRepo:     contractRepo,
		employeeRepo:     employeeRepo,
		contractTypeRepo: contractTypeRepo,
		countryParamRepo: countryParamRepo,
		payrollRepo:      payrollRepo,
		overtimeRepo:     overtimeRepo,
	}
}

func (s *EmployeeContractService) Create(ctx context.Context, ec *domain.EmployeeContract) error {
	if ec == nil {
		return errors.New("contract cannot be nil")
	}
	ec.Normalize()
	if err := ec.ValidateAll(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}
	_, err := s.employeeRepo.GetByID(ctx, ec.EmployeeID)
	if err != nil {
		if errors.Is(err, domain.ErrEmployeeNotFound) {
			return domain.ErrEmployeeNotFound
		}
		return fmt.Errorf("error validating employee: %w", err)
	}
	if ec.ContractTypeID != 0 {
		_, err := s.contractTypeRepo.GetByID(ctx, ec.ContractTypeID)
		if err != nil {
			if errors.Is(err, domain.ErrContractTypeNotFound) {
				return domain.ErrContractTypeNotFound
			}
			return fmt.Errorf("error validating contract type: %w", err)
		}
	}
	return s.contractRepo.Create(ctx, ec)
}

func (s *EmployeeContractService) GetByID(ctx context.Context, id uint) (*domain.EmployeeContract, error) {
	if id == 0 {
		return nil, errors.New("invalid contract id")
	}
	return s.contractRepo.GetByID(ctx, id)
}

func (s *EmployeeContractService) GetByEmployeeID(ctx context.Context, employeeID uint) ([]domain.EmployeeContract, error) {
	if employeeID == 0 {
		return nil, errors.New("invalid employee id")
	}
	return s.contractRepo.GetByEmployeeID(ctx, employeeID)
}

func (s *EmployeeContractService) List(ctx context.Context, page, limit int) ([]domain.EmployeeContract, int64, error) {
	return s.contractRepo.List(ctx, page, limit)
}

func (s *EmployeeContractService) Update(ctx context.Context, ec *domain.EmployeeContract) error {
	if ec == nil || ec.ID == 0 {
		return errors.New("contract cannot be nil or have zero id")
	}
	if ec.ContractTypeID != 0 {
		_, err := s.contractTypeRepo.GetByID(ctx, ec.ContractTypeID)
		if err != nil {
			if errors.Is(err, domain.ErrContractTypeNotFound) {
				return domain.ErrContractTypeNotFound
			}
			return fmt.Errorf("error validating contract type: %w", err)
		}
	}
	return s.contractRepo.Update(ctx, ec)
}

func (s *EmployeeContractService) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid contract id")
	}
	return s.contractRepo.Delete(ctx, id)
}

func (s *EmployeeContractService) CalculatePayroll(ctx context.Context, contractID uint, periodStart, periodEnd time.Time, bonus, commissions int64) (*domain.PayrollRecord, error) {
	if periodStart.IsZero() || periodEnd.IsZero() {
		return nil, domain.ErrInvalidPeriod
	}
	if periodEnd.Before(periodStart) {
		return nil, domain.ErrInvalidPeriod
	}

	contract, err := s.contractRepo.GetByID(ctx, contractID)
	if err != nil {
		return nil, err
	}

	params, err := s.countryParamRepo.GetByCountryCode(ctx, contract.CountryCode)
	if err != nil {
		return nil, fmt.Errorf("error loading country params: %w", err)
	}

	calc, err := calculator.GetCalculator(contract.CountryCode)
	if err != nil {
		return nil, domain.ErrCalculatorNotFound
	}

	overtimes, err := s.overtimeRepo.ListByEmployeeAndPeriod(ctx, contract.EmployeeID, periodStart, periodEnd)
	if err != nil {
		return nil, fmt.Errorf("error loading overtime: %w", err)
	}

	var otEntries []calculator.OvertimeEntry
	for _, o := range overtimes {
		otEntries = append(otEntries, calculator.OvertimeEntry{Hours: o.Hours})
	}

	input := &calculator.PayrollInput{
		Bonus:       bonus,
		Commissions: commissions,
		Overtimes:   otEntries,
	}

	breakdown, err := calc.Calculate(contract, params, periodStart, periodEnd, input)
	if err != nil {
		return nil, fmt.Errorf("error calculating payroll: %w", err)
	}

	record := &domain.PayrollRecord{
		EmployeeContractID: contractID,
		PeriodStart:        periodStart,
		PeriodEnd:          periodEnd,
		BaseSalary:         breakdown.BaseSalary,
		Currency:           contract.Currency,
		CountryCode:        contract.CountryCode,
		GrossSalary:        breakdown.GrossSalary,
		Bonus:              breakdown.Bonus,
		Commissions:        breakdown.Commissions,
		ExtraHourPay:       breakdown.ExtraHourPay,
		HealthContribution:  breakdown.HealthContribution,
		PensionContribution: breakdown.PensionContribution,
		SolidarityFund:     breakdown.SolidarityFund,
		TransportAllowance:  breakdown.TransportAllowance,
		HousingAllowance:    breakdown.HousingAllowance,
		OtherDeductions:     breakdown.OtherDeductions,
		TotalEarnings:      breakdown.TotalEarnings,
		TotalDeductions:    breakdown.TotalDeductions,
		NetSalary:           breakdown.NetSalary,
		ServiceBonus:        breakdown.ServiceBonus,
		Severance:           breakdown.Severance,
		SeveranceInterest:   breakdown.SeveranceInterest,
		VacationProvision:   breakdown.VacationProvision,
		CalculatedAt:        time.Now(),
	}

	if err := s.payrollRepo.Create(ctx, record); err != nil {
		return nil, fmt.Errorf("error saving payroll record: %w", err)
	}

	return record, nil
}

func (s *EmployeeContractService) GetPayrollRecords(ctx context.Context, contractID uint, page, limit int) ([]domain.PayrollRecord, int64, error) {
	return s.payrollRepo.GetByContractID(ctx, contractID, page, limit)
}

package service

import (
	"context"
	"fmt"
	"time"

	"github.com/arrase21/crm/internal/calculator"
	"github.com/arrase21/crm/internal/domain"
)

type PayrollService struct {
	contractRepo     domain.EmployeeContractRepo
	countryParamRepo domain.CountryParamRepo
	payrollRepo      domain.PayrollRecordRepo
	overtimeRepo     domain.OvertimeRepo
}

func NewPayrollService(
	contractRepo domain.EmployeeContractRepo,
	countryParamRepo domain.CountryParamRepo,
	payrollRepo domain.PayrollRecordRepo,
	overtimeRepo domain.OvertimeRepo,
) *PayrollService {
	return &PayrollService{
		contractRepo:     contractRepo,
		countryParamRepo: countryParamRepo,
		payrollRepo:      payrollRepo,
		overtimeRepo:     overtimeRepo,
	}
}

func (s *PayrollService) CalculatePayroll(ctx context.Context, contractID uint, periodStart, periodEnd time.Time, bonus, commissions int64) (*domain.PayrollRecord, error) {
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

func (s *PayrollService) GetPayrollRecords(ctx context.Context, contractID uint, page, limit int) ([]domain.PayrollRecord, int64, error) {
	return s.payrollRepo.GetByContractID(ctx, contractID, page, limit)
}

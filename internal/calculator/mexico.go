package calculator

import (
	"math"
	"time"

	"github.com/arrase21/crm/internal/domain"
)

type MexicoCalculator struct{}

func init() {
	Register(&MexicoCalculator{})
}

func (c *MexicoCalculator) CountryCode() string {
	return "MX"
}

func (c *MexicoCalculator) Calculate(contract *domain.EmployeeContract, params *domain.CountryParam, periodStart, periodEnd time.Time, input *PayrollInput) (*PayrollBreakdown, error) {
	days := int64(daysBetween(periodStart, periodEnd))

	grossSalary := contract.BaseSalary * days / 30

	salaryForContributions := contract.BaseSalary
	if salaryForContributions < params.MinWage {
		salaryForContributions = params.MinWage
	}

	healthBPS := int64(math.Round(params.HealthRate * 10000))
	pensionBPS := int64(math.Round(params.PensionRate * 10000))

	health := salaryForContributions * healthBPS * days / 10000 / 30
	pension := salaryForContributions * pensionBPS * days / 10000 / 30

	totalDeductions := health + pension
	netSalary := grossSalary - totalDeductions

	return &PayrollBreakdown{
		BaseSalary:          contract.BaseSalary,
		GrossSalary:         grossSalary,
		HealthContribution:  health,
		PensionContribution: pension,
		TransportAllowance:  0,
		HousingAllowance:    0,
		OtherDeductions:     0,
		NetSalary:           netSalary,
	}, nil
}

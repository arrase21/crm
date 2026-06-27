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

	hourlyRate := contract.BaseSalary / 30 / 8
	var extraHourPay int64
	for _, ot := range input.Overtimes {
		totalHours := int64(math.Round(ot.Hours))
		extraHourPay += totalHours * hourlyRate * 200 / 100
	}

	bonus := input.Bonus
	commissions := input.Commissions

	salaryForContributions := contract.BaseSalary
	if salaryForContributions < params.MinWage {
		salaryForContributions = params.MinWage
	}

	healthBPS := int64(math.Round(params.HealthRate * 10000))
	pensionBPS := int64(math.Round(params.PensionRate * 10000))

	health := salaryForContributions * healthBPS * days / 10000 / 30
	pension := salaryForContributions * pensionBPS * days / 10000 / 30

	totalEarnings := grossSalary + extraHourPay + bonus + commissions
	totalDeductions := health + pension
	netSalary := totalEarnings - totalDeductions

	return &PayrollBreakdown{
		BaseSalary:          contract.BaseSalary,
		GrossSalary:         grossSalary,
		Bonus:               bonus,
		Commissions:         commissions,
		ExtraHourPay:        extraHourPay,
		HealthContribution:  health,
		PensionContribution: pension,
		SolidarityFund:      0,
		TransportAllowance:  0,
		HousingAllowance:    0,
		OtherDeductions:     0,
		TotalEarnings:       totalEarnings,
		TotalDeductions:     totalDeductions,
		NetSalary:           netSalary,
	}, nil
}

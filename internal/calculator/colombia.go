package calculator

import (
	"math"
	"time"

	"github.com/arrase21/crm/internal/domain"
)

type ColombiaCalculator struct{}

func init() {
	Register(&ColombiaCalculator{})
}

func (c *ColombiaCalculator) CountryCode() string {
	return "CO"
}

func (c *ColombiaCalculator) Calculate(contract *domain.EmployeeContract, params *domain.CountryParam, periodStart, periodEnd time.Time, input *PayrollInput) (*PayrollBreakdown, error) {
	days := int64(daysBetween(periodStart, periodEnd))

	grossSalary := contract.BaseSalary * days / 30

	hourlyRate := contract.BaseSalary / 30 / 8
	var extraHourPay int64
	for _, ot := range input.Overtimes {
		totalHours := int64(math.Round(ot.Hours))
		extraHourPay += totalHours * hourlyRate * 150 / 100
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

	solidarityFund := int64(0)
	if contract.BaseSalary > params.MinWage*4 {
		solidarityFund = salaryForContributions * 100 * days / 10000 / 30
	}

	transport := int64(0)
	if contract.BaseSalary <= params.MinWage*2 {
		transport = params.TransportSubs * days / 30
	}

	housing := int64(0)
	if params.HousingSubs > 0 && contract.BaseSalary <= params.MinWage*2 {
		housing = params.HousingSubs * days / 30
	}

	totalEarnings := grossSalary + extraHourPay + bonus + commissions + transport + housing
	totalDeductions := health + pension + solidarityFund
	netSalary := totalEarnings - totalDeductions

	serviceBonus := grossSalary * 180 / 360
	severance := grossSalary * days / 360
	severanceInterest := severance * 12 / 100 * days / 360
	vacationProvision := grossSalary * 15 / 360

	return &PayrollBreakdown{
		BaseSalary:          contract.BaseSalary,
		GrossSalary:         grossSalary,
		Bonus:               bonus,
		Commissions:         commissions,
		ExtraHourPay:        extraHourPay,
		HealthContribution:  health,
		PensionContribution: pension,
		SolidarityFund:      solidarityFund,
		TransportAllowance:  transport,
		HousingAllowance:    housing,
		OtherDeductions:     0,
		TotalEarnings:       totalEarnings,
		TotalDeductions:     totalDeductions,
		NetSalary:           netSalary,
		ServiceBonus:        serviceBonus,
		Severance:           severance,
		SeveranceInterest:   severanceInterest,
		VacationProvision:   vacationProvision,
	}, nil
}

func daysBetween(start, end time.Time) float64 {
	return end.Sub(start).Hours() / 24
}

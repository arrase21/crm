package calculator

import (
	"fmt"
	"time"

	"github.com/arrase21/crm/internal/domain"
)

type OvertimeEntry struct {
	Hours float64
}

type PayrollInput struct {
	Bonus       int64
	Commissions int64
	Overtimes   []OvertimeEntry
}

type PayrollBreakdown struct {
	BaseSalary          int64
	GrossSalary         int64
	Bonus               int64
	Commissions         int64
	ExtraHourPay        int64
	HealthContribution  int64
	PensionContribution int64
	SolidarityFund      int64
	TransportAllowance  int64
	HousingAllowance    int64
	OtherDeductions     int64
	TotalEarnings       int64
	TotalDeductions     int64
	NetSalary           int64
	ServiceBonus        int64
	Severance           int64
	SeveranceInterest   int64
	VacationProvision   int64
}

type PayrollCalculator interface {
	Calculate(contract *domain.EmployeeContract, params *domain.CountryParam, periodStart, periodEnd time.Time, input *PayrollInput) (*PayrollBreakdown, error)
	CountryCode() string
}

var registry = make(map[string]PayrollCalculator)

func Register(c PayrollCalculator) {
	registry[c.CountryCode()] = c
}

func GetCalculator(countryCode string) (PayrollCalculator, error) {
	c, ok := registry[countryCode]
	if !ok {
		return nil, fmt.Errorf("calculator not found for country: %s", countryCode)
	}
	return c, nil
}

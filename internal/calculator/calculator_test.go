package calculator

import (
	"testing"
	"time"

	"github.com/arrase21/crm/internal/domain"
)

func TestRegisterAndGetCalculator(t *testing.T) {
	registry = make(map[string]PayrollCalculator)

	mx := &MexicoCalculator{}
	Register(mx)

	calc, err := GetCalculator("MX")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calc.CountryCode() != "MX" {
		t.Errorf("expected MX, got %s", calc.CountryCode())
	}
}

func TestGetCalculatorNotFound(t *testing.T) {
	registry = make(map[string]PayrollCalculator)

	_, err := GetCalculator("BR")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
}

func TestRegisterOverwrites(t *testing.T) {
	registry = make(map[string]PayrollCalculator)

	Register(&MexicoCalculator{})
	Register(&MexicoCalculator{})

	calc, err := GetCalculator("MX")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calc == nil {
		t.Fatal("expected calculator but got nil")
	}
}

func TestMexicoCalculator_CountryCode(t *testing.T) {
	c := &MexicoCalculator{}
	if got := c.CountryCode(); got != "MX" {
		t.Errorf("MexicoCalculator.CountryCode() = %s, want MX", got)
	}
}

func TestColombiaCalculator_CountryCode(t *testing.T) {
	c := &ColombiaCalculator{}
	if got := c.CountryCode(); got != "CO" {
		t.Errorf("ColombiaCalculator.CountryCode() = %s, want CO", got)
	}
}

func TestMexicoCalculator_Calculate(t *testing.T) {
	c := &MexicoCalculator{}
	contract := &domain.EmployeeContract{
		BaseSalary: 20000,
		Currency:   "MXN",
	}
	params := &domain.CountryParam{
		CountryCode:  "MX",
		MinWage:      8000,
		HealthRate:   0.05,
		PensionRate:  0.065,
	}
	periodStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	periodEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)

	result, err := c.Calculate(contract, params, periodStart, periodEnd, &PayrollInput{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.GrossSalary <= 0 {
		t.Errorf("expected positive gross salary, got %d", result.GrossSalary)
	}
	if result.NetSalary <= 0 {
		t.Errorf("expected positive net salary, got %d", result.NetSalary)
	}
	if result.NetSalary >= result.GrossSalary {
		t.Errorf("net salary (%d) should be less than gross (%d)", result.NetSalary, result.GrossSalary)
	}
}

func TestMexicoCalculator_BelowMinWage(t *testing.T) {
	c := &MexicoCalculator{}
	contract := &domain.EmployeeContract{
		BaseSalary: 5000,
		Currency:   "MXN",
	}
	params := &domain.CountryParam{
		CountryCode:  "MX",
		MinWage:      8000,
		HealthRate:   0.05,
		PensionRate:  0.065,
	}
	periodStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	periodEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)

	result, err := c.Calculate(contract, params, periodStart, periodEnd, &PayrollInput{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.GrossSalary != 5000.0 {
		t.Errorf("gross salary should be based on contract, not min wage, got %d", result.GrossSalary)
	}
}

func TestMexicoCalculator_ExactlyMinWage(t *testing.T) {
	c := &MexicoCalculator{}
	contract := &domain.EmployeeContract{
		BaseSalary: 8000,
		Currency:   "MXN",
	}
	params := &domain.CountryParam{
		CountryCode:  "MX",
		MinWage:      8000,
		HealthRate:   0.05,
		PensionRate:  0.065,
	}
	periodStart := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
	periodEnd := time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC)

	result, err := c.Calculate(contract, params, periodStart, periodEnd, &PayrollInput{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.GrossSalary <= 0 {
		t.Errorf("expected positive gross salary, got %d", result.GrossSalary)
	}
}

func TestColombiaCalculator_Calculate(t *testing.T) {
	c := &ColombiaCalculator{}
	contract := &domain.EmployeeContract{
		BaseSalary: 2000000,
		Currency:   "COP",
	}
	params := &domain.CountryParam{
		CountryCode:   "CO",
		MinWage:       1000000,
		HealthRate:   0.04,
		PensionRate:  0.04,
		TransportSubs: 162000,
	}
	periodStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	periodEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)

	result, err := c.Calculate(contract, params, periodStart, periodEnd, &PayrollInput{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.GrossSalary <= 0 {
		t.Errorf("expected positive gross salary, got %d", result.GrossSalary)
	}
	if result.TransportAllowance <= 0 {
		t.Errorf("expected transport allowance for salary <= 2x min wage, got %d", result.TransportAllowance)
	}
	if result.NetSalary <= 0 {
		t.Errorf("expected positive net salary, got %d", result.NetSalary)
	}
}

func TestColombiaCalculator_AboveTransportThreshold(t *testing.T) {
	c := &ColombiaCalculator{}
	contract := &domain.EmployeeContract{
		BaseSalary: 3000000,
		Currency:   "COP",
	}
	params := &domain.CountryParam{
		CountryCode:   "CO",
		MinWage:       1000000,
		HealthRate:   0.04,
		PensionRate:  0.04,
		TransportSubs: 162000,
	}
	periodStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	periodEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)

	result, err := c.Calculate(contract, params, periodStart, periodEnd, &PayrollInput{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.TransportAllowance != 0 {
		t.Errorf("expected no transport allowance for salary > 2x min wage, got %d", result.TransportAllowance)
	}
}

func TestColombiaCalculator_WithHousingSubsidy(t *testing.T) {
	c := &ColombiaCalculator{}
	contract := &domain.EmployeeContract{
		BaseSalary: 1500000,
		Currency:   "COP",
	}
	params := &domain.CountryParam{
		CountryCode:   "CO",
		MinWage:       1000000,
		HealthRate:   0.04,
		PensionRate:  0.04,
		TransportSubs: 162000,
		HousingSubs:   120000,
	}
	periodStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	periodEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)

	result, err := c.Calculate(contract, params, periodStart, periodEnd, &PayrollInput{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.HousingAllowance <= 0 {
		t.Errorf("expected housing allowance when HousingSubs > 0, got %d", result.HousingAllowance)
	}
}

func TestColombiaCalculator_BelowMinWage(t *testing.T) {
	c := &ColombiaCalculator{}
	contract := &domain.EmployeeContract{
		BaseSalary: 500000,
		Currency:   "COP",
	}
	params := &domain.CountryParam{
		CountryCode:   "CO",
		MinWage:       1000000,
		HealthRate:   0.04,
		PensionRate:  0.04,
		TransportSubs: 162000,
	}
	periodStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	periodEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)

	result, err := c.Calculate(contract, params, periodStart, periodEnd, &PayrollInput{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.GrossSalary != 500000.0 {
		t.Errorf("gross salary should be based on contract salary, got %d", result.GrossSalary)
	}
}

func TestColombiaCalculator_TransportOnlyAboveMinWage(t *testing.T) {
	c := &ColombiaCalculator{}
	contract := &domain.EmployeeContract{
		BaseSalary: 1000000,
		Currency:   "COP",
	}
	params := &domain.CountryParam{
		CountryCode:   "CO",
		MinWage:       1000000,
		HealthRate:   0.04,
		PensionRate:  0.04,
		TransportSubs: 162000,
	}
	periodStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	periodEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)

	result, err := c.Calculate(contract, params, periodStart, periodEnd, &PayrollInput{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.TransportAllowance <= 0 {
		t.Errorf("expected transport allowance at exactly 1x min wage, got %d", result.TransportAllowance)
	}
}

func TestColombiaCalculator_AtTwoTimesMinWage(t *testing.T) {
	c := &ColombiaCalculator{}
	contract := &domain.EmployeeContract{
		BaseSalary: 2000000,
		Currency:   "COP",
	}
	params := &domain.CountryParam{
		CountryCode:   "CO",
		MinWage:       1000000,
		HealthRate:   0.04,
		PensionRate:  0.04,
		TransportSubs: 162000,
	}
	periodStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	periodEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)

	result, err := c.Calculate(contract, params, periodStart, periodEnd, &PayrollInput{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.TransportAllowance <= 0 {
		t.Errorf("expected transport allowance at exactly 2x min wage, got %d", result.TransportAllowance)
	}
}

func TestDaysBetween(t *testing.T) {
	tests := []struct {
		name  string
		start time.Time
		end   time.Time
		want  float64
	}{
		{
			name:  "same day",
			start: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			end:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			want:  0,
		},
		{
			name:  "one day",
			start: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			end:   time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			want:  1,
		},
		{
			name:  "full month January",
			start: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			end:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
			want:  30,
		},
		{
			name:  "leap year February",
			start: time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
			end:   time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
			want:  29,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := daysBetween(tt.start, tt.end)
			if got != tt.want {
				t.Errorf("daysBetween() = %v, want %v", got, tt.want)
			}
		})
	}
}

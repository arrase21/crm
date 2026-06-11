package domain

import (
	"errors"
	"strings"
	"time"
)

func (u *User) Validate() error {
	if u.Gender != "M" && u.Gender != "F" {
		return errors.New("Invalid option")
	}
	if u.BirthDay.IsZero() || u.BirthDay.After(time.Now()) {
		return errors.New("invalid birthday")
	}
	if u.IsMinor() {
		return errors.New("user must be over 18")
	}
	return nil
}

func (u *User) Normalize() {
	u.FirstName = strings.TrimSpace(u.FirstName)
	u.LastName = strings.TrimSpace(u.LastName)
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	u.Dni = strings.TrimSpace(u.Dni)
	u.Phone = strings.TrimSpace(u.Phone)
}

func (u *User) Required() error {
	u.Normalize()

	switch {
	case u.FirstName == "":
		return errors.New("first name is required")
	case u.LastName == "":
		return errors.New("last name is required")
	case u.Email == "":
		return errors.New("email is required")
	case u.Dni == "":
		return errors.New("dni is required")
	case u.Phone == "":
		return errors.New("phone is required")
	}

	return nil
}

func (u *User) IsMinor() bool {
	now := time.Now()
	age := now.Year() - u.BirthDay.Year()
	if now.Month() < u.BirthDay.Month() || (now.Month() == u.BirthDay.Month() && now.Day() < u.BirthDay.Day()) {
		age--
	}
	return age < 18
}

func (u *User) ValidateAll() error {
	if err := u.Required(); err != nil {
		return err
	}
	if err := u.Validate(); err != nil {
		return err
	}
	return nil
}

func (d *Department) Normalize() {
	d.Name = strings.TrimSpace(d.Name)
	d.Code = strings.TrimSpace(d.Code)
}
func (d *Department) Required() error {
	d.Normalize()

	if d.Name == "" {
		return errors.New("name is required")
	}
	if d.Code == "" {
		return errors.New("code is required")
	}
	return nil
}
func (d *Department) Validate() error {
	if len(d.Name) > 100 {
		return errors.New("name must be at most 100 characters")
	}
	if len(d.Code) > 20 {
		return errors.New("code must be at most 20 characters")
	}
	return nil
}
func (d *Department) ValidateAll() error {
	d.Normalize() // Normalizar UNA sola vez al inicio
	if err := d.Required(); err != nil {
		return err
	}
	if err := d.Validate(); err != nil {
		return err
	}
	return nil
}

func (p *Position) Normalize() {
	p.Name = strings.TrimSpace(p.Name)
	p.Description = strings.TrimSpace(p.Description)
}

func (p *Position) Required() error {
	p.Normalize()

	if p.Name == "" {
		return errors.New("name is required")
	}
	return nil
}

func (p *Position) Validate() error {
	if len(p.Name) > 100 {
		return errors.New("name must be at most 100 characters")
	}
	if len(p.Description) > 255 {
		return errors.New("description must be at most 255 characters")
	}
	return nil
}

func (p *Position) ValidateAll() error {
	p.Normalize()
	if err := p.Required(); err != nil {
		return err
	}
	if err := p.Validate(); err != nil {
		return err
	}
	return nil
}

func (e *Employee) Normalize() {}

func (e *Employee) Required() error {
	if e.UserID == 0 {
		return errors.New("user is required")
	}
	return nil
}

func (e *Employee) Validate() error {
	if e.DepartmentID == 0 {
		return errors.New("department is required")
	}
	if e.PositionID == 0 {
		return errors.New("position is required")
	}
	return nil
}

func (e *Employee) ValidateAll() error {
	if err := e.Required(); err != nil {
		return err
	}
	if err := e.Validate(); err != nil {
		return err
	}
	return nil
}

func (ct *ContractType) Normalize() {
	ct.Name = strings.TrimSpace(ct.Name)
	ct.Description = strings.TrimSpace(ct.Description)
	ct.CountryCode = strings.ToUpper(strings.TrimSpace(ct.CountryCode))
}

func (ct *ContractType) Required() error {
	ct.Normalize()

	if ct.Name == "" {
		return errors.New("name is required")
	}
	if ct.CountryCode == "" {
		return errors.New("country code is required")
	}
	return nil
}

func (ct *ContractType) Validate() error {
	if len(ct.Name) > 100 {
		return errors.New("name must be at most 100 characters")
	}
	if len(ct.CountryCode) != 2 {
		return ErrInvalidCountryCode
	}
	return nil
}

func (ct *ContractType) ValidateAll() error {
	ct.Normalize()
	if err := ct.Required(); err != nil {
		return err
	}
	if err := ct.Validate(); err != nil {
		return err
	}
	return nil
}

func (ec *EmployeeContract) Normalize() {
	ec.CountryCode = strings.ToUpper(strings.TrimSpace(ec.CountryCode))
	ec.Currency = strings.ToUpper(strings.TrimSpace(ec.Currency))
}

func (ec *EmployeeContract) Required() error {
	ec.Normalize()

	if ec.EmployeeID == 0 {
		return errors.New("employee is required")
	}
	if ec.CountryCode == "" {
		return errors.New("country code is required")
	}
	if ec.Currency == "" {
		return errors.New("currency is required")
	}
	return nil
}

func (ec *EmployeeContract) Validate() error {
	if ec.BaseSalary <= 0 {
		return errors.New("base salary must be greater than zero")
	}
	if len(ec.CountryCode) != 2 {
		return ErrInvalidCountryCode
	}
	if len(ec.Currency) != 3 {
		return ErrInvalidCurrency
	}
	if ec.StartDate.IsZero() {
		return errors.New("start date is required")
	}
	if ec.EndDate != nil && !ec.EndDate.IsZero() && ec.EndDate.Before(ec.StartDate) {
		return errors.New("end date must be after start date")
	}
	if ec.WorkHoursPerDay <= 0 || ec.WorkHoursPerDay > 24 {
		return errors.New("work hours per day must be between 1 and 24")
	}
	if ec.WorkDaysPerWeek <= 0 || ec.WorkDaysPerWeek > 7 {
		return errors.New("work days per week must be between 1 and 7")
	}
	return nil
}

func (ec *EmployeeContract) ValidateAll() error {
	ec.Normalize()
	if err := ec.Required(); err != nil {
		return err
	}
	if err := ec.Validate(); err != nil {
		return err
	}
	return nil
}

func (cp *CountryParam) Normalize() {
	cp.CountryCode = strings.ToUpper(strings.TrimSpace(cp.CountryCode))
	cp.Name = strings.TrimSpace(cp.Name)
	cp.Currency = strings.ToUpper(strings.TrimSpace(cp.Currency))
}

func (cp *CountryParam) Required() error {
	cp.Normalize()

	if cp.CountryCode == "" {
		return errors.New("country code is required")
	}
	if cp.Name == "" {
		return errors.New("name is required")
	}
	if cp.Currency == "" {
		return errors.New("currency is required")
	}
	return nil
}

func (cp *CountryParam) Validate() error {
	if len(cp.CountryCode) != 2 {
		return ErrInvalidCountryCode
	}
	if len(cp.Currency) != 3 {
		return ErrInvalidCurrency
	}
	if cp.MinWage <= 0 {
		return errors.New("minimum wage must be greater than zero")
	}
	if cp.HealthRate <= 0 || cp.HealthRate >= 100 {
		return errors.New("health rate must be between 0 and 100")
	}
	if cp.PensionRate <= 0 || cp.PensionRate >= 100 {
		return errors.New("pension rate must be between 0 and 100")
	}
	return nil
}

func (cp *CountryParam) ValidateAll() error {
	cp.Normalize()
	if err := cp.Required(); err != nil {
		return err
	}
	if err := cp.Validate(); err != nil {
		return err
	}
	return nil
}

func (pr *PayrollRecord) Validate() error {
	if pr.EmployeeContractID == 0 {
		return errors.New("contract is required")
	}
	if pr.PeriodStart.IsZero() || pr.PeriodEnd.IsZero() {
		return errors.New("period is required")
	}
	if pr.PeriodEnd.Before(pr.PeriodStart) {
		return ErrInvalidPeriod
	}
	return nil
}

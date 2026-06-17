package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	TenantID  uint           `gorm:"not null;index;uniqueIndex:idx_users_tenant_dni,idx_users_tenant_phone,idx_users_tenant_email" json:"tenant_id"`
	FirstName string         `gorm:"size:30;not null" json:"first_name"`
	LastName  string         `gorm:"size:40;not null" json:"last_name"`
	Dni       string         `gorm:"size:20;not null;uniqueIndex:idx_users_tenant_dni" json:"dni"`
	Gender    string         `gorm:"size:1;not null;check:gender IN ('M', 'F')" json:"gender"`
	Phone     string         `gorm:"size:15;not null;uniqueIndex:idx_users_tenant_phone" json:"phone"`
	Email     string         `gorm:"size:50;not null;uniqueIndex:idx_users_tenant_email" json:"email"`
	BirthDay  time.Time      `gorm:"not null" json:"birth_day"`
	Password     string         `gorm:"-" json:"-"`
	PasswordHash string         `gorm:"size:255;not null" json:"-"`
	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index:idx_users_deleted_at" json:"deleted_at,omitzero"`
}

type Department struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	TenantID  uint           `gorm:"not null;index;uniqueIndex:idx_dept_tenant_code" json:"tenant_id"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	Code      string         `gorm:"size:20;uniqueIndex:idx_dept_tenant_code" json:"code"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitzero"`
	Positions []Position     `gorm:"foreignKey:DepartmentID" json:"positions,omitempty"`
}

type Position struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	TenantID     uint           `gorm:"not null;index" json:"tenant_id"`
	DepartmentID uint           `gorm:"index"`
	Name         string         `gorm:"size:100;not null" json:"name"`
	Description  string         `gorm:"size:255;not null" json:"description"`
	IsActive     bool           `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitzero"`
	Department   Department     `gorm:"foreignKey:DepartmentID"`
}

type Employee struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	TenantID     uint           `gorm:"not null;index" json:"tenant_id"`
	UserID       uint           `gorm:"not null;uniqueIndex" json:"user_id"`
	DepartmentID uint           `gorm:"index" json:"department_id"`
	PositionID   uint           `gorm:"index" json:"position_id"`
	SupervisorID *uint          `gorm:"index" json:"supervisor_id,omitempty"`
	IsActive     bool           `gorm:"default:true;index" json:"is_active"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitzero"`
	User         User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Department   Department     `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	Position     Position       `gorm:"foreignKey:PositionID" json:"position,omitempty"`
	Supervisor   *Employee      `gorm:"foreignKey:SupervisorID" json:"supervisor,omitempty"`
	Contracts    []EmployeeContract `gorm:"foreignKey:EmployeeID" json:"contracts,omitempty"`
}

type ContractType struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Description string    `gorm:"size:255" json:"description"`
	CountryCode string    `gorm:"size:2;not null;index" json:"country_code"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type EmployeeContract struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	TenantID       uint      `gorm:"not null;index" json:"tenant_id"`
	EmployeeID     uint      `gorm:"not null;index" json:"employee_id"`
	ContractTypeID uint      `gorm:"index" json:"contract_type_id"`

	BaseSalary  int64  `gorm:"not null" json:"base_salary"`
	Currency    string `gorm:"size:3" json:"currency"`
	CountryCode string `gorm:"size:2;not null;index" json:"country_code"`

	StartDate time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	IsActive  bool      `gorm:"index" json:"is_active"`

	WorkHoursPerDay float64 `json:"work_hours_per_day"`
	WorkDaysPerWeek float64 `json:"work_days_per_week"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Employee     Employee     `gorm:"foreignKey:EmployeeID" json:"employee,omitempty"`
	ContractType ContractType `gorm:"foreignKey:ContractTypeID" json:"contract_type,omitempty"`
}

type CountryParam struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	CountryCode   string    `gorm:"size:2;not null;uniqueIndex:idx_country_code" json:"country_code"`
	Name          string    `gorm:"size:100;not null" json:"name"`
	Currency      string    `gorm:"size:3;not null" json:"currency"`
	MinWage       int64     `json:"min_wage"`
	HealthRate    float64   `json:"health_rate"`
	PensionRate   float64   `json:"pension_rate"`
	TransportSubs int64     `json:"transport_subs"`
	HousingSubs   int64     `json:"housing_subs"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type Attendance struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	TenantID   uint           `gorm:"not null;index" json:"tenant_id"`
	EmployeeID uint           `gorm:"not null;index;uniqueIndex:idx_att_emp_date" json:"employee_id"`
	Date       time.Time      `gorm:"not null;uniqueIndex:idx_att_emp_date" json:"date"`
	CheckIn    *time.Time     `json:"check_in,omitempty"`
	CheckOut   *time.Time     `json:"check_out,omitempty"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitzero"`
	Employee   Employee       `gorm:"foreignKey:EmployeeID" json:"employee,omitempty"`
}

type Overtime struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	TenantID   uint           `gorm:"not null;index" json:"tenant_id"`
	EmployeeID uint           `gorm:"not null;index;uniqueIndex:idx_ot_emp_date" json:"employee_id"`
	Date       time.Time      `gorm:"not null;uniqueIndex:idx_ot_emp_date" json:"date"`
	Hours      float64        `json:"hours"`
	Reason     string         `gorm:"size:255" json:"reason"`
	Approved   bool           `gorm:"default:false" json:"approved"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitzero"`
	Employee   Employee       `gorm:"foreignKey:EmployeeID" json:"employee,omitempty"`
}

type PayrollRecord struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	TenantID           uint      `gorm:"not null;index" json:"tenant_id"`
	EmployeeContractID uint      `gorm:"not null;index" json:"employee_contract_id"`
	PeriodStart        time.Time `gorm:"not null" json:"period_start"`
	PeriodEnd          time.Time `gorm:"not null" json:"period_end"`

	BaseSalary          int64  `json:"base_salary"`
	Currency            string `gorm:"size:3" json:"currency"`
	CountryCode         string `gorm:"size:2;not null" json:"country_code"`
	GrossSalary         int64  `json:"gross_salary"`
	Bonus               int64  `gorm:"default:0" json:"bonus"`
	Commissions         int64  `gorm:"default:0" json:"commissions"`
	ExtraHourPay        int64  `gorm:"default:0" json:"extra_hour_pay"`
	HealthContribution  int64  `json:"health_contribution"`
	PensionContribution int64  `json:"pension_contribution"`
	SolidarityFund      int64  `gorm:"default:0" json:"solidarity_fund"`
	TransportAllowance  int64  `json:"transport_allowance"`
	HousingAllowance    int64  `json:"housing_allowance"`
	OtherDeductions     int64  `json:"other_deductions"`
	TotalEarnings       int64  `gorm:"default:0" json:"total_earnings"`
	TotalDeductions     int64  `gorm:"default:0" json:"total_deductions"`
	NetSalary           int64  `json:"net_salary"`
	ServiceBonus        int64  `gorm:"default:0" json:"service_bonus"`
	Severance           int64  `gorm:"default:0" json:"severance"`
	SeveranceInterest   int64  `gorm:"default:0" json:"severance_interest"`
	VacationProvision   int64  `gorm:"default:0" json:"vacation_provision"`

	CalculatedAt time.Time      `gorm:"not null" json:"calculated_at"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitzero"`

	EmployeeContract EmployeeContract `gorm:"foreignKey:EmployeeContractID" json:"employee_contract,omitempty"`
}

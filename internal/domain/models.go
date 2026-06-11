package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	TenantID  uint           `gorm:"not null;index" json:"tenant_id"`
	FirstName string         `gorm:"size:30;not null" json:"first_name"`
	LastName  string         `gorm:"size:40;not null" json:"last_name"`
	Dni       string         `gorm:"size:20;not null;uniqueIndex:idx_users_tenant_dni,composite:tenant_dni" json:"dni"`
	Gender    string         `gorm:"size:1;not null;check:gender IN ('M', 'F')" json:"gender"`
	Phone     string         `gorm:"size:15;not null;uniqueIndex:idx_users_tenant_phone,composite:tenant_phone" json:"phone"`
	Email     string         `gorm:"size:50;not null;uniqueIndex:idx_users_tenant_email,composite:tenant_email" json:"email"`
	BirthDay  time.Time      `gorm:"not null" json:"birth_day"`
	Password     string         `gorm:"-" json:"-"`
	PasswordHash string         `gorm:"size:255;not null" json:"-"`
	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index:idx_users_deleted_at" json:"deleted_at,omitzero"`
}

type Department struct {
	ID        uint           `gorm:"primaryKey"`
	TenantID  uint           `gorm:"not null;index"`
	Name      string         `gorm:"size:100;not null"`
	Code      string         `gorm:"size:20;uniqueIndex:idx_dept_tenant_code,composite:tenant_code"`
	IsActive  bool           `gorm:"default:true"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Positions []Position     `gorm:"foreignKey:DepartmentID"`
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
	ID           uint `gorm:"primaryKey"`
	TenantID     uint `gorm:"not null;index"`
	UserID       uint `gorm:"not null;uniqueIndex"`
	DepartmentID uint `gorm:"index"`
	PositionID   uint `gorm:"index"`
	SupervisorID *uint      `gorm:"index"`
	IsActive     bool `gorm:"default:true;index"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt     `gorm:"index"`
	User         User               `gorm:"foreignKey:UserID"`
	Department   Department         `gorm:"foreignKey:DepartmentID"`
	Position     Position           `gorm:"foreignKey:PositionID"`
	Supervisor   *Employee          `gorm:"foreignKey:SupervisorID"`
	Contracts    []EmployeeContract `gorm:"foreignKey:EmployeeID"`
}

type ContractType struct {
	ID          uint      `gorm:"primaryKey"`
	Name        string    `gorm:"size:100;not null"`
	Description string    `gorm:"size:255"`
	CountryCode string    `gorm:"size:2;not null;index"`
	IsActive    bool      `gorm:"default:true"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

type EmployeeContract struct {
	ID             uint `gorm:"primaryKey"`
	TenantID       uint `gorm:"not null;index"`
	EmployeeID     uint `gorm:"not null;index"`
	ContractTypeID uint `gorm:"index"`

	BaseSalary  int64 `gorm:"not null"`
	Currency    string `gorm:"size:3"`
	CountryCode string `gorm:"size:2;not null;index"`

	StartDate time.Time
	EndDate   *time.Time
	IsActive  bool `gorm:"index"`

	WorkHoursPerDay float64
	WorkDaysPerWeek float64

	CreatedAt time.Time
	UpdatedAt time.Time

	Employee     Employee     `gorm:"foreignKey:EmployeeID"`
	ContractType ContractType `gorm:"foreignKey:ContractTypeID"`
}

type CountryParam struct {
	ID            uint   `gorm:"primaryKey"`
	CountryCode   string `gorm:"size:2;not null;uniqueIndex:idx_country_code"`
	Name          string `gorm:"size:100;not null"`
	Currency      string `gorm:"size:3;not null"`
	MinWage       int64
	HealthRate    float64
	PensionRate   float64
	TransportSubs int64
	HousingSubs   int64
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

type Attendance struct {
	ID         uint           `gorm:"primaryKey"`
	TenantID   uint           `gorm:"not null;index"`
	EmployeeID uint           `gorm:"not null;index"`
	Date       time.Time      `gorm:"not null;index:idx_att_emp_date,composite:emp_date"`
	CheckIn    *time.Time
	CheckOut   *time.Time
	CreatedAt  time.Time      `gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	Employee   Employee       `gorm:"foreignKey:EmployeeID"`
}

type Overtime struct {
	ID         uint           `gorm:"primaryKey"`
	TenantID   uint           `gorm:"not null;index"`
	EmployeeID uint           `gorm:"not null;index"`
	Date       time.Time      `gorm:"not null;index:idx_ot_emp_date,composite:emp_date"`
	Hours      float64
	Reason     string         `gorm:"size:255"`
	Approved   bool           `gorm:"default:false"`
	CreatedAt  time.Time      `gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	Employee   Employee       `gorm:"foreignKey:EmployeeID"`
}

type PayrollRecord struct {
	ID                 uint      `gorm:"primaryKey"`
	TenantID           uint      `gorm:"not null;index"`
	EmployeeContractID uint      `gorm:"not null;index"`
	PeriodStart        time.Time `gorm:"not null"`
	PeriodEnd          time.Time `gorm:"not null"`

	BaseSalary          int64
	Currency            string `gorm:"size:3"`
	CountryCode         string `gorm:"size:2;not null"`
	GrossSalary         int64
	Bonus               int64 `gorm:"default:0"`
	Commissions         int64 `gorm:"default:0"`
	ExtraHourPay        int64 `gorm:"default:0"`
	HealthContribution  int64
	PensionContribution int64
	SolidarityFund      int64 `gorm:"default:0"`
	TransportAllowance  int64
	HousingAllowance    int64
	OtherDeductions     int64
	TotalEarnings       int64 `gorm:"default:0"`
	TotalDeductions     int64 `gorm:"default:0"`
	NetSalary           int64
	ServiceBonus        int64 `gorm:"default:0"`
	Severance           int64 `gorm:"default:0"`
	SeveranceInterest   int64 `gorm:"default:0"`
	VacationProvision   int64 `gorm:"default:0"`

	CalculatedAt time.Time      `gorm:"not null"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`

	EmployeeContract EmployeeContract `gorm:"foreignKey:EmployeeContractID"`
}

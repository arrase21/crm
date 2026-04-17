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
	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index:idx_users_deleted_at" json:"deleted_at,omitzero"`
	// Roles     []Role         `gorm:"many2many:user_roles;" json:"roles,omitzero"`
}

// type Department struct {
// 	ID        uint           `gorm:"primaryKey"`
// 	TenantID  uint           `gorm:"not null;index"`
// 	Name      string         `gorm:"size:100;not null"`
// 	Code      string         `gorm:"size:20;uniqueIndex:idx_dept_tenant_code,composite:tenant_code"`
// 	IsActive  bool           `gorm:"default:true"`
// 	CreatedAt time.Time      `gorm:"autoCreateTime"`
// 	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
// 	DeletedAt gorm.DeletedAt `gorm:"index"`
// 	Positions []Position     `gorm:"foreignKey:DepartmentID"`
// }
//
// type Position struct {
// 	ID           uint           `gorm:"primaryKey" json:"id"`
// 	TenantID     uint           `gorm:"not null;index" json:"tenant_id"`
// 	DepartmentID uint           `gorm:"index"`
// 	NamePosition string         `gorm:"size:100;not null" json:"name_position"`
// 	Description  string         `gorm:"size:255;not null" json:"description"`
// 	IsActive     bool           `gorm:"default:true"`
// 	CreatedAt    time.Time      `gorm:"not null" json:"created_at"`
// 	UpdatedAt    time.Time      `gorm:"not null" json:"updated_at"`
// 	DeletedAt    gorm.DeletedAt `gorm:"index:idx_users_deleted_at" json:"deleted_at,omitzero"`
// 	Deparment    Department     `gorm:"foreignKey:DepartmentID"`
// }
//
// type Employee struct {
// 	ID           uint `gorm:"primaryKey"`
// 	TenantID     uint `gorm:"not null;index"`
// 	UserID       uint `gorm:"not null;uniqueIndex"`
// 	DepartmentID uint `gorm:"index"`
// 	PositionID   uint `gorm:"index"`
// 	IsActive     bool `gorm:"default:true;index"`
// 	CreatedAt    time.Time
// 	UpdatedAt    time.Time
// 	DeletedAt    gorm.DeletedAt     `gorm:"index"`
// 	User         User               `gorm:"foreignKey:UserID"`
// 	Department   Department         `gorm:"foreignKey:DepartmentID"`
// 	Position     Position           `gorm:"foreignKey:PositionID"`
// 	Contracts    []EmployeeContract `gorm:"foreignKey:EmployeeID"`
// }
// type EmployeeContract struct {
// 	ID             uint `gorm:"primaryKey"`
// 	TenantID       uint `gorm:"not null;index"`
// 	EmployeeID     uint `gorm:"not null;index"`
// 	ContractTypeID uint `gorm:"index"`
//
// 	BaseSalary float64
// 	Currency   string `gorm:"size:3"`
//
// 	StartDate time.Time
// 	EndDate   *time.Time
// 	IsActive  bool `gorm:"index"`
//
// 	WorkHoursPerDay     float64
// 	WorkDaysPerWeek     float64
// 	HealthContribution  float64
// 	PensionContribution float64
// 	TransportAllowance  float64
// 	HousingAllowance    float64
//
// 	CreatedAt time.Time
// 	UpdatedAt time.Time
//
// 	Employee     Employee     `gorm:"foreignKey:EmployeeID"`
// 	ContractType ContractType `gorm:"foreignKey:ContractTypeID"`
// }
//
// type ContractType struct {
// 	ID          uint   `gorm:"primaryKey"`
// 	Name        string `gorm:"size:50;not null"`
// 	Description string `gorm:"size:255"`
// }

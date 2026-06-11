package domain

import (
	"context"
	"time"
)

// Tenant Context
type contextKey string

const TenantIDKey contextKey = "tenant_id"

type UserRepo interface {
	Create(ctx context.Context, usr *User) error
	GetByID(ctx context.Context, id uint) (*User, error)
	GetByDNI(ctx context.Context, dni string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByPhone(ctx context.Context, phone string) (*User, error)
	List(ctx context.Context, page, limit int) ([]User, int64, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uint) error
}

type DepartmentRepo interface {
	Create(ctx context.Context, dept *Department) error
	GetByID(ctx context.Context, id uint) (*Department, error)
	GetByCode(ctx context.Context, code string) (*Department, error)
	GetByName(ctx context.Context, name string) (*Department, error)
	List(ctx context.Context, page, limit int) ([]Department, int64, error)
	Update(ctx context.Context, dept *Department) error
	Delete(ctx context.Context, id uint) error
}

type PositionRepo interface {
	Create(ctx context.Context, pstn *Position) error
	GetByID(ctx context.Context, id uint) (*Position, error)
	GetByIDWithDepartment(ctx context.Context, id uint) (*Position, error)
	GetByName(ctx context.Context, name string) (*Position, error)
	List(ctx context.Context, page, limit int) ([]Position, int64, error)
	ListByDepartment(ctx context.Context, departmentID uint) ([]Position, int64, error)
	CountByDepartment(ctx context.Context, departmentID uint) (int64, error)
	Update(ctx context.Context, pstn *Position) error
	Delete(ctx context.Context, id uint) error
}

type EmployeeRepo interface {
	Create(ctx context.Context, emp *Employee) error
	GetByID(ctx context.Context, id uint) (*Employee, error)
	GetByUserID(ctx context.Context, userID uint) (*Employee, error)
	List(ctx context.Context, page, limit int) ([]Employee, int64, error)
	ListActive(ctx context.Context, page, limit int) ([]Employee, int64, error)
	ListBySupervisor(ctx context.Context, supervisorEmployeeID uint, page, limit int) ([]Employee, int64, error)
	Update(ctx context.Context, emp *Employee) error
	Delete(ctx context.Context, id uint) error
}

type EmployeeContractRepo interface {
	Create(ctx context.Context, empct *EmployeeContract) error
	GetByID(ctx context.Context, id uint) (*EmployeeContract, error)
	GetByEmployeeID(ctx context.Context, employeeID uint) ([]EmployeeContract, error)
	List(ctx context.Context, page, limit int) ([]EmployeeContract, int64, error)
	Update(ctx context.Context, empct *EmployeeContract) error
	Delete(ctx context.Context, id uint) error
}

type ContractTypeRepo interface {
	Create(ctx context.Context, ct *ContractType) error
	GetByID(ctx context.Context, id uint) (*ContractType, error)
	GetByCountry(ctx context.Context, countryCode string) ([]ContractType, error)
	List(ctx context.Context, page, limit int) ([]ContractType, int64, error)
	Update(ctx context.Context, ct *ContractType) error
	Delete(ctx context.Context, id uint) error
}

type CountryParamRepo interface {
	Create(ctx context.Context, cp *CountryParam) error
	GetByID(ctx context.Context, id uint) (*CountryParam, error)
	GetByCountryCode(ctx context.Context, countryCode string) (*CountryParam, error)
	List(ctx context.Context, page, limit int) ([]CountryParam, int64, error)
	Update(ctx context.Context, cp *CountryParam) error
	Delete(ctx context.Context, id uint) error
}

type AttendanceRepo interface {
	Create(ctx context.Context, a *Attendance) error
	GetByID(ctx context.Context, id uint) (*Attendance, error)
	ListByEmployee(ctx context.Context, employeeID uint, page, limit int) ([]Attendance, int64, error)
	ListByDepartment(ctx context.Context, departmentID uint, page, limit int) ([]Attendance, int64, error)
	List(ctx context.Context, page, limit int) ([]Attendance, int64, error)
	Update(ctx context.Context, a *Attendance) error
	Delete(ctx context.Context, id uint) error
}

type OvertimeRepo interface {
	Create(ctx context.Context, o *Overtime) error
	GetByID(ctx context.Context, id uint) (*Overtime, error)
	ListByEmployee(ctx context.Context, employeeID uint, page, limit int) ([]Overtime, int64, error)
	ListByEmployeeAndPeriod(ctx context.Context, employeeID uint, start, end time.Time) ([]Overtime, error)
	ListByDepartment(ctx context.Context, departmentID uint, page, limit int) ([]Overtime, int64, error)
	List(ctx context.Context, page, limit int) ([]Overtime, int64, error)
	Update(ctx context.Context, o *Overtime) error
	Delete(ctx context.Context, id uint) error
}



type RoleRepo interface {
	GetByName(ctx context.Context, name string) (*Role, error)
	List(ctx context.Context) ([]Role, error)
	GetUserRoles(ctx context.Context, userID uint) ([]string, error)
	AssignRole(ctx context.Context, userID, roleID uint) error
	UnassignRole(ctx context.Context, userID, roleID uint) error
}

type PayrollRecordRepo interface {
	Create(ctx context.Context, pr *PayrollRecord) error
	GetByID(ctx context.Context, id uint) (*PayrollRecord, error)
	GetByContractID(ctx context.Context, contractID uint, page, limit int) ([]PayrollRecord, int64, error)
	List(ctx context.Context, page, limit int) ([]PayrollRecord, int64, error)
}

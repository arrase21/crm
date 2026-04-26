package domain

import (
	"context"
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
	// GetByDNI(ctx context.Context, dni string) (*Employee, error)
	List(ctx context.Context, page, limit int) ([]Employee, int64, error)
	ListActive(ctx context.Context, page, limit int) ([]Employee, int64, error)
	// ListByDepartment(ctx context.Context, deptID uint, page, limit int) ([]Employee, int64, error)
	// ListByPosition(ctx context.Context, positionID uint, page, limit int) ([]Employee, int64, error)
	Update(ctx context.Context, emp *Employee) error
	Delete(ctx context.Context, id uint) error
}

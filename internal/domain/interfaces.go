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
	GetByName(ctx context.Context, name string) (*Position, error)
	List(ctx context.Context, page, limit int) ([]Position, int64, error)
	Update(ctx context.Context, pstn *Position) error
	Delete(ctx context.Context, id uint) error
}

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

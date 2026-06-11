package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/arrase21/crm/internal/domain"
)

func TestTenantFromCtx(t *testing.T) {
	t.Run("tenant in context - success", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), domain.TenantIDKey, uint(1))
		tenantID, err := tenantFromCtx(ctx)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if tenantID != 1 {
			t.Errorf("expected tenant id 1, got %d", tenantID)
		}
	})

	t.Run("no tenant in context - error", func(t *testing.T) {
		_, err := tenantFromCtx(context.Background())
		if err == nil {
			t.Error("expected error but got nil")
		}
	})

	t.Run("tenant is zero - error", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), domain.TenantIDKey, uint(0))
		_, err := tenantFromCtx(ctx)
		if err == nil {
			t.Error("expected error for zero tenant id")
		}
	})

	t.Run("wrong type in context - error", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), domain.TenantIDKey, "not-a-uint")
		_, err := tenantFromCtx(ctx)
		if err == nil {
			t.Error("expected error for wrong type")
		}
	})
}

func TestIsDuplicateError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
		{
			name: "unique constraint error",
			err:  errors.New("UNIQUE constraint failed: users.dni"),
			want: true,
		},
		{
			name: "duplicate key error",
			err:  errors.New(`ERROR: duplicate key value violates unique constraint "idx_users_dni"`),
			want: true,
		},
		{
			name: "other error",
			err:  errors.New("record not found"),
			want: false,
		},
		{
			name: "connection error",
			err:  errors.New("connection refused"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isDuplicateError(tt.err)
			if got != tt.want {
				t.Errorf("isDuplicateError(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}

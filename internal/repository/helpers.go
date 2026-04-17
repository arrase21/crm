package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/arrase21/crm/internal/domain"
)

// tenant validations
func tenatFromctx(ctx context.Context) (uint, error) {
	tenantID, ok := ctx.Value(domain.TenantIDKey).(uint)
	if !ok || tenantID == 0 {
		return 0, errors.New("tenant not fund in context")
	}
	return tenantID, nil
}

// Duplicate validate
func isDuplicateError(err error) bool {
	if err == nil {
		return false
	}
	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "duplicate") ||
		strings.Contains(errMsg, "unique constraint")
}

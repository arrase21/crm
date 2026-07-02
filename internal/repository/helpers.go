package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/arrase21/crm/internal/domain"
	"github.com/jackc/pgx/v5/pgconn"
)

func isDuplicateError(err error) (bool, string) {
	if err == nil {
		return false, ""
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return true, pgErr.ConstraintName
	}

	errMsg := strings.ToLower(err.Error())
	if strings.Contains(errMsg, "duplicate") ||
		strings.Contains(errMsg, "unique constraint") {
		return true, ""
	}

	return false, ""
}

func tenantFromCtx(ctx context.Context) (uint, error) {
	tenantID, ok := ctx.Value(domain.TenantIDKey).(uint)
	if !ok || tenantID == 0 {
		return 0, errors.New("tenant not found in context")
	}
	return tenantID, nil
}



package middleware

import (
	"context"
	"net/http"
	"testing"

	"github.com/arrase21/crm/internal/domain"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupPermissionTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to db: %v", err)
	}
	db.AutoMigrate(&domain.Permission{}, &domain.Role{}, &domain.UserRole{}, &domain.RolePermission{})
	return db
}

func TestRequirePermission_NoClaims(t *testing.T) {
	db := setupPermissionTestDB(t)
	c, w := setupGin()

	middleware := RequirePermission(db, "users", "read")
	middleware(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
	if c.IsAborted() != true {
		t.Error("expected context to be aborted")
	}
}

func TestRequirePermission_SuperAdminBypass(t *testing.T) {
	db := setupPermissionTestDB(t)
	c, w := setupGin()

	ctx := context.WithValue(c.Request.Context(), domain.ClaimsKey, &domain.Claims{
		UserID:   1,
		TenantID: 1,
		Roles:    []string{"super_admin"},
	})
	c.Request = c.Request.WithContext(ctx)

	middleware := RequirePermission(db, "users", "read")
	middleware(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	if c.IsAborted() {
		t.Error("expected context not to be aborted")
	}
}

func TestRequirePermission_NoPermission(t *testing.T) {
	db := setupPermissionTestDB(t)
	db.Create(&domain.Permission{Resource: "users", Action: "read"})

	c, w := setupGin()
	ctx := context.WithValue(c.Request.Context(), domain.ClaimsKey, &domain.Claims{
		UserID:   1,
		TenantID: 1,
		Roles:    []string{"employee"},
	})
	c.Request = c.Request.WithContext(ctx)

	middleware := RequirePermission(db, "users", "read")
	middleware(c)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected status 403, got %d", w.Code)
	}
	if c.IsAborted() != true {
		t.Error("expected context to be aborted")
	}
}

func TestRequirePermission_WithPermission(t *testing.T) {
	db := setupPermissionTestDB(t)

	perm := domain.Permission{Resource: "users", Action: "read"}
	db.Create(&perm)

	role := domain.Role{Name: "employee", TenantID: 1}
	db.Create(&role)

	userRole := domain.UserRole{UserID: 1, RoleID: role.ID}
	db.Create(&userRole)

	rolePerm := domain.RolePermission{RoleID: role.ID, PermissionID: perm.ID}
	db.Create(&rolePerm)

	c, w := setupGin()
	ctx := context.WithValue(c.Request.Context(), domain.ClaimsKey, &domain.Claims{
		UserID:   1,
		TenantID: 1,
		Roles:    []string{"employee"},
	})
	c.Request = c.Request.WithContext(ctx)

	middleware := RequirePermission(db, "users", "read")
	middleware(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	if c.IsAborted() {
		t.Error("expected context not to be aborted")
	}
}

func TestScopeSupervisor_NoClaims(t *testing.T) {
	db := setupPermissionTestDB(t)
	c, w := setupGin()

	middleware := ScopeSupervisor(db)
	middleware(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestScopeSupervisor_AdminBypass(t *testing.T) {
	db := setupPermissionTestDB(t)
	c, w := setupGin()

	ctx := context.WithValue(c.Request.Context(), domain.ClaimsKey, &domain.Claims{
		UserID:   1,
		TenantID: 1,
		Roles:    []string{"super_admin"},
	})
	c.Request = c.Request.WithContext(ctx)

	middleware := ScopeSupervisor(db)
	middleware(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestScopeSupervisor_PayrollManagerBypass(t *testing.T) {
	db := setupPermissionTestDB(t)
	c, w := setupGin()

	ctx := context.WithValue(c.Request.Context(), domain.ClaimsKey, &domain.Claims{
		UserID:   1,
		TenantID: 1,
		Roles:    []string{"payroll_manager"},
	})
	c.Request = c.Request.WithContext(ctx)

	middleware := ScopeSupervisor(db)
	middleware(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

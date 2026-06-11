package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arrase21/crm/internal/domain"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupRoleTest() (*gorm.DB, *gin.Engine) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&domain.User{}, &domain.Role{}, &domain.UserRole{}, &domain.Permission{}, &domain.RolePermission{})

	gin.SetMode(gin.TestMode)
	router := gin.New()

	handler := NewRoleHandler(db)
	router.POST("/roles/assign", handler.Assign)
	router.DELETE("/roles/assign", handler.Unassign)
	router.GET("/roles", handler.ListRoles)
	router.GET("/roles/me", handler.UserRoles)

	return db, router
}

func TestRoleHandler_Assign_Success(t *testing.T) {
	db, router := setupRoleTest()

	db.Create(&domain.User{ID: 1, TenantID: 1, Email: "test@example.com", Dni: "12345678", Phone: "1234567890", FirstName: "Test", LastName: "User", Gender: "M"})
	db.Create(&domain.Role{ID: 1, Name: "employee", TenantID: 0})

	body := map[string]interface{}{
		"user_id":   1,
		"role_name": "employee",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/roles/assign", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), domain.ClaimsKey, &domain.Claims{UserID: 1, TenantID: 1, Roles: []string{"super_admin"}})
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestRoleHandler_Assign_InvalidJSON(t *testing.T) {
	_, router := setupRoleTest()

	req, _ := http.NewRequest("POST", "/roles/assign", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestRoleHandler_Assign_Unauthorized(t *testing.T) {
	_, router := setupRoleTest()

	body := map[string]interface{}{
		"user_id":   1,
		"role_name": "employee",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/roles/assign", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestRoleHandler_Assign_UserNotFound(t *testing.T) {
	db, router := setupRoleTest()

	db.Create(&domain.Role{ID: 1, Name: "employee", TenantID: 0})

	body := map[string]interface{}{
		"user_id":   999,
		"role_name": "employee",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/roles/assign", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), domain.ClaimsKey, &domain.Claims{UserID: 1, TenantID: 1, Roles: []string{"super_admin"}})
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestRoleHandler_Assign_RoleNotFound(t *testing.T) {
	db, router := setupRoleTest()

	db.Create(&domain.User{ID: 1, TenantID: 1, Email: "test@example.com", Dni: "12345678", Phone: "1234567890", FirstName: "Test", LastName: "User", Gender: "M"})

	body := map[string]interface{}{
		"user_id":   1,
		"role_name": "nonexistent_role",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/roles/assign", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), domain.ClaimsKey, &domain.Claims{UserID: 1, TenantID: 1, Roles: []string{"super_admin"}})
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestRoleHandler_Unassign_Success(t *testing.T) {
	db, router := setupRoleTest()

	db.Create(&domain.User{ID: 1, TenantID: 1, Email: "test@example.com", Dni: "12345678", Phone: "1234567890", FirstName: "Test", LastName: "User", Gender: "M"})
	db.Create(&domain.Role{ID: 1, Name: "employee", TenantID: 0})
	db.Create(&domain.UserRole{UserID: 1, RoleID: 1})

	body := map[string]interface{}{
		"user_id":   1,
		"role_name": "employee",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("DELETE", "/roles/assign", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), domain.ClaimsKey, &domain.Claims{UserID: 1, TenantID: 1, Roles: []string{"super_admin"}})
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestRoleHandler_Unassign_NotAssigned(t *testing.T) {
	db, router := setupRoleTest()

	db.Create(&domain.User{ID: 1, TenantID: 1, Email: "test@example.com", Dni: "12345678", Phone: "1234567890", FirstName: "Test", LastName: "User", Gender: "M"})
	db.Create(&domain.Role{ID: 1, Name: "employee", TenantID: 0})

	body := map[string]interface{}{
		"user_id":   1,
		"role_name": "employee",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("DELETE", "/roles/assign", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), domain.ClaimsKey, &domain.Claims{UserID: 1, TenantID: 1, Roles: []string{"super_admin"}})
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestRoleHandler_ListRoles(t *testing.T) {
	db, router := setupRoleTest()

	db.Create(&domain.Role{ID: 1, Name: "admin", TenantID: 0})
	db.Create(&domain.Role{ID: 2, Name: "employee", TenantID: 0})

	req, _ := http.NewRequest("GET", "/roles", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var roles []domain.Role
	if err := json.Unmarshal(w.Body.Bytes(), &roles); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if len(roles) != 2 {
		t.Errorf("expected 2 roles, got %d", len(roles))
	}
}

func TestRoleHandler_UserRoles(t *testing.T) {
	db, router := setupRoleTest()

	db.Create(&domain.Role{ID: 1, Name: "employee", TenantID: 0})
	db.Create(&domain.UserRole{UserID: 1, RoleID: 1})

	req, _ := http.NewRequest("GET", "/roles/me", nil)
	ctx := context.WithValue(req.Context(), domain.ClaimsKey, &domain.Claims{UserID: 1, TenantID: 1, Roles: []string{"employee"}})
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestRoleHandler_UserRoles_Unauthorized(t *testing.T) {
	_, router := setupRoleTest()

	req, _ := http.NewRequest("GET", "/roles/me", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

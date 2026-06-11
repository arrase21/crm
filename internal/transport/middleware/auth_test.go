package middleware

import (
	"net/http"
	"os"
	"testing"

	"github.com/arrase21/crm/internal/domain"
	"github.com/arrase21/crm/internal/service"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupAuthMiddlewareTest(t *testing.T) *service.AuthService {
	os.Setenv("JWT_SECRET", "test-secret-key")
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to db: %v", err)
	}
	db.AutoMigrate(&domain.User{})
	return service.NewAuthService(db)
}

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	svc := setupAuthMiddlewareTest(t)
	defer os.Unsetenv("JWT_SECRET")

	c, w := setupGin()
	middleware := AuthMiddleware(svc)
	middleware(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
	if c.IsAborted() != true {
		t.Error("expected context to be aborted")
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	svc := setupAuthMiddlewareTest(t)
	defer os.Unsetenv("JWT_SECRET")

	c, w := setupGin()
	c.Request.Header.Set("Authorization", "Bearer invalid-token")

	middleware := AuthMiddleware(svc)
	middleware(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
	if c.IsAborted() != true {
		t.Error("expected context to be aborted")
	}
}

func TestAuthMiddleware_WithoutBearerPrefix(t *testing.T) {
	svc := setupAuthMiddlewareTest(t)
	defer os.Unsetenv("JWT_SECRET")

	c, w := setupGin()
	c.Request.Header.Set("Authorization", "Basic some-token")

	middleware := AuthMiddleware(svc)
	middleware(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
	if c.IsAborted() != true {
		t.Error("expected context to be aborted")
	}
}

func TestAuthMiddleware_Success(t *testing.T) {
	svc := setupAuthMiddlewareTest(t)
	defer os.Unsetenv("JWT_SECRET")

	token, err := generateTestToken(svc)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	c, w := setupGin()
	c.Request.Header.Set("Authorization", "Bearer "+token)

	middleware := AuthMiddleware(svc)
	middleware(c)
	c.Next()

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	if c.IsAborted() {
		t.Error("expected context not to be aborted")
	}

	claims, ok := c.Request.Context().Value(domain.ClaimsKey).(*domain.Claims)
	if !ok {
		t.Fatal("expected claims in context")
	}
	if claims.UserID != 1 {
		t.Errorf("expected user id 1, got %d", claims.UserID)
	}
}

func generateTestToken(svc *service.AuthService) (string, error) {
	return svc.GenerateToken(1, 1, []string{"admin"})
}

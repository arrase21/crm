package middleware

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/arrase21/crm/internal/domain"
	"github.com/arrase21/crm/internal/repository"
	"github.com/arrase21/crm/internal/service"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupAuthMiddlewareTest(t *testing.T) (*service.AuthService, *gorm.DB) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to db: %v", err)
	}
	db.AutoMigrate(&domain.User{}, &domain.Role{}, &domain.UserRole{})
	userRepo := repository.NewGormUserRepository(db)
	roleRepo := repository.NewGormRoleRepository(db)
	return service.NewAuthService(userRepo, roleRepo), db
}

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	svc, _ := setupAuthMiddlewareTest(t)
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
	svc, _ := setupAuthMiddlewareTest(t)
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
	svc, _ := setupAuthMiddlewareTest(t)
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
	svc, db := setupAuthMiddlewareTest(t)
	defer os.Unsetenv("JWT_SECRET")

	token, err := generateTestToken(svc, db)
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

func generateTestToken(svc *service.AuthService, db *gorm.DB) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte("test123"), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	db.Create(&domain.User{
		TenantID:     1,
		FirstName:    "Test",
		LastName:     "User",
		Dni:          "12345678",
		Gender:       "M",
		Phone:        "1234567890",
		Email:        "test@example.com",
		PasswordHash: string(hash),
	})

	ctx := context.WithValue(context.Background(), domain.TenantIDKey, uint(1))
	resp, err := svc.Login(ctx, domain.LoginRequest{
		Email:    "test@example.com",
		Password: "test123",
	})
	if err != nil {
		return "", err
	}
	return resp.Token, nil
}

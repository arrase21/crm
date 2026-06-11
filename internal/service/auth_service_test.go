package service

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/arrase21/crm/internal/domain"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupAuthTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to test db: %v", err)
	}
	if err := db.AutoMigrate(&domain.User{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	return db
}

func createTestUser(db *gorm.DB, email, password string) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	db.Create(&domain.User{
		TenantID:     1,
		FirstName:    "Test",
		LastName:     "User",
		Dni:          "12345678",
		Gender:       "M",
		Phone:        "1234567890",
		Email:        email,
		BirthDay:     time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		PasswordHash: string(hash),
	})
}

func TestAuthService_Login_Success(t *testing.T) {
	db := setupAuthTestDB(t)
	createTestUser(db, "test@example.com", "secret123")

	svc := NewAuthService(db)
	ctx := context.Background()

	resp, err := svc.Login(ctx, domain.LoginRequest{
		Email:    "test@example.com",
		Password: "secret123",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Token == "" {
		t.Error("expected token to be non-empty")
	}
	if resp.UserID == 0 {
		t.Error("expected user id to be set")
	}
	if resp.TenantID != 1 {
		t.Errorf("expected tenant id 1, got %d", resp.TenantID)
	}
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	db := setupAuthTestDB(t)
	createTestUser(db, "test@example.com", "secret123")

	svc := NewAuthService(db)
	ctx := context.Background()

	_, err := svc.Login(ctx, domain.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	})
	if err == nil {
		t.Fatal("expected error for wrong password")
	}
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	db := setupAuthTestDB(t)
	svc := NewAuthService(db)
	ctx := context.Background()

	_, err := svc.Login(ctx, domain.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "secret123",
	})
	if err == nil {
		t.Fatal("expected error for non-existent user")
	}
}

func TestAuthService_GenerateAndValidateToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer os.Unsetenv("JWT_SECRET")

	db := setupAuthTestDB(t)
	svc := NewAuthService(db)

	token, err := svc.GenerateToken(1, 1, []string{"admin"})
	if err != nil {
		t.Fatalf("expected no error generating token, got %v", err)
	}

	claims, err := svc.ValidateToken(token)
	if err != nil {
		t.Fatalf("expected no error validating token, got %v", err)
	}
	if claims.UserID != 1 {
		t.Errorf("expected user id 1, got %d", claims.UserID)
	}
	if claims.TenantID != 1 {
		t.Errorf("expected tenant id 1, got %d", claims.TenantID)
	}
	if len(claims.Roles) != 1 || claims.Roles[0] != "admin" {
		t.Errorf("expected roles ['admin'], got %v", claims.Roles)
	}
}

func TestAuthService_ValidateToken_Invalid(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer os.Unsetenv("JWT_SECRET")

	db := setupAuthTestDB(t)
	svc := NewAuthService(db)

	_, err := svc.ValidateToken("invalid-token")
	if err == nil {
		t.Fatal("expected error for invalid token")
	}
}

func TestAuthService_ValidateToken_Empty(t *testing.T) {
	db := setupAuthTestDB(t)
	svc := NewAuthService(db)

	_, err := svc.ValidateToken("")
	if err == nil {
		t.Fatal("expected error for empty token")
	}
}

func TestToStringSlice(t *testing.T) {
	input := []interface{}{"admin", "user", "viewer"}
	result := toStringSlice(input)

	if len(result) != 3 {
		t.Errorf("expected 3 items, got %d", len(result))
	}
	if result[0] != "admin" {
		t.Errorf("expected 'admin', got '%s'", result[0])
	}
	if result[1] != "user" {
		t.Errorf("expected 'user', got '%s'", result[1])
	}
	if result[2] != "viewer" {
		t.Errorf("expected 'viewer', got '%s'", result[2])
	}
}

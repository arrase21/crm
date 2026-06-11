package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/arrase21/crm/internal/domain"
	"github.com/arrase21/crm/internal/service"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupAuthTest(t *testing.T) (*service.AuthService, *gin.Engine) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to db: %v", err)
	}
	db.AutoMigrate(&domain.User{})

	svc := service.NewAuthService(db)
	handler := NewAuthHandler(svc)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/auth/login", handler.Login)

	return svc, router
}

func createLoginUser(db *gorm.DB) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.DefaultCost)
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
}

func TestAuthHandler_Login_Success(t *testing.T) {
	svc, router := setupAuthTest(t)
	defer os.Unsetenv("JWT_SECRET")

	createLoginUser(svc.GetDB())

	body := map[string]interface{}{
		"email":    "test@example.com",
		"password": "secret123",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}

	var resp domain.LoginResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if resp.Token == "" {
		t.Error("expected token to be non-empty")
	}
}

func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	svc, router := setupAuthTest(t)
	defer os.Unsetenv("JWT_SECRET")

	createLoginUser(svc.GetDB())

	body := map[string]interface{}{
		"email":    "test@example.com",
		"password": "wrongpassword",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestAuthHandler_Login_UserNotFound(t *testing.T) {
	_, router := setupAuthTest(t)
	defer os.Unsetenv("JWT_SECRET")

	body := map[string]interface{}{
		"email":    "nonexistent@example.com",
		"password": "secret123",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestAuthHandler_Login_InvalidJSON(t *testing.T) {
	_, router := setupAuthTest(t)
	defer os.Unsetenv("JWT_SECRET")

	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

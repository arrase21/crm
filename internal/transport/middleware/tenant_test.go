package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arrase21/crm/internal/domain"
	"github.com/gin-gonic/gin"
)

func setupGin() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	return c, w
}

func TestTenantMiddleware_MissingHeader(t *testing.T) {
	c, w := setupGin()

	middleware := TenantMiddleware()
	middleware(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
	expected := `{"error":"X-Tenant-ID header is required"}`
	if w.Body.String() != expected {
		t.Errorf("expected body '%s', got '%s'", expected, w.Body.String())
	}
	if c.IsAborted() != true {
		t.Error("expected context to be aborted")
	}
}

func TestTenantMiddleware_InvalidHeader(t *testing.T) {
	c, w := setupGin()
	c.Request.Header.Set("X-Tenant-ID", "abc")

	middleware := TenantMiddleware()
	middleware(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
	expected := `{"error":"invalid X-Tenant-ID"}`
	if w.Body.String() != expected {
		t.Errorf("expected body '%s', got '%s'", expected, w.Body.String())
	}
	if c.IsAborted() != true {
		t.Error("expected context to be aborted")
	}
}

func TestTenantMiddleware_ZeroHeader(t *testing.T) {
	c, w := setupGin()
	c.Request.Header.Set("X-Tenant-ID", "0")

	middleware := TenantMiddleware()
	middleware(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
	if c.IsAborted() != true {
		t.Error("expected context to be aborted")
	}
}

func TestTenantMiddleware_Success(t *testing.T) {
	c, w := setupGin()
	c.Request.Header.Set("X-Tenant-ID", "42")

	middleware := TenantMiddleware()
	middleware(c)
	c.Next()

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	tenantID, ok := c.Request.Context().Value(domain.TenantIDKey).(uint)
	if !ok {
		t.Fatal("expected tenant ID in context")
	}
	if tenantID != 42 {
		t.Errorf("expected tenant ID 42, got %d", tenantID)
	}
	if c.IsAborted() {
		t.Error("expected context not to be aborted")
	}
}

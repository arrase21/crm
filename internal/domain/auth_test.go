package domain

import (
	"testing"
)

func TestLoginRequest_Fields(t *testing.T) {
	req := LoginRequest{
		Email:    "test@example.com",
		Password: "secret123",
	}
	if req.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got '%s'", req.Email)
	}
	if req.Password != "secret123" {
		t.Errorf("expected password 'secret123', got '%s'", req.Password)
	}
}

func TestLoginResponse_Fields(t *testing.T) {
	resp := LoginResponse{
		Token:    "jwt-token",
		UserID:   1,
		TenantID: 1,
		Roles:    []string{"admin", "user"},
	}
	if resp.Token != "jwt-token" {
		t.Errorf("expected token 'jwt-token', got '%s'", resp.Token)
	}
	if resp.UserID != 1 {
		t.Errorf("expected user id 1, got %d", resp.UserID)
	}
	if resp.TenantID != 1 {
		t.Errorf("expected tenant id 1, got %d", resp.TenantID)
	}
	if len(resp.Roles) != 2 {
		t.Errorf("expected 2 roles, got %d", len(resp.Roles))
	}
}

func TestClaims_Fields(t *testing.T) {
	claims := Claims{
		UserID:   1,
		TenantID: 1,
		Roles:    []string{"super_admin"},
	}
	if claims.UserID != 1 {
		t.Errorf("expected user id 1, got %d", claims.UserID)
	}
	if claims.TenantID != 1 {
		t.Errorf("expected tenant id 1, got %d", claims.TenantID)
	}
	if len(claims.Roles) != 1 || claims.Roles[0] != "super_admin" {
		t.Errorf("expected roles ['super_admin'], got %v", claims.Roles)
	}
}

func TestClaimsKey(t *testing.T) {
	if ClaimsKey != "claims" {
		t.Errorf("expected ClaimsKey 'claims', got '%s'", ClaimsKey)
	}
}

package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/arrase21/crm/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo domain.UserRepo
	roleRepo domain.RoleRepo
}

func NewAuthService(userRepo domain.UserRepo, roleRepo domain.RoleRepo) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

func (s *AuthService) Login(ctx context.Context, req domain.LoginRequest) (*domain.LoginResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	roleNames, err := s.roleRepo.GetUserRoles(ctx, user.ID)
	if err != nil {
		return nil, errors.New("failed to get user roles")
	}

	token, err := s.generateToken(user.ID, user.TenantID, roleNames)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &domain.LoginResponse{
		Token:    token,
		UserID:   user.ID,
		TenantID: user.TenantID,
		Roles:    roleNames,
	}, nil
}

func (s *AuthService) generateToken(userID, tenantID uint, roles []string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "change-me-in-production"
	}

	claims := jwt.MapClaims{
		"user_id":   userID,
		"tenant_id": tenantID,
		"roles":     roles,
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (s *AuthService) ValidateToken(tokenStr string) (*domain.Claims, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "change-me-in-production"
	}

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("invalid user_id claim")
	}
	tenantIDFloat, ok := claims["tenant_id"].(float64)
	if !ok {
		return nil, errors.New("invalid tenant_id claim")
	}
	rolesRaw, ok := claims["roles"].([]interface{})
	if !ok {
		return nil, errors.New("invalid roles claim")
	}

	return &domain.Claims{
		UserID:   uint(userIDFloat),
		TenantID: uint(tenantIDFloat),
		Roles:    toStringSlice(rolesRaw),
	}, nil
}

func toStringSlice(v []interface{}) []string {
	out := make([]string, len(v))
	for i, s := range v {
		out[i] = s.(string)
	}
	return out
}

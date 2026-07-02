package domain

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=4"`
}

type LoginResponse struct {
	Token        string   `json:"token"`
	RefreshToken string   `json:"refresh_token"`
	UserID       uint     `json:"user_id"`
	TenantID     uint     `json:"tenant_id"`
	Roles        []string `json:"roles"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type Claims struct {
	UserID   uint
	TenantID uint
	Roles    []string
}

const ClaimsKey contextKey = "claims"

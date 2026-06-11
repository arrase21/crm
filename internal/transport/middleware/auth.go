package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/arrase21/crm/internal/domain"
	"github.com/arrase21/crm/internal/service"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid token"})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")
		claims, err := authService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		ctx := context.WithValue(c.Request.Context(), domain.ClaimsKey, claims)
		ctx = context.WithValue(ctx, domain.TenantIDKey, claims.TenantID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

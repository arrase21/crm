package middleware

import (
	"context"
	"net/http"
	"strconv"

	"github.com/arrase21/crm/internal/domain"
	"github.com/gin-gonic/gin"
)

func TenantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantHeader := c.GetHeader("X-Tenant-ID")
		if tenantHeader == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "X-Tenant-ID header is required"})
			c.Abort()
			return
		}
		tenantID, err := strconv.ParseUint(tenantHeader, 10, 32)
		if err != nil || tenantID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid X-Tenant-ID"})
			c.Abort()
			return
		}
		ctx := context.WithValue(c.Request.Context(), domain.TenantIDKey, uint(tenantID))
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

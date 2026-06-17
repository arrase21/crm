package middleware

import (
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/arrase21/crm/internal/cache"
	"github.com/arrase21/crm/internal/domain"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PermissionMiddleware struct {
	db    *gorm.DB
	cache *cache.Cache
}

func NewPermissionMiddleware(db *gorm.DB) *PermissionMiddleware {
	return &PermissionMiddleware{
		db:    db,
		cache: cache.New(5 * time.Minute),
	}
}

func (pm *PermissionMiddleware) RequirePermission(resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := c.Request.Context().Value(domain.ClaimsKey).(*domain.Claims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		if slices.Contains(claims.Roles, "super_admin") {
			c.Next()
			return
		}

		key := fmt.Sprintf("perm:%d:%d:%s:%s", claims.TenantID, claims.UserID, resource, action)
		if _, found := pm.cache.Get(key); found {
			c.Next()
			return
		}

		var count int64
		pm.db.Table("role_permissions").
			Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
			Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
			Where("user_roles.user_id = ?", claims.UserID).
			Where("permissions.resource = ? AND permissions.action = ?", resource, action).
			Count(&count)

		if count == 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}

		pm.cache.Set(key, true)
		c.Next()
	}
}
func (pm *PermissionMiddleware) InvalidateAll() {
	pm.cache.Clear()
}
func (pm *PermissionMiddleware) InvalidateUser(tenantID, userID uint) {
	prefix := fmt.Sprintf("perm:%d:%d:", tenantID, userID)
	pm.cache.DeleteByPrefix(prefix)
}

func RequirePermission(db *gorm.DB, resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := c.Request.Context().Value(domain.ClaimsKey).(*domain.Claims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		if slices.Contains(claims.Roles, "super_admin") {
			c.Next()
			return
		}

		var count int64
		db.Table("role_permissions").
			Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
			Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
			Where("user_roles.user_id = ?", claims.UserID).
			Where("permissions.resource = ? AND permissions.action = ?", resource, action).
			Count(&count)

		if count == 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func ScopeSupervisor(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := c.Request.Context().Value(domain.ClaimsKey).(*domain.Claims)
		if !ok {
			c.Next()
			return
		}

		for _, role := range claims.Roles {
			if role == "super_admin" || role == "payroll_manager" {
				c.Next()
				return
			}
		}

		var deptID uint
		db.Table("employees").
			Where("user_id = ?", claims.UserID).
			Pluck("department_id", &deptID)

		c.Set("scope_department_id", deptID)
		c.Next()
	}
}

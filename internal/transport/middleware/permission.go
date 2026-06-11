package middleware

import (
	"net/http"

	"github.com/arrase21/crm/internal/domain"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RequirePermission(db *gorm.DB, resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := c.Request.Context().Value(domain.ClaimsKey).(*domain.Claims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		for _, role := range claims.Roles {
			if role == "super_admin" {
				c.Next()
				return
			}
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

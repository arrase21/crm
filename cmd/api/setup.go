package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/arrase21/crm/internal/database"
	"github.com/arrase21/crm/internal/domain"
	"github.com/arrase21/crm/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type setupRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=4"`
}

func newSetupHandler(db *gorm.DB, userSvc *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var count int64
		db.WithContext(c.Request.Context()).Model(&domain.User{}).Where("tenant_id = 1").Count(&count)
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": database.ErrAlreadySetup})
			return
		}

		var req setupRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		admin := &domain.User{
			TenantID:  1,
			FirstName: "Admin",
			LastName:  "System",
			Dni:       "ADMIN-001",
			Gender:    "M",
			Phone:     "0000000000",
			Email:     req.Email,
			BirthDay:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		if err := userSvc.Create(c.Request.Context(), admin, req.Password); err != nil {
			slog.Error("setup: failed to create admin user", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		if err := database.AssignSuperAdminUser(db, admin.ID); err != nil {
			slog.Error("setup: failed to assign super admin role", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "admin user created, you can now login"})
	}
}

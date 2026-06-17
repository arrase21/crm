package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/arrase21/crm/internal/domain"
	"github.com/arrase21/crm/internal/service"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var errAlreadySetup = errors.New("the system is already set up")

func seedPermissionsAndRoles(db *gorm.DB) {
	for i := range domain.Permissions {
		p := &domain.Permissions[i]
		db.Where("resource = ? AND action = ?", p.Resource, p.Action).FirstOrCreate(p)
	}

	roles := []struct {
		Name   string
		Perms  []struct{ resource, action string }
		Tenant bool
	}{
		{
			Name: "super_admin", Tenant: false,
			Perms: []struct{ resource, action string }{
				{"users", "create"}, {"users", "read"}, {"users", "update"}, {"users", "delete"},
				{"employees", "create"}, {"employees", "read"}, {"employees", "update"}, {"employees", "delete"},
				{"departments", "create"}, {"departments", "read"}, {"departments", "update"}, {"departments", "delete"},
				{"positions", "create"}, {"positions", "read"}, {"positions", "update"}, {"positions", "delete"},
				{"contracts", "create"}, {"contracts", "read"}, {"contracts", "update"}, {"contracts", "delete"},
				{"payroll", "calculate"}, {"payroll", "read"},
				{"attendance", "create"}, {"attendance", "read"}, {"attendance", "update"}, {"attendance", "delete"},
				{"overtime", "create"}, {"overtime", "read"}, {"overtime", "update"}, {"overtime", "delete"},
				{"roles", "assign"},
			},
		},
		{
			Name: "payroll_manager", Tenant: true,
			Perms: []struct{ resource, action string }{
				{"payroll", "calculate"}, {"payroll", "read"},
				{"contracts", "create"}, {"contracts", "read"}, {"contracts", "update"},
				{"employees", "read"},
				{"attendance", "read"}, {"attendance", "update"},
				{"overtime", "read"}, {"overtime", "update"},
			},
		},
		{
			Name: "supervisor", Tenant: true,
			Perms: []struct{ resource, action string }{
				{"employees", "read"}, {"employees", "create"}, {"employees", "update"},
				{"payroll", "read"},
				{"attendance", "create"}, {"attendance", "read"},
				{"overtime", "create"}, {"overtime", "read"},
			},
		},
		{
			Name: "employee", Tenant: true,
			Perms: []struct{ resource, action string }{
				{"users", "read"}, {"employees", "read"}, {"payroll", "read"},
			},
		},
	}

	for _, r := range roles {
		role := domain.Role{Name: r.Name}
		if r.Tenant {
			role.TenantID = 1
		}
		db.Where("name = ?", role.Name).FirstOrCreate(&role)

		for _, p := range r.Perms {
			var perm domain.Permission
			db.Where("resource = ? AND action = ?", p.resource, p.action).First(&perm)
			db.Where("role_id = ? AND permission_id = ?", role.ID, perm.ID).
				FirstOrCreate(&domain.RolePermission{RoleID: role.ID, PermissionID: perm.ID})
		}
	}
}

func assignSuperAdminUser(db *gorm.DB, userID uint) error {
	var role domain.Role
	if err := db.Where("name = ?", "super_admin").First(&role).Error; err != nil {
		return err
	}

	return db.Where("user_id = ? AND role_id = ?", userID, role.ID).
		FirstOrCreate(&domain.UserRole{UserID: userID, RoleID: role.ID}).Error
}

func seedMasterData(db *gorm.DB) {
	adminEmail := os.Getenv("ADMIN_EMAIL")
	if adminEmail == "" {
		adminEmail = "admin@crm.com"
	}
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		adminPassword = "admin123"
	}

	var count int64
	db.Model(&domain.User{}).Where("tenant_id = 1").Count(&count)
	if count == 0 {
		hash, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Warning: could not hash admin password: %v", err)
			return
		}

		admin := domain.User{
			TenantID:     1,
			FirstName:    "Admin",
			LastName:     "System",
			Dni:          "ADMIN-001",
			Gender:       "M",
			Phone:        "0000000000",
			Email:        adminEmail,
			BirthDay:     time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			PasswordHash: string(hash),
		}
		if err := db.Create(&admin).Error; err != nil {
			log.Printf("Warning: could not create admin user: %v", err)
		} else if err := assignSuperAdminUser(db, admin.ID); err != nil {
			log.Printf("Warning: could not assign super admin role: %v", err)
		} else {
			log.Printf("Admin user created: %s", adminEmail)
		}
	}

	contractTypes := []domain.ContractType{
		{Name: "Indefinido", Description: "Contrato por tiempo indefinido", CountryCode: "MX", IsActive: true},
		{Name: "Temporal", Description: "Contrato por tiempo determinado", CountryCode: "MX", IsActive: true},
		{Name: "Por obra o tiempo determinado", Description: "Contrato por obra determinada", CountryCode: "MX", IsActive: true},
		{Name: "Término Indefinido", Description: "Contrato a término indefinido", CountryCode: "CO", IsActive: true},
		{Name: "Término Fijo", Description: "Contrato a término fijo", CountryCode: "CO", IsActive: true},
		{Name: "Obra o Labor", Description: "Contrato por obra o labor", CountryCode: "CO", IsActive: true},
	}
	for _, ct := range contractTypes {
		db.Where("name = ? AND country_code = ?", ct.Name, ct.CountryCode).FirstOrCreate(&ct)
	}

	countries := []domain.CountryParam{
		{
			CountryCode: "MX", Name: "México", Currency: "MXN",
			MinWage: 37489, HealthRate: 0.125, PensionRate: 0.1667,
			TransportSubs: 0, HousingSubs: 0,
		},
		{
			CountryCode: "CO", Name: "Colombia", Currency: "COP",
			MinWage: 130000000, HealthRate: 0.04, PensionRate: 0.04,
			TransportSubs: 16200000, HousingSubs: 0,
		},
	}
	for _, c := range countries {
		db.Where("country_code = ?", c.CountryCode).FirstOrCreate(&c)
	}

	log.Println("Master data seeded successfully")
}

type setupRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=4"`
}

func newSetupHandler(db *gorm.DB, userSvc *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var count int64
		db.Model(&domain.User{}).Where("tenant_id = 1").Count(&count)
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": errAlreadySetup.Error()})
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
			Password:  req.Password,
		}

		if err := userSvc.Create(c.Request.Context(), admin); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := assignSuperAdminUser(db, admin.ID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "admin user created, you can now login"})
	}
}

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/arrase21/crm/internal/domain"
	"github.com/arrase21/crm/internal/repository"
	"github.com/arrase21/crm/internal/service"
	"github.com/arrase21/crm/internal/transport/http"
	"github.com/arrase21/crm/internal/transport/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	godotenv.Load()
	dsn := getDSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database connection: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Connected to PostgreSQL successfully")

	// Auto-migrate
	// Nota: Optional
	if err := db.AutoMigrate(
		&domain.User{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migration completed")

	userRepo := repository.NewGormUserRepository(db)
	userSvc := service.NewUserService(userRepo)
	userHandler := http.NewUserHandler(userSvc)
	//Department
	deptRepo := repository.NewGormDepartmentRepository(db)
	deptSvc := service.NewDepartmentService(deptRepo)
	deptHandler := http.NewDepartmentHandler(deptSvc)

	router := gin.Default()

	api := router.Group("/api/v1")
	api.Use(middleware.TenantMiddleware())

	{
		users := api.Group("/users")
		{
			users.POST("", userHandler.Create)
			users.GET("/:id", userHandler.GetByID)
			users.GET("", userHandler.GetByDni) // ?dni=xxx
			users.GET("/list", userHandler.List)
			users.PUT("/:id", userHandler.Update)
			users.DELETE("/:id", userHandler.Delete)
		}
		departments := api.Group("/departments")
		{
			departments.POST("", deptHandler.Create)
			departments.GET("/:id", deptHandler.GetByID)
			departments.GET("", deptHandler.GetByCode) // ?code=xxx
			departments.GET("/list", deptHandler.List)
			departments.PUT("/:id", deptHandler.Update)
			departments.DELETE("/:id", deptHandler.Delete)
		}
	}

	// Iniciar servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getDSN() string {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")
	// host := getEnv("DB_HOST", "")
	// port := getEnv("DB_PORT", "5432")
	// user := getEnv("DB_USER", "postgres")
	// password := getEnv("DB_PASSWORD", "postgres")
	// dbname := getEnv("DB_NAME", "crm")
	// sslmode := getEnv("DB_SSLMODE", "disable")

	if host == "" || user == "" || password == "" || dbname == "" {
		log.Fatal("Environment missing")
	}

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

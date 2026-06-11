package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/arrase21/crm/internal/repository"
	"github.com/arrase21/crm/internal/service"
	httptransport "github.com/arrase21/crm/internal/transport/http"
	"github.com/arrase21/crm/internal/transport/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func main() {
	godotenv.Load()

	logLevel := &slog.LevelVar{}
	logLevel.Set(slog.LevelInfo)
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})))

	dsn := getDSN()

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Warn),
	})
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	db.Exec("SET search_path TO crm,public")

	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("failed to get database connection", "error", err)
		os.Exit(1)
	}

	sqlDB.SetMaxOpenConns(getEnvInt("DB_MAX_OPEN_CONNS", 25))
	sqlDB.SetMaxIdleConns(getEnvInt("DB_MAX_IDLE_CONNS", 10))
	sqlDB.SetConnMaxLifetime(getEnvDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute))

	if err := sqlDB.Ping(); err != nil {
		slog.Error("failed to ping database", "error", err)
		os.Exit(1)
	}
	slog.Info("connected to database")

	m := gormigrate.New(db, gormigrate.DefaultOptions, getMigrations())
	if err := m.Migrate(); err != nil {
		slog.Error("migration failed", "error", err)
		os.Exit(1)
	}
	slog.Info("migrations completed")

	seedPermissionsAndRoles(db)

	seedEnabled := getEnv("SEED_ENABLED", "true")
	if seedEnabled == "true" {
		seedMasterData(db)
	}

	userRepo := repository.NewGormUserRepository(db)
	userSvc := service.NewUserService(userRepo)
	userHandler := httptransport.NewUserHandler(userSvc)

	pstnRepo := repository.NewGormPositionRepository(db)
	deptRepo := repository.NewGormDepartmentRepository(db)
	deptSvc := service.NewDepartmentService(deptRepo, pstnRepo)
	deptHandler := httptransport.NewDepartmentHandler(deptSvc)
	pstnSvc := service.NewPositionService(pstnRepo, deptRepo)
	pstnHandler := httptransport.NewPositionHandler(pstnSvc)

	empRepo := repository.NewGormEmployeeRepository(db)
	empSvc := service.NewEmployeeService(empRepo, userRepo, deptRepo, pstnRepo)
	empHandler := httptransport.NewEmployeeHandler(empSvc)

	ctRepo := repository.NewGormContractTypeRepository(db)
	countryRepo := repository.NewGormCountryParamRepository(db)
	contractRepo := repository.NewGormEmployeeContractRepository(db)
	payrollRepo := repository.NewGormPayrollRecordRepository(db)
	otRepo := repository.NewGormOvertimeRepository(db)
	contractSvc := service.NewEmployeeContractService(contractRepo, empRepo, ctRepo, countryRepo, payrollRepo, otRepo)
	contractHandler := httptransport.NewEmployeeContractHandler(contractSvc)

	attRepo := repository.NewGormAttendanceRepository(db)
	attSvc := service.NewAttendanceService(attRepo, empRepo)
	attHandler := httptransport.NewAttendanceHandler(attSvc)

	otSvc := service.NewOvertimeService(otRepo, empRepo)
	otHandler := httptransport.NewOvertimeHandler(otSvc)

	roleRepo := repository.NewGormRoleRepository(db)
	authSvc := service.NewAuthService(userRepo, roleRepo)
	authHandler := httptransport.NewAuthHandler(authSvc)
	roleSvc := service.NewRoleService(roleRepo, userRepo)
	roleHandler := httptransport.NewRoleHandler(roleSvc)

	ginMode := getEnv("GIN_MODE", "release")
	gin.SetMode(ginMode)
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	rl := middleware.NewRateLimiter(rate.Limit(100), 200, 10*time.Minute)
	router.Use(rl.Middleware())

	router.GET("/health", func(c *gin.Context) {
		if err := sqlDB.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	router.GET("/ready", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Tenant-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := router.Group("/api/v1")
	api.Use(middleware.TenantMiddleware())

	api.POST("/auth/login", authHandler.Login)

	if seedEnabled != "true" {
		api.POST("/setup", newSetupHandler(db, userSvc))
	}

	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(authSvc))
	{
		roles := protected.Group("/roles")
		{
			roles.POST("/assign", middleware.RequirePermission(db, "roles", "assign"), roleHandler.Assign)
			roles.DELETE("/assign", middleware.RequirePermission(db, "roles", "assign"), roleHandler.Unassign)
			roles.GET("", roleHandler.ListRoles)
			roles.GET("/me", roleHandler.UserRoles)
		}
		users := protected.Group("/users")
		{
			users.POST("", middleware.RequirePermission(db, "users", "create"), userHandler.Create)
			users.GET("/:id", middleware.RequirePermission(db, "users", "read"), userHandler.GetByID)
			users.GET("", middleware.RequirePermission(db, "users", "read"), userHandler.GetByDni)
			users.GET("/list", middleware.RequirePermission(db, "users", "read"), userHandler.List)
			users.PUT("/:id", middleware.RequirePermission(db, "users", "update"), userHandler.Update)
			users.DELETE("/:id", middleware.RequirePermission(db, "users", "delete"), userHandler.Delete)
		}
		departments := protected.Group("/departments")
		{
			departments.POST("", middleware.RequirePermission(db, "departments", "create"), deptHandler.Create)
			departments.GET("/:id", middleware.RequirePermission(db, "departments", "read"), deptHandler.GetByID)
			departments.GET("", middleware.RequirePermission(db, "departments", "read"), deptHandler.GetByCode)
			departments.GET("/list", middleware.RequirePermission(db, "departments", "read"), deptHandler.List)
			departments.PUT("/:id", middleware.RequirePermission(db, "departments", "update"), deptHandler.Update)
			departments.DELETE("/:id", middleware.RequirePermission(db, "departments", "delete"), deptHandler.Delete)
		}
		positions := protected.Group("/positions")
		{
			positions.POST("", middleware.RequirePermission(db, "positions", "create"), pstnHandler.Create)
			positions.GET("/:id", middleware.RequirePermission(db, "positions", "read"), pstnHandler.GetByID)
			positions.GET("/search", middleware.RequirePermission(db, "positions", "read"), pstnHandler.GetByName)
			positions.GET("/list", middleware.RequirePermission(db, "positions", "read"), pstnHandler.List)
			positions.PUT("/:id", middleware.RequirePermission(db, "positions", "update"), pstnHandler.Update)
			positions.DELETE("/:id", middleware.RequirePermission(db, "positions", "delete"), pstnHandler.Delete)
		}
		employees := protected.Group("/employees")
		{
			employees.POST("", middleware.RequirePermission(db, "employees", "create"), empHandler.Create)
			employees.GET("/:id", middleware.RequirePermission(db, "employees", "read"), empHandler.GetByID)
			employees.GET("", middleware.RequirePermission(db, "employees", "read"), empHandler.GetByUserID)
			employees.GET("/list", middleware.RequirePermission(db, "employees", "read"), empHandler.List)
			employees.PUT("/:id", middleware.RequirePermission(db, "employees", "update"), empHandler.Update)
			employees.DELETE("/:id", middleware.RequirePermission(db, "employees", "delete"), empHandler.Delete)
		}
		contracts := protected.Group("/contracts")
		{
			contracts.POST("", middleware.RequirePermission(db, "contracts", "create"), contractHandler.Create)
			contracts.GET("/:id", middleware.RequirePermission(db, "contracts", "read"), contractHandler.GetByID)
			contracts.GET("", middleware.RequirePermission(db, "contracts", "read"), contractHandler.GetByEmployeeID)
			contracts.GET("/list", middleware.RequirePermission(db, "contracts", "read"), contractHandler.List)
			contracts.PUT("/:id", middleware.RequirePermission(db, "contracts", "update"), contractHandler.Update)
			contracts.DELETE("/:id", middleware.RequirePermission(db, "contracts", "delete"), contractHandler.Delete)
		}
		payroll := protected.Group("/payroll")
		{
			payroll.POST("/calculate", middleware.RequirePermission(db, "payroll", "calculate"), contractHandler.CalculatePayroll)
			payroll.GET("/records", middleware.RequirePermission(db, "payroll", "read"), contractHandler.GetPayrollRecords)
		}
		attendance := protected.Group("/attendance")
		{
			attendance.POST("", middleware.RequirePermission(db, "attendance", "create"), attHandler.Create)
			attendance.GET("/:id", middleware.RequirePermission(db, "attendance", "read"), attHandler.GetByID)
			attendance.GET("", middleware.RequirePermission(db, "attendance", "read"), attHandler.List)
			attendance.PUT("/:id", middleware.RequirePermission(db, "attendance", "update"), attHandler.Update)
			attendance.DELETE("/:id", middleware.RequirePermission(db, "attendance", "delete"), attHandler.Delete)
		}
		overtime := protected.Group("/overtime")
		{
			overtime.POST("", middleware.RequirePermission(db, "overtime", "create"), otHandler.Create)
			overtime.GET("/:id", middleware.RequirePermission(db, "overtime", "read"), otHandler.GetByID)
			overtime.GET("", middleware.RequirePermission(db, "overtime", "read"), otHandler.List)
			overtime.PUT("/:id", middleware.RequirePermission(db, "overtime", "update"), otHandler.Update)
			overtime.DELETE("/:id", middleware.RequirePermission(db, "overtime", "delete"), otHandler.Delete)
		}
	}

	port := getEnv("PORT", "8080")
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		slog.Info("server starting", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced shutdown", "error", err)
		os.Exit(1)
	}

	if err := sqlDB.Close(); err != nil {
		slog.Error("database connection close error", "error", err)
	}
	slog.Info("server stopped")
}

func getDSN() string {
	host := getEnv("DB_HOST", "")
	port := getEnv("DB_PORT", "")
	user := getEnv("DB_USER", "")
	password := getEnv("DB_PASSWORD", "")
	dbname := getEnv("DB_NAME", "")
	sslmode := getEnv("DB_SSLMODE", "disable")

	if host == "" || user == "" || password == "" || dbname == "" {
		slog.Error("missing required database environment variables")
		os.Exit(1)
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

func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	n, err := strconv.Atoi(value)
	if err != nil {
		slog.Warn("invalid env value, using default", "key", key, "value", value, "default", defaultValue)
		return defaultValue
	}
	return n
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	d, err := time.ParseDuration(value)
	if err != nil {
		slog.Warn("invalid env duration, using default", "key", key, "value", value, "default", defaultValue)
		return defaultValue
	}
	return d
}

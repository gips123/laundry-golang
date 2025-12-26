package main

import (
	"fmt"
	"laundry-go/internal/config"
	"laundry-go/internal/database"
	"laundry-go/internal/handlers"
	"laundry-go/internal/middleware"
	"laundry-go/internal/models"
	"laundry-go/internal/repository"
	"laundry-go/internal/service"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Set Gin mode
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Connect to database
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate
	err = db.AutoMigrate(
		&models.User{},
		&models.Laundry{},
		&models.Service{},
		&models.Order{},
		&models.OrderService{},
		&models.Review{},
	)
	if err != nil {
		// Check if error is just "relation already exists" - this is OK
		errStr := err.Error()
		if contains(errStr, "already exists") || contains(errStr, "42P07") {
			log.Println("Tables already exist, skipping migration...")
		} else {
			log.Fatalf("Failed to migrate database: %v", err)
		}
	} else {
		log.Println("Database migration completed successfully")
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	laundryRepo := repository.NewLaundryRepository(db)
	serviceRepo := repository.NewServiceRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	orderServiceRepo := repository.NewOrderServiceRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg)
	laundryService := service.NewLaundryService(laundryRepo, serviceRepo, userRepo)
	orderService := service.NewOrderService(orderRepo, orderServiceRepo, serviceRepo, laundryRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	laundryHandler := handlers.NewLaundryHandler(laundryService)
	orderHandler := handlers.NewOrderHandler(orderService)

	// Setup router
	router := gin.Default()

	// CORS middleware
	router.Use(middleware.CORSMiddleware(cfg))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/me", middleware.AuthMiddleware(cfg), authHandler.GetMe)
			auth.PATCH("/update-location", middleware.AuthMiddleware(cfg), authHandler.UpdateLocation)
		}

		// Laundry routes
		laundries := api.Group("/laundries")
		{
			laundries.GET("", laundryHandler.GetAll)
			laundries.GET("/:id", laundryHandler.GetByID)
		}

		// Order routes (protected)
		orders := api.Group("/orders")
		orders.Use(middleware.AuthMiddleware(cfg))
		{
			orders.POST("", orderHandler.Create)
			orders.GET("", orderHandler.GetAll)
			orders.GET("/:id", orderHandler.GetByID)
			orders.PATCH("/:id/cancel", orderHandler.Cancel)
			orders.PATCH("/:id/status", middleware.RequireRole("laundry_owner"), orderHandler.UpdateStatus)
		}
	}

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

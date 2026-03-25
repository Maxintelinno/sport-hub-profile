package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/maxintelinno/sport-hub-profile/internal/handlers"
	customMiddleware "github.com/maxintelinno/sport-hub-profile/internal/middleware"
	"github.com/maxintelinno/sport-hub-profile/internal/models"
	"github.com/maxintelinno/sport-hub-profile/internal/repositories"
	"github.com/maxintelinno/sport-hub-profile/internal/services"
	"github.com/maxintelinno/sport-hub-profile/pkg/config"
	"github.com/maxintelinno/sport-hub-profile/pkg/database"
)

func main() {
	cfg := config.LoadConfig()

	e := echo.New()

	// Initialize DB
	db, err := database.InitDB()
	if err != nil {
		e.Logger.Fatal(err)
	}

	// Migrations
	db.AutoMigrate(&models.OTPRequest{})

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Repositories
	userRepo := repositories.NewUserRepository(db)
	reportRepo := repositories.NewReportRepository(db)
	dashboardRepo := repositories.NewDashboardRepository(db)

	// Services
	profileService := services.NewProfileService(userRepo)
	reportService := services.NewReportService(reportRepo)
	dashboardService := services.NewDashboardService(userRepo, dashboardRepo)
	authService := services.NewAuthService(userRepo)

	// Handlers
	healthHandler := handlers.NewHealthHandler()
	profileHandler := handlers.NewProfileHandler(profileService)
	reportHandler := handlers.NewReportHandler(reportService)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)
	authHandler := handlers.NewAuthHandler(authService)

	// Routes
	e.GET("/health", healthHandler.Check)
	
	// Public routes
	e.POST("/v1/auth/forgot-password", authHandler.ForgotPassword)
	e.POST("/v1/auth/check-phone", authHandler.CheckPhone)
	
	// Authenticated routes
	v1Auth := e.Group("/v1")
	v1Auth.Use(customMiddleware.JWTMiddleware(cfg.JwtSecret))
	v1Auth.GET("/profile", profileHandler.GetProfile)
	v1Auth.GET("/reports/revenue", reportHandler.GetRevenueReport)
	v1Auth.GET("/owner/dashboard", dashboardHandler.GetDashboard)

	// Start server
	go func() {
		if err := e.Start(fmt.Sprintf(":%s", cfg.AppPort)); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

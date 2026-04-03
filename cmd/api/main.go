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
	db.AutoMigrate(&models.OTPRequest{}, &models.OwnerBankAccount{}, &models.OwnerSettlement{}, &models.OwnerPayout{}, &models.OwnerStaff{})

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Repositories
	userRepo := repositories.NewUserRepository(db)
	reportRepo := repositories.NewReportRepository(db)
	dashboardRepo := repositories.NewDashboardRepository(db)
	bankRepo := repositories.NewBankRepository(db)
	payoutRepo := repositories.NewPayoutRepository(db)
	ownerStaffRepo := repositories.NewOwnerStaffRepository(db)

	// Services
	profileService := services.NewProfileService(userRepo)
	reportService := services.NewReportService(reportRepo)
	dashboardService := services.NewDashboardService(userRepo, dashboardRepo)
	authService := services.NewAuthService(userRepo)
	bankService := services.NewBankService(bankRepo, userRepo, payoutRepo)
	ownerStaffService := services.NewOwnerStaffService(ownerStaffRepo)

	// Handlers
	healthHandler := handlers.NewHealthHandler()
	profileHandler := handlers.NewProfileHandler(profileService)
	reportHandler := handlers.NewReportHandler(reportService)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)
	authHandler := handlers.NewAuthHandler(authService)
	bankHandler := handlers.NewBankHandler(bankService)
	ownerStaffHandler := handlers.NewOwnerStaffHandler(ownerStaffService)

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
	
	// Bank Account routes
	v1Auth.GET("/owner/bank-accounts", bankHandler.GetBankAccounts)
	v1Auth.POST("/owner/bank-accounts", bankHandler.AddBankAccount)
	v1Auth.PUT("/owner/bank-accounts/:id", bankHandler.UpdateBankAccount)
	v1Auth.DELETE("/owner/bank-accounts/:id", bankHandler.DeleteBankAccount)
	v1Auth.POST("/owner/bank-accounts/:id/set-default", bankHandler.SetDefaultBankAccount)
	v1Auth.GET("/owner/staff", ownerStaffHandler.GetStaff)

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

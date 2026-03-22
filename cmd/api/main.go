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

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Repositories
	userRepo := repositories.NewUserRepository(db)

	// Services
	profileService := services.NewProfileService(userRepo)

	// Handlers
	healthHandler := handlers.NewHealthHandler()
	profileHandler := handlers.NewProfileHandler(profileService)

	// Routes
	e.GET("/health", healthHandler.Check)
	
	v1 := e.Group("/v1")
	v1.Use(customMiddleware.JWTMiddleware(cfg.JwtSecret))
	v1.GET("/profile", profileHandler.GetProfile)

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

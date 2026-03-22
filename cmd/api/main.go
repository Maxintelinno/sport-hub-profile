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
	"github.com/maxintelinno/sport-hub-profile/pkg/config"
)

func main() {

	cfg := config.LoadConfig()

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Handlers
	healthHandler := handlers.NewHealthHandler()

	// Routes
	e.GET("/health", healthHandler.Check)

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

	// Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	e.Logger.Fatal(e.Start(":" + port))
}

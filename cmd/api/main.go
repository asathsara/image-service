package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"image-service/internal/config"
	"image-service/internal/handler"
	"image-service/internal/middleware"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Initialize Configuration
	cfg := config.LoadConfig()

	// 1. Initialize VIPS
	// This is important because govips uses the C libvips engine
	vips.Startup(nil)
	defer vips.Shutdown()

	// 2. Initialize Fiber App
	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024, // 10MB limit (matches Fastify)
	})

	// Global Middleware
	app.Use(logger.New())

	// Health Check (Public)
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	// Protected Routes
	api := app.Group("/", middleware.Auth(cfg))
	api.Post("/compress", handler.CompressHandler)

	// Handle Graceful Shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Println("Shutting down server...")
		_ = app.Shutdown()
	}()

	// Start Server
	log.Printf("Server running on port %s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

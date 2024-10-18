package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/mysterybee07/result-distribution-system/initializers"
	"github.com/mysterybee07/result-distribution-system/routes"
)

func init() {
	initializers.Connect()
	initializers.LoadEnvironment()
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("The port is taken by another process.")
		port = "8080"
	}
	log.Println("Starting the server on port " + port + "..........")

	// Load templates
	// engine := html.New("./resources/views", ".html")
	// engine.AddFunc("add", func(values ...int) int {
	// 	sum := 0
	// 	for _, v := range values {
	// 		sum += v
	// 	}
	// 	return sum
	// })

	app := fiber.New(fiber.Config{
		// Views: engine,
	})

	// Initialize session store
	store := session.New()

	// Use session middleware
	app.Use(func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			return err
		}
		c.Locals("session", sess)
		return c.Next()
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000, http://localhost:5173/",
		// AllowOrigins:     "",
		AllowMethods:     "GET, POST, PUT, DELETE, PATCH, OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
	}))

	// // Use flash messages middleware
	// app.Use(middleware.FlashMessages)
	// // Loading static files
	// app.Static("/", "./static")
	// // Loading images
	// app.Static("/static", "./static")

	// Authentication routes
	routes.Home(app)

	// Protected routes
	// app.Use(middleware.AuthRequired)

	routes.SetupRoutes(app)

	// Start the server and handle graceful shutdown
	go func() {
		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Graceful shutdown handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Create a context with a timeout to allow for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the Fiber app gracefully
	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}

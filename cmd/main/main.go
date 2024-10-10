package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"
	"github.com/mysterybee07/result-distribution-system/initializers"
	"github.com/mysterybee07/result-distribution-system/middleware"
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
	engine := html.New("./resources/views", ".html")
	engine.AddFunc("add", func(values ...int) int {
		sum := 0
		for _, v := range values {
			sum += v
		}
		return sum
	})

	app := fiber.New(fiber.Config{
		Views: engine,
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
	// Use flash messages middleware
	app.Use(middleware.FlashMessages)
	// Loading static files
	app.Static("/", "./static")
	// Loading images
	app.Static("/static", "./static")

	// Authentication routes
	routes.Home(app)

	// Protected routes
	// app.Use(middleware.AuthRequired)

	// Routes
	// routes.Profile(app)
	// routes.Dashboard(app)
	// routes.Student(app)
	// routes.Batch(app)
	// routes.Program(app)
	// routes.Semester(app)
	// routes.Course(app)
	// routes.Mark(app)
	// routes.Result(app)
	// routes.Error(app)
	routes.SetupRoutes(app)

	err := app.Listen(":" + port)
	if err != nil {
		log.Fatalf("Server failed to listen: %v", err)
	}
	log.Println("Server exited")
}

package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
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
		log.Fatal("The port:" + port + " is taken by another process.")
		port = "8080"
	}
	log.Println("Starting the server on port " + port + "..........")

	//Load templates
	engine := html.New("./resources/views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
		// ViewsLayout:""
	})

	app.Static("/", "./static")
	routes.Home(app)
	routes.Profile(app)
	routes.Dashboard(app)
	routes.Student(app)
	routes.Batch(app)
	routes.Program(app)
	err := app.Listen(":" + port)
	if err != nil {
		log.Fatalf("Server failed to listen: %v", err)

	}
	log.Println("Server exited")

}

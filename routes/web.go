package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/controllers"
)

func Home(app *fiber.App) {
	app.Get("", controllers.Home)
}

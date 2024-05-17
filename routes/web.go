package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/controllers"
)

func Home(app *fiber.App) {
	app.Get("", controllers.Home)
	app.Get("/register", controllers.Register)
	app.Post("/register", controllers.StoreRegister)
	app.Get("/login", controllers.Login)
	app.Post("/login", controllers.LoginUser)

}

func Profile(app *fiber.App) {
	app.Get("/profile/", controllers.Profile)
}

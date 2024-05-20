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
	app.Get("/forgot-password", controllers.ForgotPassword)
	// app.Get("/users/:id", controllers.GetUserById)

}

func Profile(app *fiber.App) {
	app.Get("/profile/", controllers.Profile)
}

func Admin(app *fiber.App) {
	app.Get("/dashboard", controllers.Index)
}

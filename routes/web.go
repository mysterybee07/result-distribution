package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/controllers"
	"github.com/mysterybee07/result-distribution-system/middleware"
)

func Home(app *fiber.App) {
	app.Get("", controllers.Home)
	app.Get("/register", controllers.Register)
	app.Post("/register", controllers.StoreRegister)
	app.Get("/login", controllers.Login)
	app.Post("/login", controllers.LoginUser)
	app.Get("/logout", controllers.LogoutUser)
	app.Get("/forgot-password", controllers.ForgotPassword)
}

func Profile(app *fiber.App) {
	app.Get("/profile", middleware.AuthRequired, controllers.GetUserProfile)
}

func Dashboard(app *fiber.App) {
	app.Get("/dashboard", middleware.AuthRequired, middleware.AdminRequired, controllers.Index)
}

func Student(app *fiber.App) {
	app.Get("/students/add", middleware.AuthRequired, middleware.AdminRequired, controllers.AddStudent)
	app.Post("/students/add", middleware.AuthRequired, middleware.AdminRequired, controllers.StoreStudent)
	app.Get("/students", middleware.AuthRequired, middleware.AdminRequired, controllers.GetStudents)
	// app.Get("/students/edit/:id", middleware.AuthRequired, middleware.AdminRequired, controllers.EditStudentForm)
	app.Put("/students/edit/:id", middleware.AuthRequired, middleware.AdminRequired, controllers.UpdateStudent)
}

func Batch(app *fiber.App) {
	app.Get("/batches/add", middleware.AuthRequired, middleware.AdminRequired, controllers.AddBatch)
	app.Post("/batches/add", middleware.AuthRequired, middleware.AdminRequired, controllers.CreateBatch)
}

func Program(app *fiber.App) {
	app.Get("/programs/add", middleware.AuthRequired, middleware.AdminRequired, controllers.AddProgram)
	app.Post("/programs/add", middleware.AuthRequired, middleware.AdminRequired, controllers.StoreProgram)
}

func Semester(app *fiber.App) {
	app.Get("/semesters/add", middleware.AuthRequired, middleware.AdminRequired, controllers.AddSemester)
	app.Post("/semesters/add", middleware.AuthRequired, middleware.AdminRequired, controllers.StoreSemester)
}

func Subject(app *fiber.App) {
	app.Get("/courses/add", middleware.AuthRequired, middleware.AdminRequired, controllers.AddCourse)
	app.Post("/courses/add", middleware.AuthRequired, middleware.AdminRequired, controllers.StoreCourse)
}

func Mark(app *fiber.App) {
	app.Post("/marks/add", middleware.AuthRequired, middleware.AdminRequired, controllers.CreateMarks)
	app.Put("/marks/edit/:id", middleware.AuthRequired, middleware.AdminRequired, controllers.UpdateMarks)
	app.Get("/marks/:symbolNumber", middleware.AuthRequired, middleware.AdminRequired, controllers.GetMarksBySymbolNumber)
	app.Post("/publish-results", middleware.AuthRequired, middleware.AdminRequired, controllers.PublishResults)
}

func Result(app *fiber.App) {
	app.Get("/results", middleware.AuthRequired, middleware.AdminRequired, controllers.AddResult)
	app.Post("/results", middleware.AuthRequired, middleware.AdminRequired, controllers.PublishResults)
}

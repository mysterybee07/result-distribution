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
	app.Get("/forgot-password", controllers.ForgotPassword)
	// app.Get("/users/:id", controllers.GetUserById)

}

func Profile(app *fiber.App) {
	// app.Get("/profile/", controllers.Profile)
	app.Get("/profile", middleware.AuthRequired, controllers.GetUserProfile)
}

func Dashboard(app *fiber.App) {
	app.Get("/dashboard", controllers.Index)
}

func Student(app *fiber.App) {
	app.Get("/students/add", controllers.AddStudent)
	app.Post("/students/add", controllers.StoreStudent)
	app.Get("/students", controllers.GetStudents)
	app.Get("/students/edit/:id", controllers.EditStudent)
	app.Put("/students/edit/:id", controllers.UpdateStudent)
}

func Batch(app *fiber.App) {
	app.Get("/batches/add", controllers.AddBatch)
	app.Post("/batches/add", controllers.CreateBatch)
}

func Program(app *fiber.App) {
	app.Get("/programs/add", controllers.AddProgram)
	app.Post("/programs/add", controllers.StoreProgram)
}

func Semester(app *fiber.App) {
	app.Get("/semesters/add", controllers.AddSemester)
	app.Post("/semesters/add", controllers.StoreSemester)
}

func Subject(app *fiber.App) {
	app.Get("/courses/add", controllers.AddCourse)
	app.Post("/courses/add", controllers.StoreCourse)
}

func Mark(app *fiber.App) {
	// app.Get("/marks/add", controllers.AddMark)
	app.Post("/marks/add", controllers.CreateMarks)
	app.Put("/marks/edit/:id", controllers.UpdateMarks)
	app.Get("/marks/:symbolNumber", controllers.GetMarksBySymbolNumber)
}

package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mysterybee07/result-distribution-system/controllers"
	"github.com/mysterybee07/result-distribution-system/middleware"
)

func Home(app *fiber.App) {
	app.Get("", controllers.Home)
	app.Get("/user/register", controllers.Register)
	// app.Get("/register-admin", middleware.AuthRequired, middleware.AdminRequired, controllers.RegisterAdmin)
	app.Post("/user/register", controllers.StoreRegister)
	app.Get("/user/login", controllers.Login)
	app.Post("/user/login", controllers.LoginUser)
	app.Post("/user/logout", controllers.LogoutUser)
	app.Get("/user/forgot-password", controllers.ForgotPassword)
	app.Put("/user/update/:id", controllers.UpdateUser)
	app.Get("/users", controllers.GetLoginUser)
	// app.Get("/logout", controllers.LogoutUser)
}

func Profile(app *fiber.App) {
	app.Get("/profile", middleware.AuthRequired, controllers.GetUserProfile)
}

// func Dashboard(app *fiber.App) {
// 	app.Get("/dashboard", middleware.AuthRequired, middleware.AdminRequired, controllers.Index)
// }

func Student(app *fiber.App) {
	app.Get("/students/add", middleware.AuthRequired, middleware.AdminRequired, controllers.AddStudent)
	app.Post("/students/add", middleware.AuthRequired, middleware.AdminRequired, controllers.StoreStudents)
	app.Get("/students", middleware.AuthRequired, middleware.AdminRequired, controllers.GetStudents)
	// app.Get("/students/edit/:id", middleware.AuthRequired, middleware.AdminRequired, controllers.EditStudentForm)
	app.Put("/students/edit/:id", middleware.AuthRequired, middleware.AdminRequired, controllers.UpdateStudent)
	app.Get("/students/:id", controllers.GetStudentById)
	app.Get("/students/edit/:id", middleware.AuthRequired, middleware.AdminRequired, controllers.EditStudent)
}

func Batch(app *fiber.App) {
	app.Get("/batches", middleware.AuthRequired, middleware.SuperadminRequired, controllers.AddBatch)
	app.Post("/batches", middleware.AuthRequired, middleware.SuperadminRequired, controllers.CreateBatch)
}

func Program(app *fiber.App) {
	app.Get("/programs", middleware.AuthRequired, middleware.SuperadminRequired, controllers.AddProgram)
	app.Post("/programs", middleware.AuthRequired, middleware.SuperadminRequired, controllers.StoreProgram)
}

func Semester(app *fiber.App) {
	app.Get("/semesters", middleware.AuthRequired, middleware.SuperadminRequired, controllers.AddSemester)
	app.Post("/semesters", middleware.AuthRequired, middleware.SuperadminRequired, controllers.StoreSemester)
}

func Subject(app *fiber.App) {
	app.Get("/courses", middleware.AuthRequired, middleware.SuperadminRequired, controllers.AddCourse)
	app.Post("/courses", middleware.AuthRequired, middleware.SuperadminRequired, controllers.StoreCourse)
}

func Mark(app *fiber.App) {
	app.Post("/marks/add", middleware.AuthRequired, middleware.AdminRequired, controllers.CreateMarks)
	app.Put("/marks/edit/:id", middleware.AuthRequired, middleware.AdminRequired, controllers.UpdateMarks)
	app.Get("/marks/:symbolNumber", middleware.AuthRequired, middleware.AdminRequired, controllers.GetMarksBySymbolNumber)
	// app.Post("/publish-results", middleware.AuthRequired, middleware.AdminRequired, controllers.PublishResults)
	app.Get("/getstudents", controllers.GetFilteredStudents)
	app.Get("/getfiltercourses", controllers.GetFilteredCourses)
	app.Get("/getfiltersemesters", controllers.GetFilteredSemesters)
}

func Result(app *fiber.App) {
	app.Get("/results", middleware.AuthRequired, middleware.SuperadminRequired, controllers.AddResult)
	app.Post("/results", middleware.AuthRequired, middleware.SuperadminRequired, controllers.PublishResults)
}

func Error(app *fiber.App) {
	app.Get("/404", controllers.NotFound)
	app.Get("/500", controllers.ServerError)
}

// func Subject(app *fiber.App) {
// 	app.Get("/courses", controllers.AddCourse)
// 	app.Post("/courses", controllers.StoreCourse)
// 	app.Get("/api/semesters/:programID", controllers.GetSemestersByProgramID)
// }

// func Mark(app *fiber.App) {
// 	app.Get("/marks/add", controllers.AddMarks)
// 	app.Post("/marks/add", controllers.CreateMarks)
// 	app.Put("/marks/edit/:id", controllers.UpdateMarks)
// 	app.Get("/marks/:symbolNumber", controllers.GetMarksBySymbolNumber)
// 	app.Get("/getstudents", controllers.GetFilteredStudents)
// 	app.Get("/getfiltercourses", controllers.GetFilteredCourses)
// 	app.Get("/getfiltersemesters", controllers.GetFilteredSemesters)

// }
func Dashboard(app *fiber.App) {
	app.Get("/dashboard", controllers.Index)
	app.Get("/failstudents", controllers.FailStudents)
}

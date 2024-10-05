package routes

import (
	"github.com/gofiber/fiber/v2"
	controllers "github.com/mysterybee07/result-distribution-system/controllers/admin"
	errorController "github.com/mysterybee07/result-distribution-system/controllers/error"
	homeController "github.com/mysterybee07/result-distribution-system/controllers/home"
	userController "github.com/mysterybee07/result-distribution-system/controllers/user"
	"github.com/mysterybee07/result-distribution-system/middleware"
)

func Home(app *fiber.App) {
	app.Get("", homeController.Home)
	app.Get("/user/register", homeController.Register)
	// app.Get("/register-admin", middleware.AuthRequired, middleware.AdminRequired, controllers.RegisterAdmin)
	app.Post("/user/register", homeController.StoreRegister)
	app.Get("/login", homeController.Login)
	app.Post("/user/login", homeController.LoginUser)
	app.Post("/user/logout", homeController.LogoutUser)
	app.Get("/user/forgot-password", homeController.ForgotPassword)
	app.Put("/user/update/:id", homeController.UpdateUser)
	app.Get("/users", homeController.GetLoginUser)
	app.Get("/all-users", homeController.GetAllUsers)
	app.Get("/user/:id", homeController.GetUserById)
	// app.Get("/logout", controllers.LogoutUser)
}

func Profile(app *fiber.App) {
	app.Get("/profile", middleware.AuthRequired, userController.GetUserProfile)
}

// func Dashboard(app *fiber.App) {
// 	app.Get("/dashboard", middleware.AuthRequired, middleware.AdminRequired, controllers.Index)
// }

func Student(app *fiber.App) {
	app.Get("/students/add", middleware.AuthRequired, middleware.AdminRequired, controllers.Student)
	// app.Post("/students/add", middleware.AuthRequired, middleware.AdminRequired, controllers.StoreStudents)
	app.Get("/students", middleware.AuthRequired, middleware.AdminRequired, controllers.GetStudents)
	// app.Get("/students/edit/:id", middleware.AuthRequired, middleware.AdminRequired, controllers.EditStudentForm)
	// app.Put("/students/update/:id", middleware.AuthRequired, middleware.AdminRequired, controllers.UpdateStudent)
	app.Put("/student/update/:id", controllers.UpdateStudent)
	app.Get("/students/:id", controllers.GetStudentById)
	app.Get("/students/edit/:id", middleware.AuthRequired, middleware.AdminRequired, controllers.EditStudent)
	app.Post("/students/create", controllers.CreateStudents)
	app.Get("/getstudents", controllers.GetFilteredStudents)
}

func Batch(app *fiber.App) {
	app.Get("/batches", middleware.AuthRequired, middleware.SuperadminRequired, controllers.Batch)
	// app.Post("/batches", middleware.AuthRequired, middleware.SuperadminRequired, controllers.CreateBatch)
	app.Post("/batches/create", controllers.CreateBatch)
}

func Program(app *fiber.App) {
	app.Get("/programs", middleware.AuthRequired, middleware.SuperadminRequired, controllers.Program)
	// app.Post("/programs/add", middleware.AuthRequired, middleware.SuperadminRequired, controllers.CreateProgram)
	app.Post("/programs/create", controllers.CreateProgram)
	app.Put("/programs/update/:id", controllers.UpdateProgram)
}

func Semester(app *fiber.App) {
	app.Get("/semesters", middleware.AuthRequired, middleware.SuperadminRequired, controllers.Semester)
	// app.Post("/semesters", middleware.AuthRequired, middleware.SuperadminRequired, controllers.CreateSemester)
	app.Post("/semester/create", controllers.CreateSemester)
	app.Put("/semester/update/:id", controllers.UpdateSemester)
	app.Get("/semester/sem-by-program", controllers.GetSemestersByProgramID)
}

func Course(app *fiber.App) {
	app.Get("/courses", middleware.AuthRequired, middleware.SuperadminRequired, controllers.Course)
	// app.Post("/courses/create", middleware.AuthRequired, middleware.SuperadminRequired, controllers.CreateCourse)
	app.Post("/courses/create", controllers.CreateCourses)
	app.Put("/course/update/:id", controllers.UpdateCourse)
	app.Get("/getfiltercourses", controllers.GetFilteredCourses)

}

func Mark(app *fiber.App) {
	app.Get("/marks", controllers.Marks)
	// app.Post("/marks/create", middleware.AuthRequired, middleware.AdminRequired, controllers.CreateMarks)
	// app.Put("/marks/update/:id", middleware.AuthRequired, middleware.AdminRequired, controllers.UpdateMarks)
	app.Post("/marks/create", controllers.CreateMarks)
	app.Put("/marks/update/:id", controllers.UpdateMarks)
	// app.Get("/marks/:symbolNumber", middleware.AuthRequired, middleware.AdminRequired, controllers.GetMarksBySymbolNumber)
	app.Get("/marks/:symbolNumber", controllers.GetMarksBySymbolNumber)
	// app.Post("/publish-results", middleware.AuthRequired, middleware.AdminRequired, controllers.PublishResults)

}

func Result(app *fiber.App) {
	app.Get("/result", middleware.AuthRequired, middleware.SuperadminRequired, controllers.Result)
	app.Post("/result/publish", middleware.AuthRequired, middleware.SuperadminRequired, controllers.PublishResults)
}

func Error(app *fiber.App) {
	app.Get("/404", errorController.NotFound)
	app.Get("/500", errorController.ServerError)
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

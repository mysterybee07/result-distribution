// package routes

// import (
// 	"github.com/gofiber/fiber/v2"
// 	"github.com/mysterybee07/result-distribution-system/controllers"
// 	"github.com/mysterybee07/result-distribution-system/middleware"
// )

// func Home(app *fiber.App) {
// 	app.Get("", controllers.Home)
// 	app.Get("/register", controllers.Register)
// 	app.Post("/register", controllers.StoreRegister)
// 	app.Get("/login", controllers.Login)
// 	app.Post("/login", controllers.LoginUser)
// 	app.Get("/logout", controllers.LogoutUser)
// 	app.Get("/forgot-password", controllers.ForgotPassword)
// }

// func Profile(app *fiber.App) {
// 	app.Get("/profile", middleware.AuthRequired, controllers.GetUserProfile)
// }

// func Dashboard(app *fiber.App) {
// 	app.Get("/dashboard", middleware.AuthRequired, middleware.AdminRequired, controllers.Index)
// }

// func Student(app *fiber.App) {
// 	app.Get("/students/add", middleware.AuthRequired, middleware.AdminRequired, controllers.AddStudent)
// 	app.Post("/students/add", middleware.AuthRequired, middleware.AdminRequired, controllers.StoreStudent)
// 	app.Get("/students", middleware.AuthRequired, middleware.AdminRequired, controllers.GetStudents)
// 	// app.Get("/students/edit/:id", middleware.AuthRequired, middleware.AdminRequired, controllers.EditStudentForm)
// 	app.Put("/students/edit/:id", middleware.AuthRequired, middleware.AdminRequired, controllers.UpdateStudent)
// }

// func Batch(app *fiber.App) {
// 	app.Get("/batches", middleware.AuthRequired, middleware.SuperadminRequired, controllers.AddBatch)
// 	app.Post("/batches", middleware.AuthRequired, middleware.SuperadminRequired, controllers.CreateBatch)
// }

// func Program(app *fiber.App) {
// 	app.Get("/programs", middleware.AuthRequired, middleware.SuperadminRequired, controllers.AddProgram)
// 	app.Post("/programs", middleware.AuthRequired, middleware.SuperadminRequired, controllers.StoreProgram)
// }

// func Semester(app *fiber.App) {
// 	app.Get("/semesters", middleware.AuthRequired, middleware.SuperadminRequired, controllers.AddSemester)
// 	app.Post("/semesters", middleware.AuthRequired, middleware.SuperadminRequired, controllers.StoreSemester)
// }

// func Subject(app *fiber.App) {
// 	app.Get("/courses", middleware.AuthRequired, middleware.SuperadminRequired, controllers.AddCourse)
// 	app.Get("/programs/:id/semesters", controllers.GetSemestersByProgram)
// 	app.Post("/courses", middleware.AuthRequired, middleware.SuperadminRequired, controllers.StoreCourse)
// }

// func Mark(app *fiber.App) {
// 	app.Post("/marks/add", middleware.AuthRequired, middleware.AdminRequired, controllers.CreateMarks)
// 	app.Put("/marks/edit/:id", middleware.AuthRequired, middleware.AdminRequired, controllers.UpdateMarks)
// 	app.Get("/marks/:symbolNumber", middleware.AuthRequired, middleware.AdminRequired, controllers.GetMarksBySymbolNumber)
// 	app.Post("/publish-results", middleware.AuthRequired, middleware.AdminRequired, controllers.PublishResults)
// }

// func Result(app *fiber.App) {
// 	app.Get("/results", middleware.AuthRequired, middleware.SuperadminRequired, controllers.AddResult)
// 	app.Post("/results", middleware.AuthRequired, middleware.SuperadminRequired, controllers.PublishResults)
// }

// func Error(app *fiber.App) {
// 	app.Get("/404", controllers.NotFound)
// 	app.Get("/500", controllers.ServerError)
// }

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
	app.Put("/edit-user/:id", controllers.UpdateUser)
	// app.Post("/upload", controllers.Upload)
	// app.Static("/uploads", "./static/images/uploads")
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
	app.Post("/students/add", controllers.StoreStudents)
	// app.Post("/students/add", controllers.PostStudents)
	app.Get("/students", controllers.GetStudents)
	app.Get("/students/edit/:id", controllers.EditStudent)
	app.Put("/students/edit/:id", controllers.UpdateStudent)
}

func Batch(app *fiber.App) {
	app.Get("/batches", controllers.AddBatch)
	app.Post("/batches", controllers.CreateBatch)
}

func Program(app *fiber.App) {
	app.Get("/programs", controllers.AddProgram)
	app.Post("/programs", controllers.StoreProgram)
}

func Semester(app *fiber.App) {
	app.Get("/semesters", controllers.AddSemester)
	app.Post("/semesters", controllers.StoreSemester)
}

func Subject(app *fiber.App) {
	app.Get("/courses", controllers.AddCourse)
	app.Post("/courses", controllers.StoreCourse)
}

func Mark(app *fiber.App) {
	// app.Get("/marks/add", controllers.AddMark)
	app.Post("/marks/add", controllers.CreateMarks)
	app.Put("/marks/edit/:id", controllers.UpdateMarks)
	app.Get("/marks/:symbolNumber", controllers.GetMarksBySymbolNumber)
}

func Result(app *fiber.App) {
	app.Get("/results", middleware.AuthRequired, middleware.SuperadminRequired, controllers.AddResult)
	app.Post("/results", controllers.PublishResults)
}

func Error(app *fiber.App) {
	app.Get("/404", controllers.NotFound)
	app.Get("/500", controllers.ServerError)
}

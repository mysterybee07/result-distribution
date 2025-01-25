// package routes

// import (
// 	"github.com/gofiber/fiber/v2"
// 	controllers "github.com/mysterybee07/result-distribution-system/controllers/admin"
// 	errorController "github.com/mysterybee07/result-distribution-system/controllers/error"
// 	homeController "github.com/mysterybee07/result-distribution-system/controllers/home"
// 	userController "github.com/mysterybee07/result-distribution-system/controllers/user"
// 	"github.com/mysterybee07/result-distribution-system/middleware"
// )

// func Home(app *fiber.App) {
// 	app.Get("", homeController.Home)
// 	app.Get("/user/register", homeController.Register)
// 	// app.Get("/register-admin", middleware.AuthRequired, middleware.AdminRequired, controllers.RegisterAdmin)
// 	app.Post("/user/register", homeController.StoreRegister)
// 	app.Get("/login", homeController.Login)
// 	app.Post("/user/login", homeController.LoginUser)
// 	app.Post("/user/logout", homeController.LogoutUser)
// 	app.Get("/user/forgot-password", homeController.ForgotPassword)
// 	app.Put("/user/update/:id", homeController.UpdateUser)
// 	app.Get("/users", homeController.GetLoginUser)
// 	app.Get("/all-users", homeController.GetAllUsers)
// 	app.Get("/user/:id", homeController.GetUserById)
// 	// app.Get("/logout", controllers.LogoutUser)
// }

// func Profile(app *fiber.App) {
// 	app.Get("/profile", middleware.AuthRequired, userController.GetUserProfile)
// }

// // func Dashboard(app *fiber.App) {
// // 	app.Get("/dashboard", middleware.AuthRequired, middleware.AdminRequired, controllers.Index)
// // }

// func Student(app *fiber.App) {
// 	app.Get("/students/add", middleware.AuthRequired, middleware.AdminRequired, controllers.Student)
// 	// app.Post("/students/add", middleware.AuthRequired, middleware.AdminRequired, controllers.StoreStudents)
// 	app.Get("/students", middleware.AuthRequired, middleware.AdminRequired, controllers.GetStudents)
// 	// app.Get("/students/edit/:id", middleware.AuthRequired, middleware.AdminRequired, controllers.EditStudentForm)
// 	// app.Put("/students/update/:id", middleware.AuthRequired, middleware.AdminRequired, controllers.UpdateStudent)
// 	app.Put("/student/update/:id", controllers.UpdateStudent)
// 	app.Get("/students/:id", controllers.GetStudentById)
// 	app.Get("/students/edit/:id", middleware.AuthRequired, middleware.AdminRequired, controllers.EditStudent)
// 	app.Post("/students/create", controllers.CreateStudents)
// 	app.Get("/getstudents", controllers.GetFilteredStudents)
// }

// func Batch(app *fiber.App) {
// 	app.Get("/batches", middleware.AuthRequired, middleware.SuperadminRequired, controllers.Batch)
// 	// app.Post("/batches", middleware.AuthRequired, middleware.SuperadminRequired, controllers.CreateBatch)
// 	app.Post("/batches/create", controllers.CreateBatch)
// }

// func Program(app *fiber.App) {
// 	app.Get("/programs", middleware.AuthRequired, middleware.SuperadminRequired, controllers.Program)
// 	// app.Post("/programs/add", middleware.AuthRequired, middleware.SuperadminRequired, controllers.CreateProgram)
// 	app.Post("/programs/create", controllers.CreateProgram)
// 	app.Put("/programs/update/:id", controllers.UpdateProgram)
// }

// func Semester(app *fiber.App) {
// 	app.Get("/semesters", middleware.AuthRequired, middleware.SuperadminRequired, controllers.Semester)
// 	// app.Post("/semesters", middleware.AuthRequired, middleware.SuperadminRequired, controllers.CreateSemester)
// 	app.Post("/semester/create", controllers.CreateSemester)
// 	app.Put("/semester/update/:id", controllers.UpdateSemester)
// 	app.Get("/semester/sem-by-program", controllers.GetSemestersByProgramID)
// }

// func Course(app *fiber.App) {
// 	app.Get("/courses", middleware.AuthRequired, middleware.SuperadminRequired, controllers.Course)
// 	// app.Post("/courses/create", middleware.AuthRequired, middleware.SuperadminRequired, controllers.CreateCourse)
// 	app.Post("/courses/create", controllers.CreateCourses)
// 	app.Put("/course/update/:id", controllers.UpdateCourse)
// 	app.Get("/getfiltercourses", controllers.GetFilteredCourses)

// }

// func Mark(app *fiber.App) {
// 	app.Get("/marks", controllers.Marks)
// 	// app.Post("/marks/create", middleware.AuthRequired, middleware.AdminRequired, controllers.CreateMarks)
// 	// app.Put("/marks/update/:id", middleware.AuthRequired, middleware.AdminRequired, controllers.UpdateMarks)
// 	app.Post("/marks/create", controllers.CreateMarks)
// 	app.Put("/marks/update/:id", controllers.UpdateMarks)
// 	// app.Get("/marks/:symbolNumber", middleware.AuthRequired, middleware.AdminRequired, controllers.GetMarksBySymbolNumber)
// 	app.Get("/marks/:symbolNumber", controllers.GetMarksBySymbolNumber)
// 	// app.Post("/publish-results", middleware.AuthRequired, middleware.AdminRequired, controllers.PublishResults)

// }

// func Result(app *fiber.App) {
// 	app.Get("/result", middleware.AuthRequired, middleware.SuperadminRequired, controllers.Result)
// 	app.Post("/result/publish", middleware.AuthRequired, middleware.SuperadminRequired, controllers.PublishResults)
// }

// func Error(app *fiber.App) {
// 	app.Get("/404", errorController.NotFound)
// 	app.Get("/500", errorController.ServerError)
// }

// // func Subject(app *fiber.App) {
// // 	app.Get("/courses", controllers.AddCourse)
// // 	app.Post("/courses", controllers.StoreCourse)
// // 	app.Get("/api/semesters/:programID", controllers.GetSemestersByProgramID)
// // }

// // func Mark(app *fiber.App) {
// // 	app.Get("/marks/add", controllers.AddMarks)
// // 	app.Post("/marks/add", controllers.CreateMarks)
// // 	app.Put("/marks/edit/:id", controllers.UpdateMarks)
// // 	app.Get("/marks/:symbolNumber", controllers.GetMarksBySymbolNumber)
// // 	app.Get("/getstudents", controllers.GetFilteredStudents)
// // 	app.Get("/getfiltercourses", controllers.GetFilteredCourses)
// // 	app.Get("/getfiltersemesters", controllers.GetFilteredSemesters)

// // }
// func Dashboard(app *fiber.App) {
// 	app.Get("/dashboard", controllers.Index)
// 	app.Get("/failstudents", controllers.FailStudents)
// }

package routes

import (
	"github.com/gofiber/fiber/v2"
	adminController "github.com/mysterybee07/result-distribution-system/controllers/admin"
	authController "github.com/mysterybee07/result-distribution-system/controllers/auth"
	errorController "github.com/mysterybee07/result-distribution-system/controllers/error"
	examController "github.com/mysterybee07/result-distribution-system/controllers/exam"
	noticeController "github.com/mysterybee07/result-distribution-system/controllers/notice"
	userController "github.com/mysterybee07/result-distribution-system/controllers/user"
	"github.com/mysterybee07/result-distribution-system/middleware"
)

// func Home(app *fiber.App) {
// 	app.Get("", authController.Home)
// }

func SetupRoutes(app *fiber.App) {
	// Home/User Routes
	user := app.Group("/user")

	// user.Get("/register", authController.Register)
	user.Post("/register", authController.StoreRegister)
	// user.Get("/login", authController.Login)
	user.Post("/login", authController.LoginUser)
	user.Post("/logout", authController.LogoutUser)
	user.Get("/forgot-password", authController.ForgotPassword)
	user.Put("/update/:id", authController.UpdateUser)
	user.Get("/active", middleware.AuthRequired, authController.AuthorizedUser)
	user.Get("", authController.GetAllUsers)
	user.Get("/:id", authController.GetUserById)
	// user.Get("/logout", controllers.LogoutUser)

	// Profile Routes
	profile := app.Group("/profile")
	profile.Get("/", middleware.AuthRequired, userController.GetUserProfile)

	// Student Routes
	student := app.Group("/students")
	// student := app.Group("/students", middleware.AuthRequired, middleware.AdminRequired)
	// student.Get("/add", adminController.Student)
	// student.Post("/add", adminController.StoreStudents)
	student.Get("", middleware.AuthRequired, adminController.GetStudents)
	student.Put("/update/:id", adminController.UpdateStudent)
	student.Get("/:id", adminController.GetStudentById)
	student.Get("/edit/:id", adminController.EditStudent)
	student.Post("/create", adminController.CreateStudents)
	student.Get("/filter", adminController.GetFilteredStudents)
	student.Delete("/delete", adminController.DeleteStudent)
	student.Get("/pass-students-by-semester", adminController.PassingStudentsBySemester)
	student.Get("/fail-students-by-course", adminController.FailedStudentsByCourse)

	// Batch Routes
	batch := app.Group("/batch")
	// batch := app.Group("/batches", middleware.AuthRequired, middleware.SuperadminRequired)
	batch.Get("", adminController.GetBatches)
	batch.Post("/create", adminController.CreateBatch)
	batch.Put("/update/:id", adminController.UpdateBatch)
	// batch.Get("/")

	// Program Routes
	program := app.Group("/program")
	// program := app.Group("/programs", middleware.AuthRequired, middleware.SuperadminRequired)
	program.Get("", middleware.AuthRequired, adminController.GetPrograms)
	// program.Post("/add", adminController.CreateProgram)
	program.Post("/create", middleware.AuthRequired, adminController.CreateProgram)
	program.Put("/update/:id", adminController.UpdateProgram)

	// Semester Routes
	semester := app.Group("/semester")
	// semester := app.Group("/semesters", middleware.AuthRequired, middleware.SuperadminRequired)
	semester.Get("/", adminController.Semester)
	// semester.Post("/", adminController.CreateSemester)
	semester.Post("/create", adminController.CreateSemester)
	semester.Put("/update/:id", adminController.UpdateSemester)
	semester.Get("/by-program/:id", adminController.GetSemestersByProgramID)

	// Course Routes
	course := app.Group("/courses")
	// course := app.Group("/courses", middleware.AuthRequired, middleware.SuperadminRequired)
	// course.Get("/", adminController.Course)
	// course.Post("/create", adminController.CreateCourses)
	course.Post("/create", adminController.CreateCourses)
	course.Put("/update/:id", adminController.UpdateCourse)
	course.Get("/filter", adminController.GetFilteredCourses)
	course.Get("/:id", adminController.GetCourseById)

	// Mark Routes
	mark := app.Group("/marks")
	mark.Get("/", adminController.Marks)
	// mark.Post("/create", middleware.AuthRequired, middleware.AdminRequired, adminController.CreateMarks)
	mark.Post("/create", adminController.CreateMarks)
	// mark.Put("/update/:id", middleware.AuthRequired, middleware.AdminRequired, adminController.UpdateMarks)
	mark.Put("/update/:id", adminController.UpdateMarks)
	// mark.Get("/:symbolNumber", middleware.AuthRequired, middleware.AdminRequired, adminController.GetMarksBySymbolNumber)
	mark.Get("/:symbolNumber", adminController.GetMarksBySymbolNumber)
	// app.Post("/publish-results", middleware.AuthRequired, middleware.AdminRequired, adminController.PublishResults)

	// Result Routes
	// result := app.Group("/result", middleware.AuthRequired, middleware.SuperadminRequired)
	result := app.Group("/result")
	result.Get("", adminController.Result)
	result.Post("/publish", adminController.PublishResults)

	// Error Routes
	errorGroup := app.Group("/error")
	errorGroup.Get("/404", errorController.NotFound)
	errorGroup.Get("/500", errorController.ServerError)

	// Dashboard Routes
	// dashboard := app.Group("/dashboard")
	// dashboard.Get("/", adminController.Index)
	// dashboard.Get("/failstudents", adminController.FailStudents)

	//notice routes
	notice := app.Group("/notice")
	notice.Get("", noticeController.GetAllNotices)
	notice.Post("/create", noticeController.CreateNotice)
	notice.Delete("/delete/:id", noticeController.DeleteNotice)
	notice.Put("/update/:id", noticeController.UpdateNotice)
	notice.Get("/by-id/:id", noticeController.GetNoticeById)
	notice.Get("/by-program", noticeController.GetNoticesByProgram)
	notice.Get("/by-program-and-batch", noticeController.GetNoticesByProgramAndBatch)
	notice.Post("/publish", noticeController.PublishNotice)

	exam := app.Group("/exam")
	exam.Get("/assign-centers", examController.AssignCentersHandler)
	exam.Post("/update-center-and-capacity", adminController.AssignCenterAndCapacity)
	exam.Post("/update-capacity", adminController.UpdateCapacity)
	exam.Post("/schedule/create", adminController.CreateExamRoutine)
	exam.Post("/schedule/publish", adminController.PublishExamRoutine)

	college := app.Group("/college")
	college.Get("", adminController.GetColleges)
	college.Post("/upload-college", adminController.UploadColleges)
	college.Get("/centers-by-program-and-batch", adminController.GetCenterCollegesByProgramAndBatch)
	// college.Get("/all-centers", adminController.GetAllCenterColleges)
	college.Put("/update-college/:id", adminController.UpdateCollege)
	college.Delete("/delete-college/:id", adminController.DeleteCollege)
}

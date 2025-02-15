package routes

import (
	"github.com/gofiber/fiber/v2"
	adminController "github.com/mysterybee07/result-distribution-system/controllers/admin"

	// controllers "github.com/mysterybee07/result-distribution-system/controllers/admin"
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
	profile.Get("", middleware.AuthRequired, userController.GetUserProfile)

	// Student Routes
	student := app.Group("/students")
	// student := app.Group("/students", middleware.AuthRequired, middleware.AdminRequired)
	// student.Get("/add", adminController.Student)
	// student.Post("/add", adminController.StoreStudents)
	student.Get("", adminController.GetStudents)
	student.Put("/update/:id", adminController.UpdateStudent)
	student.Get("/filtered", adminController.FilteredStudents)
	student.Get("/:id", adminController.GetStudentById)
	student.Get("/edit/:id", adminController.EditStudent)
	student.Post("/create", adminController.CreateStudents)
	student.Get("/by-program/:id", adminController.FilterStudentsByProgram)
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
	course.Get("", adminController.GetAllCourses)
	// course.Post("/create", adminController.CreateCourses)
	// TODO: get all courses
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
	exam.Get("/routines", adminController.ListExamsRoutine)
	exam.Get("/schedules", adminController.ListExamSchedules)
	exam.Get("/schedules/by-batch-program", adminController.GetFilteredExamSchedules)
	exam.Get("/assign-centers", examController.AssignCentersHandler)
	exam.Post("/update-center-and-capacity", adminController.AssignCenterAndCapacity)
	exam.Put("/update-capacity/:id", adminController.UpdateCapacity)
	exam.Post("/schedule/create", adminController.CreateExamRoutine)
	exam.Post("/schedule/publish/:id", adminController.PublishExamRoutine)

	college := app.Group("/college")
	college.Get("", adminController.GetColleges)
	college.Post("/upload-college", adminController.UploadColleges)
	college.Get("/centers-by-program-and-batch", adminController.GetCenterCollegesByProgramAndBatch)
	// college.Get("/all-centers", adminController.GetAllCenterColleges)
	college.Put("/update-college/:id", adminController.UpdateCollege)
	college.Delete("/delete-college/:id", adminController.DeleteCollege)
}

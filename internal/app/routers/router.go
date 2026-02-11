package routers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"

	"course-reg/internal/app/domain/export"
	"course-reg/internal/app/handler"
	"course-reg/internal/app/middleware"
)

// InitRouter initialize routing information
func InitRouter(
	runMode, sessionKey string,
	h *handler.Handlers,
) *gin.Engine {
	gin.SetMode(runMode) // set gin mode (must be called before gin.New())
	r := gin.New()
	if gin.Mode() != gin.ReleaseMode {
		r.Use(gin.Logger())
	}
	r.Use(gin.Recovery())    // panic 발생시 500
	r.Use(middleware.CORS()) // CORS

	// session
	store := memstore.NewStore([]byte(sessionKey)) // authentication key for session
	r.Use(sessions.Sessions("course_reg_session", store))

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", h.Auth.Login)
			auth.POST("/logout", h.Auth.Logout)
			auth.GET("/check", h.Auth.Check)
		}

		admin := v1.Group("/admin")
		admin.Use(middleware.AuthAdmin())
		{
			admin.GET("/registration/state", h.Admin.GetRegistrationState)
			admin.POST("/registration/start", h.Admin.StartRegistration)
			admin.POST("/registration/pause", h.Admin.PauseRegistration)
			admin.PUT("/registration/period", h.Admin.SetRegistrationPeriod)
			admin.GET("/registration/period", h.Admin.GetRegistrationPeriod)

			setup := admin.Group("/setup")
			{
				// todo : reset과 init atomic하게 묶기?
				// admin.POST("/students/init", adminHandler.RegisterStudents)

				setup.POST("/students/register", h.Admin.RegisterStudents)
				setup.DELETE("/students/reset", h.Admin.ResetStudents)

				setup.POST("/courses", h.Admin.CreateCourse)
				setup.DELETE("/courses/:course_id", h.Admin.DeleteCourse)
				setup.POST("/courses/register", h.Admin.RegisterCourses)
				setup.DELETE("/courses/reset", h.Admin.ResetCourses)

				setup.DELETE("/enrollments/reset", h.Admin.ResetEnrollments)
			}

			// 수강 신청 기간 중에는 어떻게?
			// admin.POST("/enrollments", adminHandler.AddEnrollment)
			// admin.DELETE("/enrollments", adminHandler.CancelEnrollment)
		}

		user := v1.Group("/courses")
		user.Use(middleware.AuthUser())
		{
			user.StaticFile("/", export.StaticCoursesFilePath)
			user.GET("/status", h.CourseReg.GetAllCourseStatus)
			// user.GET("/enrollments", courseRegHandler.GetCoursesCapacityStatus)

		}

		courseReg := v1.Group("/course-reg")
		courseReg.Use(middleware.AuthStudent())
		{
			courseReg.POST("/enrollment", h.CourseReg.EnrollCourse)
			// courseReg.DELETE("/:course_id/enroll", courseRegHandler.CancelEnrollment)

			// courseReg.POST("/:course_id/waitlist", courseRegHandler.AddToWaitlist)
			// courseReg.DELETE("/:course_id/waitlist", courseRegHandler.DeleteToWaitlist)
		}
	}
	return r
}

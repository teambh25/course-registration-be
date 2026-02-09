package routers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"

	"course-reg/internal/app/handler"
	"course-reg/internal/app/middleware"
	"course-reg/internal/pkg/constant"
	"course-reg/internal/pkg/setting"
)

// InitRouter initialize routing information
func InitRouter(
	adminHandler *handler.AdminHandler,
	authHandler *handler.AuthHandler,
	courseRegHandler *handler.CourseRegHandler,
) *gin.Engine {
	gin.SetMode(setting.ServerSetting.RunMode) // set gin mode (must be called before gin.New())
	r := gin.New()
	if gin.Mode() != gin.ReleaseMode {
		r.Use(gin.Logger())
	}
	r.Use(gin.Recovery())    // panic 발생시 500
	r.Use(middleware.CORS()) // CORS

	// session
	store := memstore.NewStore([]byte(setting.SecretSetting.SessionKey)) // authentication key for session
	r.Use(sessions.Sessions("course_reg_session", store))

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", authHandler.Logout)
			auth.GET("/check", authHandler.Check)
		}

		admin := v1.Group("/admin")
		// admin.Use(middleware.AuthAdmin())
		{
			admin.GET("/registration/state", adminHandler.GetRegistrationState)
			admin.POST("/registration/start", adminHandler.StartRegistration)
			admin.POST("/registration/pause", adminHandler.PauseRegistration)
			admin.PUT("/registration/period", adminHandler.SetRegistrationPeriod)
			admin.GET("/registration/period", adminHandler.GetRegistrationPeriod)

			setup := admin.Group("/setup")
			{
				// todo : reset과 init atomic하게 묶기?
				// admin.POST("/students/init", adminHandler.RegisterStudents)

				setup.POST("/students/register", adminHandler.RegisterStudents)
				setup.DELETE("/students/reset", adminHandler.ResetStudents)

				setup.POST("/courses", adminHandler.CreateCourse)
				setup.DELETE("/courses/:course_id", adminHandler.DeleteCourse)
				setup.POST("/courses/register", adminHandler.RegisterCourses)
				setup.DELETE("/courses/reset", adminHandler.ResetCourses)

				setup.DELETE("/enrollments/reset", adminHandler.ResetEnrollments)
			}

			// 수강 신청 기간 중에는 어떻게?
			// admin.POST("/enrollments", adminHandler.AddEnrollment)
			// admin.DELETE("/enrollments", adminHandler.CancelEnrollment)
		}

		user := v1.Group("/courses")
		user.Use(middleware.AuthUser())
		{
			user.StaticFile("/", constant.StaticCoursesFilePath)
			user.GET("/status", courseRegHandler.GetAllCourseStatus)
			// user.GET("/enrollments", courseRegHandler.GetCoursesCapacityStatus)

		}

		courseReg := v1.Group("/course-reg")
		courseReg.Use(middleware.AuthStudent())
		{
			courseReg.POST("/enrollment", courseRegHandler.EnrollCourse)
			// courseReg.DELETE("/:course_id/enroll", courseRegHandler.CancelEnrollment)

			// courseReg.POST("/:course_id/waitlist", courseRegHandler.AddToWaitlist)
			// courseReg.DELETE("/:course_id/waitlist", courseRegHandler.DeleteToWaitlist)
		}
	}
	return r
}

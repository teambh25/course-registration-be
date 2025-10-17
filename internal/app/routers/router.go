package routers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"

	"course-reg/internal/app/handler"
	"course-reg/internal/app/middleware"
	"course-reg/internal/pkg/setting"
	"course-reg/internal/pkg/util"
)

// InitRouter initialize routing information
func InitRouter(adminHandler *handler.AdminHandler, authHandler *handler.AuthHandler, courseRegHandler *handler.CourseRegHandler, timeProvider util.TimeProvider) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery()) // panic 발생시 500

	// CORS
	r.Use(middleware.CORS())

	// session
	store := memstore.NewStore([]byte(setting.SecretSetting.SessionKey)) // authentication key for session
	r.Use(sessions.Sessions("course_reg_session", store))

	// r.StaticFS("/export", http.Dir(export.GetExcelFullPath()))

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

			admin.POST("students/register", adminHandler.RegisterStudents)
			admin.DELETE("students/reset", adminHandler.ResetStudents)

			admin.POST("/courses", adminHandler.CreateCourse)
			admin.DELETE("/courses/:course_id", adminHandler.DeleteCourse)
			admin.POST("/courses/register", adminHandler.RegisterCourses)
			admin.DELETE("courses/reset", adminHandler.ResetCourses)

			admin.POST("/registration-period", adminHandler.SetRegistrationPeriod)
			admin.GET("/registration-period", adminHandler.GetRegistrationPeriod)

			// admin.POST("/enrollments", adminHandler.AddEnrollment)
			// admin.DELETE("/enrollments", adminHandler.CancelEnrollment)
		}

		user := v1.Group("/user")
		// user.Use(middleware.Auth())
		{
			user.GET("/courses", courseRegHandler.GetAllCourses)
			courseReg := user.Group("/course-reg")
			// courseReg.Use(middleware.CheckRegistrationPeriod(timeProvider))
			{
				courseReg.POST("/:course_id/enroll", courseRegHandler.EnrollCourse)
				courseReg.DELETE("/:course_id/enroll", courseRegHandler.CancelEnrollment)
				// courseReg.GET("/capacity", courseRegHandler.GetCoursesCapacityStatus)

				// courseReg.POST("/:course_id/waitlist", courseRegHandler.AddToWaitlist)
				// courseReg.DELETE("/:course_id/waitlist", courseRegHandler.DeleteToWaitlist)
			}
		}
	}
	return r
}

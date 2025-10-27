package routers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"

	"course-reg/internal/app/handler"
	"course-reg/internal/app/middleware"
	"course-reg/internal/pkg/constant"
	"course-reg/internal/pkg/setting"
	"course-reg/internal/pkg/util"
)

// InitRouter initialize routing information
func InitRouter(adminHandler *handler.AdminHandler, authHandler *handler.AuthHandler, courseRegHandler *handler.CourseRegHandler, timeProvider util.TimeProvider) *gin.Engine {
	gin.SetMode(setting.ServerSetting.RunMode) // set gin mode (must be called before gin.New())
	r := gin.New()
	r.Use(gin.Logger())
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

			// todo : 수강 신청 기간 중엔 학생/강의 수정 못하도록 변경
			admin.POST("students/register", adminHandler.RegisterStudents)
			admin.DELETE("students/reset", adminHandler.ResetStudents)

			admin.POST("/courses", adminHandler.CreateCourse)
			admin.DELETE("/courses/:course_id", adminHandler.DeleteCourse)
			admin.POST("/courses/register", adminHandler.RegisterCourses)
			admin.DELETE("courses/reset", adminHandler.ResetCourses)

			// todo : 수강신청 start, pause 기능 만들기
			admin.PUT("/registration-period", adminHandler.SetRegistrationPeriod)
			admin.GET("/registration-period", adminHandler.GetRegistrationPeriod)

			// admin.POST("/enrollments", adminHandler.AddEnrollment)
			// admin.DELETE("/enrollments", adminHandler.CancelEnrollment)
		}

		user := v1.Group("/user")
		user.Use(middleware.Auth())
		{
			user.StaticFile("/courses", constant.StaticCoursesFilePath)

			courseReg := user.Group("/course-reg")
			// courseReg.Use(middleware.CheckRegistrationPeriod(timeProvider))
			{
				courseReg.POST("/:course_id/enroll", courseRegHandler.EnrollCourse)
				// courseReg.DELETE("/:course_id/enroll", courseRegHandler.CancelEnrollment)

				// courseReg.GET("/capacity", courseRegHandler.GetCoursesCapacityStatus)

				// courseReg.POST("/:course_id/waitlist", courseRegHandler.AddToWaitlist)
				// courseReg.DELETE("/:course_id/waitlist", courseRegHandler.DeleteToWaitlist)
			}
		}
	}
	return r
}

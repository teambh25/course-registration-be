package routers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"

	"course-reg/handler"
	authmiddleware "course-reg/middleware"
	"course-reg/pkg/setting"
)

// InitRouter initialize routing information
func InitRouter(adminHandler *handler.AdminHandler, authHandler *handler.AuthHandler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery()) // panic 발생시 500

	// CORS 설정
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3001"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	store := memstore.NewStore([]byte(setting.SecretSetting.SessionKey)) // authentication key for session
	r.Use(sessions.Sessions("course_reg_session", store))

	// r.StaticFS("/export", http.Dir(export.GetExcelFullPath()))

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", authHandler.Logout)
		}

		admin := v1.Group("/admin")
		admin.Use(authmiddleware.AuthAdmin())
		{

			admin.POST("students/register", adminHandler.RegisterStudents)
			admin.POST("/courses", adminHandler.CreateCourse)
			// admin.DELETE("/courses", admin.Handlers.DeleteCourse)

			// admin.GET("/courses/:course_id/students", admin.Handlers.GetEnrolledStudentsByCourse)
			// admin.POST("/enrollments/force", admin.Handlers.ForceEnrollStudent)
			// admin.POST("/enrollments/cancel", admin.Handlers.CancelEnrollment)
			// admin.GET("/students/:student_id/courses", admin.Handlers.GetCoursesByStudent)
		}

		student := v1.Group("/courses")
		student.Use(authmiddleware.AuthStudent()) // amdin도 쓸 수 있어야되나..?
		{

		}

	}
	return r
}

package routers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"

	"course-reg/internal/app/handler"
	"course-reg/internal/app/middleware"
	"course-reg/internal/pkg/setting"
)

// InitRouter initialize routing information
func InitRouter(adminHandler *handler.AdminHandler, authHandler *handler.AuthHandler, studentHandler *handler.StudentHandler) *gin.Engine {
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
		}

		admin := v1.Group("/admin")
		admin.Use(middleware.AuthAdmin())
		{

			admin.POST("students/register", adminHandler.RegisterStudents)
			admin.DELETE("students/reset", adminHandler.ResetStudents)

			admin.POST("/courses", adminHandler.CreateCourse)
			admin.DELETE("/courses/:id", adminHandler.DeleteCourse)

			// admin.POST("/enrollments", adminHandler.AddEnrollment)
			// admin.DELETE("/enrollments", adminHandler.CancelEnrollment)
		}

		student := v1.Group("/courses")
		student.Use(middleware.AuthStudent()) // amdin도 쓸 수 있어야되나..?
		{

		}

	}
	return r
}

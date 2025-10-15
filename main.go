package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"course-reg/handler"
	"course-reg/models"
	"course-reg/pkg/setting"
	"course-reg/pkg/util"
	"course-reg/repository"
	"course-reg/routers"
	"course-reg/service"
)

func init() {
	setting.Setup()
	// logging.Setup() // logging 할 수 있는 환경일 떄 다시 사용
	util.Setup()
}

func main() {
	db := models.Setup()
	studentRepo := repository.NewStudentRepository(db)
	courseRepo := repository.NewCourseRepository(db)
	enrollRepo := repository.NewEnrollmentRepository(db)

	authService := service.NewAuthService(studentRepo)
	authHandler := handler.NewAuthHandler(authService)

	adminService := service.NewAdminService(studentRepo, courseRepo, enrollRepo)
	adminHandler := handler.NewAdminHandler(adminService)

	gin.SetMode(setting.ServerSetting.RunMode)
	routersInit := routers.InitRouter(adminHandler, authHandler)

	readTimeout := setting.ServerSetting.ReadTimeout
	writeTimeout := setting.ServerSetting.WriteTimeout
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)
	maxHeaderBytes := 1 << 20

	server := &http.Server{
		Addr:           endPoint,
		Handler:        routersInit,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	log.Printf("[info] start http server listening %s", endPoint)

	server.ListenAndServe()
}

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"course-reg/internal/app/handler"
	"course-reg/internal/app/repository"
	"course-reg/internal/app/routers"
	"course-reg/internal/app/service"
	"course-reg/internal/app/worker"
	"course-reg/internal/pkg/database"
	"course-reg/internal/pkg/setting"
	"course-reg/internal/pkg/util"
)

func init() {
	setting.Setup()
	// logging.Setup() // logging 할 수 있는 환경일 떄 다시 사용
	util.Setup()
}

func main() {
	db := database.Setup()
	studentRepo := repository.NewStudentRepository(db)
	courseRepo := repository.NewCourseRepository(db)
	enrollRepo := repository.NewEnrollmentRepository(db)

	// todo : log 파일을 통한 배치 쓰기로 변겨할 예정
	// dbWorker := worker.NewDBWorker(enrollRepo)
	// dbWorker.Start()

	enrollmentWorker := worker.NewEnrollmentWorker(1000)

	// todo : 이전의 db 데이터 로딩
	// courses, err := courseRepo.FetchAllCourses()
	// if err != nil {
	// 	log.Printf("[warning] failed to load courses: %v", err)
	// 	courses = []models.Course{}
	// }
	// enrollments, err := enrollRepo.LoadAllEnrollments()
	// if err != nil {
	// 	log.Printf("[warning] failed to load enrollments: %v", err)
	// 	enrollments = []models.Enrollment{}
	// }

	enrollmentWorker.Start()

	authService := service.NewAuthService(studentRepo)
	adminService := service.NewAdminService(studentRepo, courseRepo, enrollRepo, enrollmentWorker)
	courseRegService := service.NewCourseRegService(courseRepo, enrollRepo, enrollmentWorker)

	authHandler := handler.NewAuthHandler(authService)
	adminHandler := handler.NewAdminHandler(adminService)
	courseRegHandler := handler.NewCourseRegHandler(courseRegService)

	gin.SetMode(setting.ServerSetting.RunMode)
	routersInit := routers.InitRouter(adminHandler, authHandler, courseRegHandler)

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

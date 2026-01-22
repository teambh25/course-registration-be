package main

import (
	"fmt"
	"log"
	"net/http"

	"course-reg/internal/app"
	"course-reg/internal/pkg/setting"
)

func init() {
	setting.Setup()
	// logging.Setup() // todo : 파일에 로깅할 수 있는 환경이 아니라 삭제해도 될거 같은데?
}

func main() {
	// Initialize application with all dependencies
	application, err := app.NewApplication()
	if err != nil {
		log.Fatalf("failed to initialize application: %v", err)
	}

	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.ServerSetting.HttpPort),
		Handler:        application.Router,
		ReadTimeout:    setting.ServerSetting.ReadTimeout,
		WriteTimeout:   setting.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("[info] start http server listening %s", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

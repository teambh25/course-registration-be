package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"course-reg/internal/app"
	"course-reg/internal/pkg/setting"

	"github.com/gin-contrib/pprof"
)

var profEnabled bool

func init() {
	flag.BoolVar(&profEnabled, "prof", false, "Enable pprof profiling endpoints")
	flag.Parse()
}

func main() {
	cfg := setting.Load()

	// Initialize application with all dependencies
	application, err := app.NewApplication(cfg)
	if err != nil {
		log.Fatalf("failed to initialize application: %v", err)
	}

	// Enable pprof if --prof flag is set
	if profEnabled {
		pprof.Register(application.Router)
		log.Println("[info] pprof profiling enabled at /debug/pprof/")
	}

	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", cfg.Server.HttpPort),
		Handler:        application.Router,
		ReadTimeout:    cfg.Server.ReadTimeout,
		WriteTimeout:   cfg.Server.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("[info] start http server listening %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("[info] received shutdown signal")

	// Create a deadline to wait for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // < Docker stop timeout (default: 10s)
	defer cancel()

	// Shutdown HTTP server
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("[error] server shutdown error: %v", err)
	}

	// Shutdown application (worker, database connections, etc.)
	if err := application.Shutdown(); err != nil {
		log.Printf("[error] application shutdown error: %v", err)
	}

	log.Println("[info] server exited")
}

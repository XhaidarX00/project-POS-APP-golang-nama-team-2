package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"project_pos_app/cmd"
	_ "project_pos_app/docs"
	"project_pos_app/infra"
	"project_pos_app/router"
	"syscall"
	"time"
)

// @title App pos Team 2
// @version 1.0
// @description API for Point Of Sale
// @termsOfService http://example.com/terms/
// @contact.name API Support
// @contact.url http://example.com/support
// @contact.email support@example.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey Authentication
// @in header
// @name Authorization
// @securityDefinitions.apikey UserID
// @in header
// @name User-ID
// @securityDefinitions.apikey UserRole
// @in header
// @name User-Role

func main() {
	// Test()
	ctx, err := infra.NewIntegrateContext()
	if err != nil {
		log.Fatal("can't init service context %w", err)
	}

	r := router.NewRoutes(ctx)

	// init cron
	cmd.CronJob(ctx)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		// Start the server
		log.Printf("Server running on port 8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	// Create a timeout context for graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Attempt to gracefully shutdown the server
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	// Catching context timeout
	select {
	case <-shutdownCtx.Done():
		log.Println("Timeout of 5 seconds.")
	}

	log.Println("Server exiting")
}

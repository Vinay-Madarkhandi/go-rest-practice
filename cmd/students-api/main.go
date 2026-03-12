package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Vinay-Madarkhandi/go-rest-practice/internal/config"
	"github.com/Vinay-Madarkhandi/go-rest-practice/internal/http/handlers/student"
	"github.com/Vinay-Madarkhandi/go-rest-practice/internal/storage/mysql"
)

func main() {

	// Load config
	cfg := config.MustLoad()

	// Database connection
	db, err := mysql.New(cfg)
	if err != nil {
		slog.Error("database connection failed", "error", err)
		os.Exit(1)

	}
	slog.Info("database connection established")

	// Setup router
	router := http.NewServeMux()

	router.HandleFunc("POST /api/v1/students", student.NewStudent(db))
	router.HandleFunc("GET /api/v1/students/{id}", student.GetStudentByID(db))

	// Setup server
	server := http.Server{
		Addr:    cfg.HTTPServerConfig.Address,
		Handler: router,
	}

	slog.Info("starting http server at ", slog.String("address", cfg.HTTPServerConfig.Address))

	// Gracefully shutdown should always happen over here this is the best practice
	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("Failed to start server", err)
		}
	}()

	<-done

	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		slog.Error("Failed to shutdown server:", slog.String("error", err.Error()))
	}

	slog.Info("Server shutdown complete")

}

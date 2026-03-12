package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Vinay-Madarkhandi/go-rest-practice/internal/config"
)

func main() {
	// Load config
	cfg := config.MustLoad()
	// Database connection
	// Setup router
	router := http.NewServeMux()

	router.HandleFunc("GET /api/v1/students", func(writer http.ResponseWriter, request *http.Request) {
		_, err := writer.Write([]byte("students"))
		if err != nil {
			return
		}

	})
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
		if err != nil {
			log.Fatal("Failed to start server", err)
		}
	}()

	<-done

	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		slog.Error("Failed to shutdown server:", slog.String("error", err.Error()))
	}

	slog.Info("Server shutdown complete")

}

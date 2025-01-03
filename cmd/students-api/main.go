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

	"github.com/tushar0305/students-api/internal/config"
	"github.com/tushar0305/students-api/internal/http/handlers/student"
)

func main() {
	// load config
	cfg := config.MustLoad()


	// database setup

	// setup router
	router := http.NewServeMux()
	
	router.HandleFunc("POST /api/students",  student.New())

	// http server setup
	server := http.Server{
		Addr:    cfg.HTTPServer.Address,
		Handler: router,
	}

	slog.Info("Starting Server", slog.String("address", cfg.HTTPServer.Address))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err:= server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to start the server", err)
		}
	} ()

	<-done

	slog.Info("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		log.Fatal("Server Shutdown Failed", slog.String("error", err.Error()))
	}

	slog.Info("Server Exited Properly")
	
}
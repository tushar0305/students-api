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
	"github.com/tushar0305/students-api/internal/storage/sqlite"
)

func main() {
	// load config
	cfg := config.MustLoad()

	// create database file if it does not exist
    if _, err := os.Stat(cfg.StoragePath); os.IsNotExist(err) {
        file, err := os.Create(cfg.StoragePath)
        if err != nil {
            log.Fatal("Failed to create database file", slog.String("error", err.Error()))
        }
        file.Close()
    }


	// database setup

	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal("Failed to create storage", slog.String("error", err.Error()))
	}

	slog.Info("Storage Created", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	// setup router
	router := http.NewServeMux()
	
	router.HandleFunc("POST /api/students",  student.New(storage))

	router.HandleFunc("GET /api/students/{id}",  student.GetById(storage))

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

	err = server.Shutdown(ctx)
	if err != nil {
		log.Fatal("Server Shutdown Failed", slog.String("error", err.Error()))
	}

	slog.Info("Server Exited Properly")
	
}
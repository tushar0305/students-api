package main

import (
	"fmt"
	"net/http"
	"github.com/tushar0305/students-api/internal/config"
)

func main() {
	// load config
	cfg := config.MustLoad()


	// database setup

	// setup router
	router := http.NewServeMux()
	
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	// http server setup
	server := http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}
	fmt.Printf("Server Started %s", cfg.HTTPServer.Address)

	// start the server
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("server failed to start: %v\n", err)
	}

	
}
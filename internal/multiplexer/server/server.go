package server

import (
	"github.com/F7icK/multiplexer/internal/multiplexer/server/handlers"
	"log"
	"net/http"
	"time"
)

func StartServer(handlers *handlers.Handlers, port string) {
	router := NewRouter(handlers)

	srv := &http.Server{
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       10 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		Handler:           router,
		Addr:              port,
	}

	log.Println("Start server")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

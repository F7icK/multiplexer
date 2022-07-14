package server

import (
	"github.com/F7icK/multiplexer/internal/multiplexer/server/handlers"
	"log"
	"net/http"
	"time"
)

func StartServer(handlers *handlers.Handlers, timeoutIncoming time.Duration, port string) {
	router := NewRouter(handlers)

	srv := &http.Server{
		WriteTimeout: timeoutIncoming * time.Second,
		Handler:      router,
		Addr:         port,
	}

	log.Println("Start server")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

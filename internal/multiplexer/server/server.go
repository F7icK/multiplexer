package server

import (
	"log"
	"net/http"
	"time"

	"github.com/F7icK/multiplexer/internal/multiplexer/server/handlers"
	"github.com/F7icK/multiplexer/internal/multiplexer/types"
)

func StartServer(handler *handlers.Handlers, cfg *types.Config) {
	router := NewRouter(handler)

	srv := &http.Server{
		WriteTimeout: cfg.TimeoutIncoming * time.Second,
		Handler:      router,
		Addr:         cfg.Port,
	}

	log.Printf("Server running on port %s", cfg.Port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

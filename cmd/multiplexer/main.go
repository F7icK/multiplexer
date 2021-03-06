package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/F7icK/multiplexer/config"
	"github.com/F7icK/multiplexer/internal/multiplexer/server"
	"github.com/F7icK/multiplexer/internal/multiplexer/server/handlers"
	"github.com/F7icK/multiplexer/internal/multiplexer/service"
)

func main() {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	defer close(signalCh)

	cfg := config.NewConfig()

	srv, err := service.NewService(cfg)
	if err != nil {
		log.Printf("err with NewService in main: %s", err)
		return
	}

	handler := handlers.NewHandlers(srv)

	go func(signalCh <-chan os.Signal) {
		select {
		case sig := <-signalCh:
			log.Printf("stopped with signal: %s", sig)
			os.Exit(0)
		}
	}(signalCh)

	server.StartServer(handler, cfg)
}

package main

import (
	"github.com/F7icK/multiplexer/internal/multiplexer/server"
	"github.com/F7icK/multiplexer/internal/multiplexer/server/handlers"
	"github.com/F7icK/multiplexer/internal/multiplexer/service"
	"log"
	"os"
	"os/signal"
)

func main() {

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	defer close(signalCh)

	srvce := service.NewService(100, 1)

	endpoints := handlers.NewHandlers(srvce)

	go func(signalCh <-chan os.Signal) {
		select {
		case sig := <-signalCh:
			log.Printf("stopped with signal: %s", sig)
			os.Exit(0)
		}
	}(signalCh)

	server.StartServer(endpoints, ":8080")
}

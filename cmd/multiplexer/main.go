package main

import (
	"github.com/F7icK/multiplexer/config"
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

	cfg := config.NewConfig()

	srvce := service.NewService(cfg.LimitConnection, cfg.TimeoutOutgoing, cfg.LimitGoRoutines)

	endpoints := handlers.NewHandlers(srvce)

	go func(signalCh <-chan os.Signal) {
		select {
		case sig := <-signalCh:
			log.Printf("stopped with signal: %s", sig)
			os.Exit(0)
		}
	}(signalCh)

	server.StartServer(endpoints, cfg.TimeoutIncoming, ":8080")
}

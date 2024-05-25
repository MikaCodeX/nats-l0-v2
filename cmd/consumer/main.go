package main

import (
	"os"
	"os/signal"
	"syscall"

	"awesome/internal/rest"
	"awesome/internal/store"
	"awesome/internal/sub"

	"github.com/rs/zerolog/log"
)

func main() {
	storage := store.New()
	storage.GetOrderfromDB()
	sub := sub.New(storage)
	sub.Run()

	endpoint := rest.New(storage)
	go endpoint.Run()
	ctrlC()
	storage.Close()
	sub.Close()

}

func ctrlC() os.Signal {
	log.Print("starting http service")
	signals := make(chan os.Signal, 2)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)
	sign := <-signals
	log.Print("stopped http service")
	return sign

}

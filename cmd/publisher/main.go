package main

import (
	"log"
	"time"

	"awesome/internal/generateOrder"

	"github.com/nats-io/stan.go"
)

const (
	clusterID = "test-cluster"
	clientId  = "publisher"
)

func main() {
	connectedToServer, err := stan.Connect(clusterID, clientId, stan.NatsURL("nats://localhost:4223")) // Переменная для подключения
	if err != nil {
		log.Fatalf("Nats-connection failed :%s", err)
	}
	defer connectedToServer.Close()

	ackHandler := func(ackedNuid string, err error) {}

	t := time.NewTicker(time.Second)
	for range t.C {
		_, msg := generateOrder.RandomOrder()
		_, err := connectedToServer.PublishAsync("Orders", msg, ackHandler)
		if err != nil {
			log.Fatalf("Error :%s", err)
		}
		log.Println("Выдаем инфу в канал")

	}

}

package main

import (
	"fmt"
	"github.com/nats-io/go-nats"
	"time"
)

func main() {
	nc, _ := nats.Connect("nats://0.0.0.0:4222")

	nc.Subscribe("Supermercados.>", func(m *nats.Msg) {
		fmt.Printf("Canal: %s - Mensaje recibido: %s\n", m.Subject, string(m.Data))
	})
	for {
		// Simple Publisher
		nc.Publish("Supermercados.Murcia", []byte("Buenos días Murcia!"))
		nc.Publish("Supermercados.Andalucia", []byte("Buenos días Andalucía!"))

		time.Sleep(1 * time.Minute)
	}
	nc.Close()
}

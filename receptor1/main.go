package main

import (
	"fmt"
	"github.com/nats-io/go-nats"
	"time"
)

func main() {
	nc, _ := nats.Connect("nats://0.0.0.0:4222")
	// Simple Async Subscriber
	nc.Subscribe("Supermercados.Murcia", func(m *nats.Msg) {
		fmt.Printf("Mensaje recibido: %s\n", string(m.Data))
	})
	for {
		fmt.Println("--------------------")
		nc.Publish("Supermercados.Murcia", []byte("Comercial de Murcia!"))
		time.Sleep(24 * time.Second)
	}
}

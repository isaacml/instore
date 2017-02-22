package main

import (
	"github.com/nats-io/go-nats"
	"time"
	"github.com/isaacml/instore/libs"
)

func main() {
	nc, _ := nats.Connect("nats://0.0.0.0:4222")
	// Simple Async Subscriber
	nc.Subscribe("Imagenes.Google", handler)

	for {
		time.Sleep(1 * time.Minute)
	}
	nc.Close()
}

func handler(m *nats.Msg) {
	libs.GetFile("/home/isaac/PruebasNats/Gooogle/", string(m.Data))

}
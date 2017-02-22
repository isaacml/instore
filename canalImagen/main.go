package main

import (
	"fmt"
	"github.com/nats-io/go-nats"
	"github.com/isaacml/instore/libs"
	"time"
)

func main() {
	nc, _ := nats.Connect("nats://0.0.0.0:4222")
	nc.Subscribe("Imagenes.>", handler)

	//publicacion
	for {
		nc.Publish("Imagenes.Twitter", libs.SendFile("/home/isaac/Documentos/", "logo.png"))
		nc.Publish("Imagenes.Google", libs.SendFile("/home/isaac/Documentos/", "biwenger.png"))
		time.Sleep(1 * time.Minute)
	}
	nc.Close()
}

func handler(m *nats.Msg) {
	fmt.Printf("Canal: %s - Mensaje recibido: %s\n", m.Subject, string(m.Data))
}
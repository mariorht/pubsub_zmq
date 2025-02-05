package main

import (
	"fmt"
	"time"

    "go_pubsub_zmq"
)

func main() {
	endpoint := "tcp://127.0.0.1:5559"
	pub, err := go_pubsub_zmq.NewPublisher(endpoint)
	if err != nil {
		fmt.Printf("Error al crear Publisher: %v\n", err)
		return
	}
	defer pub.Close()

	for i := 0; i < 10; i++ {
		msg := fmt.Sprintf("Mensaje %d", i)
		if err := pub.Publish(msg); err != nil {
			fmt.Printf("Error al publicar mensaje: %v\n", err)
		} else {
			fmt.Printf("Publicado: %s\n", msg)
		}
		time.Sleep(1 * time.Second)
	}
}

package main

import (
	"fmt"

    "go_pubsub_zmq"
)

func main() {
	endpoint := "tcp://127.0.0.1:5559"
	sub, err := go_pubsub_zmq.NewSubscriber(endpoint, "")
	if err != nil {
		fmt.Printf("Error al crear Subscriber: %v\n", err)
		return
	}
	defer sub.Close()

	for {
		msg, err := sub.Receive()
		if err != nil {
			fmt.Printf("Error al recibir mensaje: %v\n", err)
			break
		}
		fmt.Printf("Recibido: %s\n", msg)
	}
}

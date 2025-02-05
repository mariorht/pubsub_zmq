package main

import (
	"fmt"

	"go_pubsub_zmq"
)

func main() {
	endpoint := "tcp://127.0.0.1:5559"
	topic := "test"

	// Crear Subscriber
	sub, err := go_pubsub_zmq.NewSubscriber(endpoint, topic)
	if err != nil {
		fmt.Printf("❌ Error al crear Subscriber: %v\n", err)
		return
	}
	defer sub.Close()

	fmt.Println("📥 Subscriber conectado. Esperando mensajes...")

	for {
		msg, err := sub.ReceiveMessage()
		if err != nil {
			fmt.Printf("❌ Error al recibir mensaje: %v\n", err)
			break
		}

		// Mostrar los datos JSON recibidos
		fmt.Printf("📥 Mensaje recibido: %+v\n", msg.Data)

		// Verificar si el mensaje tiene imágenes
		if len(msg.Images) > 0 {
			fmt.Println("⚠️ Advertencia: Se recibió metadata de imágenes, pero el procesamiento aún no está implementado.")
		}
	}
}

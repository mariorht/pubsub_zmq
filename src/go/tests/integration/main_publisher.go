package main

import (
	"fmt"
	"time"

	"go_pubsub_zmq"
)

func main() {
	endpoint := "tcp://*:5559"
	topic := "test"
	chunkSize := 1024

	pub, err := go_pubsub_zmq.NewPublisher(endpoint, topic, chunkSize)
	if err != nil {
		fmt.Printf("❌ Error al crear Publisher: %v\n", err)
		return
	}
	defer pub.Close()

	fmt.Println("✅ Publisher (Go) iniciado. Esperando suscriptores...")

	time.Sleep(2 * time.Second) // Para evitar el Slow Joiner Syndrome

	for i := 0; i < 1; i++ {
		data := map[string]interface{}{
			"key":   "probando",
			"index": i,
		}

		messageBytes, err := pub.BuildMessage(nil, data) // nil porque no enviamos imágenes reales
		if err != nil {
			fmt.Printf("❌ Error al construir mensaje: %v\n", err)
			continue
		}

		if err := pub.PublishMessage(messageBytes); err != nil {
			fmt.Printf("❌ Error al publicar mensaje: %v\n", err)
		} else {
			fmt.Printf("📤 Mensaje %d enviado con data: %v\n", i, data)
		}

		time.Sleep(1 * time.Second)
	}

	fmt.Println("✅ Publisher (Go) finalizado tras enviar 5 mensajes.")
}

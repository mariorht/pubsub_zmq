package main

import (
	"fmt"
	"time"

	"go_pubsub_zmq"
)

func main() {
	endpoint := "tcp://127.0.0.1:5555"
	topic := "camera/image_raw"
	chunkSize := 1024

	// Crear Publisher
	pub, err := go_pubsub_zmq.NewPublisher(endpoint, topic, chunkSize)
	if err != nil {
		fmt.Printf("❌ Error al crear Publisher: %v\n", err)
		return
	}
	defer pub.Close()

	fmt.Println("📤 Publisher iniciado. Enviando mensajes...")

	for i := 0; i < 10; i++ {
		// Construir mensaje sin imágenes y con datos JSON
		data := map[string]interface{}{
			"index":   i,
			"message": fmt.Sprintf("Mensaje %d enviado desde el publisher", i),
			"timestamp": time.Now().Format(time.RFC3339),
		}

		messageBytes, err := pub.BuildMessage(nil, data)
		if err != nil {
			fmt.Printf("❌ Error al construir mensaje %d: %v\n", i, err)
			continue
		}

		// Publicar el mensaje
		if err := pub.PublishMessage(messageBytes); err != nil {
			fmt.Printf("❌ Error al publicar mensaje %d: %v\n", i, err)
		} else {
			fmt.Printf("✅ Publicado: %s\n", data["message"])
		}

		time.Sleep(1 * time.Second) // Simula un intervalo entre mensajes
	}

	fmt.Println("📤 Publicación finalizada.")
}

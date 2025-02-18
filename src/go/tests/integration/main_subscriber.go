package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"go_pubsub_zmq"
)

func main() {
	endpoint := "tcp://integration_python_pub:5559" //El nombre del docker
	topic := "test"

	sub, err := go_pubsub_zmq.NewSubscriber(endpoint, topic)
	if err != nil {
		log.Fatalf("❌ Error al crear Subscriber: %v", err)
	}
	defer sub.Close()

	var lastMessage map[string]interface{}

	for i := 0; i < 1; i++ {
		msg, err := sub.ReceiveMessage()
		if err != nil {
			log.Fatalf("❌ Error al recibir mensaje: %v", err)
		}

		fmt.Printf("📥 Mensaje %d recibido: %+v\n", i, msg.Data)

		if len(msg.Images) > 0 {
			fmt.Printf("⚠️ Advertencia: Recibidas %d imágenes. Procesamiento no implementado.\n", len(msg.Images))
		}

		lastMessage = msg.Data
	}

	// Escribir el último mensaje recibido en un JSON en el volumen compartido
	file, err := os.Create("/shared/result.json")
	if err != nil {
		log.Fatalf("❌ Error al crear el archivo JSON: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(lastMessage)
	if err != nil {
		log.Fatalf("❌ Error al escribir JSON: %v", err)
	}

	fmt.Println("✅ Resultado guardado en /shared/result.json")
}

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"go_pubsub_zmq"
)

func main() {
	endpoint := "tcp://*:5555"
	topic := "test"
	chunkSize := 1024

	pub, err := go_pubsub_zmq.NewPublisher(endpoint, topic, chunkSize)
	if err != nil {
		fmt.Printf("❌ Error al crear Publisher: %v\n", err)
		return
	}
	defer pub.Close()

	fmt.Println("✅ Publisher (Go) iniciado. Esperando suscriptores...")
	time.Sleep(2 * time.Second)

	messages := 1
	for i := 0; i < messages; i++ {
		data := map[string]interface{}{
			"key":   "probando",
			"index": i,
		}

		messageBytes, err := pub.BuildMessage(nil, data, "")
		if err != nil {
			fmt.Printf("❌ Error al construir mensaje: %v\n", err)
			continue
		}

		if err := pub.PublishMessage(messageBytes); err != nil {
			fmt.Printf("❌ Error al publicar mensaje: %v\n", err)
		} else {
			fmt.Printf("📤 Mensaje %d enviado con data: %v\n", i, data)

			lastMessage := map[string]interface{}{
				"type":   "images",
				"count":  0,
				"images": []interface{}{},
				"data":   data,
			}

			// Guardar el último mensaje enviado
			file, err := os.Create("/shared/result_publisher.json")
			if err != nil {
				fmt.Printf("❌ Error al guardar JSON en publisher: %v\n", err)
				continue
			}
			encoder := json.NewEncoder(file)
			_ = encoder.Encode(lastMessage)
			file.Close()
		}

		time.Sleep(1 * time.Second)
	}

	fmt.Println("✅ Publisher (Go) finalizado tras enviar ", messages,  " mensajes.")
}

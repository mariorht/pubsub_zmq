package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"

	"go_pubsub_zmq"
)

func saveImage(img image.Image, filename string, format string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	switch format {
	case "jpeg":
		err = jpeg.Encode(file, img, nil)
	case "png":
		err = png.Encode(file, img)
	default:
		return fmt.Errorf("Formato de imagen no soportado: %s", format)
	}

	return err
}

func main() {
	endpoint := os.Getenv("PUBSUB_ENDPOINT")
	if endpoint == "" {
		endpoint = "tcp://localhost:5555"
	}

	topic := "test"
	sub, err := go_pubsub_zmq.NewSubscriber(endpoint, topic)
	if err != nil {
		log.Fatalf("❌ Error al crear Subscriber: %v", err)
	}
	defer sub.Close()

	var lastMessage map[string]interface{}
	var savedImages []string

	
	msg, images, err := sub.ReceiveMessage()
	if err != nil {
		log.Fatalf("❌ Error al recibir mensaje: %v", err)
	}

	fmt.Printf("📥 Mensaje recibido: %+v\n", msg.Data)

	// Guardar imágenes si se recibieron
	if len(images) > 0 {
		fmt.Printf("⚠️ Recibidas %d imágenes. Guardándolas...\n", len(images))
		for idx, img := range images {
			filename := fmt.Sprintf("/shared/image_%d.png", idx)
			err := saveImage(img, filename, "png")
			if err != nil {
				log.Printf("❌ Error al guardar imagen %d: %v", idx, err)
				continue
			}
			savedImages = append(savedImages, filepath.Base(filename))
		}
	}

	lastMessage = msg.Data
	

	// Escribir el último mensaje recibido en un JSON en el volumen compartido
	var imagesMetadata []map[string]interface{}

	for idx, img := range images {
		filename := fmt.Sprintf("/shared/image_%d.png", idx)
		err := saveImage(img, filename, "png")
		if err != nil {
			log.Printf("❌ Error al guardar imagen %d: %v", idx, err)
			continue
		}
	
		// 🔹 Agregar metadatos en lugar del nombre del archivo
		meta := map[string]interface{}{
			"width":    img.Bounds().Dx(),
			"height":   img.Bounds().Dy(),
			"channels": 3,
			"dtype":    "uint8",
			"size":     img.Bounds().Dx() * img.Bounds().Dy() * 3, // RGB
		}
		imagesMetadata = append(imagesMetadata, map[string]interface{}{"metadata": meta})
	}

	result := map[string]interface{}{
		"type":   "images",
		"count":  len(savedImages),
		"images": imagesMetadata,
		"data":   lastMessage,
	}

	file, err := os.Create("/shared/result.json")
	if err != nil {
		log.Fatalf("❌ Error al crear el archivo JSON: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(result)
	if err != nil {
		log.Fatalf("❌ Error al escribir JSON: %v", err)
	}

	fmt.Println("✅ Resultado guardado en /shared/result.json")
}

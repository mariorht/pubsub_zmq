package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"os"
	"time"

	"go_pubsub_zmq"
)

// createTestImage genera una imagen de prueba en memoria.
func createTestImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for x := 0; x < 100; x++ {
		for y := 0; y < 100; y++ {
			img.Set(x, y, color.RGBA{uint8(x * 2), uint8(y * 2), 200, 255})
		}
	}
	return img
}

// loadOrCreateImage intenta cargar una imagen desde un archivo, y si no la encuentra, la genera.
func loadOrCreateImage(path string) image.Image {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("⚠️ No se encontró la imagen, generando una de prueba.")
		return createTestImage()
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("⚠️ Error al decodificar la imagen, generando una de prueba.")
		return createTestImage()
	}
	return img
}

// getImageMetadata obtiene los metadatos de una imagen.
func getImageMetadata(img image.Image) map[string]interface{} {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	channels := 3 // PNG y JPEG son típicamente RGB
	dtype := "uint8"
	size := width * height * channels // Tamaño estimado en bytes

	return map[string]interface{}{
		"metadata": map[string]interface{}{
			"width":    width,
			"height":   height,
			"channels": channels,
			"dtype":    dtype,
			"size":     size,
		},
	}
}

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

		// Cargar imagen o generar una si no existe
		img := loadOrCreateImage("/assets/pong.png")

		// Obtener metadatos de la imagen
		imageMetadata := getImageMetadata(img)

		// Construir mensaje con la imagen
		messageBytes, err := pub.BuildMessage([]image.Image{img}, data, "png")
		if err != nil {
			fmt.Printf("❌ Error al construir mensaje: %v\n", err)
			continue
		}

		// Publicar mensaje
		if err := pub.PublishMessage(messageBytes); err != nil {
			fmt.Printf("❌ Error al publicar mensaje: %v\n", err)
		} else {
			fmt.Printf("📤 Mensaje %d enviado con imagen y data: %v\n", i, data)

			// Crear JSON con metadatos de la imagen
			lastMessage := map[string]interface{}{
				"type":   "images",
				"count":  1,
				"images": []interface{}{imageMetadata}, // 🔹 Guardar metadatos en lugar de nombres de archivos
				"data":   data,
			}

			// Guardar el último mensaje enviado
			jsonFile, err := os.Create("/shared/result_publisher.json")
			if err != nil {
				fmt.Printf("❌ Error al guardar JSON en publisher: %v\n", err)
				continue
			}
			encoder := json.NewEncoder(jsonFile)
			_ = encoder.Encode(lastMessage)
			jsonFile.Close()
		}

		time.Sleep(1 * time.Second)
	}

	fmt.Println("✅ Publisher (Go) finalizado tras enviar ", messages, " mensajes.")
}

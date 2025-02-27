package go_pubsub_zmq_test

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"image/jpeg"
	"testing"
	"time"
	"strconv"

	"go_pubsub_zmq"
)

// createTestImage genera una imagen de prueba de 10x10 con un patrón simple.
func createTestImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			img.Set(x, y, color.RGBA{uint8(x * 10), uint8(y * 10), 100, 255})
		}
	}
	return img
}

// Test para verificar envío y recepción de imágenes PNG
func TestPubSubPNG(t *testing.T) {
	endpoint := "tcp://127.0.0.1:5555"
	topic := "test"
	chunkSize := 1024

	pub, err := go_pubsub_zmq.NewPublisher(endpoint, topic, chunkSize)
	if err != nil {
		t.Fatalf("Error inesperado: %v", err)
	}
	defer pub.Close()

	sub, err := go_pubsub_zmq.NewSubscriber(endpoint, topic)
	if err != nil {
		t.Fatalf("Error inesperado: %v", err)
	}	
	defer sub.Close()

	time.Sleep(500 * time.Millisecond) // Esperar conexión

	img := createTestImage()
	data := map[string]interface{}{"msg": "test_png"}

	msgBytes, err := pub.BuildMessage([]image.Image{img}, data, "png")
	if err != nil {
		t.Fatalf("Error inesperado: %v", err)
	}
	

	err = pub.PublishMessage(msgBytes)
	if err != nil {
		t.Fatalf("Error inesperado: %v", err)
	}
	

	receivedMsg, images, err := sub.ReceiveMessage()
	if err != nil {
		t.Fatalf("Error inesperado: %v", err)
	}
	

	if receivedMsg.Data["msg"] != "test_png" {
		t.Errorf("Esperaba 'test_png', pero recibí: %v", receivedMsg.Data["msg"])
	}
	
	if len(images) != 1 {
		t.Fatalf("Esperaba 1 imagen y se recibieron: %d", len(images))
	}
	

	// Verificar si la imagen decodificada es válida
	var buf bytes.Buffer
	err = png.Encode(&buf, images[0])
	if err != nil {
		t.Fatalf("Error inesperado: %v", err)
	}
	
	if buf.Len() <= 0 {
		t.Fatalf("La imagen PNG recibida no es válida, tamaño del buffer: %d", buf.Len())
	}
	}

// Test para verificar envío y recepción de imágenes JPEG
func TestPubSubJPEG(t *testing.T) {
	endpoint := "tcp://127.0.0.1:5556"
	topic := "test"
	chunkSize := 1024

	pub, err := go_pubsub_zmq.NewPublisher(endpoint, topic, chunkSize)
	if err != nil {
		t.Fatalf("Error inesperado: %v", err)
	}	
	defer pub.Close()

	sub, err := go_pubsub_zmq.NewSubscriber(endpoint, topic)
	if err != nil {
		t.Fatalf("Error inesperado: %v", err)
	}	
	defer sub.Close()

	time.Sleep(500 * time.Millisecond) // Esperar conexión

	img := createTestImage()
	data := map[string]interface{}{"msg": "test_jpeg"}

	msgBytes, err := pub.BuildMessage([]image.Image{img}, data, "jpeg")
	if err != nil {
		t.Fatalf("Error inesperado: %v", err)
	}
	err = pub.PublishMessage(msgBytes)
	if err != nil {
		t.Fatalf("Error inesperado: %v", err)
	}


	receivedMsg, images, err := sub.ReceiveMessage()
	if err != nil {
		t.Fatalf("Error inesperado: %v", err)
	}
	
	if receivedMsg.Data["msg"] != "test_jpeg" {
		t.Errorf("Esperaba 'test_jpeg', pero recibí: %v", receivedMsg.Data["msg"])
	}

	if len(images) != 1 {
		t.Fatalf("Esperaba 1 imagen y se recibieron: %d", len(images))
	}


	// Verificar si la imagen decodificada es válida
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, images[0], nil)
	if err != nil {
		t.Fatalf("Error inesperado: %v", err)
	}
	if buf.Len() <= 0 {
		t.Fatalf("La imagen JPG recibida no es válida, tamaño del buffer: %d", buf.Len())
	}
	}

func TestPubSubMensajeUnico(t *testing.T) {
	endpoint := "tcp://127.0.0.1:5555"
	topic := "test"
	chunkSize := 1024

	// Crea el Publisher usando el package importado.
	pub, err := go_pubsub_zmq.NewPublisher(endpoint, topic, chunkSize)
	if err != nil {
		t.Fatalf("Error al crear Publisher: %v", err)
	}
	defer pub.Close()

	// Crea el Subscriber suscrito al mismo topic.
	sub, err := go_pubsub_zmq.NewSubscriber(endpoint, topic)
	if err != nil {
		t.Fatalf("Error al crear Subscriber: %v", err)
	}
	defer sub.Close()

	testMsg := "Mensaje de prueba"

	// Asegurar que el subscriber tiene tiempo para conectarse
	time.Sleep(500 * time.Millisecond)

	// Publica el mensaje después de un pequeño retardo.
	go func() {
		time.Sleep(100 * time.Millisecond) // Asegura que el subscriber esté listo.
		messageBytes, err := pub.BuildMessage(nil, map[string]interface{}{"msg": testMsg}, "")
		if err != nil {
			t.Errorf("Error al construir el mensaje: %v", err)
			return
		}
		if err := pub.PublishMessage(messageBytes); err != nil {
			t.Errorf("Error al publicar mensaje: %v", err)
		}
	}()

	// Reintentos para asegurar la recepción del mensaje
	var msgReceived *go_pubsub_zmq.MessageReceived // Declaramos fuera para que sea accesible
	for retry := 0; retry < 5; retry++ {
		msgReceived, _, err = sub.ReceiveMessage() // Ahora usamos la variable ya declarada
		if err == nil {
			break
		}
		time.Sleep(100 * time.Millisecond) // Espera antes del siguiente intento
	}

	if err != nil {
		t.Fatalf("Error al recibir mensaje después de varios intentos: %v", err)
	}

	// Validar que el mensaje recibido coincida con el enviado
	if received, ok := msgReceived.Data["msg"].(string); !ok || received != testMsg {
		t.Errorf("El mensaje recibido (%v) no coincide con el esperado (%v)", msgReceived.Data["msg"], testMsg)
	}
}

func TestPubSubMultiplesMensajes(t *testing.T) {
	endpoint := "tcp://127.0.0.1:5560"
	topic := "test"
	chunkSize := 1024

	pub, err := go_pubsub_zmq.NewPublisher(endpoint, topic, chunkSize)
	if err != nil {
		t.Fatalf("Error al crear Publisher: %v", err)
	}
	defer pub.Close()

	sub, err := go_pubsub_zmq.NewSubscriber(endpoint, topic)
	if err != nil {
		t.Fatalf("Error al crear Subscriber: %v", err)
	}
	defer sub.Close()

	totalMensajes := 5
	msgChan := make(chan string, totalMensajes)

	// Asegurar que el subscriber tiene tiempo para conectarse
	time.Sleep(500 * time.Millisecond)

	// Goroutine que recibe los mensajes.
	go func() {
		for i := 0; i < totalMensajes; i++ {
			var msgReceived *go_pubsub_zmq.MessageReceived
			var err error

			// Reintentos de recepción para cada mensaje
			for retry := 0; retry < 5; retry++ {
				msgReceived, _, err = sub.ReceiveMessage()
				if err == nil {
					break
				}
				time.Sleep(100 * time.Millisecond)
			}

			if err != nil {
				t.Logf("Error al recibir mensaje después de varios intentos: %v", err)
				msgChan <- "" // Enviar vacío para detectar fallo
				continue
			}

			if m, ok := msgReceived.Data["msg"].(string); ok {
				msgChan <- m
			} else {
				msgChan <- ""
			}
		}
	}()

	// Publica una serie de mensajes.
	for i := 0; i < totalMensajes; i++ {
		msg := "Mensaje " + strconv.Itoa(i) + " at " + time.Now().Format("15:04:05.000")

		messageBytes, err := pub.BuildMessage(nil, map[string]interface{}{"msg": msg}, "")
		if err != nil {
			t.Fatalf("Error al construir el mensaje %d: %v", i, err)
		}
		if err := pub.PublishMessage(messageBytes); err != nil {
			t.Fatalf("Error al publicar mensaje %d: %v", i, err)
		}
		time.Sleep(100 * time.Millisecond)
	}

	// Verifica que se hayan recibido todos los mensajes.
	for i := 0; i < totalMensajes; i++ {
		select {
		case received := <-msgChan:
			if received == "" {
				t.Errorf("El mensaje %d recibido está vacío", i)
			}
		case <-time.After(3 * time.Second): // Aumentar timeout para evitar fallos intermitentes
			t.Fatalf("No se recibió el mensaje %d en el tiempo esperado", i)
		}
	}
}

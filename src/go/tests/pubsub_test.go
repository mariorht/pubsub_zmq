package go_pubsub_zmq_test

import (
	"strconv"
	"testing"
	"time"

	"go_pubsub_zmq"
)

func TestPubSubMensajeUnico(t *testing.T) {
	endpoint := "tcp://127.0.0.1:5559"
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
		messageBytes, err := pub.BuildMessage(nil, map[string]interface{}{"msg": testMsg})
		if err != nil {
			t.Errorf("Error al construir el mensaje: %v", err)
			return
		}
		if err := pub.PublishMessage(messageBytes); err != nil {
			t.Errorf("Error al publicar mensaje: %v", err)
		}
	}()

	// Reintentos para asegurar la recepción del mensaje
	var receivedMsg *go_pubsub_zmq.MessageReceived
	for i := 0; i < 5; i++ {
		receivedMsg, err = sub.ReceiveMessage()
		if err == nil {
			break
		}
		time.Sleep(100 * time.Millisecond) // Espera antes del siguiente intento
	}

	if err != nil {
		t.Fatalf("Error al recibir mensaje después de varios intentos: %v", err)
	}

	// Validar que el mensaje recibido coincida con el enviado
	if received, ok := receivedMsg.Data["msg"].(string); !ok || received != testMsg {
		t.Errorf("El mensaje recibido (%v) no coincide con el esperado (%v)", receivedMsg.Data["msg"], testMsg)
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
				msgReceived, err = sub.ReceiveMessage()
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

		messageBytes, err := pub.BuildMessage(nil, map[string]interface{}{"msg": msg})
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

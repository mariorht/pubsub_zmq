package go_pubsub_zmq_test

import (
    "testing"
    "time"

    "go_pubsub_zmq"
)

func TestPubSubMensajeUnico(t *testing.T) {
    endpoint := "tcp://127.0.0.1:5559"

    // Crea el Publisher usando el package importado.
    pub, err := go_pubsub_zmq.NewPublisher(endpoint)
    if err != nil {
        t.Fatalf("Error al crear Publisher: %v", err)
    }
    defer pub.Close()

    // Crea el Subscriber.
    sub, err := go_pubsub_zmq.NewSubscriber(endpoint, "")
    if err != nil {
        t.Fatalf("Error al crear Subscriber: %v", err)
    }
    defer sub.Close()

    testMsg := "Mensaje de prueba"

    go func() {
        time.Sleep(100 * time.Millisecond) // Asegura que el subscriber esté listo.
        if err := pub.Publish(testMsg); err != nil {
            t.Errorf("Error al publicar mensaje: %v", err)
        }
    }()

    received, err := sub.Receive()
    if err != nil {
        t.Fatalf("Error al recibir mensaje: %v", err)
    }

    if received != testMsg {
        t.Errorf("El mensaje recibido (%q) no coincide con el esperado (%q)", received, testMsg)
    }
}

func TestPubSubMultiplesMensajes(t *testing.T) {
    endpoint := "tcp://127.0.0.1:5560"

    pub, err := go_pubsub_zmq.NewPublisher(endpoint)
    if err != nil {
        t.Fatalf("Error al crear Publisher: %v", err)
    }
    defer pub.Close()

    sub, err := go_pubsub_zmq.NewSubscriber(endpoint, "")
    if err != nil {
        t.Fatalf("Error al crear Subscriber: %v", err)
    }
    defer sub.Close()

    totalMensajes := 5
    msgChan := make(chan string, totalMensajes)

    go func() {
        for i := 0; i < totalMensajes; i++ {
            msg, err := sub.Receive()
            if err != nil {
                t.Logf("Error al recibir mensaje: %v", err)
                return
            }
            msgChan <- msg
        }
    }()

    for i := 0; i < totalMensajes; i++ {
        msg := "Mensaje " + time.Now().Format("15:04:05.000")
        if err := pub.Publish(msg); err != nil {
            t.Fatalf("Error al publicar mensaje %d: %v", i, err)
        }
        time.Sleep(100 * time.Millisecond)
    }

    for i := 0; i < totalMensajes; i++ {
        select {
        case received := <-msgChan:
            if received == "" {
                t.Errorf("El mensaje %d recibido está vacío", i)
            }
        case <-time.After(2 * time.Second):
            t.Fatalf("No se recibió el mensaje %d en el tiempo esperado", i)
        }
    }
}


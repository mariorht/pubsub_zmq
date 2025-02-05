package go_pubsub_zmq

import (
	"time"

	"github.com/pebbe/zmq4"
)

// Subscriber encapsula un socket de tipo SUB.
type Subscriber struct {
	socket *zmq4.Socket
}

// NewSubscriber crea un nuevo Subscriber, lo conecta al endpoint y establece el filtro.
// Utiliza una cadena vacía en filter para recibir todos los mensajes.
func NewSubscriber(endpoint string, filter string) (*Subscriber, error) {
	socket, err := zmq4.NewSocket(zmq4.SUB)
	if err != nil {
		return nil, err
	}
	if err = socket.SetSubscribe(filter); err != nil {
		socket.Close()
		return nil, err
	}
	if err = socket.Connect(endpoint); err != nil {
		socket.Close()
		return nil, err
	}
	// Espera un momento para asegurar la conexión.
	time.Sleep(500 * time.Millisecond)
	return &Subscriber{socket: socket}, nil
}

// Receive bloquea hasta recibir un mensaje y lo retorna.
func (s *Subscriber) Receive() (string, error) {
	return s.socket.Recv(0)
}

// Close cierra el socket del Subscriber.
func (s *Subscriber) Close() {
	s.socket.Close()
}

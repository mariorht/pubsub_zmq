package go_pubsub_zmq

import (
	"time"

	"github.com/pebbe/zmq4"
)

// Publisher encapsula un socket de tipo PUB.
type Publisher struct {
	socket *zmq4.Socket
}

// NewPublisher crea un nuevo Publisher y lo vincula al endpoint especificado.
func NewPublisher(endpoint string) (*Publisher, error) {
	socket, err := zmq4.NewSocket(zmq4.PUB)
	if err != nil {
		return nil, err
	}
	if err = socket.Bind(endpoint); err != nil {
		socket.Close()
		return nil, err
	}
	// Espera un momento para que los subscribers se conecten.
	time.Sleep(500 * time.Millisecond)
	return &Publisher{socket: socket}, nil
}

// Publish envía un mensaje utilizando el socket PUB.
func (p *Publisher) Publish(msg string) error {
	_, err := p.socket.Send(msg, 0)
	return err
}

// Close cierra el socket del Publisher.
func (p *Publisher) Close() {
	p.socket.Close()
}

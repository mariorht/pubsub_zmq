package go_pubsub_zmq

import (
	"bytes"
	"encoding/json"
	"log"
	"strconv"

	"github.com/pebbe/zmq4"
)

// Subscriber encapsula un socket SUB para recibir mensajes.
type Subscriber struct {
	socket  *zmq4.Socket
	context *zmq4.Context
	Topic   []byte
}

// NewSubscriber crea un nuevo Subscriber, conecta al address y se suscribe al topic dado.
func NewSubscriber(address string, topic string) (*Subscriber, error) {
	context, err := zmq4.NewContext()
	if err != nil {
		return nil, err
	}
	socket, err := context.NewSocket(zmq4.SUB)
	if err != nil {
		return nil, err
	}
	// Se suscribe al topic indicado (o se puede usar cadena vacía para recibir todo).
	err = socket.SetSubscribe(topic)
	if err != nil {
		return nil, err
	}
	if err = socket.Connect(address); err != nil {
		return nil, err
	}
	return &Subscriber{
		socket:  socket,
		context: context,
		Topic:   []byte(topic),
	}, nil
}

// MessageReceived es la estructura que se obtiene al decodificar el JSON recibido.
type MessageReceived struct {
	Type   string                 `json:"type"`
	Count  int                    `json:"count"`
	Images []ImageMetadataWrapper `json:"images"`
	Data   map[string]interface{} `json:"data"`
}

// ReceiveMessage reensambla los fragmentos recibidos y decodifica el mensaje JSON.
// Si se recibe datos de imagen, se emite un warning y se ignoran.
func (s *Subscriber) ReceiveMessage() (*MessageReceived, error) {
	var assembled bytes.Buffer
	var totalFragments int

	// Usamos RecvMessageBytes para obtener las partes sin perder datos binarios.
	for {
		parts, err := s.socket.RecvMessageBytes(0)
		if err != nil {
			return nil, err
		}
		
		// Se esperan 4 partes: [topic, fragment index, total fragments, chunk]
		if len(parts) < 4 {
			log.Println("Received incomplete message parts; skipping")
			continue
		}
		fragIndex, err := strconv.Atoi(string(parts[1]))
		if err != nil {
			return nil, err
		}
		tot, err := strconv.Atoi(string(parts[2]))
		if err != nil {
			return nil, err
		}
		totalFragments = tot
		assembled.Write(parts[3])
		if fragIndex == totalFragments-1 {
			break
		}
	}

	fullMessage := assembled.Bytes()
	// Se espera que el mensaje sea: JSON [ + null byte + image data ]
	parts := bytes.SplitN(fullMessage, []byte{0}, 2)
	var msgJSON []byte
	var imageData []byte
	if len(parts) >= 1 {
		msgJSON = parts[0]
	}
	if len(parts) == 2 {
		imageData = parts[1]
	}
	var message MessageReceived
	if err := json.Unmarshal(msgJSON, &message); err != nil {
		return nil, err
	}
	if len(imageData) > 0 {
		log.Println("WARNING: Received image data but image processing is not implemented; ignoring image data.")
	}
	return &message, nil
}

// Close cierra el socket y termina el contexto ZeroMQ.
func (s *Subscriber) Close() {
	s.socket.Close()
	s.context.Term()
}

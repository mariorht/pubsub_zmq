package go_pubsub_zmq

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/pebbe/zmq4"
)

// Frame representa una “imagen” de forma abstracta (sin utilizar OpenCV).
// Los campos son ejemplos; en una implementación completa se extraerían propiedades reales.
type Frame struct {
	Width    int
	Height   int
	Channels int
	Dtype    string
	// Se omite el contenido real de la imagen.
}

// Publisher encapsula un socket PUB y configura el protocolo.
type Publisher struct {
	socket         *zmq4.Socket
	context        *zmq4.Context
	Topic          []byte // Tema de publicación, por ejemplo "camera/image_raw"
	ChunkSize      int    // Tamaño de fragmento en bytes
	TotalBytesSent int    // Para seguimiento (total de bytes enviados)
}

// NewPublisher crea un nuevo Publisher en la dirección y con el tema indicados.
func NewPublisher(address string, topic string, chunkSize int) (*Publisher, error) {
	context, err := zmq4.NewContext()
	if err != nil {
		return nil, err
	}
	socket, err := context.NewSocket(zmq4.PUB)
	if err != nil {
		return nil, err
	}
	if err = socket.Bind(address); err != nil {
		return nil, err
	}
	return &Publisher{
		socket:    socket,
		context:   context,
		Topic:     []byte(topic),
		ChunkSize: chunkSize,
	}, nil
}

// ImageMetadata contiene información de una “imagen”.
type ImageMetadata struct {
	Format   string `json:"format"`   // Ej.: "raw"
	Width    int    `json:"width"`    // Ancho de la imagen
	Height   int    `json:"height"`   // Alto de la imagen
	Channels int    `json:"channels"` // Número de canales
	Dtype    string `json:"dtype"`    // Tipo de dato (ej.: "uint8")
	Size     int    `json:"size"`     // Tamaño en bytes (si se procesasen los datos)
}

// ImageMetadataWrapper envuelve la metadata de una imagen.
type ImageMetadataWrapper struct {
	Metadata ImageMetadata `json:"metadata"`
}

// Message representa la estructura del mensaje (JSON) que se envía junto a datos binarios.
type Message struct {
	Type   string                 `json:"type"`   // Ej.: "images"
	Count  int                    `json:"count"`  // Número de “imágenes”
	Images []ImageMetadataWrapper `json:"images"` // Lista de metadatos para cada imagen
	Data   map[string]interface{} `json:"data"`   // Datos adicionales
}

// BuildMessage construye el mensaje (JSON) a partir de los frames y datos adicionales.
// Si se proveen frames, se emite un warning y se omite la parte binaria de imagen (queda sin datos).
func (p *Publisher) BuildMessage(frames []Frame, data map[string]interface{}) ([]byte, error) {
	var imagesMetadata []ImageMetadataWrapper

	if len(frames) > 0 {
		log.Println("WARNING: Image processing is not implemented; ignoring image data.")
		for _, frame := range frames {
			meta := ImageMetadata{
				Format:   "raw",
				Width:    frame.Width,
				Height:   frame.Height,
				Channels: frame.Channels,
				Dtype:    frame.Dtype,
				Size:     0, // No se procesan los datos de imagen
			}
			imagesMetadata = append(imagesMetadata, ImageMetadataWrapper{Metadata: meta})
		}
	}

	msg := Message{
		Type:   "images",
		Count:  len(frames),
		Images: imagesMetadata,
		Data:   data,
	}
	// Se genera el JSON. Dado que no se concatenan datos de imagen, no se agrega el separador nulo.
	return json.Marshal(msg)
}

// PublishMessage envía el mensaje binario fragmentado en partes de tamaño ChunkSize.
// Cada fragmento se envía como un mensaje multipart con las partes:
// [topic, fragment index, total fragments, chunk]
func (p *Publisher) PublishMessage(messageBytes []byte) error {
	totalLen := len(messageBytes)
	numChunks := totalLen / p.ChunkSize
	if totalLen%p.ChunkSize != 0 {
		numChunks++
	}
	p.TotalBytesSent = 0

	for i := 0; i < numChunks; i++ {
		start := i * p.ChunkSize
		end := start + p.ChunkSize
		if end > totalLen {
			end = totalLen
		}
		chunk := messageBytes[start:end]
		_, err := p.socket.SendMessage(p.Topic, []byte(strconv.Itoa(i)), []byte(strconv.Itoa(numChunks)), chunk)
		if err != nil {
			return err
		}
		p.TotalBytesSent += len(chunk)
	}
	return nil
}

// Close cierra el socket y termina el contexto ZeroMQ.
func (p *Publisher) Close() {
	p.socket.Close()
	p.context.Term()
}

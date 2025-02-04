package publisher

import (
    "encoding/json"
    "fmt"
    "github.com/pebbe/zmq4"
    "github.com/disintegration/imaging"
    "image"
    "image/jpeg"
    "os"
    "time"
)

type Publisher struct {
    socket *zmq4.Socket
    topic  string
}

type ImageMetadata struct {
    Format string `json:"format"`
    Width  int    `json:"width"`
    Height int    `json:"height"`
    Size   int    `json:"size"`
}

type Message struct {
    Type   string          `json:"type"`
    Count  int             `json:"count"`
    Images []ImageMetadata `json:"images"`
    Data   map[string]interface{} `json:"data"`
}

func NewPublisher(address, topic string) (*Publisher, error) {
    context, err := zmq4.NewContext()
    if err != nil {
        return nil, err
    }
    socket, err := context.NewSocket(zmq4.PUB)
    if err != nil {
        return nil, err
    }
    err = socket.Bind(address)
    if err != nil {
        return nil, err
    }
    return &Publisher{socket: socket, topic: topic}, nil
}

func (p *Publisher) BuildMessage(frames []image.Image, data map[string]interface{}) ([]byte, error) {
    var imagesMetadata []ImageMetadata
    var imagesData []byte

    for _, frame := range frames {
        buffer := new(bytes.Buffer)
        err := jpeg.Encode(buffer, frame, nil)
        if err != nil {
            return nil, err
        }
        imageBytes := buffer.Bytes()
        imageMetadata := ImageMetadata{
            Format: "jpg",
            Width:  frame.Bounds().Dx(),
            Height: frame.Bounds().Dy(),
            Size:   len(imageBytes),
        }
        imagesMetadata = append(imagesMetadata, imageMetadata)
        imagesData = append(imagesData, imageBytes...)
    }

    message := Message{
        Type:   "images",
        Count:  len(frames),
        Images: imagesMetadata,
        Data:   data,
    }
    messageJSON, err := json.Marshal(message)
    if err != nil {
        return nil, err
    }
    return append(messageJSON, append([]byte{0}, imagesData...)...), nil
}

func (p *Publisher) PublishMessage(messageBytes []byte) error {
    _, err := p.socket.SendMessage(p.topic, messageBytes)
    return err
}

func (p *Publisher) Close() {
    p.socket.Close()
}

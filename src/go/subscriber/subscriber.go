package subscriber

import (
    "encoding/json"
    "fmt"
    "github.com/pebbe/zmq4"
    "image"
    "image/jpeg"
    "bytes"
)

type Subscriber struct {
    socket *zmq4.Socket
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

func NewSubscriber(address, topic string) (*Subscriber, error) {
    context, err := zmq4.NewContext()
    if (err != nil) {
        return nil, err
    }
    socket, err := context.NewSocket(zmq4.SUB)
    if (err != nil) {
        return nil, err
    }
    err = socket.Connect(address)
    if (err != nil) {
        return nil, err
    }
    err = socket.SetSubscribe(topic)
    if (err != nil) {
        return nil, err
    }
    return &Subscriber{socket: socket}, nil
}

func (s *Subscriber) ReceiveMessage() ([]image.Image, map[string]interface{}, error) {
    msg, err := s.socket.RecvMessageBytes(0)
    if (err != nil) {
        return nil, nil, err
    }

    if (len(msg) < 2) {
        return nil, nil, fmt.Errorf("invalid message format")
    }

    jsonPart := msg[1]
    var message Message
    err = json.Unmarshal(jsonPart, &message)
    if (err != nil) {
        return nil, nil, err
    }

    imagesBytes := msg[2]
    var images []image.Image
    offset := 0
    for _, imageInfo := range message.Images {
        size := imageInfo.Size
        imageBytes := imagesBytes[offset : offset+size]
        img, err := jpeg.Decode(bytes.NewReader(imageBytes))
        if (err != nil) {
            return nil, nil, err
        }
        images = append(images, img)
        offset += size
    }

    return images, message.Data, nil
}

func (s *Subscriber) Close() {
    s.socket.Close()
}

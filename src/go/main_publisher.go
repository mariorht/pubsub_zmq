package main

import (
    "fmt"
    "image"
    "time"
    "github.com/disintegration/imaging"
    // Import the publisher package using a relative path
    "./publisher/publisher"
)

func main() {
    pub, err := publisher.NewPublisher("tcp://*:5555", "camera/image_raw")
    if err != nil {
        fmt.Println("Error creating publisher:", err)
        return
    }
    defer pub.Close()

    img, err := imaging.Open("../../assets/pong.png")
    if err != nil {
        fmt.Println("Error loading image:", err)
        return
    }

    angle := 0.0
    for {
        frames := []image.Image{
            imaging.Rotate(img, angle, image.Black),
            imaging.Rotate(img, angle+10, image.Black),
            imaging.Rotate(img, angle+20, image.Black),
        }
        data := map[string]interface{}{
            "key": "probando",
        }

        messageBytes, err := pub.BuildMessage(frames, data)
        if err != nil {
            fmt.Println("Error building message:", err)
            return
        }

        err = pub.PublishMessage(messageBytes)
        if err != nil {
            fmt.Println("Error publishing message:", err)
            return
        }

        angle += 10
        time.Sleep(100 * time.Millisecond)
    }
}

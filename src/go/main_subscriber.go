package main

import (
    "fmt"
    "time"
    // Import the subscriber package using a relative path
    "./subscriber/subscriber"
)

func main() {
    sub, err := subscriber.NewSubscriber("tcp://127.0.0.1:5555", "camera/image_raw")
    if err != nil {
        fmt.Println("Error creating subscriber:", err)
        return
    }
    defer sub.Close()

    for {
        images, data, err := sub.ReceiveMessage()
        if err != nil {
            fmt.Println("Error receiving message:", err)
            return
        }

        fmt.Println("Received data:", data)
        for i, img := range images {
            fmt.Printf("Received image %d: %v\n", i+1, img.Bounds())
        }

        time.Sleep(100 * time.Millisecond)
    }
}

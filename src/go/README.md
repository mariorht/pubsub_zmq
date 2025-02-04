# Go Subscriber

This is a Go-based subscriber for receiving messages from a ZeroMQ publisher.

## Prerequisites

- Go 1.13 or higher
- ZeroMQ library

## Installation

1. Install Go from the [official website](https://golang.org/dl/).
2. Install the ZeroMQ library. Follow the instructions on the [ZeroMQ website](https://zeromq.org/download/).
3. Initialize a new Go module and get the required package:
    ```sh
    go mod init zmqpubsub
    go get gopkg.in/zeromq/goczmq.v4
    ```

## Running the Subscriber

1. Navigate to the `go` directory:
    ```sh
    cd /opt/projects/pubsub_zmq/src/go
    ```
2. Run the Go subscriber:
    ```sh
    go run main_subscriber.go
    ```

The subscriber will connect to the publisher at `tcp://127.0.0.1:5555` and subscribe to the `camera/image_raw` topic. It will print the received data and image bounds to the console.

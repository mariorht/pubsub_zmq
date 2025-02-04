# PubSub ZMQ

PubSub ZMQ is a Python library that facilitates the publication and subscription of messages containing multiple images and additional data using ZeroMQ. This library is designed to handle large binary data efficiently by fragmenting messages into smaller chunks.

## Features

- Publish and subscribe to messages containing multiple images and additional data.
- Efficient handling of large binary data using message fragmentation.
- JSON-based metadata for images and additional data.
- Real-time image display using OpenCV.

## Installation

To install the required dependencies, run:

```bash
pip install -r requirements.txt
```

## Usage

### Publisher

The publisher reads an image, rotates it to create multiple frames, and publishes the frames along with additional data.

```python
from publisher import Publisher
import time
import cv2

if __name__ == "__main__":
    pub = Publisher()
    angle = 0
    try:
        frame = cv2.imread("pong.png")
        assert frame is not None, "❌ No se pudo cargar la imagen: pong.png"
            
        while True:
            frames = []
            for _ in range(3):
                angle = (angle + 10) % 360
                center = (frame.shape[1] // 2, frame.shape[0] // 2)
                matrix = cv2.getRotationMatrix2D(center, angle, 1.0)
                rotated_frame = cv2.warpAffine(frame, matrix, (frame.shape[1], frame.shape[0]))
                frames.append(rotated_frame)

            data = {"key": "probando"}
            message_bytes = pub.build_message(frames, data)
            pub.publish_message(message_bytes)
            time.sleep(.1)
            
    except KeyboardInterrupt:
        print("\n🛑 Publicador detenido.")
        pub.close()
```

### Subscriber

The subscriber receives the fragmented messages, reconstructs the images, and displays them using OpenCV.

```python
from subscriber import Subscriber
import time
import cv2

if __name__ == "__main__":
    sub = Subscriber()
    try:
        while True:
            start_time = time.time()
            sub.total_bytes_received = 0
            images, data = sub.receive_message()
            if images is not None:
                elapsed_time = time.time() - start_time
                mbps = (sub.total_bytes_received * 8) / (elapsed_time * 1_000_000)
                print(f"🚀 Velocidad de recepción media: {mbps:.2f} Mbps")
                
                for i, img in enumerate(images):
                    cv2.imshow(f"Imagen Recibida {i+1}", img)
                cv2.waitKey(1)
                
                print(f"📄 Datos recibidos: {data}")
    except KeyboardInterrupt:
        print("\n🛑 Suscriptor detenido.")
        sub.close()
        cv2.destroyAllWindows()
```

## License

This project is licensed under the MIT License.
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
import sys
import os
import time
import cv2
import json

sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), "../..")))

from publisher import Publisher


if __name__ == "__main__":
    pub = Publisher(address="tcp://*:5555", topic="test") #Nombre de red del docker
    angle = 0

    try:
        frame = cv2.imread("/assets/pong.png")
        assert frame is not None, f"❌ No se pudo cargar la imagen: /assets/pong.png"

        print("✅ Publicador creado. Esperando que se conecten los suscriptores...")
        time.sleep(5)

        for i in range(1):
            frames = []

            for _ in range(3):
                angle = (angle + 10) % 360
                center = (frame.shape[1] // 2, frame.shape[0] // 2)
                matrix = cv2.getRotationMatrix2D(center, angle, 1.0)
                rotated_frame = cv2.warpAffine(frame, matrix, (frame.shape[1], frame.shape[0]))
                frames.append(rotated_frame)

            data = {"key": "probando", "index": i}

            message_bytes = pub.build_message(frames, data, "png")
            pub.publish_message(message_bytes)

            # 🔄 Guardar el mensaje completo
            result_publisher = {
                "type": "images",
                "count": len(frames),
                "images": [{
                    "metadata": {
                        "width": f.shape[1],
                        "height": f.shape[0],
                        "channels": f.shape[2] if len(f.shape) > 2 else 1,
                        "dtype": str(f.dtype),
                        "size": f.nbytes,
                    }
                } for f in frames],
                "data": data
            }

            with open("/shared/result_publisher.json", "w") as f:
                json.dump(result_publisher, f, indent=4)

            print(f"📤 Mensaje {i} enviado con {len(frames)} imágenes y data {data}")

            time.sleep(1)

    except KeyboardInterrupt:
        print("\n🛑 Publicador detenido.")
    finally:
        pub.close()

```

### Subscriber

The subscriber receives the fragmented messages, reconstructs the images, and displays them using OpenCV.

```python
import sys
import os
import time
import json
import cv2

sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), "../..")))

from subscriber import Subscriber


if __name__ == "__main__":
    endpoint = os.getenv("PUBSUB_ENDPOINT", "tcp://localhost:5555")

    sub = Subscriber(address=endpoint, topic="test")

    try:
        print("✅ Subscriber (Python) esperando mensajes...")

        last_data = None

        for i in range(1):
            frames, data = sub.receive_message()
            print(f"📥 Mensaje {i} recibido: {data}")

            if frames:
                print(f"⚠️ Recibidas {len(frames)} imágenes.")
                cv2.imwrite(f"/shared/result_{i}.png", frames[0])

            last_data = data

        if last_data is not None:
            result = {
                "type": "images",
                "count": len(frames),
                "images": [{
                    "metadata": {
                        "width": f.shape[1],
                        "height": f.shape[0],
                        "channels": f.shape[2] if len(f.shape) > 2 else 1,
                        "dtype": str(f.dtype),
                        "size": f.nbytes,
                    }
                } for f in frames],
                "data": data
            }

            with open("/shared/result.json", "w") as f:
                json.dump(result, f, indent=4)

            print("✅ Resultado guardado en /shared/result.json")

    except Exception as e:
        print(f"❌ Error en el subscriber: {e}")

    finally:
        sub.close()


## Running Tests

To run the tests using pytest, navigate to the root directory of the project and execute:

```bash
pytest
```

This will automatically discover and run all the tests in the `tests` directory.

## License

This project is licensed under the MIT License.
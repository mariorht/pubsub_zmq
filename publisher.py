import zmq
import cv2
import numpy as np
import time
import json  # Add import for JSON handling

class Publisher:
    def __init__(self, address="tcp://*:5555", topic="camera/image_raw", chunk_size=100000):
        self.context = zmq.Context()
        self.socket = self.context.socket(zmq.PUB)
        self.socket.bind(address)
        
        self.topic = topic.encode()
        self.chunk_size = chunk_size  # Tamaño de fragmento en bytes

    def build_message(self, frames, data):
        """ Construye un mensaje con un número indeterminado de imágenes y datos adicionales. """
        images_metadata = []
        images_data = []
        for frame in frames:
            _, buffer = cv2.imencode(".jpg", frame)
            image_bytes = buffer.tobytes()
            image_metadata = {
                "format": "jpg",
                "width": frame.shape[1],
                "height": frame.shape[0],
                "size": len(image_bytes)
            }
            images_metadata.append({"metadata": image_metadata})
            images_data.append(image_bytes)

        message = {
            "type": "images",
            "count": len(frames),
            "images": images_metadata,
            "data": data
        }
        message_json = json.dumps(message).encode('utf-8')
        return message_json + b'\x00' + b''.join(images_data)

    def publish_message(self, message_bytes):
        """ Publica un mensaje binario fragmentado. """
        num_chunks = len(message_bytes) // self.chunk_size + 1

        self.total_bytes_sent = 0  # Initialize total bytes sent

        for i in range(num_chunks):
            chunk = message_bytes[i * self.chunk_size: (i + 1) * self.chunk_size]
            self.socket.send_multipart([self.topic, str(i).encode(), str(num_chunks).encode(), chunk])
            self.total_bytes_sent += len(chunk)

        print(f"📤 Mensaje binario publicado en {num_chunks} fragmentos.")

    def close(self):
        """ Cierra la conexión ZeroMQ. """
        self.socket.close()
        self.context.term()

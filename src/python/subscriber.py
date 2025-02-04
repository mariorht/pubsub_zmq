import zmq
import numpy as np
import cv2
import time
import json  # Add import for JSON handling

class Subscriber:
    def __init__(self, address="tcp://127.0.0.1:5555", topic="camera/image_raw"):
        self.context = zmq.Context()
        self.socket = self.context.socket(zmq.SUB)
        self.socket.connect(address)
        self.socket.setsockopt_string(zmq.SUBSCRIBE, topic)
        
        self.image_chunks = {}

    def receive_message(self, timeout=None):
        """ Recibe mensajes binarios que contienen múltiples imágenes y datos adicionales. """
        self.total_bytes_received = 0  # Initialize total bytes received
        start_time = time.time()  # Start time for reception

        if timeout:
            self.socket.setsockopt(zmq.RCVTIMEO, timeout)

        while True:
            try:
                print("Esperando recibir un mensaje...")
                topic, index, total, chunk = self.socket.recv_multipart()
                print(f"Mensaje recibido: index={index}, total={total}")
            except zmq.Again:
                print("❌ Timeout: No se recibieron mensajes en el tiempo esperado.")
                return None, None

            index = int(index.decode())
            total = int(total.decode())

            self.image_chunks[index] = chunk
            self.total_bytes_received += len(chunk)

            if len(self.image_chunks) == total:  # Si todos los fragmentos llegaron
                ordered_chunks = [self.image_chunks[i] for i in sorted(self.image_chunks.keys())]
                message_bytes = b"".join(ordered_chunks)

                # Split the JSON part and the image part
                if b'\x00' in message_bytes:
                    json_part, images_bytes = message_bytes.split(b'\x00', 1)
                else:
                    json_part = message_bytes
                    images_bytes = b''

                try:
                    message = json.loads(json_part.decode('utf-8'))
                    images = []
                    offset = 0
                    for image_info in message.get("images", []):
                        size = image_info["metadata"]["size"]
                        image_bytes = images_bytes[offset:offset + size]
                        image_array = np.frombuffer(image_bytes, dtype=np.uint8)
                        frame = cv2.imdecode(image_array, cv2.IMREAD_COLOR)
                        images.append(frame)
                        offset += size
                    print(f"✅ Imágenes embebidas en binario reconstruidas correctamente.")
                    self.image_chunks.clear()  # Limpiar para el siguiente mensaje
                    return images, message["data"]
                except json.JSONDecodeError:
                    print("❌ Error decodificando JSON.")
                    self.image_chunks.clear()
                    return None, None

    def close(self):
        """ Cierra la conexión ZeroMQ. """
        self.socket.close()
        self.context.term()

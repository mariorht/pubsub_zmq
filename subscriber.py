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

    def receive_message(self):
        """ Recibe mensajes y determina si es una imagen o un JSON. """
        self.total_bytes_received = 0  # Initialize total bytes received
        start_time = time.time()  # Start time for reception

        while True:
            topic, index, total, chunk = self.socket.recv_multipart()
            index = int(index.decode())
            total = int(total.decode())

            self.image_chunks[index] = chunk
            self.total_bytes_received += len(chunk)

            if len(self.image_chunks) == total:  # Si todos los fragmentos llegaron
                ordered_chunks = [self.image_chunks[i] for i in sorted(self.image_chunks.keys())]
                message_bytes = b"".join(ordered_chunks)

                # Check if the message is JSON
                if message_bytes.startswith(b'{'):
                    try:
                        message = json.loads(message_bytes.decode('utf-8'))
                        print(f"✅ JSON recibido: {message}")
                        self.image_chunks.clear()  # Limpiar para el siguiente mensaje
                        return message
                    except json.JSONDecodeError:
                        print("❌ Error decodificando JSON.")
                        self.image_chunks.clear()
                        return None
                else:
                    # Reconstruir la imagen
                    image_array = np.frombuffer(message_bytes, dtype=np.uint8)
                    frame = cv2.imdecode(image_array, cv2.IMREAD_COLOR)
                    print(f"✅ Imagen reconstruida correctamente. Tamaño: {frame.shape[1]}x{frame.shape[0]}")
                    self.image_chunks.clear()  # Limpiar para la siguiente imagen
                    return frame

    def close(self):
        """ Cierra la conexión ZeroMQ. """
        self.socket.close()
        self.context.term()

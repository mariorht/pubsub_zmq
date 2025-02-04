import zmq
import numpy as np
import cv2
import time

class Subscriber:
    def __init__(self, address="tcp://127.0.0.1:5555", topic="camera/image_raw"):
        self.context = zmq.Context()
        self.socket = self.context.socket(zmq.SUB)
        self.socket.connect(address)
        self.socket.setsockopt_string(zmq.SUBSCRIBE, topic)
        
        self.image_chunks = {}

    def receive_image(self):
        """ Recibe fragmentos y reconstruye la imagen cuando estén todos. """
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
                image_bytes = b"".join(ordered_chunks)

                # Reconstruir la imagen
                image_array = np.frombuffer(image_bytes, dtype=np.uint8)
                frame = cv2.imdecode(image_array, cv2.IMREAD_COLOR)


                print(f"✅ Imagen reconstruida correctamente. Tamaño: {frame.shape[1]}x{frame.shape[0]}")
                self.image_chunks.clear()  # Limpiar para la siguiente imagen
                return frame

    def close(self):
        """ Cierra la conexión ZeroMQ. """
        self.socket.close()
        self.context.term()

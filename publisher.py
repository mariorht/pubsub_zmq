import zmq
import cv2
import numpy as np
import time

class Publisher:
    def __init__(self, address="tcp://*:5555", topic="camera/image_raw", chunk_size=100000):
        self.context = zmq.Context()
        self.socket = self.context.socket(zmq.PUB)
        self.socket.bind(address)
        
        self.topic = topic.encode()
        self.chunk_size = chunk_size  # Tamaño de fragmento en bytes

    def publish_image(self, image_path=None, use_camera=False):
        """ Captura una imagen de un archivo o una cámara y la publica fragmentada. """
        if use_camera:
            cap = cv2.VideoCapture(0)
            ret, frame = cap.read()
            cap.release()
            if not ret:
                print("❌ Error capturando imagen de la cámara.")
                return
        else:
            frame = cv2.imread(image_path)
            if frame is None:
                print(f"❌ No se pudo cargar la imagen: {image_path}")
                return

        _, buffer = cv2.imencode(".jpg", frame)
        image_bytes = buffer.tobytes()

        num_chunks = len(image_bytes) // self.chunk_size + 1

        for i in range(num_chunks):
            chunk = image_bytes[i * self.chunk_size: (i + 1) * self.chunk_size]
            self.socket.send_multipart([self.topic, str(i).encode(), str(num_chunks).encode(), chunk])
            time.sleep(0.001)  # Pequeño delay para evitar congestión

        print(f"📤 Imagen publicada en {num_chunks} fragmentos. Tamaño: {frame.shape[1]}x{frame.shape[0]}")

    def close(self):
        """ Cierra la conexión ZeroMQ. """
        self.socket.close()
        self.context.term()

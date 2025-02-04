import pytest
import zmq
import cv2
import time
import json
import numpy as np
import os

from publisher import Publisher
from subscriber import Subscriber

@pytest.fixture
def setup_pubsub():
    pub = Publisher(address="tcp://127.0.0.1:5555")
    sub = Subscriber(address="tcp://127.0.0.1:5555")
    yield pub, sub
    pub.close()
    sub.close()

def test_pubsub(setup_pubsub):
    pub, sub = setup_pubsub

    # Set environment variable for image path
    base_dir = os.path.dirname(os.path.abspath(__file__))
    image_path = os.path.join(base_dir, "../../../assets/pong.png")
    os.environ["PONG_IMAGE_PATH"] = image_path

    # Load image
    frame = cv2.imread(image_path)
    assert frame is not None, f"❌ No se pudo cargar la imagen: {image_path}"

    frames = [frame, frame, frame]  # Use the same frame 3 times for simplicity
    data = {"key": "probando"}

    # Build and publish message
    message_bytes = pub.build_message(frames, data)
    pub.publish_message(message_bytes)

    # Receive and verify message
    images, received_data = sub.receive_message(timeout=5000)
    assert images is not None, "❌ No se recibieron imágenes."
    assert received_data == data, "❌ Los datos recibidos no coinciden."

    for img in images:
        assert img.shape == frame.shape, "❌ La imagen recibida no tiene las mismas dimensiones que la original."

def test_pubsub_text_only(setup_pubsub):
    pub, sub = setup_pubsub

    print("Enviando mensaje sin imágenes...")
    # No images, only text data
    frames = []
    data = {"key": "probando", "message": "Este es un mensaje de prueba sin imágenes."}

    # Build and publish message
    message_bytes = pub.build_message(frames, data)
    pub.publish_message(message_bytes)
    print("Mensaje publicado.")

    # Receive and verify message
    images, received_data = sub.receive_message(timeout=5000)
    print("Mensaje recibido.")
    assert images is not None, "❌ No se recibieron imágenes."
    assert images == [], "❌ Se recibieron imágenes cuando no se esperaban."
    assert received_data == data, "❌ Los datos recibidos no coinciden."

if __name__ == "__main__":
    pytest.main()

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
    image_path = os.path.join(base_dir, "../../../../assets/pong.png")
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

    for img, original in zip(images, frames):
        assert np.array_equal(img, original), "❌ La imagen recibida no es idéntica a la original."

def test_pubsub_different_formats(setup_pubsub):
    pub, sub = setup_pubsub

    # Set environment variable for image path
    base_dir = os.path.dirname(os.path.abspath(__file__))
    image_path = os.path.join(base_dir, "../../../../assets/pong.png")

    # Load image
    frame = cv2.imread(image_path)
    assert frame is not None, f"❌ No se pudo cargar la imagen: {image_path}"

    # Convert image to different formats
    frame_float32 = frame.astype(np.float32)
    frame_uint16 = frame.astype(np.uint16)

    formats = [frame, frame_float32, frame_uint16]
    data = {"key": "probando"}

    for fmt in formats:
        frames = [fmt, fmt, fmt]  # Use the same frame 3 times for simplicity

        # Build and publish message
        message_bytes = pub.build_message(frames, data)
        pub.publish_message(message_bytes)

        # Receive and verify message
        images, received_data = sub.receive_message(timeout=5000)
        assert images is not None, "❌ No se recibieron imágenes."
        assert received_data == data, "❌ Los datos recibidos no coinciden."

        for img, original in zip(images, frames):
            assert np.array_equal(img, original), f"❌ La imagen recibida no es idéntica a la original para el formato {fmt.dtype}."



def test_pubsub_multiple_channels(setup_pubsub):
    pub, sub = setup_pubsub

    # Set environment variable for image path
    base_dir = os.path.dirname(os.path.abspath(__file__))
    image_path = os.path.join(base_dir, "../../../../assets/pong.png")

    # Load image
    frame = cv2.imread(image_path)
    assert frame is not None, f"❌ No se pudo cargar la imagen: {image_path}"

    # Convert image to different channel formats
    frame_gray = cv2.cvtColor(frame, cv2.COLOR_BGR2GRAY)
    frame_rgba = cv2.cvtColor(frame, cv2.COLOR_BGR2RGBA)

    formats = [frame_gray, frame_rgba]
    data = {"key": "probando"}

    for fmt in formats:
        frames = [fmt, fmt, fmt]  # Use the same frame 3 times for simplicity

        # Build and publish message
        message_bytes = pub.build_message(frames, data, "raw")
        pub.publish_message(message_bytes)

        # Receive and verify message
        images, received_data = sub.receive_message(timeout=5000)
        assert images is not None, "❌ No se recibieron imágenes."
        assert received_data == data, "❌ Los datos recibidos no coinciden."

        for img, original in zip(images, frames):
            if len(original.shape) == 2:  # Grayscale image
                original = original[:, :, np.newaxis]
            assert np.array_equal(img, original), f"❌ La imagen recibida no es idéntica a la original para el formato {fmt.shape}."



def test_pubsub_jpeg(setup_pubsub):
    pub, sub = setup_pubsub

    # Set environment variable for image path
    base_dir = os.path.dirname(os.path.abspath(__file__))
    image_path = os.path.join(base_dir, "../../../../assets/pong.png")

    # Load image
    frame = cv2.imread(image_path)
    assert frame is not None, f"❌ No se pudo cargar la imagen: {image_path}"

    frames = [frame, frame, frame]  # Use the same frame 3 times for simplicity
    data = {"key": "probando"}

    # Build and publish message in JPEG format
    message_bytes = pub.build_message(frames, data, format="jpeg")
    pub.publish_message(message_bytes)

    # Receive and verify message
    images, received_data = sub.receive_message(timeout=5000)
    assert images is not None, "❌ No se recibieron imágenes."
    assert received_data == data, "❌ Los datos recibidos no coinciden."

    for img in images:
        assert img is not None, "❌ La imagen decodificada es nula."
        assert img.shape == frame.shape, f"❌ La imagen recibida tiene dimensiones incorrectas: {img.shape} != {frame.shape}"


def test_pubsub_jpeg_quality(setup_pubsub):
    pub, sub = setup_pubsub

    base_dir = os.path.dirname(os.path.abspath(__file__))
    image_path = os.path.join(base_dir, "../../../../assets/pong.png")

    # Load image
    frame = cv2.imread(image_path)
    assert frame is not None, f"❌ No se pudo cargar la imagen: {image_path}"

    frames = [frame]  # Enviar solo una imagen para comparar con precisión
    data = {"key": "probando"}

    # Build and publish message in JPEG format
    message_bytes = pub.build_message(frames, data, format="jpeg")
    pub.publish_message(message_bytes)

    # Receive and verify message
    images, received_data = sub.receive_message(timeout=5000)
    assert images is not None, "❌ No se recibieron imágenes."
    assert received_data == data, "❌ Los datos recibidos no coinciden."

    received_image = images[0]
    assert received_image.shape == frame.shape, "❌ La imagen recibida tiene dimensiones incorrectas."

    # Calcular la diferencia absoluta entre la imagen original y la recibida
    difference = cv2.absdiff(frame, received_image)
    mean_diff = np.mean(difference)

    print(f"📊 Diferencia promedio entre imágenes: {mean_diff}")

    # La diferencia debería ser baja (JPEG introduce ligeras pérdidas)
    assert mean_diff < 5, f"❌ La pérdida de calidad de JPEG es demasiado alta: {mean_diff}"


def test_pubsub_png(setup_pubsub):
    pub, sub = setup_pubsub

    # Set environment variable for image path
    base_dir = os.path.dirname(os.path.abspath(__file__))
    image_path = os.path.join(base_dir, "../../../../assets/pong.png")

    # Load image
    frame = cv2.imread(image_path)
    assert frame is not None, f"❌ No se pudo cargar la imagen: {image_path}"

    frames = [frame, frame, frame]  # Usar 3 copias de la imagen
    data = {"key": "probando"}

    # Enviar imágenes en formato PNG
    message_bytes = pub.build_message(frames, data, format="png")
    pub.publish_message(message_bytes)

    # Recibir y verificar mensaje
    images, received_data = sub.receive_message(timeout=5000)
    assert images is not None, "❌ No se recibieron imágenes."
    assert received_data == data, "❌ Los datos recibidos no coinciden."

    for img, original in zip(images, frames):
        assert img.shape == original.shape, "❌ La imagen recibida tiene dimensiones incorrectas."

    print("✅ Test de imágenes PNG pasado correctamente.")


# def test_pubsub_text_only(setup_pubsub):
#     pub, sub = setup_pubsub

#     print("📨 Enviando mensaje sin imágenes...")
    
#     # Mensaje sin imágenes, solo con datos
#     frames = []  # Lista vacía de imágenes
#     data = {
#         "key": "probando",
#         "message": "Este es un mensaje de prueba sin imágenes."
#     }

#     # Construir y enviar el mensaje
#     message_bytes = pub.build_message(frames, data)
#     pub.publish_message(message_bytes)
#     print("✅ Mensaje publicado.")

#     # Recibir y verificar mensaje
#     images, received_data = sub.receive_message(timeout=5000)
#     print("📩 Mensaje recibido.")

#     # Verificar que no hay imágenes y que los datos coinciden
#     assert images is not None, "❌ No se recibieron imágenes (debería ser una lista vacía)."
#     assert images == [], "❌ Se recibieron imágenes cuando no se esperaban."
#     assert received_data == data, "❌ Los datos recibidos no coinciden."

#     print("✅ Test de mensaje sin imágenes completado exitosamente.")




if __name__ == "__main__":
    pytest.main()

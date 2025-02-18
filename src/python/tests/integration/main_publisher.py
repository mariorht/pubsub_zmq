import sys
import os
import time
import cv2

sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), "../..")))

from publisher import Publisher


if __name__ == "__main__":
    pub = Publisher(address="tcp://integration_python_pub:5559", topic="test") #Nombre de red del docker
    angle = 0

    try:
        frame = cv2.imread("/assets/pong.png")
        assert frame is not None, f"❌ No se pudo cargar la imagen: /assets/pong.png"

        time.sleep(2)

        print("✅ Publicador creado. Enviando mensajes...")

        for i in range(1):  # Envía 1 mensajes y termina
            frames = []

            for _ in range(3):  # Genera 3 imágenes rotadas
                angle = (angle + 10) % 360
                center = (frame.shape[1] // 2, frame.shape[0] // 2)
                matrix = cv2.getRotationMatrix2D(center, angle, 1.0)
                rotated_frame = cv2.warpAffine(frame, matrix, (frame.shape[1], frame.shape[0]))
                frames.append(rotated_frame)

            data = {"key": "probando", "index": i}

            message_bytes = pub.build_message(frames, data)
            pub.publish_message(message_bytes)

            print(f"📤 Mensaje {i} enviado con {len(frames)} imágenes y data {data}")

            time.sleep(1)

    except KeyboardInterrupt:
        print("\n🛑 Publicador detenido.")
    finally:
        pub.close()

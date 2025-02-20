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

            message_bytes = pub.build_message(frames, data)
            pub.publish_message(message_bytes)

            # Guardar el último mensaje enviado
            with open("/shared/result_publisher.json", "w") as f:
                json.dump(data, f)

            print(f"📤 Mensaje {i} enviado con {len(frames)} imágenes y data {data}")

            time.sleep(1)

    except KeyboardInterrupt:
        print("\n🛑 Publicador detenido.")
    finally:
        pub.close()

import sys
import os
import time
import json

sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), "../..")))

from subscriber import Subscriber


if __name__ == "__main__":
    sub = Subscriber(address="tcp://integration_go_pub:5559", topic="test")

    try:
        print("✅ Subscriber (Python) esperando mensajes...")

        last_data = None

        for i in range(1):
            images, data = sub.receive_message()
            print(f"📥 Mensaje {i} recibido: {data}")

            if images:
                print(f"⚠️ Recibidas {len(images)} imágenes (ignorado por ahora).")

            last_data = data

        if last_data is not None:
            with open("/shared/result.json", "w") as f:
                json.dump(last_data, f)
            print("✅ Resultado guardado en /shared/result.json")

    except Exception as e:
        print(f"❌ Error en el subscriber: {e}")

    finally:
        sub.close()

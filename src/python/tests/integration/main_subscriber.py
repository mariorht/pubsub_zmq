import sys
import os
import time
import json
import cv2

sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), "../..")))

from subscriber import Subscriber


if __name__ == "__main__":
    endpoint = os.getenv("PUBSUB_ENDPOINT", "tcp://localhost:5555")

    sub = Subscriber(address=endpoint, topic="test")

    try:
        print("✅ Subscriber (Python) esperando mensajes...")

        last_data = None

        for i in range(1):
            frames, data = sub.receive_message()
            print(f"📥 Mensaje {i} recibido: {data}")

            if frames:
                print(f"⚠️ Recibidas {len(frames)} imágenes.")
                cv2.imwrite(f"/shared/result_{i}.png", frames[0])

            last_data = data

        if last_data is not None:
            result = {
                "type": "images",
                "count": len(frames),
                "images": [{
                    "metadata": {
                        "width": f.shape[1],
                        "height": f.shape[0],
                        "channels": f.shape[2] if len(f.shape) > 2 else 1,
                        "dtype": str(f.dtype),
                        "size": f.nbytes,
                    }
                } for f in frames],
                "data": data
            }

            with open("/shared/result.json", "w") as f:
                json.dump(result, f, indent=4)

            print("✅ Resultado guardado en /shared/result.json")

    except Exception as e:
        print(f"❌ Error en el subscriber: {e}")

    finally:
        sub.close()

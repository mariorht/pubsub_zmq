from subscriber import Subscriber
import time
import cv2  # Add import for OpenCV
import numpy as np  # Add import for numpy

if __name__ == "__main__":
    sub = Subscriber( topic="camera/image_raw")
    try:
        while True:
            start_time = time.time()  # Reset start time for each message
            sub.total_bytes_received = 0  # Reset total bytes received for each message
            images, data = sub.receive_message()
            if images is not None:
                elapsed_time = time.time() - start_time
                mbps = (sub.total_bytes_received * 8) / (elapsed_time * 1_000_000)
                print(f"🚀 Velocidad de recepción media: {mbps:.2f} Mbps")
                
                # Mostrar las imágenes si el mensaje contiene imágenes
                for i, img in enumerate(images):
                    cv2.imshow(f"Imagen Recibida {i+1}", img)
                cv2.waitKey(1)  # Esperar 1 ms para actualizar la ventana
                
                # Imprimir el JSON recibido en data
                print(f"📄 Datos recibidos: {data}")
    except KeyboardInterrupt:
        print("\n🛑 Suscriptor detenido.")
        sub.close()
        cv2.destroyAllWindows()  # Cerrar todas las ventanas de OpenCV

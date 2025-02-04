# PubSub ZMQ Project

Este proyecto contiene una biblioteca de Python que facilita la publicación y suscripción de mensajes que contienen múltiples imágenes y datos adicionales utilizando ZeroMQ. La biblioteca está diseñada para manejar grandes datos binarios de manera eficiente mediante la fragmentación de mensajes en partes más pequeñas.

## Estructura del Proyecto

El proyecto está organizado de la siguiente manera:

- `src/python`: Contiene el código fuente de la biblioteca y los scripts de ejemplo.
  - `publisher.py`: Implementación del publicador.
  - `subscriber.py`: Implementación del suscriptor.
  - `main_publisher.py`: Script de ejemplo para el publicador.
  - `main_subscriber.py`: Script de ejemplo para el suscriptor.
  - `tests`: Contiene las pruebas unitarias para la biblioteca.
- `assets`: Contiene los recursos utilizados en los ejemplos, como imágenes.
- `protocol.md`: Documentación del protocolo utilizado para codificar y decodificar mensajes.
- `requirements.txt`: Lista de dependencias necesarias para ejecutar el proyecto.
- `.gitignore`: Lista de archivos y directorios que deben ser ignorados por Git.
- `pytest.ini`: Configuración de pytest para el proyecto.

## Instalación

Para instalar las dependencias necesarias, ejecuta:

```bash
pip install -r src/python/requirements.txt
```

## Uso

### Publicador

El publicador lee una imagen, la rota para crear múltiples fotogramas y publica los fotogramas junto con datos adicionales.

```python
from publisher import Publisher
import time
import cv2

if __name__ == "__main__":
    pub = Publisher()
    angle = 0
    try:
        frame = cv2.imread("../assets/pong.png")  # Ruta relativa actualizada
        assert frame is not None, "❌ No se pudo cargar la imagen: pong.png"
            
        while True:
            frames = []
            for _ in range(3):
                angle = (angle + 10) % 360
                center = (frame.shape[1] // 2, frame.shape[0] // 2)
                matrix = cv2.getRotationMatrix2D(center, angle, 1.0)
                rotated_frame = cv2.warpAffine(frame, matrix, (frame.shape[1], frame.shape[0]))
                frames.append(rotated_frame)

            data = {"key": "probando"}
            message_bytes = pub.build_message(frames, data)
            pub.publish_message(message_bytes)
            time.sleep(.1)
            
    except KeyboardInterrupt:
        print("\n🛑 Publicador detenido.")
        pub.close()
```

### Suscriptor

El suscriptor recibe los mensajes fragmentados, reconstruye las imágenes y las muestra utilizando OpenCV.

```python
from subscriber import Subscriber
import time
import cv2

if __name__ == "__main__":
    sub = Subscriber()
    try:
        while True:
            start_time = time.time()
            sub.total_bytes_received = 0
            images, data = sub.receive_message()
            if images is not None:
                elapsed_time = time.time() - start_time
                mbps = (sub.total_bytes_received * 8) / (elapsed_time * 1_000_000)
                print(f"🚀 Velocidad de recepción media: {mbps:.2f} Mbps")
                
                for i, img in enumerate(images):
                    cv2.imshow(f"Imagen Recibida {i+1}", img)
                cv2.waitKey(1)
                
                print(f"📄 Datos recibidos: {data}")
    except KeyboardInterrupt:
        print("\n🛑 Suscriptor detenido.")
        sub.close()
        cv2.destroyAllWindows()
```

## Ejecución de Pruebas

Para ejecutar las pruebas utilizando pytest, navega al directorio raíz del proyecto y ejecuta:

```bash
pytest
```

Esto descubrirá y ejecutará automáticamente todas las pruebas en el directorio `tests`.

## Licencia

Este proyecto está licenciado bajo la Licencia MIT.

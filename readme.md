# PubSub ZMQ Project

Este proyecto contiene librerías que facilitan la publicación y suscripción de mensajes que contienen múltiples imágenes y datos adicionales utilizando ZeroMQ. La biblioteca está diseñada para manejar grandes datos binarios de manera eficiente mediante la fragmentación de mensajes en partes más pequeñas.

## Estado del Proyecto

- **Python:**  
  La implementación en Python está completada. La biblioteca permite la publicación y suscripción de mensajes, reconstrucción de imágenes y procesamiento de datos adicionales utilizando un protocolo basado en ZeroMQ.

- **Go:**  
  La implementación en Go se encuentra en desarrollo. Actualmente se dispone de una versión muy sencilla que permite la comunicación básica a través de ZeroMQ, pero aún no soporta completamente el protocolo definido en la versión de Python.

- **C++:**  
  Se tiene previsto desarrollar una versión en C++ del proyecto, lo que permitirá ampliar la compatibilidad y explorar mejoras en el rendimiento.

## Estructura del Proyecto

El proyecto está organizado de la siguiente manera:

- **src/python:**  
  Contiene el código fuente de la biblioteca y los scripts de ejemplo.
  - `publisher.py`: Implementación del publicador.
  - `subscriber.py`: Implementación del suscriptor.
  - `main_publisher.py`: Script de ejemplo para el publicador.
  - `main_subscriber.py`: Script de ejemplo para el suscriptor.
  - `tests`: Contiene las pruebas unitarias para la biblioteca.
  - `pytest.ini`: Configuración de pytest para el proyecto.
  - `requirements.txt`: Lista de dependencias necesarias para ejecutar el proyecto.

- **src/go:**  
  Contiene la implementación en Go (actualmente en proceso).
  - `publisher.go`: Implementación básica del publicador.
  - `subscriber.go`: Implementación básica del suscriptor.
  - `examples`: Scripts de ejemplo para el publicador y suscriptor.
  - `tests`: Pruebas unitarias para la implementación en Go.
  - `go.mod`, `go.sum`: Configuración del módulo Go.

- **assets:**  
  Contiene recursos utilizados en los ejemplos, como imágenes (por ejemplo, `pong.png`).

- **Otros archivos:**
  - `protocol.md`: Documentación del protocolo utilizado para codificar y decodificar mensajes.
  - `docker-compose.yml`: Configuración de Docker Compose para ejecutar los tests en entornos de Python y Go.
  - `run_tests.sh`: Script con logs coloridos para ejecutar secuencialmente los tests de Go y Python.
  - `readme.md`: Este archivo.

## Instalación

Para instalar las dependencias necesarias en la parte de Python, ejecuta:

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

El proyecto está preparado para ejecutar las pruebas utilizando Docker. El archivo **docker-compose.yml** configura dos servicios (Go y Python) que, al ejecutarse, corren los tests de cada entorno.

Para ejecutar las pruebas, puedes utilizar:

```bash
docker compose up
```

Esto iniciará los servicios y ejecutará los tests. También puedes usar el script **run_tests.sh** para ejecutarlos de forma secuencial (primero los tests de Go y, de ser exitosos, luego los de Python):

```bash
./run_tests.sh
```

## Pendiente

- **Python:**  
  Resolver el problema con el test corto que no envía imágenes, el cual falla en la recepción en el entorno Docker (aunque fuera funciona correctamente).

- **Go:**  
  Terminar la implementación de la librería en Go para soportar completamente el protocolo utilizado en la versión de Python.

- **C++:**  
  Desarrollar una versión en C++ del proyecto para ampliar la compatibilidad y explorar mejoras de rendimiento.

## Licencia

Este proyecto está licenciado bajo la Licencia MIT.


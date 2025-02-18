# PubSub ZMQ Project

Este proyecto contiene librerías que facilitan la publicación y suscripción de mensajes que contienen múltiples imágenes y datos adicionales utilizando ZeroMQ. La biblioteca está diseñada para manejar grandes datos binarios de manera eficiente mediante la fragmentación de mensajes en partes más pequeñas.

## Estado del Proyecto

- **Python:**  
  La implementación en Python está completada. La biblioteca permite la publicación y suscripción de mensajes, reconstrucción de imágenes y procesamiento de datos adicionales utilizando un protocolo basado en ZeroMQ.

- **Go:**  
  La implementación en Go se encuentra en desarrollo. Actualmente se dispone de una versión sencilla que soporta el protocolo de mensajes pero obvia las imágenes. Se han añadido pruebas de integración entre Go y Python, que verifican que los mensajes enviados por uno son recibidos correctamente por el otro.

- **C++:**  
  La implementación en C++ está en progreso. Ya se cuenta con una versión básica del publicador y suscriptor que soporta el protocolo, aunque actualmente no se realiza el manejo de imágenes. Se han creado pruebas de integración entre C++ y Python, y entre C++ y Go. Los publicadores y suscriptores en C++ generan archivos JSON con los mensajes enviados y recibidos para validación.

## Estructura del Proyecto

El proyecto está organizado de la siguiente manera:

- **src/python:**  
  Contiene el código fuente de la biblioteca y los scripts de ejemplo.
  - `publisher.py`: Implementación del publicador.
  - `subscriber.py`: Implementación del suscriptor.
  - `main_publisher.py`: Script de ejemplo para el publicador.
  - `main_subscriber.py`: Script de ejemplo para el suscriptor.
  - `tests`: Contiene las pruebas unitarias e integración.
  - `pytest.ini`: Configuración de pytest para el proyecto.
  - `requirements.txt`: Lista de dependencias necesarias para ejecutar el proyecto.

- **src/go:**  
  Contiene la implementación en Go.
  - `publisher.go`: Implementación básica del publicador.
  - `subscriber.go`: Implementación básica del suscriptor.
  - `examples`: Scripts de ejemplo para el publicador y suscriptor.
  - `tests`: Pruebas unitarias y de integración entre Go y Python.
  - `go.mod`, `go.sum`: Configuración del módulo Go.

- **src/cpp:**  
  Contiene la implementación en C++.
  - `include`: Cabeceras de las clases `Publisher` y `Subscriber`.
  - `src`: Implementación de `Publisher` y `Subscriber`.
  - `tests`: Scripts de integración para verificar el funcionamiento entre C++, Go y Python.
  - `Dockerfile`: Imagen Docker para compilar y ejecutar los ejemplos y tests en C++.
  - `CMakeLists.txt`: Configuración de CMake para compilar el código.

- **assets:**  
  Contiene recursos utilizados en los ejemplos, como imágenes (por ejemplo, `pong.png`).

- **Otros archivos:**
  - `protocol.md`: Documentación del protocolo utilizado para codificar y decodificar mensajes.
  - `docker-compose.yml`: Configuración de Docker Compose para ejecutar los tests en entornos de Python, Go y C++.
  - `run_tests.sh`: Script con logs coloridos para ejecutar secuencialmente los tests de C++, Go y Python.
  - `readme.md`: Este archivo.

## Instalación

Para instalar las dependencias necesarias en la parte de Python, ejecuta:

```bash
pip install -r src/python/requirements.txt
```

## Uso

### Publicador

Ejecuta el publicador de ejemplo en Python:

```bash
python src/python/main_publisher.py
```

Ejecuta el publicador de ejemplo en Go:

```bash
go run src/go/examples/main_publisher.go
```

Ejecuta el publicador de ejemplo en C++ dentro del contenedor Docker:

```bash
docker-compose run --rm integration_cpp_pub
```

### Suscriptor

Ejecuta el suscriptor de ejemplo en Python:

```bash
python src/python/main_subscriber.py
```

Ejecuta el suscriptor de ejemplo en Go:

```bash
go run src/go/examples/main_subscriber.go
```

Ejecuta el suscriptor de ejemplo en C++ dentro del contenedor Docker:

```bash
docker-compose run --rm integration_cpp_sub
```

## Ejecución de Pruebas

El proyecto está preparado para ejecutar las pruebas utilizando Docker. El archivo **docker-compose.yml** configura servicios para C++, Go y Python, que al ejecutarse corren los tests de cada entorno y pruebas de integración entre ellos.

Para ejecutar las pruebas, puedes utilizar:

```bash
docker-compose up
```

Esto iniciará los servicios y ejecutará los tests. También puedes usar el script **run_tests.sh** para ejecutarlos de forma secuencial (primero C++, luego Go y finalmente Python):

```bash
./run_tests.sh
```

## Licencia

Este proyecto está licenciado bajo la Licencia MIT.


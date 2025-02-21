# PubSub ZMQ Project

# 🏗️🚧 WORK IN PROGRESS 🚧🏗️

Este proyecto contiene librerías que facilitan la publicación y suscripción de mensajes que contienen múltiples imágenes y datos adicionales utilizando ZeroMQ.  
La biblioteca está diseñada para manejar grandes datos binarios de manera eficiente mediante la fragmentación de mensajes en partes más pequeñas.

---

## 🚀 **Estado del Proyecto**

- **Python:**  
  ✅ Implementación completada. Se pueden publicar y suscribir mensajes, reconstruir imágenes y procesar datos adicionales con ZeroMQ.  
  ✅ Genera un archivo JSON con el mensaje enviado (`/shared/result_publisher.json`) y otro con el mensaje recibido (`/shared/result.json`).  

- **Go:**  
  🔧 Implementación en progreso. Se soporta el protocolo de mensajes, pero aún no se procesan imágenes.  
  ✅ Se han añadido pruebas de integración con Python y C++ para verificar la compatibilidad.  
  ✅ Ahora los suscriptores usan `PUBSUB_ENDPOINT` para configurar la dirección del publicador.  

- **C++:**  
  🔧 Implementación en progreso. Ya se soporta el protocolo de mensajes, pero aún no se procesan imágenes.  
  ✅ Se han creado pruebas de integración entre C++ y Python, y entre C++ y Go.  
  ✅ Se han estandarizado los mensajes JSON para asegurar compatibilidad entre lenguajes.  

---

## 📂 **Estructura del Proyecto**

```
.
├── src/
│   ├── python/
│   │   ├── publisher.py  # Implementación del publicador
│   │   ├── subscriber.py  # Implementación del suscriptor
│   │   ├── main_publisher.py  # Script de prueba del publicador
│   │   ├── main_subscriber.py  # Script de prueba del suscriptor
│   │   ├── tests/  # Pruebas unitarias e integración
│   │   ├── pytest.ini  # Configuración de pytest
│   │   ├── requirements.txt  # Dependencias
│   ├── go/
│   │   ├── publisher.go  # Implementación del publicador
│   │   ├── subscriber.go  # Implementación del suscriptor
│   │   ├── main_publisher.go  # Script de prueba del publicador
│   │   ├── main_subscriber.go  # Script de prueba del suscriptor
│   │   ├── tests/  # Pruebas unitarias e integración
│   │   ├── go.mod / go.sum  # Configuración del módulo Go
│   ├── cpp/
│   │   ├── include/  # Cabeceras de las clases Publisher y Subscriber
│   │   ├── src/  # Implementaciones de Publisher y Subscriber
│   │   ├── tests/  # Pruebas de integración entre C++, Go y Python
│   │   ├── CMakeLists.txt  # Configuración de CMake
│   │   ├── Dockerfile  # Imagen Docker para compilar y ejecutar C++
├── assets/  # Contiene archivos como imágenes utilizadas en las pruebas
├── protocol.md  # Documentación del protocolo de mensajes
├── docker-compose.yml  # Configuración de Docker Compose para pruebas
├── run_tests.sh  # Script para ejecutar todas las pruebas de integración
└── readme.md  # Este archivo
```

---

## 📜 **Uso de Publicadores y Suscriptores**

### ▶️ **Ejecutar un Publicador**
📌 **Python**
```bash
python src/python/main_publisher.py
```
📌 **Go**
```bash
go run src/go/main_publisher.go
```
📌 **C++ (Dentro de Docker)**
```bash
docker-compose run --rm integration_cpp_pub
```

---

### 🎧 **Ejecutar un Suscriptor**
📌 **Python**
```bash
python src/python/main_subscriber.py
```
📌 **Go**
```bash
go run src/go/main_subscriber.go
```
📌 **C++ (Dentro de Docker)**
```bash
docker-compose run --rm integration_cpp_sub
```

---

## 🔧 **Configuración de Endpoint con `PUBSUB_ENDPOINT`**
Por defecto, los suscriptores se conectan a `tcp://localhost:5555`,  
pero puedes definir un **endpoint diferente** usando la variable de entorno `PUBSUB_ENDPOINT`.

📌 **Ejemplo en Python**
```bash
PUBSUB_ENDPOINT=tcp://integration_cpp_pub:5555 python src/python/main_subscriber.py
```

📌 **Ejemplo en Go**
```bash
PUBSUB_ENDPOINT=tcp://integration_go_pub:5555 go run src/go/main_subscriber.go
```

---

## 🛠️ **Ejecución de Pruebas**
El proyecto usa **Docker Compose** para ejecutar los tests en Python, Go y C++.

🔹 **Para ejecutar todas las pruebas unitarias e integración:**
```bash
docker-compose up
```

🔹 **Para ejecutarlas secuencialmente con logs detallados:**
```bash
./run_tests.sh
```

---

## 📊 **Estructura del Mensaje**
Todos los lenguajes manejan la misma estructura JSON.

```json
{
  "type": "images",
  "count": 1,
  "images": [
    {
      "metadata": {
        "width": 1222,
        "height": 945,
        "channels": 3,
        "dtype": "uint8",
        "size": 3464370
      }
    }
  ],
  "data": {
    "key": "probando",
    "index": 0
  }
}
```
🔹 **Cada publicador guarda el mensaje enviado en:**  
📌 `/shared/result_publisher.json`  

🔹 **Cada suscriptor guarda el mensaje recibido en:**  
📌 `/shared/result.json`

Si los archivos **son diferentes, la prueba fallará**.

---

## 📌 **Estado Actual de las Pruebas de Integración**
✅ **Python ↔ Go** (Mensajes enviados y recibidos coinciden correctamente).  
✅ **Python ↔ C++** (Mensajes enviados y recibidos coinciden correctamente).  
✅ **Go ↔ C++** (Mensajes enviados y recibidos coinciden correctamente).  

---

## 📝 **Licencia**
Este proyecto está licenciado bajo la **Licencia MIT**.

# PubSub ZMQ Project

## 🏗️🚧 Estado del Proyecto 🚧🏗️

Este proyecto implementa **publicadores y suscriptores ZeroMQ** en **Python, Go y C++**  
para transmitir mensajes con imágenes y datos en formato JSON.

Todos los lenguajes usan el mismo **protocolo JSON**,  
permitiendo la interoperabilidad entre diferentes implementaciones.

---

## 🚀 **Estado Actual**
| Lenguaje  | Estado |
|-----------|--------|
| **Python** | ✅ **Completado.** Publicador y suscriptor funcionando con imágenes PNG/JPEG/Raw Opencv. |
| **Go** | ✅ **Completado.** Publicador y suscriptor con soporte para imágenes PNG/JPEG. **NO SOPORTA RAW DE OPENCV** |
| **C++** | ✅ **Completado.** Publicador y suscriptor con imágenes con imágenes PNG/JPEG/Raw Opencv. |

✅ **Las pruebas de integración confirman la compatibilidad entre Python, Go y C++.**

---
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
        "width": 100,
        "height": 100,
        "channels": 3,
        "dtype": "uint8",
        "size": 30000
      }
    }
  ],
  "data": {
    "key": "probando",
    "index": 0
  }
}

```


## 📝 **Licencia**
Este proyecto está licenciado bajo la **Licencia MIT**.

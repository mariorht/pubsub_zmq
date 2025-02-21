### 📜 **PubSub ZMQ Protocol**  

Este documento describe el protocolo utilizado en **PubSub ZMQ** para codificar y decodificar mensajes que contienen **múltiples imágenes** y **datos adicionales** mediante **ZeroMQ**.  
El protocolo está diseñado para **fragmentar y reconstruir** los mensajes correctamente, asegurando compatibilidad entre **Python, Go y C++**.

---

## 🏗️ **Estructura del Mensaje**  

Un mensaje en **PubSub ZMQ** está compuesto por **dos partes principales**:  

1. **Metadatos JSON** → Contiene información de las imágenes y datos adicionales.  
2. **Datos Binarios de las Imágenes** → Contiene las imágenes codificadas en el formato especificado.  

Ambas partes están concatenadas y separadas por un **byte nulo (`\x00`)**.  

---

### 📌 **1️⃣ Metadatos JSON**  

El JSON contiene **información de las imágenes** y **datos adicionales**.  

📌 **Ejemplo de JSON estructurado:**  
```json
{
    "type": "images",
    "count": 2,
    "images": [
        {
            "metadata": {
                "format": "jpg",
                "width": 1222,
                "height": 945,
                "size": 3464370,
                "dtype": "uint8"
            }
        },
        {
            "metadata": {
                "format": "jpg",
                "width": 800,
                "height": 600,
                "size": 2145689,
                "dtype": "uint8"
            }
        }
    ],
    "data": {
        "key": "prueba",
        "index": 0
    }
}
```

🔹 **Explicación de los campos:**  
| Campo       | Descripción |
|-------------|------------|
| `"type"`    | Tipo de mensaje (`"images"` para indicar que contiene imágenes). |
| `"count"`   | Cantidad de imágenes en el mensaje. |
| `"images"`  | Lista de metadatos de cada imagen. |
| `"format"`  | Formato de codificación de la imagen (`"jpg"`, `"png"`, etc.). |
| `"width"`   | Ancho de la imagen en píxeles. |
| `"height"`  | Alto de la imagen en píxeles. |
| `"size"`    | Tamaño de la imagen en bytes. |
| `"dtype"`   | Tipo de datos de la imagen (`"uint8"`, `"float32"`, etc.). |
| `"data"`    | Información adicional enviada con el mensaje (puede ser cualquier JSON válido). |

---

### 📌 **2️⃣ Datos Binarios de las Imágenes**  

Después del JSON, se concatenan **las imágenes codificadas en binario** en el formato especificado (`JPEG`, `PNG`, etc.).

#### **Ejemplo de flujo de datos en memoria:**
```
<JSON METADATA> + \x00 + <BINARY IMAGE 1> + <BINARY IMAGE 2> + ...
```
Donde:
- `<JSON METADATA>` es la información de las imágenes y datos adicionales.
- `\x00` (NULL BYTE) es el separador.
- `<BINARY IMAGE 1>`, `<BINARY IMAGE 2>` son las imágenes codificadas.

---

## 📦 **Fragmentación del Mensaje**  

Para manejar **mensajes grandes**, se fragmentan en **múltiples partes**.  
Cada fragmento es enviado como un **mensaje multipart** con los siguientes elementos:

1️⃣ **Topic** → Tema de publicación (`"camera/image_raw"`, `"test"`, etc.).  
2️⃣ **Índice del Fragmento** → Número de fragmento dentro del mensaje.  
3️⃣ **Total de Fragmentos** → Cantidad total de fragmentos en el mensaje.  
4️⃣ **Datos del Fragmento** → Parte del mensaje en binario.  

---

### **🔄 Ejemplo de Envío**
Si un mensaje pesa **300 KB** y el tamaño máximo de fragmento es **100 KB**, se enviará en **3 fragmentos**:

| Topic              | Índice | Total | Datos |
|--------------------|--------|-------|-------|
| `"camera/image_raw"` | `0` | `3` | `Fragmento 1` |
| `"camera/image_raw"` | `1` | `3` | `Fragmento 2` |
| `"camera/image_raw"` | `2` | `3` | `Fragmento 3` |

Cada **Subscriber** debe **reensamblar** los fragmentos en el orden correcto antes de procesar el mensaje.

---

## 📥 **Recepción y Reconstrucción del Mensaje**  

El **Subscriber** recibe los fragmentos y reconstruye el mensaje en **6 pasos**:

1️⃣ **Recibe todos los fragmentos** en el orden correcto.  
2️⃣ **Une los fragmentos** en un solo flujo de bytes.  
3️⃣ **Separa el JSON de los datos binarios** usando `\x00`.  
4️⃣ **Decodifica el JSON** para obtener los metadatos de las imágenes.  
5️⃣ **Extrae cada imagen** según su tamaño y la reconstruye.  
6️⃣ **Procesa los datos adicionales** enviados en `"data"`.  

---

## 🔧 **Manejo de Errores**  

Si el mensaje **no puede ser reconstruido**, se descarta y se genera un **log de error**.  

| Error Detectado | Acción |
|-----------------|--------|
| **Fragmentos perdidos** | Descarta el mensaje. |
| **Separador `\x00` no encontrado** | Descarta el mensaje. |
| **JSON inválido** | Descarta el mensaje. |
| **Datos binarios incompletos** | Descarta el mensaje. |

---

## ✅ **Resumen del Protocolo**
✔️ **Mensajes compuestos por JSON + Imágenes en binario.**  
✔️ **Fragmentación en múltiples partes si el mensaje es grande.**  
✔️ **Cada fragmento contiene índice y total de fragmentos.**  
✔️ **Recepción y reconstrucción garantizan la entrega confiable.**  

---

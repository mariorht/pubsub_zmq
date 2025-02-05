# src/python/Dockerfile
FROM python:3.8

# Instala la dependencia de OpenCV (libGL)
RUN apt-get update && apt-get install -y libgl1-mesa-glx

WORKDIR /app

# Copia el código Python y los archivos de configuración
COPY . .

# Instala las dependencias de Python
RUN pip install --no-cache-dir -r requirements.txt

# Comando por defecto para ejecutar pytest
CMD ["pytest"]

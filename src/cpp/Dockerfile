FROM ubuntu:24.04

# Instalar dependencias
RUN apt-get update && apt-get install -y \
    cmake \
    g++ \
    libzmq3-dev \
    make

# Crear directorio de trabajo
WORKDIR /app

# Copiar el código fuente
COPY . .

# Compilar usando CMake
RUN mkdir -p build && cd build && cmake .. && make

# CMD ["/app/build/main_publisher"]

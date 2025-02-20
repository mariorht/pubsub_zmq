FROM ubuntu:24.04

# Instalar dependencias
RUN apt-get update && apt-get install -y \
    cmake \
    g++ \
    libzmq3-dev \
    libopencv-dev \
    pkg-config \
    make \
    libgtest-dev

# Compilar GoogleTest ya que no trae los binarios precompilados
RUN cd /usr/src/gtest && cmake . && make && cp lib/*.a /usr/lib/

# Crear directorio de trabajo
WORKDIR /app

# Copiar el código fuente
COPY . .

# Compilar usando CMake
RUN mkdir -p build && cd build && cmake .. && make

# CMD ["/app/build/main_publisher"]

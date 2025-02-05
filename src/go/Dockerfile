# src/go/Dockerfile
FROM golang:1.13

# Instala libzmq3-dev (dependencia de ZeroMQ)
RUN apt-get update && apt-get install -y libzmq3-dev

WORKDIR /app

# Copia el código al contenedor
COPY . .

# Comando por defecto para ejecutar los tests
CMD ["go", "test", "-v", "./tests"]

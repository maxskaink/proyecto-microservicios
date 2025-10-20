FROM golang:1.25-alpine

# Instalar dependencias y Air
RUN apk add --no-cache git bash curl ca-certificates tzdata && \
    go install github.com/air-verse/air@latest

# Directorio de trabajo
WORKDIR /app

# Copiar los archivos de módulos y descargar dependencias
COPY go.mod go.sum* ./
RUN go mod download

# Copiar el resto del código
COPY . .

# Variables de entorno
ENV CGO_ENABLED=0 \
    GOPROXY=https://proxy.golang.org,direct

# Crear archivo dummy si no existe
RUN mkdir -p /app/secrets && \
    if [ ! -f /app/secrets/serviceAccountKey.json ]; then \
        echo '{"type":"service_account"}' > /app/secrets/serviceAccountKey.json; \
    fi

# Exponer puerto
EXPOSE 8080

# Comando por defecto
CMD ["air", "-c", ".air.toml"]

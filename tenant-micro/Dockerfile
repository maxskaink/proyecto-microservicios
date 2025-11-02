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

# Exponer puerto
EXPOSE 8082

# Comando por defecto
CMD ["air", "-c", ".air.toml"]

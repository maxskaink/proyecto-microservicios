FROM golang:1.25-alpine

# Instalar dependencias básicas
RUN apk add --no-cache git bash curl ca-certificates tzdata

# Directorio de trabajo
WORKDIR /app

# Copiar solo los archivos de módulos primero (para mejor caché)
COPY go.mod go.sum* ./
RUN go mod download

# Ahora copia todo el código fuente
COPY . .

# Configurar variables Go
ENV CGO_ENABLED=0 \
    GOPROXY=https://proxy.golang.org,direct

# Crear estructura para secrets si no existe
RUN mkdir -p /app/secrets && \
    if [ ! -f /app/secrets/serviceAccountKey.json ]; then \
        echo '{"type":"service_account"}' > /app/secrets/serviceAccountKey.json; \
    fi

# Exponer puerto
EXPOSE 8080

# Comando de ejecución
CMD ["go", "run", "./cmd/server"]
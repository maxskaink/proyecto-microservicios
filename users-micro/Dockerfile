##############################
# Etapa base (dependencias)
##############################
FROM golang:1.25-alpine AS base
RUN apk add --no-cache git bash curl ca-certificates tzdata
WORKDIR /app
COPY users-micro/go.mod users-micro/go.sum* ./
COPY pkg/observability /pkg/observability
RUN go mod download

##############################
# Builder (compila binario)
##############################
FROM base AS builder
ARG BUILD_DIR=./cmd/server
ENV CGO_ENABLED=0 \
    GOPROXY=https://proxy.golang.org,direct
COPY users-micro/ .
RUN mkdir -p /app/secrets && \
    if [ ! -f /app/secrets/serviceAccountKey.json ]; then echo '{"type":"service_account"}' > /app/secrets/serviceAccountKey.json; fi
RUN go build -ldflags "-s -w" -o /out/app "$BUILD_DIR"

##############################
# Dev (hot reload con Air)
##############################
FROM base AS dev
RUN go install github.com/air-verse/air@latest
ARG APP_PORT=8080
ENV CGO_ENABLED=0 PORT=$APP_PORT
COPY users-micro/ .
RUN mkdir -p /app/secrets && \
    if [ ! -f /app/secrets/serviceAccountKey.json ]; then echo '{"type":"service_account"}' > /app/secrets/serviceAccountKey.json; fi
EXPOSE $APP_PORT
CMD ["air", "-c", ".air.toml"]

##############################
# Prod (imagen mínima)
##############################
FROM alpine:3.20 AS prod
RUN apk add --no-cache ca-certificates tzdata && adduser -D -H appuser
WORKDIR /app
ARG APP_PORT=8080
ENV PORT=$APP_PORT
COPY --from=builder /out/app ./app
COPY --from=builder /app/secrets ./secrets
EXPOSE $APP_PORT
USER appuser
ENTRYPOINT ["./app"]

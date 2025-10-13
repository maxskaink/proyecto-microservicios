# syntax=docker/dockerfile:1.6

FROM golang:1.22-alpine AS builder
WORKDIR /app

# Habilitar módulos y descargar dependencias (vacío inicialmente)
COPY go.mod ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download || true

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/users-micro ./cmd/server

FROM gcr.io/distroless/base-debian12
WORKDIR /
COPY --from=builder /out/users-micro /users-micro

ENV PORT=8080
EXPOSE 8080

USER nonroot:nonroot
ENTRYPOINT ["/users-micro"]

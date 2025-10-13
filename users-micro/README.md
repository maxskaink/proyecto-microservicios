# users-micro

Microservicio (esqueleto) en Go con Gin + GORM (Postgres), arquitectura en 4 capas: controladores, servicios, dominio (DDD) y repositorios. Autenticación prevista con Firebase (pendiente de implementación).

## Estructura

- cmd/server: entrypoint
- internal/
  - controllers: HTTP handlers y mapeo de rutas
  - services: lógica de negocio (interfaces)
  - domain: entidades DDD
  - repositories: contratos de persistencia (GORM + Postgres)
  - server: router, middlewares, bootstrapping
  - config: configuración
- pkg/logger: logger simple
- test/: pruebas iniciales

## Ejecutar (dev)

Requiere Go 1.22+. Aún sin dependencias inicializadas.

## Docker

Se incluye Dockerfile multi-stage (build + runtime). Variables de entorno típicas:

- PORT
- DATABASE_URL
- FIREBASE_PROJECT_ID

## Notas

- Este proyecto es un esqueleto: no hay lógica implementada, solo contratos y estructura.
- La auth con Firebase se integrará con un middleware que verifica el token de Authorization.
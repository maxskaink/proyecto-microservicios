# API Gateway con Traefik

Este componente sirve como punto de entrada único para todos los microservicios del proyecto.

## Características

- Enrutamiento automático basado en rutas (`/api/users`, etc.)
- Descubrimiento dinámico de servicios mediante integración con Consul
- Dashboard para monitoreo visual en http://localhost:8080/dashboard/
- Balanceo de carga automático entre múltiples instancias del mismo servicio

## Configuración

La configuración de Traefik se divide en dos partes:

1. **Configuración estática** (`traefik.yml`): Define los proveedores, puntos de entrada y configuración general.
2. **Configuración dinámica** (`dynamic/services.yml`): Define rutas, middleware y servicios.

## Acceso a los servicios

Una vez que el contenedor esté en funcionamiento, puedes acceder a los servicios a través de:

- API Gateway: http://localhost
  - Microservicio de usuarios: http://localhost/api/users
- Dashboard de Traefik: http://localhost:8080/dashboard/

## Configuración para nuevos microservicios

Para agregar un nuevo microservicio al gateway, puedes:

1. **Registrarlo en Consul** (recomendado): El servicio se registrará automáticamente y Traefik lo detectará.
2. **Configuración estática**: Añadir una nueva entrada en `dynamic/services.yml`.

## Notas importantes

- Este gateway está configurado para entornos de desarrollo. Para producción, desactiva `insecure: true` y configura TLS.
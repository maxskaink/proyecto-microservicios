# Servicio de Descubrimiento con Consul

Este componente permite el registro y descubrimiento dinámico de servicios en la arquitectura de microservicios.

## Características

- Registro automático de servicios
- UI de administración en http://localhost:8500/ui/
- Health checks para monitoreo de servicios
- Integración con Traefik para enrutamiento dinámico
- Almacén de clave-valor para configuración distribuida

## Configuración

La configuración principal se encuentra en `config.json` y define:

1. Configuración del servidor Consul
2. Servicios registrados estáticamente (para respaldo)
3. Health checks para los servicios

## Acceso a la interfaz

Una vez que el contenedor esté en funcionamiento, puedes acceder a:

- UI de Consul: http://localhost:8500/ui/

## Integración con microservicios

Para que un microservicio se registre en Consul, debe:

1. Implementar un endpoint de health check (típicamente `/health`)
2. Registrarse en Consul al iniciar
3. Desregistrarse al apagar

Este registro puede realizarse de dos formas:
1. A través de la API HTTP de Consul
2. Usando un cliente de Consul en tu código Go

## Ejemplo de registro manual via HTTP

```bash
curl -X PUT -d '{
  "ID": "my-service-1",
  "Name": "my-service",
  "Tags": ["traefik.enable=true"],
  "Address": "my-service-container",
  "Port": 8080,
  "Check": {
    "HTTP": "http://my-service-container:8080/health",
    "Interval": "10s"
  }
}' http://localhost:8500/v1/agent/service/register
```
# Kong Gateway - Configuración

Este directorio contiene la configuración de Kong Gateway como API Gateway para el proyecto de microservicios.

## Estructura

```
api-gateway/
├── Dockerfile              # Imagen de Kong con plugin personalizado
├── init-kong.sh           # Script de inicialización e integración con Consul
├── kong/
│   └── kong.yml          # Configuración declarativa de Kong (servicios y rutas)
└── plugins/
    └── tenant-extractor/  # Plugin Lua para multitenencia
        ├── handler.lua    # Lógica del plugin
        └── schema.lua     # Esquema de configuración
```

## Características

### 1. API Gateway con Kong
- Kong Gateway 3.5 en modo declarativo (sin base de datos)
- Configuración mediante archivo YAML (`kong.yml`)
- Proxy en puerto 8000 (mapeado a 80 en el host)
- Admin API en puerto 8001 (mapeado a 8090 en el host)

### 2. Integración con Consul
- Kong se registra automáticamente en Consul al iniciar
- Script `init-kong.sh` maneja el registro del servicio
- Health checks configurados para monitoreo

### 3. Plugin de Multitenencia (tenant-extractor)
El plugin Lua personalizado extrae el tenant desde la URL y lo agrega como header.

#### Funcionamiento:
- **URL sin tenant**: `/api/users` → No agrega header `X-Tenant-ID`
- **URL con tenant**: `/tenant1/api/users` → Agrega header `X-Tenant-ID: tenant1`

#### Implementación:
- El plugin analiza el path de la request
- Extrae el primer segmento si coincide con el patrón `/[tenant]/api/...`
- Agrega el header `X-Tenant-ID` con el valor del tenant
- Los microservicios backend reciben el header y procesan según el tenant

## Configuración de Rutas

Kong está configurado con las siguientes rutas (definidas en `kong.yml`):

### Usuarios
- `GET /api/users` → `users-service` (sin tenant)
- `GET /tenant1/api/users` → `users-service` (con tenant)

### Productos
- `GET /api/products` → `products-service` (sin tenant)
- `GET /tenant1/api/products` → `products-service` (con tenant)

### Tenants
- `GET /api/tenants` → `tenant-service` (sin tenant)
- `GET /tenant1/api/tenants` → `tenant-service` (con tenant)

## Servicios Backend

Kong hace proxy a los siguientes servicios:
- **users-service**: `http://users_api:8080`
- **products-service**: `http://product_api:8081`
- **tenant-service**: `http://tenant_api:8082`

## Admin API

Kong Admin API está disponible en `http://localhost:8090` para:
- Consultar servicios: `GET http://localhost:8090/services`
- Consultar rutas: `GET http://localhost:8090/routes`
- Consultar plugins: `GET http://localhost:8090/plugins`
- Ver status: `GET http://localhost:8090/status`

## Desarrollo del Plugin

El plugin `tenant-extractor` está escrito en Lua siguiendo la estructura estándar de Kong:

### handler.lua
Contiene la lógica principal:
- Fase `access`: Se ejecuta antes de hacer proxy al backend
- Extrae el tenant del path usando pattern matching
- Agrega el header `X-Tenant-ID` a la request

### schema.lua
Define el esquema de configuración:
- `header_name`: Nombre del header (default: "X-Tenant-ID")

## Logs y Debugging

Ver logs de Kong:
```bash
docker logs api_gateway -f
```

Los logs incluyen:
- Access logs del proxy
- Admin API access logs
- Mensajes de debug del plugin (cuando se encuentra o no un tenant)

## Migración desde Traefik

Cambios principales:
- Puerto proxy: 80 (Kong) vs 80 (Traefik)
- Puerto admin: 8090 (Kong Admin API) vs 8090 (Traefik Dashboard)
- Plugin: Lua (Kong) vs Go (Traefik)
- Configuración: YAML declarativo (Kong) vs YAML dinámico (Traefik)

Las rutas y servicios mantienen la misma funcionalidad que con Traefik.

# proyecto-microservicios
Modular microservices-based application enabling independent development, deployment, and scaling of services. Each service manages a specific business function, ensuring flexibility, resilience, and streamlined integration in distributed systems.

## Arquitectura

El proyecto utiliza:
- **Kong Gateway**: API Gateway con soporte para multitenencia mediante plugin Lua personalizado
- **Consul**: Service Discovery para registro y descubrimiento de servicios
- **RabbitMQ**: Message broker para comunicación asíncrona entre microservicios
- **PostgreSQL**: Base de datos para cada microservicio

### Microservicios
- **users-micro**: Gestión de usuarios (Puerto 8080)
- **product-micro**: Gestión de productos (Puerto 8081)
- **tenant-micro**: Gestión de tenants/inquilinos (Puerto 8082)

### API Gateway - Kong
Kong Gateway actúa como punto de entrada único para todos los microservicios, proporcionando:
- Enrutamiento de requests a los servicios apropiados
- Plugin de multitenencia que extrae el tenant de la URL y lo agrega como header
- Integración con Consul para descubrimiento de servicios
- Load balancing y health checks

#### Rutas soportadas:
- Sin tenant: `/api/users`, `/api/products`, `/api/tenants`
- Con tenant: `/tenant1/api/users`, `/tenant1/api/products`, `/tenant1/api/tenants`

El plugin `tenant-extractor` extrae el tenant de la URL y lo agrega como header `X-Tenant-ID` para que los microservicios puedan aislar datos por tenant.

## Despliegue local (Postgres + users-api)

1. Copia el archivo de variables de entorno:
	- `cp .env.example .env`
2. Levanta los servicios con Docker Compose (desde la raíz del repo):
	- `docker compose -f deploy/docker-compose.yml --env-file .env up --build`
3. Prueba el endpoint de salud del micro de usuarios:
	- `curl http://localhost:$PORT/health`

Variables clave (puedes editarlas en `.env`):
- `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`, `POSTGRES_PORT`
- `PORT` (puerto del micro)
- `DATABASE_URL` (ej: `postgres://users:userspass@db:5432/users?sslmode=disable`)
- `FIREBASE_PROJECT_ID`

Notas:
- El compose construye la imagen usando `users-micro/Dockerfile`.
- El contenedor `db` expone Postgres en el puerto `${POSTGRES_PORT}` local.

## Inicio rápido con Docker Compose

Para iniciar toda la infraestructura (Kong, Consul, RabbitMQ y microservicios):

```bash
./start.sh
```

Este script:
1. Configura variables de entorno necesarias
2. Inicia todos los contenedores con `docker-compose up -d`
3. Espera a que los servicios estén disponibles
4. Muestra información de acceso a los servicios

### Acceso a servicios:
- **API Gateway (Kong)**: http://localhost
- **Kong Admin API**: http://localhost:8090
- **Consul UI**: http://localhost:8500/ui/
- **RabbitMQ Management**: http://localhost:15672 (usuario: guest, contraseña: guest)

### Ejemplos de requests:

#### Sin tenant:
```bash
# Usuarios
curl http://localhost/api/users

# Productos
curl http://localhost/api/products
```

#### Con tenant:
```bash
# Usuarios del tenant1
curl http://localhost/tenant1/api/users

# Productos del tenant1
curl http://localhost/tenant1/api/products
```

El header `X-Tenant-ID` será agregado automáticamente por Kong y enviado a los microservicios.

### Consultar configuración de Kong:
```bash
# Ver servicios configurados
curl http://localhost:8090/services

# Ver rutas configuradas
curl http://localhost:8090/routes

# Ver plugins activos
curl http://localhost:8090/plugins

# Ver status de Kong
curl http://localhost:8090/status
```

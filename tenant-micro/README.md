# Tenant Microservice

Microservicio encargado de la gestión de tenants en el sistema. Permite crear, listar, obtener y eliminar tenants. Cuando se crea o elimina un tenant, se publican eventos a RabbitMQ para que otros microservicios (users-micro y product-micro) creen o eliminen los esquemas correspondientes.

## Estructura del Proyecto

```
tenant-micro/
├── cmd/
│   └── main.go              # Punto de entrada de la aplicación
├── internal/
│   ├── controllers/         # Capa HTTP - manejo de requests
│   ├── services/            # Capa de lógica de negocio
│   ├── repositories/        # Capa de acceso a datos
│   ├── models/              # Modelos de datos (GORM)
│   └── dto/                 # Data Transfer Objects (request/response)
├── .air.toml               # Configuración de hot reload (Air)
├── Dockerfile              # Definición del contenedor
├── go.mod                  # Dependencias del proyecto
├── go.sum                  # Lock de versiones
└── README.md               # Este archivo
```

## Variables de Entorno

```env
# Base de datos
DB_HOST=tenant_db
DB_USER=tenant
DB_PASSWORD=tenantpass
DB_NAME=tenant_db
DB_PORT=5432

# RabbitMQ
RABBITMQ_USER=guest
RABBITMQ_PASSWORD=guest
RABBITMQ_HOST=rabbitmq
RABBITMQ_PORT=5672

# Servidor
PORT=8082
```

## Endpoints

### Crear Tenant
```bash
POST /api/tenants
Content-Type: application/json

{
  "tenant_id": "acme-corp",
  "tenant_name": "ACME Corporation"
}

Response: 201 Created
{
  "id": "uuid",
  "tenant_id": "acme-corp",
  "tenant_name": "ACME Corporation",
  "created_at": "2024-01-15T10:30:00Z"
}
```

### Obtener todos los Tenants
```bash
GET /api/tenants

Response: 200 OK
[
  {
    "id": "uuid",
    "tenant_id": "acme-corp",
    "tenant_name": "ACME Corporation",
    "created_at": "2024-01-15T10:30:00Z"
  }
]
```

### Obtener Tenant por ID
```bash
GET /api/tenants/{tenant_id}

Response: 200 OK
{
  "id": "uuid",
  "tenant_id": "acme-corp",
  "tenant_name": "ACME Corporation",
  "created_at": "2024-01-15T10:30:00Z"
}
```

### Eliminar Tenant
```bash
DELETE /api/tenants/{tenant_id}

Response: 200 OK
{
  "message": "Tenant eliminado correctamente"
}
```

### Health Check
```bash
GET /health

Response: 200 OK
{
  "status": "ok"
}
```

## Eventos Publicados

El microservicio publica eventos en RabbitMQ con exchange `tenants_events` (tipo topic):

### tenant.created
Se publica cuando se crea un nuevo tenant.
```json
{
  "event_type": "tenant.created",
  "data": {
    "tenant_id": "acme-corp",
    "tenant_name": "ACME Corporation"
  }
}
```

### tenant.deleted
Se publica cuando se elimina un tenant.
```json
{
  "event_type": "tenant.deleted",
  "data": {
    "tenant_id": "acme-corp"
  }
}
```

## Desarrollo Local

### Con Docker Compose
```bash
cd proyecto-microservicios
docker-compose up tenant_api
```

El contenedor utilizará Air para hot reload.

### Sin Docker
```bash
# Instalar Air para hot reload
go install github.com/air-verse/air@latest

# Ejecutar con variables de entorno
export DB_HOST=localhost
export DB_USER=tenant
export DB_PASSWORD=tenantpass
export DB_NAME=tenant_db
export DB_PORT=5432
export RABBITMQ_HOST=localhost
export RABBITMQ_PORT=5672
export PORT=8082

air -c .air.toml
```

## Arquitectura

### 3 Capas

1. **Controller (HTTP Handler)**
   - Maneja requests HTTP
   - Valida entrada (DTOs)
   - Retorna respuestas HTTP

2. **Service (Business Logic)**
   - Implementa lógica de negocio
   - Coordina operaciones
   - Publica eventos a RabbitMQ

3. **Repository (Data Persistence)**
   - Acceso directo a base de datos con GORM
   - Operaciones CRUD en tabla `tenants`

## Base de Datos

### Tabla: tenants

```sql
CREATE TABLE tenants (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id VARCHAR(255) UNIQUE NOT NULL,
  tenant_name VARCHAR(255) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_tenant_id ON tenants(tenant_id);
```

## Integración con Otros Microservicios

Cuando se publica un evento `tenant.created`, otros microservicios (users-micro y product-micro) escuchan y crean automáticamente un esquema PostgreSQL con el nombre `tenant_{tenant_id}` para aislar datos por tenant.

Cuando se publica un evento `tenant.deleted`, los otros microservicios eliminan el esquema correspondiente.

## Rutas en Traefik

- `GET /api/tenants` → tenant-micro
- `POST /api/tenants` → tenant-micro
- `GET /api/tenants/{id}` → tenant-micro
- `DELETE /api/tenants/{id}` → tenant-micro

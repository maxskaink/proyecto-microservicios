# Shipping Microservice

Microservicio para la gestión de carritos, órdenes y envíos.

## Características

- **Gestión de Carrito**: Agregar, actualizar y eliminar productos del carrito de compras
- **Gestión de Órdenes**: Crear órdenes desde el carrito, actualizar estados (pending, paid, cancelled)
- **Gestión de Envíos**: Crear y gestionar envíos automáticamente cuando una orden es pagada
- **Multi-tenancy**: Soporte completo para múltiples tenants con schemas separados en PostgreSQL
- **Event-Driven**: Consume eventos de RabbitMQ para sincronizar información de productos, usuarios y tenants
- **Service Discovery**: Registro automático en Consul
- **API Gateway**: Integración con Kong para enrutamiento

## Arquitectura

```
shipping-micro/
├── cmd/server/           # Punto de entrada de la aplicación
├── internal/
│   ├── config/           # Configuración
│   ├── controllers/      # Controladores HTTP
│   ├── db/               # Acceso a base de datos
│   │   ├── models/       # Modelos de BD
│   │   ├── mappers/      # Conversión entre modelos y DTOs
│   │   ├── repositories/ # Repositorios
│   │   └── tenant/       # Gestión de schemas multi-tenant
│   ├── domain/           # Entidades de dominio y errores
│   ├── dto/              # Data Transfer Objects
│   ├── events/           # Handlers de eventos
│   ├── messaging/        # Integración con RabbitMQ
│   ├── middleware/       # Middlewares (tenant, auth, etc.)
│   ├── server/           # Configuración del servidor y rutas
│   │   └── discovery/    # Service Discovery (Consul)
│   └── services/         # Lógica de negocio
└── pkg/logger/           # Logger

```

## Entidades

### Cart Item
- Artículos en el carrito de un usuario
- Relación con Product y User

### Order
- Órdenes/Pedidos creados desde el carrito
- Estados: pending, paid, cancelled
- Contiene OrderItems

### Shipping
- Envíos creados automáticamente al pagar una orden
- Estados: pending, in_transit, delivered, cancelled
- Número de tracking único

## Eventos Consumidos

### product.created
```json
{
  "id": "uuid",
  "name": "Producto X",
  "description": "Descripción",
  "price": 99.99,
  "stock": 100
}
```

### user.created
```json
{
  "id": "uuid",
  "email": "user@example.com",
  "name": "Usuario"
}
```

### tenant.created
```json
{
  "tenant_id": "tenant_123",
  "name": "Tenant Name"
}
```

### tenant.deleted
```json
{
  "tenant_id": "tenant_123"
}
```

## Eventos Publicados

### order.created
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "total_price": 199.98,
  "status": "pending"
}
```

### order.status_changed
```json
{
  "id": "uuid",
  "status": "paid",
  "previous_status": "pending"
}
```

### shipping.created
```json
{
  "id": "uuid",
  "order_id": "uuid",
  "tracking_number": "TRACK123",
  "status": "pending"
}
```

## API Endpoints

### Cart

- `POST /cart/items` - Agregar producto al carrito
- `GET /cart` - Obtener carrito del usuario
- `PUT /cart/items/:id` - Actualizar cantidad de un item
- `DELETE /cart/items/:id` - Eliminar item del carrito
- `DELETE /cart` - Vaciar carrito

### Orders

- `POST /orders` - Crear orden desde carrito
- `GET /orders` - Listar órdenes del usuario
- `GET /orders/:id` - Obtener detalle de orden
- `PUT /orders/:id/status` - Actualizar estado de orden (Admin)

### Shipping

- `GET /shippings` - Listar envíos
- `GET /shippings/:id` - Obtener detalle de envío
- `GET /shippings/order/:order_id` - Obtener envío por orden
- `PUT /shippings/:id/status` - Actualizar estado de envío (Admin)

## Variables de Entorno

```env
PORT=8083
DATABASE_URL=postgres://shipping:shippingpass@shipping_db:5432/shipping_db?sslmode=disable
RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/
CONSUL_HTTP_ADDR=consul:8500
SERVICE_NAME=shipping-service
SERVICE_TAGS=shipping,api
```

## Ejecución

### Con Docker Compose
```bash
docker-compose up shipping_api
```

### Desarrollo Local
```bash
cd shipping-micro
air -c .air.toml
```

## Testing

```bash
go test ./...
```

## Dependencias

- Gin (Web Framework)
- GORM (ORM)
- PostgreSQL Driver
- RabbitMQ Client
- Consul Client
- Google UUID

## Multi-tenancy

Cada tenant tiene su propio schema en PostgreSQL. El middleware valida el header `X-Tenant-ID` en cada request y ejecuta todas las operaciones en el schema correspondiente.

## Integración

1. **RabbitMQ**: Escucha eventos de products, users y tenants
2. **Consul**: Se registra automáticamente para service discovery
3. **Kong**: El API Gateway enruta las peticiones basándose en el registro de Consul

## Flujo de Negocio

1. Usuario agrega productos al carrito (POST /cart/items)
2. Usuario confirma compra y crea una orden (POST /orders)
3. La orden se crea en estado "pending"
4. Admin/Sistema actualiza estado a "paid" (PUT /orders/:id/status)
5. Automáticamente se crea un envío cuando la orden es "paid"
6. Se puede rastrear el envío y actualizar su estado

## Próximos Pasos

- [ ] Implementar autenticación JWT
- [ ] Agregar tests unitarios e integración
- [ ] Implementar paginación en listados
- [ ] Agregar métricas y monitoreo
- [ ] Implementar circuit breakers
- [ ] Agregar cache con Redis

# Guía de Inicio Rápido - tenant-micro

## 🚀 Inicio Rápido

### 1. Clonar repositorio y posicionarse en la carpeta
```bash
cd /home/maxskaink/Universidad/Octavo/microservicios/proyecto-microservicios
```

### 2. Iniciar todos los servicios con Docker Compose
```bash
docker-compose up
```

Esto iniciará:
- ✅ API Gateway (Traefik)
- ✅ Consul (Service Discovery)
- ✅ RabbitMQ
- ✅ PostgreSQL para users-micro
- ✅ PostgreSQL para product-micro
- ✅ **PostgreSQL para tenant-micro** (puerto 5434)
- ✅ users-micro
- ✅ product-micro
- ✅ **tenant-micro** (puerto 8082)

### 3. Esperar a que todos los servicios estén listos

Puedes verificar el estado visitando:
- Traefik Dashboard: http://localhost:8090
- Consul UI: http://localhost:8500
- RabbitMQ Management: http://localhost:15672 (guest/guest)

### 4. Crear un tenant (ejemplo)

```bash
curl -X POST http://localhost/api/tenants \
  -H "Content-Type: application/json" \
  -d '{
    "tenant_id": "acme-corp",
    "tenant_name": "ACME Corporation"
  }'
```

**Respuesta esperada:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "tenant_id": "acme-corp",
  "tenant_name": "ACME Corporation",
  "created_at": "2024-01-15T10:30:00Z"
}
```

### 5. Verificar creación automática de esquemas

El evento `tenant.created` se publica automáticamente a RabbitMQ. Otros microservicios crean esquemas:
- En users-micro: `tenant_acme_corp`
- En product-micro: `tenant_acme_corp`

### 6. Listar todos los tenants

```bash
curl -X GET http://localhost/api/tenants
```

### 7. Obtener un tenant específico

```bash
curl -X GET http://localhost/api/tenants/acme-corp
```

### 8. Eliminar un tenant

```bash
curl -X DELETE http://localhost/api/tenants/acme-corp
```

---

## 🏗️ Estructura Arquitectónica

### Flujo de Multitenancia

```
┌─────────────────────────────────────────────┐
│        Traefik (API Gateway)                │
│  /{tenant}/api/users → X-Tenant-ID header   │
└─────────────────────────────────────────────┘
         ↓
┌─────────────────────────────────────────────┐
│  /api/tenants → tenant-micro                │
│  /api/users → users-micro                   │
│  /api/products → product-micro              │
└─────────────────────────────────────────────┘
         ↓
    RabbitMQ (Events)
    ├─ tenant.created → users-micro
    ├─ tenant.created → product-micro
    ├─ tenant.deleted → users-micro
    └─ tenant.deleted → product-micro
         ↓
PostgreSQL (Schema per Tenant)
├─ tenant_db (tenant-micro) - Metadata
├─ users_db with schema tenant_{id} (users-micro)
└─ product_db with schema tenant_{id} (product-micro)
```

---

## 📋 Endpoints de tenant-micro

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| POST | `/api/tenants` | Crear nuevo tenant |
| GET | `/api/tenants` | Listar todos los tenants |
| GET | `/api/tenants/{id}` | Obtener tenant por ID |
| DELETE | `/api/tenants/{id}` | Eliminar tenant |
| GET | `/health` | Health check |

---

## 🔧 Troubleshooting

### Error: "Tenant ya existe"
```
{
  "error": "unique constraint"
}
```
**Solución:** El `tenant_id` ya fue creado. Usa un identificador único.

### Error: "No se puede conectar a RabbitMQ"
```
Advertencia: No se pudo conectar a RabbitMQ
```
**Solución:** RabbitMQ puede no estar listo. Los eventos simplemente no se publican, pero el tenant se crea igual.

### Error: "Puerto 5434 ya en uso"
```
Error: cannot bind to port 5434
```
**Solución:** Cambia el puerto en `docker-compose.yml` o detén otros servicios.

---

## 🛠️ Desarrollo Local (Sin Docker)

### Prerequisitos
```bash
# Go 1.21+
# PostgreSQL 16+
# RabbitMQ 3.8+
```

### Instalación
```bash
cd tenant-micro
go mod download
```

### Variables de Entorno
Copia `.env.example` a `.env` y ajusta:
```bash
cp .env.example .env
```

### Ejecutar localmente
```bash
air -c .air.toml
```

El servidor iniciará en `http://localhost:8082`

---

## 📝 Características Principales

✅ **CRUD Completo**: Crear, leer, actualizar y eliminar tenants
✅ **Event-Driven**: Publica eventos para sincronización automática
✅ **Multitenancia**: Aislamiento de datos por tenant
✅ **Health Check**: Endpoint para verificar estado del servicio
✅ **Hot Reload**: Desarrollo con Air (cambios en tiempo real)
✅ **Docker Ready**: Contenedor listo para producción

---

## 📚 Documentación Completa

Ver `README.md` en este directorio para documentación detallada.

---

## 🤝 Integración con Otros Microservicios

### users-micro & product-micro

Ambos microservicios escuchan eventos de `tenant-micro`:

1. **tenant.created**: Crea automáticamente schema `tenant_{id}`
2. **tenant.deleted**: Elimina automáticamente schema `tenant_{id}`

Verificar listeners en:
- `users-micro/internal/messaging/handlers/tenant_event_handler.go`
- `product-micro/internal/messaging/handlers/tenant_event_handler.go`

---

## 📞 Contacto

Para preguntas o problemas, contacta al equipo de desarrollo.

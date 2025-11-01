# Configuración de Tenants - Traefik

Esta carpeta contiene configuraciones específicas de tenants para Traefik.

Aunque Traefik maneja dinámicamente los tenants mediante regex en las rutas,
este directorio puede usarse para:

1. Configuraciones personalizadas por tenant (rate limiting, autenticación, etc.)
2. Documentación de tenants disponibles
3. Ejemplos de uso

## Cómo funciona

Las rutas de tenant se configuran en `services.yml` usando patrones regex:

```yaml
routers:
  users-tenant:
    rule: "PathPrefix(`/{tenant:[a-zA-Z0-9_-]+}/api/users`)"
    middlewares:
      - tenant-extractor
```

Traefik automáticamente:
1. Captura el valor del tenant desde la ruta
2. Lo pasa al middleware tenant-extractor
3. El middleware lo agrega como header X-Tenant-ID
4. Reescribe la ruta a /api/users

## Ejemplo de Request

```bash
curl -X GET http://localhost/tenant-a/api/users

# Traefik hace:
# 1. Extrae: tenant = "tenant-a"
# 2. Reescribe: /api/users
# 3. Agrega header: X-Tenant-ID: tenant-a
# 4. Forwarda a: http://users_api:8080/api/users
```

## Tenants Disponibles (Ejemplo)

- `default` - Tenant por defecto
- `company_a` - Company A
- `company_b` - Company B

## Agregar Middleware Personalizado

Si necesitas rate limiting o autenticación por tenant,
puedes agregar configuración aquí.

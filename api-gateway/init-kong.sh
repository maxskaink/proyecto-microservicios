#!/bin/sh
# init-kong.sh
# Script de inicialización para Kong con integración a Consul

set -e

echo "Iniciando configuración de Kong con Consul..."

# Esperar a que Consul esté disponible
echo "Esperando a que Consul esté disponible..."
until curl -sf http://consul:8500/v1/status/leader > /dev/null 2>&1; do
  echo "Consul no está listo aún. Reintentando en 2 segundos..."
  sleep 2
done

echo "Consul está disponible"

# Verificar que los servicios están registrados en Consul
echo "Verificando servicios en Consul..."
SERVICES=$(curl -s http://consul:8500/v1/catalog/services | jq -r 'keys[]')
echo "Servicios registrados en Consul:"
echo "${SERVICES}"

# Registrar Kong en Consul
KONG_SERVICE_ID="kong-gateway"
KONG_SERVICE_NAME="kong"
KONG_PORT="8000"
KONG_HOST=$(hostname -i)

echo "Registrando Kong en Consul..."
curl -X PUT "http://consul:8500/v1/agent/service/register" \
  -H "Content-Type: application/json" \
  -d @- <<EOF
{
  "ID": "${KONG_SERVICE_ID}",
  "Name": "${KONG_SERVICE_NAME}",
  "Tags": ["api-gateway", "kong"],
  "Address": "${KONG_HOST}",
  "Port": ${KONG_PORT},
  "Check": {
    "HTTP": "http://${KONG_HOST}:${KONG_PORT}/",
    "Interval": "10s",
    "Timeout": "3s"
  }
}
EOF

if [ $? -eq 0 ]; then
  echo "Kong registrado exitosamente en Consul"
else
  echo "Error al registrar Kong en Consul, pero continuando..."
fi

# Verificar resolución DNS de servicios desde Consul
echo "Verificando resolución DNS de servicios..."
for service in users-service products-service tenant-service; do
  echo "Resolviendo ${service}.service.consul..."
  nslookup ${service}.service.consul consul || echo "No se pudo resolver ${service}"
done

echo "Configuración de Kong completada"
echo "Kong puede ahora descubrir servicios dinámicamente desde Consul"

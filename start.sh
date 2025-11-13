#!/bin/bash

set -euo pipefail

GREEN='\033[0;32m'
YELLOW='\033[0;33m'
RED='\033[0;31m'
NC='\033[0m'

PROJECT_NAME="proyecto-microservicios"
MODE=${1:-dev} # dev (por defecto) | prod

if [[ "$MODE" != "dev" && "$MODE" != "prod" ]]; then
  echo -e "${RED}Modo inválido: $MODE. Usa 'dev' o 'prod'.${NC}"; exit 1;
fi

echo -e "${YELLOW}Modo seleccionado: $MODE${NC}"

# Asegurar que el script se ejecuta desde la raíz del repo (donde está este script)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

echo -e "${GREEN}Preparando infraestructura de microservicios...${NC}"

# Comprobar que Docker está disponible
if ! command -v docker >/dev/null 2>&1; then
  echo -e "${RED}Docker no está instalado o no está en PATH. Instala Docker antes de continuar.${NC}"; exit 1
fi

HOST_UID=${HOST_UID:-$(id -u)}
HOST_GID=${HOST_GID:-$(id -g)}
export HOST_UID HOST_GID

echo -e "${YELLOW}HOST_UID y HOST_GID configurados: HOST_UID=$HOST_UID HOST_GID=$HOST_GID${NC}"

echo -e "${YELLOW}Deteniendo contenedores antiguos...${NC}"
docker compose down --remove-orphans || true

echo -e "${YELLOW}Eliminando contenedores detenidos del proyecto...${NC}"
docker ps -a --filter "network=${PROJECT_NAME}_micro-network" -q | xargs -r docker rm -f || true

echo -e "${YELLOW}Eliminando imágenes dangling...${NC}"
docker images -f "dangling=true" -q | xargs -r docker rmi || true

echo -e "${YELLOW}Eliminando red del proyecto si existe...${NC}"
docker network rm ${PROJECT_NAME}_micro-network 2>/dev/null || true

echo -e "${YELLOW}Limpiando redes huérfanas...${NC}"
docker network prune -f >/dev/null 2>&1 || true

echo -e "${GREEN}Iniciando servicios...${NC}"

COMPOSE_FILE="docker-compose.yml"
if [[ "$MODE" == "prod" ]]; then
  COMPOSE_FILE="docker-compose.prod.yml"
fi

echo -e "${YELLOW}Usando archivo: $COMPOSE_FILE${NC}"

# Verificar que el archivo de compose existe
if [[ ! -f "$COMPOSE_FILE" ]]; then
  echo -e "${RED}No se encontró $COMPOSE_FILE en $(pwd). Asegúrate de ejecutar el script desde la raíz del repo.${NC}"; exit 1
fi

# Si es modo prod, asegurarse de que las carpetas de contexto existen (advertencia, no fatal)
if [[ "$MODE" == "prod" ]]; then
  for d in users-micro product-micro shipping-micro tenant-micro pkg/observability api-gateway service-discovery; do
    if [[ ! -e "$d" ]]; then
      echo -e "${YELLOW}Advertencia: no se encontró '$d' en el repo. Si falta, el build puede fallar.${NC}"
    fi
  done
fi

# Intentar levantar. Si falla, intentar limpiar la red y reintentar una vez más.
if ! docker compose -f "$COMPOSE_FILE" up --build -d; then
  echo -e "${RED}⚠️ Error creando servicios en la primera tentativa. Intentando limpiar red y reintentar...${NC}"
  docker network rm ${PROJECT_NAME}_micro-network 2>/dev/null || true
  if ! docker compose -f "$COMPOSE_FILE" up --build -d; then
    echo -e "${RED}Fallo al levantar los servicios tras reintento. Mostrando últimos logs de compose para diagnóstico...${NC}"
    docker compose -f "$COMPOSE_FILE" ps || true
    docker compose -f "$COMPOSE_FILE" logs --no-color --timestamps --tail=200 || true
    exit 1
  fi
fi

echo -e "${YELLOW}Esperando inicio...${NC}"
sleep 8

echo -e "${GREEN}Servicios activos:${NC}"
docker compose -f "$COMPOSE_FILE" ps

echo -e "\n${GREEN}✅ Infraestructura levantada (${MODE})${NC}"
echo -e "- Consul: http://localhost:8500/ui/"
echo -e "- Kong Admin: http://localhost:8090"
echo -e "- RabbitMQ: http://localhost:15672 (guest/guest)"
echo -e "- Prometheus: http://localhost:9090"
echo -e "- Grafana: http://localhost:3000 (admin/admin)"
echo -e "- API Products: http://localhost/api/products"
echo -e "- API Products (tenant): http://localhost/tenant1/api/products"

#!/bin/bash

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

echo -e "${GREEN}Preparando infraestructura de microservicios...${NC}"

# Exportar UID/GID sin sobrescribir si existen
export UID=$(id -u)
export GID=$(id -g)

echo -e "${YELLOW}UID y GID configurados: UID=$UID GID=$GID${NC}"

echo -e "${YELLOW}Deteniendo contenedores antiguos...${NC}"
docker compose down --remove-orphans

echo -e "${YELLOW}Eliminando contenedores detenidos del proyecto...${NC}"
docker ps -a --filter "network=${PROJECT_NAME}_micro-network" -q | xargs -r docker rm -f

echo -e "${YELLOW}Eliminando imágenes dangling...${NC}"
docker images -f "dangling=true" -q | xargs -r docker rmi

echo -e "${YELLOW}Eliminando red del proyecto si existe...${NC}"
docker network rm ${PROJECT_NAME}_micro-network 2>/dev/null

echo -e "${YELLOW}Limpiando redes huérfanas...${NC}"
docker network prune -f >/dev/null 2>&1

echo -e "${GREEN}Iniciando servicios...${NC}"

COMPOSE_FILE="docker-compose.yml"
if [[ "$MODE" == "prod" ]]; then
  COMPOSE_FILE="docker-compose.prod.yml"
fi

echo -e "${YELLOW}Usando archivo: $COMPOSE_FILE${NC}"

if ! docker compose -f "$COMPOSE_FILE" up --build -d; then
  echo -e "${RED}⚠️ Error creando red. Intentando limpiar...${NC}"
  docker network rm ${PROJECT_NAME}_micro-network 2>/dev/null
  docker compose -f "$COMPOSE_FILE" up --build -d
fi

echo -e "${YELLOW}Esperando inicio...${NC}"
sleep 8

echo -e "${GREEN}Servicios activos:${NC}"
docker compose -f "$COMPOSE_FILE" ps

echo -e "\n${GREEN}✅ Infraestructura levantada (${MODE})${NC}"
echo -e "- Consul: http://localhost:8500/ui/"
echo -e "- Kong Admin: http://localhost:8090"
echo -e "- RabbitMQ: http://localhost:15672 (guest/guest)"
echo -e "- API Products: http://localhost/api/products"
echo -e "- API Products (tenant): http://localhost/tenant1/api/products"

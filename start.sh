#!/bin/bash

# Colores para mejorar la salida
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

PROJECT_NAME="proyecto-microservicios"

echo -e "${GREEN}Preparando infraestructura de microservicios...${NC}"

# Exportar UID y GID para evitar warnings en Docker
export UID=$(id -u)
export GID=$(id -g)
echo -e "${YELLOW}UID y GID configurados: UID=$UID GID=$GID${NC}"

echo -e "${YELLOW}Deteniendo contenedores antiguos...${NC}"
docker-compose down --remove-orphans

echo -e "${YELLOW}Eliminando contenedores detenidos...${NC}"
docker rm $(docker ps -aq) 2>/dev/null

echo -e "${YELLOW}Eliminando imágenes dangling...${NC}"
docker rmi $(docker images -f "dangling=true" -q) 2>/dev/null

echo -e "${YELLOW}Eliminando red del proyecto si existe...${NC}"
docker network rm ${PROJECT_NAME}_micro-network 2>/dev/null

echo -e "${YELLOW}Limpiando redes huérfanas...${NC}"
docker network prune -f >/dev/null 2>&1

echo -e "${YELLOW}Verificando redes del proyecto...${NC}"
if docker network inspect proyecto-microservicios_micro-network >/dev/null 2>&1; then
  echo -e "${RED}Red encontrada, eliminando...${NC}"
  docker network rm proyecto-microservicios_micro-network >/dev/null 2>&1
fi

echo -e "${GREEN}Iniciando los servicios...${NC}"
docker-compose up --build -d

echo -e "${YELLOW}Esperando a que los servicios arranquen...${NC}"
sleep 10

echo -e "${GREEN}Servicios en ejecución:${NC}"
docker-compose ps

echo -e "\n${GREEN}Acceso a los servicios:${NC}"
echo -e "- API Gateway: ${YELLOW}http://localhost${NC}"
echo -e "- Kong Admin API: ${YELLOW}http://localhost:8090${NC}"
echo -e "- UI de Consul: ${YELLOW}http://localhost:8500/ui/${NC}"
echo -e "- API Usuarios: ${YELLOW}http://localhost/api/users${NC}"
echo -e "- API Usuarios (tenant): ${YELLOW}http://localhost/tenant1/api/users${NC}"
echo -e "- API Productos: ${YELLOW}http://localhost/api/products${NC}"
echo -e "- API Productos (tenant): ${YELLOW}http://localhost/tenant1/api/products${NC}"
echo -e "- RabbitMQ: ${YELLOW}http://localhost:15672${NC} (guest / guest)"

echo -e "\n${GREEN}✅ Infraestructura iniciada correctamente!${NC}"

#!/bin/bash

# Colores para mejorar la salida
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Iniciando la infraestructura de microservicios...${NC}"

# Definir variables de entorno para docker-compose
export UID=$(id -u)
export GID=$(id -g)

# Iniciar los servicios
echo -e "${YELLOW}Iniciando Docker Compose...${NC}"
docker-compose up -d

# Esperar a que los servicios estén disponibles
echo -e "${YELLOW}Esperando a que los servicios estén disponibles...${NC}"
sleep 10

# Mostrar los contenedores en ejecución
echo -e "${GREEN}Servicios en ejecución:${NC}"
docker-compose ps

# Mostrar información de acceso
echo -e "\n${GREEN}Acceso a los servicios:${NC}"
echo -e "- API Gateway: ${YELLOW}http://localhost${NC}"
echo -e "- Dashboard Traefik: ${YELLOW}http://localhost:8090/dashboard/${NC}"
echo -e "- UI de Consul: ${YELLOW}http://localhost:8500/ui/${NC}"
echo -e "- API de usuarios: ${YELLOW}http://localhost/api/users${NC} (a través del Gateway)"
echo -e "- API de usuarios (directo): ${YELLOW}http://localhost:8090${NC} (desarrollo)"
echo -e "- rabbitmq Management: ${YELLOW}http://localhost:15672${NC} (usuario: guest, contraseña: guest)"

echo -e "\n${GREEN}¡Infraestructura iniciada correctamente!${NC}"

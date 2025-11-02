# Usar imagen oficial de Kong OSS
FROM kong:3.4

# Cambiar a usuario root temporalmente para instalar dependencias
USER root

# Instalar curl, jq y herramientas DNS necesarias para Consul
RUN apt-get update && apt-get install -y curl jq dnsutils && rm -rf /var/lib/apt/lists/*

# Crear directorios necesarios
RUN mkdir -p /kong /usr/local/share/lua/5.1/kong/plugins/tenant-extractor

# Copiar configuración de Kong
COPY kong/kong.yml /kong/kong.yml

# Copiar plugin Lua
COPY plugins/tenant-extractor /usr/local/share/lua/5.1/kong/plugins/tenant-extractor

# Copiar script de inicialización
COPY init-kong.sh /docker-entrypoint.d/init-kong.sh
RUN chmod +x /docker-entrypoint.d/init-kong.sh

# Volver al usuario kong
USER kong

EXPOSE 8000 8001 8443 8444

# El entrypoint por defecto de Kong se encarga de ejecutar los scripts en /docker-entrypoint.d/

#!/bin/sh
set -e

: "${API_URL:=http://localhost:80}"

# Ensure config dir exists inside the built dist and write runtime config
# Angular genera artefactos navegables en dist/web-microservicios/browser
mkdir -p /app/dist/web-microservicios/browser/assets/config
cat > /app/dist/web-microservicios/browser/assets/config/config.json <<EOF
{
  "apiUrl": "${API_URL}"
}
EOF

# SPA fallback: servir index.html también como 404.html para rutas como /home
if [ -f /app/dist/web-microservicios/browser/index.html ]; then
  cp /app/dist/web-microservicios/browser/index.html /app/dist/web-microservicios/browser/404.html
fi


# Exec the CMD (serve the static files)
exec "$@"

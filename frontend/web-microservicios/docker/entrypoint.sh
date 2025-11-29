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

# Exec the CMD (serve the static files)
exec "$@"

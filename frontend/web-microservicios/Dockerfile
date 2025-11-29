FROM node:20
WORKDIR /app

# Install dependencies
COPY package*.json ./
RUN npm ci --no-audit --no-fund

# Copy source and build
COPY . .
# Build using project's configured outputPath (dist/web-microservicios)
RUN npm run build -- --configuration=production

# Install a minimal static server
RUN npm i -g http-server --no-audit --no-fund

# Copy entrypoint that writes runtime config into the built assets
COPY docker/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 4200

ENTRYPOINT ["/entrypoint.sh"]
 CMD ["http-server", "/app/dist/web-microservicios/browser", "-p", "4200", "--cors"]

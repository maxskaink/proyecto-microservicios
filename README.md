# proyecto-microservicios
Modular microservices-based application enabling independent development, deployment, and scaling of services. Each service manages a specific business function, ensuring flexibility, resilience, and streamlined integration in distributed systems.

## Despliegue local (Postgres + users-api)

1. Copia el archivo de variables de entorno:
	- `cp .env.example .env`
2. Levanta los servicios con Docker Compose (desde la raíz del repo):
	- `docker compose -f deploy/docker-compose.yml --env-file .env up --build`
3. Prueba el endpoint de salud del micro de usuarios:
	- `curl http://localhost:$PORT/health`

Variables clave (puedes editarlas en `.env`):
- `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`, `POSTGRES_PORT`
- `PORT` (puerto del micro)
- `DATABASE_URL` (ej: `postgres://users:userspass@db:5432/users?sslmode=disable`)
- `FIREBASE_PROJECT_ID`

Notas:
- El compose construye la imagen usando `users-micro/Dockerfile`.
- El contenedor `db` expone Postgres en el puerto `${POSTGRES_PORT}` local.

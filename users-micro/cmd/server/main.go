package main

import (
	"log"

	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/server"
)

// @title           API Usuarios
// @version         1.0
// @description     Microservicio para la gestión de usuarios.
// @host            localhost:8080
// @BasePath        /
func main() {
	// Inicializa y levanta el servidor HTTP (solo estructura, sin implementación específica).
	if err := server.Run(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

package main

import (
	"log"

	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/server"
)

// Entry point del microservicio de usuarios.
func main() {
	// Inicializa y levanta el servidor HTTP (solo estructura, sin implementación específica).
	if err := server.Run(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

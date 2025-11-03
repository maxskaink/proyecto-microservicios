package main

import (
	"log"

	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/server"
)

// @title           API Shipping
// @version         1.0
// @description     Microservicio para la gestión de envíos, carritos y órdenes.
// @host            localhost:8080
// @BasePath        /
func main() {
	// Inicializa y levanta el servidor HTTP
	if err := server.Run(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/maxskaink/proyecto-microservicios/tenant-micro/internal/controllers"
	"github.com/maxskaink/proyecto-microservicios/tenant-micro/internal/discovery"
	"github.com/maxskaink/proyecto-microservicios/tenant-micro/internal/models"
	"github.com/maxskaink/proyecto-microservicios/tenant-micro/internal/repositories"
	"github.com/maxskaink/proyecto-microservicios/tenant-micro/internal/services"
	"github.com/streadway/amqp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var serviceRegistry *discovery.ServiceRegistration

func main() {
	godotenv.Load()

	// Conectar a BD
	db := connectDB()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// Migrar modelos
	db.AutoMigrate(&models.Tenant{})

	// Conectar a RabbitMQ
	channel := connectRabbitMQ()
	defer func() {
		if channel != nil {
			channel.Close()
		}
	}()

	// Crear capas
	repo := repositories.NewTenantRepository(db)
	service := services.NewTenantService(repo, channel)
	controller := controllers.NewTenantController(service)

	// Configurar service discovery
	setupServiceDiscovery()

	// Configurar Gin
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	// Rutas
	r.POST("/tenants", controller.CreateTenant)
	r.GET("/tenants", controller.GetAllTenants)
	r.GET("/tenants/:id", controller.GetTenant)
	r.DELETE("/tenants/:id", controller.DeleteTenant)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Manejo de señales para cierre controlado
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Iniciar servidor en una goroutine
	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8082"
		}
		addr := fmt.Sprintf(":%s", port)
		log.Printf("Tenant microservice iniciado en puerto %s\n", port)
		if err := r.Run(addr); err != nil && err != http.ErrServerClosed {
			log.Printf("Error al iniciar servidor: %v", err)
		}
	}()

	// Esperar señal de cierre
	<-quit
	log.Println("Cerrando servidor...")

	// Desregistrar del service discovery
	cleanup()
}

func cleanup() {
	log.Println("Iniciando limpieza de recursos...")

	// Desregistrar el servicio
	if serviceRegistry != nil {
		serviceRegistry.Deregister()
	}

	log.Println("Limpieza finalizada")
}

func connectDB() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error conectando a BD:", err)
	}

	return db
}

func connectRabbitMQ() *amqp.Channel {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		os.Getenv("RABBITMQ_USER"),
		os.Getenv("RABBITMQ_PASSWORD"),
		os.Getenv("RABBITMQ_HOST"),
		os.Getenv("RABBITMQ_PORT"),
	)

	conn, err := amqp.Dial(url)
	if err != nil {
		log.Println("Advertencia: No se pudo conectar a RabbitMQ:", err)
		return nil
	}

	channel, err := conn.Channel()
	if err != nil {
		log.Println("Advertencia: No se pudo crear canal RabbitMQ:", err)
		return nil
	}

	// Declarar exchange
	channel.ExchangeDeclare("tenants_events", "topic", true, false, false, false, nil)

	return channel
}

func setupServiceDiscovery() {
	var err error
	serviceRegistry, err = discovery.NewServiceRegistration()
	if err != nil {
		log.Printf("Error al configurar service discovery: %v", err)
		return
	}

	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = "tenant-service"
	}

	port := os.Getenv("PORT")
	portNum, _ := strconv.Atoi(port)
	if portNum == 0 {
		portNum = 8082
	}

	// Preparar etiquetas para Traefik
	tags := []string{"tenant", "api"}
	if os.Getenv("TRAEFIK_ENABLE") != "false" {
		tags = append(tags,
			"traefik.enable=true",
			"traefik.http.routers.tenants.rule=PathPrefix(`/api/tenants`)",
			"traefik.http.routers.tenants.entrypoints=web",
			"traefik.http.middlewares.strip-prefix.stripprefix.prefixes=/api",
			"traefik.http.routers.tenants.middlewares=strip-prefix",
		)
	}

	// Registrar el servicio
	err = serviceRegistry.Register(
		fmt.Sprintf("%s-%d", serviceName, portNum),
		serviceName,
		portNum,
		tags,
	)

	if err != nil {
		log.Printf("Error al registrar el servicio en Consul: %v", err)
	} else {
		log.Println("Servicio registrado correctamente en Consul")
	}
}

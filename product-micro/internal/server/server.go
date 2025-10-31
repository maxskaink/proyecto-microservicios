package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/db"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/db/repositories"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/events"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/server/discovery"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/server/validators"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/services"
	"github.com/maxskaink/proyecto-microservicios/product-micro/pkg/logger"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Services
var ProductService services.IProductService
var UserService services.IUserService
var serviceRegistry *discovery.ServiceRegistration
var userEventHandler *events.UserHandler

// Repositories
var ProductRepository repositories.IProductRepository
var UserRepository repositories.IUserRepository

// Midddlewares
var AuthMiddleware gin.HandlerFunc

// Run arranca el servidor HTTP con Gin.
// Solo registra una ruta de salud para validar que el contenedor responde.
func Run() error {

	configDB()
	configServices()
	InitFirebase()
	validators.RegisterValidators()

	//Configurar service discovery
	setupServiceDiscovery()
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())     //For logs request
	r.Use(CORSMiddleware()) //For manage the cors

	// Healthcheck básico
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Registrar rutas de controladores (usuarios, pedidos, perfiles)
	RegisterRoutes(r)

	// Manejo de señales para cierre controlado
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Iniciar servidor en una goroutine
	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		addr := fmt.Sprintf(":%s", port)
		logger.Info(fmt.Sprintf("Servidor iniciado en %s", addr))
		if err := r.Run(addr); err != nil && err != http.ErrServerClosed {
			logger.Error(fmt.Sprintf("Error al iniciar servidor: %v", err))
		}
	}()

	// Esperar señal de cierre
	<-quit
	logger.Info("Cerrando servidor...")

	// Ejecutar limpieza
	cleanup()

	return nil
}

func configDB() {
	providerDB, err := db.NewGormDBProvider()

	if nil != err {
		log.Fatal(err.Error())
		return
	}

	DB, _ := providerDB.DB(&gin.Context{})

	ProductRepository = repositories.NewProductRepository(DB)
	UserRepository = repositories.NewUserRepository(DB)
}

func configServices() {
	UserService = services.NewUserService(UserRepository)
	ProductService = services.NewProductService(ProductRepository, UserService)

	// Configurar el manejador de eventos de usuarios
	var err error
	userEventHandler, err = events.NewUserHandler(UserRepository)
	if err != nil {
		logger.Error(fmt.Sprintf("Error al crear manejador de eventos de usuarios: %v", err))
		return
	}

	// Iniciar el consumo de eventos en una goroutine
	go func() {
		ctx := context.Background()
		if err := userEventHandler.Start(ctx); err != nil {
			logger.Error(fmt.Sprintf("Error al iniciar consumo de eventos: %v", err))
		}
	}()

	logger.Info("Manejador de eventos de usuarios iniciado correctamente")
}

func setupServiceDiscovery() {
	var err error
	serviceRegistry, err = discovery.NewServiceRegistration()
	if err != nil {
		logger.Error(fmt.Sprintf("Error al configurar service discovery: %v", err))
		return
	}

	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = "users-service"
	}

	port := os.Getenv("PORT")
	portNum, _ := strconv.Atoi(port)
	if portNum == 0 {
		portNum = 8080
	}

	// Preparar etiquetas para Traefik
	tags := []string{"products"}
	if os.Getenv("TRAEFIK_ENABLE") != "false" {
		tags = append(tags,
			"traefik.enable=true",
			"traefik.http.routers.products.rule=PathPrefix(`/api/products`)",
			"traefik.http.routers.products.entrypoints=web",
			"traefik.http.middlewares.strip-prefix.stripprefix.prefixes=/api",
			"traefik.http.routers.products.middlewares=strip-prefix",
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
		logger.Error(fmt.Sprintf("Error al registrar el servicio en Consul: %v", err))
	} else {
		logger.Info("Servicio registrado correctamente en Consul")
	}
}

func cleanup() {
	// Desregistrar el servicio de Consul
	if serviceRegistry != nil {
		if err := serviceRegistry.Deregister(); err != nil {
			logger.Error(fmt.Sprintf("Error al desregistrar el servicio: %v", err))
		} else {
			logger.Info("Servicio desregistrado correctamente")
		}
	}

	// Cerrar manejador de eventos de usuarios
	if userEventHandler != nil {
		if err := userEventHandler.Close(); err != nil {
			logger.Error(fmt.Sprintf("Error al cerrar manejador de eventos de usuarios: %v", err))
		} else {
			logger.Info("Manejador de eventos de usuarios cerrado correctamente")
		}
	}
}

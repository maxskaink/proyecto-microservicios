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
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/db/tenant"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/messaging"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/server/discovery"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/server/validators"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/services"
	tenant_services "github.com/maxskaink/proyecto-microservicios/product-micro/internal/services/tenant"
	"github.com/maxskaink/proyecto-microservicios/product-micro/pkg/logger"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

// Services
var ProductService services.IProductService
var UserService services.IUserService
var serviceRegistry *discovery.ServiceRegistration
var eventManager *messaging.EventManager
var tenantService *tenant_services.TenantService

// Repositories
var ProductRepository repositories.IProductRepository
var UserRepository repositories.IUserRepository

// Midddlewares
var AuthMiddleware gin.HandlerFunc

// Database
var DB *gorm.DB

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

	DB, _ = providerDB.DB(&gin.Context{})
	tenantDB := tenant.NewTenantDB(DB)
	tenantService = tenant_services.NewTenantService(DB)

	ProductRepository = repositories.NewProductRepository(DB, tenantDB)
	UserRepository = repositories.NewUserRepository(DB, tenantDB)
}

func configServices() {
	UserService = services.NewUserService(UserRepository)
	// Crear publisher de eventos de productos
	productPublisher, err := messaging.NewProductEventPublisher()
	if err != nil {
		logger.Error(fmt.Sprintf("Error al crear publisher de productos: %v", err))
	}
	ProductService = services.NewProductService(ProductRepository, UserService, productPublisher)

	// Configurar el gestor de eventos (consumidor de RabbitMQ)
	eventManager, err = messaging.NewEventManager(UserRepository, tenantService)
	if err != nil {
		logger.Error(fmt.Sprintf("Error al crear gestor de eventos: %v", err))
		return
	}

	// Iniciar el consumo de eventos en una goroutine
	go func() {
		ctx := context.Background()
		if err := eventManager.Start(ctx); err != nil {
			logger.Error(fmt.Sprintf("Error al iniciar consumo de eventos: %v", err))
		}
	}()

	logger.Info("Gestor de eventos iniciado correctamente")
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

	// Cerrar gestor de eventos
	if eventManager != nil {
		if err := eventManager.Close(); err != nil {
			logger.Error(fmt.Sprintf("Error al cerrar gestor de eventos: %v", err))
		} else {
			logger.Info("Gestor de eventos cerrado correctamente")
		}
	}
}

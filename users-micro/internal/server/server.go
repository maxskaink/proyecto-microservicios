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
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/db"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/db/repositories"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/db/tenant"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/messaging"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/messaging/handlers"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/messaging/rabbitmq"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/middleware"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/server/discovery"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/services"
	tenant_services "github.com/maxskaink/proyecto-microservicios/users-micro/internal/services/tenant"
	"github.com/maxskaink/proyecto-microservicios/users-micro/pkg/logger"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

// Services
var UserService services.UserService
var TenantService *tenant_services.TenantService
var serviceRegistry *discovery.ServiceRegistration
var msgPublisher messaging.Publisher
var eventManager *messaging.EventManager

// Repositories
var UserRepository repositories.UserRepository

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

	//Configurar service discovery
	setupServiceDiscovery()
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())     //For logs request
	r.Use(CORSMiddleware()) //For manage the cors

	// Healthcheck básico (sin middleware de tenant, es público)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Registrar middleware de tenant para rutas protegidas
	// Las rutas de administración (sin tenant) se registrarán después
	protectedRoutes := r.Group("/api")
	protectedRoutes.Use(middleware.TenantMiddleware(DB))

	r.GET("/users/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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

func cleanup() {
	logger.Info("Iniciando limpieza de recursos...")

	// Cerrar gestor de eventos
	if eventManager != nil {
		if err := eventManager.Close(); err != nil {
			logger.Error(fmt.Sprintf("Error al cerrar gestor de eventos: %v", err))
		} else {
			logger.Info("Gestor de eventos cerrado correctamente")
		}
	}

	// Cerrar publicador de mensajes
	if msgPublisher != nil {
		if err := msgPublisher.Close(); err != nil {
			logger.Error(fmt.Sprintf("Error al cerrar publicador de mensajes: %v", err))
		} else {
			logger.Info("Publicador de mensajes cerrado correctamente")
		}
	}

	// Desregistrar el servicio
	if serviceRegistry != nil {
		serviceRegistry.Deregister()
	}

	logger.Info("Limpieza finalizada")
}

func configDB() {
	providerDB, err := db.NewGormDBProvider()

	if nil != err {
		log.Fatal(err.Error())
		return
	}

	DB, _ = providerDB.DB(&gin.Context{})
	tenantDB := tenant.NewTenantDB(DB)

	UserRepository = repositories.NewUserRepository(DB, *tenantDB)
}

func configServices() {
	// Obtener DB provider para TenantService
	providerDB, err := db.NewGormDBProvider()
	if err != nil {
		logger.Error(fmt.Sprintf("Error al obtener DB provider para TenantService: %v", err))
		return
	}
	dbConn, _ := providerDB.DB(&gin.Context{})

	// Inicializar el publicador de mensajes
	factory := messaging.NewFactory(rabbitmq.DefaultConfig())
	msgPublisher, err = factory.CreatePublisher()
	if err != nil {
		logger.Error(fmt.Sprintf("Error al crear publicador de mensajes: %v", err))
		msgPublisher = nil
	} else {
		logger.Info("Publicador de mensajes inicializado correctamente")
	}

	// Crear TenantService
	TenantService = tenant_services.NewTenantService(dbConn, msgPublisher)
	logger.Info("TenantService inicializado correctamente")

	UserService = services.NewUserService(UserRepository, msgPublisher)

	// Inicializar el gestor de eventos (consumidor de RabbitMQ)
	// Crear el dispatcher con handlers
	dispatcher := buildEventDispatcher()

	eventManager, err = messaging.NewEventManager(dispatcher)
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
	tags := []string{"users"}
	if os.Getenv("TRAEFIK_ENABLE") != "false" {
		tags = append(tags,
			"traefik.enable=true",
			"traefik.http.routers.users.rule=PathPrefix(`/api/users`)",
			"traefik.http.routers.users.entrypoints=web",
			"traefik.http.middlewares.strip-prefix.stripprefix.prefixes=/api",
			"traefik.http.routers.users.middlewares=strip-prefix",
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

// buildEventDispatcher construye el dispatcher de eventos
func buildEventDispatcher() messaging.Dispatcher {
	return handlers.NewEventDispatcher(TenantService)
}

package server

import (
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
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/server/discovery"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/services"
	"github.com/maxskaink/proyecto-microservicios/users-micro/pkg/logger"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Services
var UserService services.UserService
var serviceRegistry *discovery.ServiceRegistration

// Repositories
var UserRepository repositories.UserRepository

// Midddlewares
var AuthMiddleware gin.HandlerFunc

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

	UserRepository = repositories.NewUserRepository(DB)
}

func configServices() {
	UserService = services.NewUserService(UserRepository)
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

func cleanup() {
	// Desregistrar el servicio de Consul
	if serviceRegistry != nil {
		if err := serviceRegistry.Deregister(); err != nil {
			logger.Error(fmt.Sprintf("Error al desregistrar el servicio: %v", err))
		} else {
			logger.Info("Servicio desregistrado correctamente")
		}
	}
}

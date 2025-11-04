package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/controllers"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/db"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/db/repositories"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/db/tenant"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/messaging"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/messaging/rabbitmq"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/middleware"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/server/discovery"
	cart "github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/services/cart"
	order "github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/services/order"
	shipping "github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/services/shipping"
	tenant_services "github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/services/tenant"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/pkg/logger"
	"gorm.io/gorm"
)

var (
	DB              *gorm.DB
	serviceRegistry *discovery.ServiceRegistration
	tenantService   *tenant_services.TenantService
	tenantDB        *tenant.TenantDB

	// repositorios
	productRepo repositories.IProductRepository
	userRepo    repositories.IUserRepository
	cartRepo    repositories.ICartRepository
	orderRepo   repositories.IOrderRepository
	shipRepo    repositories.IShippingRepository

	// servicios
	cartService     *cart.Service
	orderService    *order.Service
	shippingService *shipping.Service

	// eventos
	eventManager *messaging.EventManager
	publisher    *rabbitmq.Publisher
)

func Run() error {
	configDB()
	configServices()
	setupServiceDiscovery()
	InitFirebase()

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(CORSMiddleware())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	RegisterRoutes(r)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8083"
		}
		addr := fmt.Sprintf(":%s", port)
		logger.Info(fmt.Sprintf("Servidor iniciado en %s", addr))
		if err := r.Run(addr); err != nil && err != http.ErrServerClosed {
			logger.Error(fmt.Sprintf("Error al iniciar servidor: %v", err))
		}
	}()

	<-quit
	logger.Info("Cerrando servidor...")
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
	tenantDB = tenant.NewTenantDB(DB)
	tenantService = tenant_services.NewTenantService(DB)

	logger.Info("Base de datos configurada correctamente")
}

func configServices() {
	// Inicializar repositorios
	productRepo = repositories.NewProductRepository(DB, tenantDB)
	userRepo = repositories.NewUserRepository(DB, tenantDB)
	cartRepo = repositories.NewCartRepository(DB, tenantDB)
	orderRepo = repositories.NewOrderRepository(DB, tenantDB)
	shipRepo = repositories.NewShippingRepository(DB, tenantDB)

	// Inicializar publisher
	cm, err := rabbitmq.NewConnectionManager(rabbitmq.DefaultConfig())
	if err != nil {
		logger.Error(fmt.Sprintf("Error creando conexión RabbitMQ: %v", err))
	} else {
		publisher, err = rabbitmq.NewPublisher(cm)
		if err != nil {
			logger.Error(fmt.Sprintf("Error creando publisher: %v", err))
		}
	}

	// Inicializar servicios
	cartService = cart.NewService(cartRepo, productRepo)
	shippingService = shipping.NewService(shipRepo)
	orderService = order.NewService(cartRepo, productRepo, orderRepo, shipRepo, publisher, userRepo)

	// Inicializar consumidor de eventos RabbitMQ
	var emErr error
	eventManager, emErr = messaging.NewEventManager(productRepo, userRepo, tenantService)
	if emErr != nil {
		logger.Error(fmt.Sprintf("Error al crear gestor de eventos: %v", emErr))
	} else {
		go func() { _ = eventManager.Start(&gin.Context{}) }()
		logger.Info("Gestor de eventos iniciado correctamente")
	}

	logger.Info("Servicios y repositorios configurados correctamente")
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
		serviceName = "shipping-service"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	portInt := 8083
	fmt.Sscanf(port, "%d", &portInt)

	tags := []string{"shipping", "api"}
	serviceID := fmt.Sprintf("%s-%d", serviceName, portInt)

	if err := serviceRegistry.Register(serviceName, serviceID, portInt, tags); err != nil {
		logger.Error(fmt.Sprintf("Error al registrar servicio: %v", err))
		return
	}

	logger.Info(fmt.Sprintf("Servicio registrado en Consul: %s", serviceID))
}

func RegisterRoutes(r *gin.Engine) {
	tenantMW := middleware.TenantMiddleware(DB)

	api := r.Group("/api/")
	api.Use(middleware.CORSMiddleware(), tenantMW, middleware.FirebaseAuthMiddleware(userRepo, FirebaseApp))
	{
		// Controladores
		cartCtrl := controllers.NewCartController(cartService)
		orderCtrl := controllers.NewOrderController(orderService)
		shipCtrl := controllers.NewShippingController(shippingService)

		cartCtrl.Register(api)
		orderCtrl.Register(api)
		shipCtrl.Register(api)
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Tenant-ID")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func cleanup() {
	if eventManager != nil {
		_ = eventManager.Close()
	}
	if publisher != nil {
		_ = publisher.Close()
	}
	if serviceRegistry != nil {
		serviceName := os.Getenv("SERVICE_NAME")
		if serviceName == "" {
			serviceName = "shipping-service"
		}
		port := os.Getenv("PORT")
		if port == "" {
			port = "8083"
		}
		portInt := 8083
		fmt.Sscanf(port, "%d", &portInt)
		serviceID := fmt.Sprintf("%s-%d", serviceName, portInt)
		serviceRegistry.Deregister(serviceID)
	}

	logger.Info("Limpieza completada")
}

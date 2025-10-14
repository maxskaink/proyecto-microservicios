package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/db"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/db/repositories"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/services"
)

// Services
var UserService services.UserService

// Repositories
var UserRepository repositories.UserRepository

// Midddlewares
var AuthMiddleware gin.HandlerFunc

// Run arranca el servidor HTTP con Gin.
// Solo registra una ruta de salud para validar que el contenedor responde.
func Run() error {

	configDB()
	configServices()
	InitFirebase() //Connect with firebase

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())     //For logs request
	r.Use(CORSMiddleware()) //For manage the cors

	// Healthcheck básico
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	//r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Registrar rutas de controladores (usuarios, pedidos, perfiles)
	RegisterRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)
	return r.Run(addr)
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

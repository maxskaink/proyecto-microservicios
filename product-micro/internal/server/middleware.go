package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/services"
	"github.com/maxskaink/proyecto-microservicios/product-micro/pkg/logger"
)

func CORSMiddleware() gin.HandlerFunc {
	cfg := cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	// Could be important read this config from a env file
	return cors.New(cfg)
}

func FirebaseAuthMiddleware(userService services.IUserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			return
		}

		idToken := strings.TrimPrefix(authHeader, "Bearer ")

		client, err := FirebaseApp.Auth(context.Background())
		if err != nil {
			fmt.Printf("Token asked: %s\n", idToken)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Auth client error"})
			return
		}

		token, err := client.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		email, _ := token.Claims["email"].(string)
		name := strings.Split(email, "@")[0]

		_, err = userService.GetUserByUUID(token.UID)

		if err != nil {
			logger.Error("User with uid " + token.UID + " no existe en la bd de productos")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "No se ha encontrado el usuario, reintentelo mas tarde"})
		}

		c.Set("name", name)
		c.Set("uid", token.UID)
		c.Set("email", email)

		c.Next()
	}
}

package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/dto"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/services"
	"github.com/maxskaink/proyecto-microservicios/users-micro/pkg/logger"
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

func FirebaseAuthMiddleware(userService services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantId := c.GetHeader("X-Tenant-ID")

		if tenantId == "" {
			logger.Error("Missing X-Tenant-ID")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Missing X-Tenant-ID"})
			return
		}

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

		user, err := userService.GetUserByUUID(token.UID, tenantId)

		if err != nil {
			// Si el usuario no existe, crearlo
			newUser := dto.UserRequest{
				FirebaseUID: token.UID,
				Email:       email,
				Name:        name,
			}
			fmt.Printf("Creando usuario con email %s, y nombre %s \n", email, name)
			_, err := userService.CreateUser(newUser, tenantId)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el usuario" + err.Error()})
				return
			}
		} else {
			name = user.Name
		}

		c.Set("name", name)
		c.Set("uid", token.UID)
		c.Set("email", email)

		c.Next()
	}
}

package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	firebase "firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
	uuid "github.com/google/uuid"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/db/repositories"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/domain"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/pkg/logger"
)

const (
	uuidKey   = "uuid"
	userIdKey = "id"
	nameKey   = "name"
	emailKey  = "email"
)

func GetUUIDFromContext(c *gin.Context) (uuid.UUID, error) {
	value := c.GetString(uuidKey)
	if value == "" {
		return uuid.Nil, fmt.Errorf("UUID not found in context")
	}

	fmt.Printf("Trying to get uuid of %s\n", value)
	uuidValid, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid UUID type")
	}

	return uuidValid, nil
}

func GetUserIDFromContext(c *gin.Context) (string, error) {
	value := c.GetString(userIdKey)
	if value == "" {
		return "", domain.UnauthorizedError{Message: "user ID not found in context"}
	}
	return value, nil
}

func FirebaseAuthMiddleware(userRepository repositories.IUserRepository, firebaseApp *firebase.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantId := GetTenantFromContext(c)

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			return
		}

		idToken := strings.TrimPrefix(authHeader, "Bearer ")

		client, err := firebaseApp.Auth(context.Background())
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

		userResponse, err := userRepository.GetByUUID(token.UID, tenantId)

		if err != nil {
			logger.Error("User with uid " + token.UID + " no existe en la bd de productos")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "No se ha encontrado el usuario, reintentelo mas tarde"})
			return
		}

		c.Set(nameKey, name)
		c.Set(uuidKey, token.UID)
		c.Set(emailKey, email)
		c.Set(userIdKey, userResponse.ID)

		c.Next()
	}
}

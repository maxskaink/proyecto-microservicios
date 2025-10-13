package server

import "github.com/gin-gonic/gin"

// AuthMiddleware retorna un middleware placeholder que más adelante verificará tokens de Firebase.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: verificar token Firebase del header Authorization
		c.Next()
	}
}

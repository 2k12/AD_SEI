// package middleware

// import (
// 	"net/http"
// 	"strings"

// 	"github.com/gin-gonic/gin"
// 	"github.com/golang-jwt/jwt/v4"
// )

// func AuthMiddleware(secret string) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		tokenString := c.GetHeader("Authorization")
// 		if tokenString == "" {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token requerido"})
// 			c.Abort()
// 			return
// 		}

// 		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

// 		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 			return []byte(secret), nil
// 		})
// 		if err != nil || !token.Valid {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
// 			c.Abort()
// 			return
// 		}

// 		c.Next()
// 	}
// }

package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token requerido"})
			c.Abort()
			return
		}

		// Eliminar el prefijo "Bearer "
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Parsear el token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		// Extraer los claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No se pudieron obtener los claims del token"})
			c.Abort()
			return
		}

		// Buscar la clave correcta ("id" en este caso)
		userID, exists := claims["id"]
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "El token no contiene un ID de usuario"})
			c.Abort()
			return
		}

		// Guardar el ID del usuario en el contexto
		c.Set("userID", userID)

		// Continuar con la solicitud
		c.Next()
	}
}

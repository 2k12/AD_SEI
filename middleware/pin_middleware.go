package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func ValidatePIN() gin.HandlerFunc {
	return func(c *gin.Context) {
		providedPin := c.GetHeader("X-API-PIN")
		requiredPin := os.Getenv("API_PIN")

		if providedPin == "" || providedPin != requiredPin {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "PIN inválido o no proporcionado"})
			c.Abort()
			return
		}

		c.Next()
	}
}

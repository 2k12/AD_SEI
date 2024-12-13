package controllers

import (
	"net/http"
	"seguridad-api/services"

	"github.com/gin-gonic/gin"
)

// Login maneja la solicitud de inicio de sesión
func Login(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Parsear los datos del cuerpo de la solicitud
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// Llamar al servicio de autenticación
	token, err := services.Authenticate(loginData.Email, loginData.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Devolver el token de acceso
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func Logout(c *gin.Context) {
	// Aquí no se necesita hacer nada en el servidor, ya que el JWT es independiente
	// Pero puedes enviar una respuesta que indique que el logout fue exitoso
	c.JSON(http.StatusOK, gin.H{"message": "Sesión cerrada exitosamente"})
}

package controllers

import (
	"net/http"
	"seguridad-api/services"

	"github.com/gin-gonic/gin"
)

// Login autentica al usuario y genera un token JWT
// @Summary Iniciar sesión
// @Description Autentica un usuario con email y contraseña, devolviendo un token JWT
// @Tags Autenticación
// @Accept json
// @Produce json
// @Param loginData body object true "Datos de inicio de sesión"
// @Success 200 {object} map[string]string "token"
// @Failure 400 {object}  map[string]string "error"
// @Failure 401 {object}  map[string]string "error"
// @Router /login [post]
func Login(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	token, err := services.Authenticate(loginData.Email, loginData.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Logout cierra la sesión del usuario
// @Summary Cerrar sesión
// @Description Invalida la sesión actual del usuario. Requiere un Bearer Token.
// @Tags Autenticación
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]string "message"
// @Failure 401 {object} map[string]string "error"
// @Router /logout [post]
func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Sesión cerrada exitosamente"})
}

package controllers

import (
	// "log"
	"net/http"
	helpers "seguridad-api/helpers"
	"seguridad-api/services"

	// email "seguridad-api/services/email"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type TokenResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type ErrorResponse struct {
	Error string `json:"error" example:"Credenciales inválidas"`
}

type LoginData struct {
	Email     string `json:"email" example:"user@example.com"`
	Password  string `json:"password" example:"securePassword123"`
	ModuleKey string `json:"module_key" example:"....."`
}

// Login autentica al usuario y genera un token JWT
// @Summary Iniciar sesión
// @Description Autentica un usuario con email, contraseña y key del módulo devolviendo un token JWT
// @Tags Autenticación
// @Accept json
// @Produce json
// @Param loginData body LoginData true "Datos de inicio de sesión (email,password y la key del módulo correspondiente)"
// @Success 200 {object} TokenResponse "token"
// @Failure 400 {object} ErrorResponse "Datos inválidos"
// @Failure 401 {object} ErrorResponse "Credenciales inválidas"
// @Router /login [post]
func Login(c *gin.Context) {
	var loginData LoginData

	currentTime := time.Now()
	ecuadorTime := helpers.AdjustToEcuadorTime(currentTime)

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	token, err := services.Authenticate(loginData.Email, loginData.Password, loginData.ModuleKey)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	parsedToken, _, err := new(jwt.Parser).ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al decodificar el token"})
		return
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || claims["id"] == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ID del usuario no encontrado en el token"})
		return
	}

	userIDUint := uint(claims["id"].(float64))

	event := "INSERT"
	description := "Se registra ingreso a la plataforma, usuario: " + loginData.Email
	originService := "SEGURIDAD"

	if auditErr := services.RegisterAudit(event, description, userIDUint, originService, ecuadorTime); auditErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Usuario creado, pero no se pudo registrar la auditoría"})
		return
	}

	// // Enviar correo de notificación
	// go func() {
	// 	subject := "Inicio de sesión exitoso"
	// 	body := "Hola,\n\nSe ha registrado un inicio de sesión en la plataforma con tu cuenta de correo: " + loginData.Email + ".\n\nFecha y hora: " + ecuadorTime.Format("02/01/2006 15:04:05") + "\n\nSi no reconoces esta actividad, por favor contacta al soporte."
	// 	err := email.SendEmail(loginData.Email, subject, body)
	// 	if err != nil {
	// 		log.Printf("Error al enviar el correo electrónico: %v", err)
	// 	}
	// }()

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

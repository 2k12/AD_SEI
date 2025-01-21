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
// @Description Este endpoint autentica un usuario utilizando su email, contraseña y la clave del módulo correspondiente. Si la autenticación es exitosa, se genera y devuelve un token JWT.
// @Tags Autenticación
// @Accept json
// @Produce json
// @Param loginData body LoginData true "Datos de inicio de sesión (email, password y la key del módulo correspondiente)"
// @Success 200 {object} TokenResponse "Token JWT generado"
// @Failure 400 {object} ErrorResponse "Datos inválidos en la solicitud"
// @Failure 401 {object} ErrorResponse "Credenciales inválidas"
// @Failure 500 {object} ErrorResponse "Error interno al procesar la autenticación"
// @Router /login [post]
func Login(c *gin.Context) {
	var loginData LoginData

	currentTime := time.Now()
	ecuadorTime := helpers.AdjustToEcuadorTime(currentTime)

	// Verificar los datos de la solicitud
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// Autenticación del usuario
	token, err := services.Authenticate(loginData.Email, loginData.Password, loginData.ModuleKey)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Procesar y validar el token
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

	// Registrar evento de auditoría
	event := "INSERT"
	description := "Se registra ingreso a la plataforma, usuario: " + loginData.Email
	originService := "SEGURIDAD"

	if auditErr := services.RegisterAudit(event, description, userIDUint, originService, ecuadorTime); auditErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Usuario creado, pero no se pudo registrar la auditoría"})
		return
	}

	// Devolver el token generado
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Logout cierra la sesión del usuario
// @Summary Cerrar sesión
// @Description Invalida la sesión actual del usuario. Requiere un Bearer Token válido. Este endpoint cierra la sesión del usuario, eliminando el acceso al sistema hasta una nueva autenticación.
// @Tags Autenticación
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]string "message"
// @Failure 401 {object} map[string]string "error"
// @Router /logout [post]
func Logout(c *gin.Context) {
	// Aquí se pueden agregar acciones para invalidar el token si fuera necesario
	c.JSON(http.StatusOK, gin.H{"message": "Sesión cerrada exitosamente"})
}

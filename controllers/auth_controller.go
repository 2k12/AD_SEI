package controllers

import (
	"log"
	"net/http"
	helpers "seguridad-api/helpers"
	"seguridad-api/services"

	// email "seguridad-api/services/email"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
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

	// // Enviar correo electrónico al usuario
	// err = sendWelcomeEmail(loginData.Email)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al enviar el correo electrónico: " + err.Error()})
	// 	return
	// }

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

func sendWelcomeEmail(userEmail string) error {
	from := mail.NewEmail("Módulo de Seguridad", "sheremypavon12@gmail.com")
	subject := "¡Bienvenido al Módulo de Seguridad!"
	to := mail.NewEmail("Usuario", userEmail)

	plainTextContent := `¡Hola! 

Gracias por iniciar sesión en el Módulo de Seguridad. 
Si no fuiste tú quien realizó esta acción, por favor comunícate de inmediato con nuestro equipo de soporte.

Saludos cordiales, 
El equipo de Seguridad`

	htmlContent := `
		<html>
			<head>
				<style>
					body {
						font-family: Arial, sans-serif;
						background-color: #f4f4f4;
						color: #333;
					}
					.container {
						background-color: #ffffff;
						border-radius: 8px;
						padding: 20px;
						box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
						max-width: 600px;
						margin: 0 auto;
					}
					.header {
						text-align: center;
						background-color: #4CAF50;
						color: white;
						padding: 10px 0;
						border-radius: 8px;
					}
					.content {
						margin-top: 20px;
						font-size: 16px;
					}
					.footer {
						margin-top: 30px;
						text-align: center;
						font-size: 12px;
						color: #888;
					}
					.button {
						background-color: #4CAF50;
						color: white;
						padding: 10px 20px;
						border-radius: 5px;
						text-decoration: none;
						font-weight: bold;
					}
					.button:hover {
						background-color: #45a049;
					}
				</style>
			</head>
			<body>
				<div class="container">
					<div class="header">
						<h2>¡Bienvenido al Módulo de Seguridad!</h2>
					</div>
					<div class="content">
						<p>Hola, <strong>usuario</strong>!</p>
						<p>Gracias por iniciar sesión en el Módulo de Seguridad.</p>
						<p><strong>Si no fuiste tú quien inició sesión, por favor comunícate con soporte inmediatamente.</strong></p>
						<p>Para más información, visita nuestra página de soporte o contáctanos directamente.</p>
					</div>
					<div class="footer">
						<p>Saludos cordiales,</p>
						<p><strong>El equipo de Seguridad</strong></p>
						<p><a href="https://sei-ad-frontend.vercel.app" class="button">Ir al soporte</a></p>
					</div>
				</div>
			</body>
		</html>`

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)

	if err != nil {
		return fmt.Errorf("error al enviar el correo: %v", err)
	}

	if response.StatusCode >= 400 {
		return fmt.Errorf("error en el envío del correo, código de estado: %d, cuerpo: %s", response.StatusCode, response.Body)
	}

	log.Printf("Correo enviado exitosamente a %s. Código de estado: %d", userEmail, response.StatusCode)
	return nil
}

func RequestPasswordReset(c *gin.Context) {
	var request struct {
		Email string `json:"email"`
	}

	fmt.Println("Solicitud recibida:", c.Request.Body)

	if err := c.ShouldBindJSON(&request); err != nil {
		fmt.Println("Error en el formato de datos:", err) // 🔍 Log de error
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de datos inválido"})
		return
	}

	fmt.Println("Email recibido:", request.Email)

	err := services.SendPasswordResetEmail(request.Email)
	if err != nil {
		fmt.Println("Error al enviar correo:", err) // 🔍 Log de error
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Correo de restablecimiento enviado"})
}

// Confirmar restablecimiento de contraseña
// Confirmar restablecimiento de contraseña
func ResetPassword(c *gin.Context) {
	var request struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de datos inválido"})
		return
	}

	err := services.ResetPassword(request.Token, request.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contraseña restablecida exitosamente"})
}

// func UnlockUser(c *gin.Context) {
// 	var user models.User
// 	email := c.Param("email")

// 	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
// 		return
// 	}

// 	user.FailedAttempts = 0
// 	user.LockedUntil = nil
// 	config.DB.Save(&user)

// 	c.JSON(http.StatusOK, gin.H{"message": "Usuario desbloqueado exitosamente"})
// }

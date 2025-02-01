package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"
	"seguridad-api/config"
	"seguridad-api/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error al cargar el archivo .env")
	}
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		log.Println("JWT_SECRET_KEY no está configurada")
	}
}

const MAX_ATTEMPTS = 3
const LOCK_DURATION = 15 * time.Minute

func Authenticate(email, password, module_key string) (string, error) {
	var user models.User

	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return "", errors.New("usuario o contraseña inválidos")
	}

	// Verificar si la cuenta está bloqueada
	if user.LockedUntil != nil && time.Now().Before(*user.LockedUntil) {
		return "", errors.New("la cuenta está bloqueada. Inténtelo más tarde")
	}

	if !user.Active {
		return "", errors.New("la cuenta está inactiva")
	}

	if user.ModuleKey != module_key {
		return "", errors.New("no dispone de acceso a este módulo")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		// Incrementar intentos fallidos
		user.FailedAttempts++
		if user.FailedAttempts >= MAX_ATTEMPTS {
			lockTime := time.Now().Add(LOCK_DURATION)
			user.LockedUntil = &lockTime
		}
		config.DB.Save(&user)
		return "", errors.New("usuario o contraseña inválidos")
	}

	// Reiniciar intentos fallidos al iniciar sesión correctamente
	user.FailedAttempts = 0
	user.LockedUntil = nil
	config.DB.Save(&user)

	var roles []models.Role
	config.DB.Model(&user).Association("Roles").Find(&roles)

	roleNames := []string{}
	permissions := []string{}

	for _, role := range roles {
		roleNames = append(roleNames, role.Name)

		var perms []models.Permission
		config.DB.Model(&role).Association("Permissions").Find(&perms)
		for _, perm := range perms {
			log.Printf("Permission: %s", perm.Name)
			permissions = append(permissions, perm.Name)
		}
	}

	log.Printf("Permissions collected: %v", permissions)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":          user.ID,
		"Name":        user.Name,
		"email":       user.Email,
		"roles":       roleNames,
		"permissions": permissions,
		"exp":         time.Now().Add(time.Hour * 1).Unix(),
	})

	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		return "", errors.New("clave secreta no definida en .env")
	}

	return token.SignedString([]byte(secretKey))
}

func generateResetToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// Enviar correo con enlace de restablecimiento
func SendPasswordResetEmail(email string) error {
	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return errors.New("usuario no encontrado")
	}

	// Generar token y guardar en la BD
	token := generateResetToken()
	expiry := time.Now().Add(30 * time.Minute) // Expira en 30 minutos
	user.ResetToken = token
	user.ResetTokenExpiry = &expiry
	config.DB.Save(&user)

	// 🔍 Verificar token generado
	fmt.Println("Token generado para", user.Email, ":", token)

	// Crear contenido del correo
	subject := "Restablecimiento de Contraseña"
	body := fmt.Sprintf(`
    <div style="font-family: Arial, sans-serif; max-width: 500px; margin: auto; padding: 20px; border: 1px solid #ddd; border-radius: 8px; background-color: #f9f9f9;">
        <h2 style="color: #333; text-align: center;">🔐 Restablecimiento de Contraseña</h2>
        <p style="color: #555; text-align: justify;">
            Hemos recibido una solicitud para restablecer tu contraseña. Para proceder, haz clic en el botón de abajo:
        </p>
        <div style="text-align: center; margin: 20px 0;">
            <a href="http://localhost:5173/reset-password?token=%s" 
               style="display: inline-block; padding: 12px 20px; font-size: 16px; color: white; background-color: #007bff; text-decoration: none; border-radius: 5px;">
               🔄 Restablecer Contraseña
            </a>
        </div>
        <p style="color: #555; text-align: justify;">
            Si no solicitaste este cambio, puedes ignorar este correo de manera segura.
        </p>
        <p style="color: #888; font-size: 12px; text-align: center;">
            ⚠️ Este enlace expirará en <strong>30 minutos</strong>. Si tienes problemas, contáctanos.
        </p>
    </div>`, token)

	err := SendEmail(user.Email, subject, body)
	if err != nil {
		fmt.Println("Error al enviar correo:", err)
		return errors.New("error al enviar el correo: " + err.Error())
	}

	fmt.Println("Correo enviado exitosamente a", user.Email)
	return nil
}

// Validar token y actualizar contraseña
func ResetPassword(token string, newPassword string) error {
	var user models.User

	// Buscar usuario por token
	if err := config.DB.Where("reset_token = ?", token).First(&user).Error; err != nil {
		return errors.New("token inválido o expirado")
	}

	// Verificar si el token ha expirado
	if user.ResetTokenExpiry == nil || time.Now().After(*user.ResetTokenExpiry) {
		return errors.New("el token ha expirado")
	}

	// Encriptar nueva contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("error al encriptar la nueva contraseña")
	}

	// Guardar nueva contraseña y limpiar token
	user.Password = string(hashedPassword)
	user.ResetToken = ""
	user.ResetTokenExpiry = nil
	if err := config.DB.Save(&user).Error; err != nil {
		return errors.New("error al actualizar la contraseña")
	}

	return nil
}

func SendEmail(to, subject, body string) error {
	// Verificar que `SENDGRID_API_KEY` está configurado
	apiKey := os.Getenv("SENDGRID_API_KEY")
	if apiKey == "" {
		log.Println("ERROR: `SENDGRID_API_KEY` no está configurado en las variables de entorno.")
		return fmt.Errorf("configuración de SendGrid incorrecta")
	}

	// Configurar correo
	from := mail.NewEmail("Security Service", "sheremypavon12@gmail.com")
	toEmail := mail.NewEmail("Usuario", to)
	plainTextContent := body
	htmlContent := "<strong>" + body + "</strong>"

	message := mail.NewSingleEmail(from, subject, toEmail, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(apiKey)

	// Enviar correo y capturar respuesta
	response, err := client.Send(message)
	if err != nil {
		log.Printf("Error al enviar correo: %v", err)
		return err
	}

	log.Printf("Correo enviado con éxito: Código %d", response.StatusCode)
	log.Printf("Cuerpo de la respuesta: %s", response.Body)

	// Verificar si SendGrid retornó un error
	if response.StatusCode >= 400 {
		log.Printf("Error en el envío: Código %d - %s", response.StatusCode, response.Body)
		return fmt.Errorf("error en el envío del correo, código %d", response.StatusCode)
	}

	return nil
}

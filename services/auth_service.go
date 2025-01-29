package services

import (
	"errors"
	"log"
	"os"
	"seguridad-api/config"
	"seguridad-api/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
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

	// Buscar usuario por email
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return "", errors.New("usuario o contraseña inválidos")
	}

	// Verificar si la cuenta está bloqueada
	if user.LockedUntil != nil && time.Now().Before(*user.LockedUntil) {
		return "", errors.New("la cuenta está bloqueada. Inténtelo más tarde")
	}

	// Verificar la contraseña
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

	// Generar token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"Name":  user.Name,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 1).Unix(),
	})

	secretKey := "tu_clave_secreta" // Cargar desde .env
	return token.SignedString([]byte(secretKey))
}

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

// Cargar las variables de entorno
func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error al cargar el archivo .env")
	}
}

// Authenticate maneja la autenticación del usuario
func Authenticate(email, password string) (string, error) {
	var user models.User

	// Buscar el usuario por email
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return "", errors.New("usuario o contraseña inválidos")
	}

	// Validar si el usuario está activo
	if !user.Active {
		return "", errors.New("la cuenta está inactiva")
	}

	// Verificar la contraseña
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("usuario o contraseña inválidos")
	}

	// Obtener roles y permisos
	var roles []models.Role
	config.DB.Model(&user).Association("Roles").Find(&roles)

	roleNames := []string{}
	permissions := []string{}

	for _, role := range roles {
		roleNames = append(roleNames, role.Name)

		var perms []models.Permission
		config.DB.Model(&role).Association("Permissions").Find(&perms)
		for _, perm := range perms {
			log.Printf("Permission: %s", perm.Name) // Depuración correcta
			permissions = append(permissions, perm.Name)
		}
	}

	log.Printf("Permissions collected: %v", permissions) // Verifica los permisos antes de generar el token

	// Crear token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":          user.ID,
		"email":       user.Email,
		"roles":       roleNames,
		"permissions": permissions, // Asegúrate de incluir los permisos
		"exp":         time.Now().Add(time.Hour * 1).Unix(),
	})

	// Obtener la clave secreta desde las variables de entorno
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		return "", errors.New("clave secreta no definida en .env")
	}

	// Firmar el token con la clave secreta
	return token.SignedString([]byte(secretKey))
}

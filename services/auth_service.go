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

// func Authenticate(email, password, module_key string) (string, error) {
// 	var user models.User

// 	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
// 		return "", errors.New("usuario o contraseña inválidos")
// 	}

// 	if !user.Active {
// 		return "", errors.New("la cuenta está inactiva")
// 	}

// 	if user.ModuleKey != module_key {
// 		return "", errors.New("no dispone de acceso a este módulo")
// 	}

// 	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
// 	if err != nil {
// 		return "", errors.New("usuario o contraseña inválidos")
// 	}

// 	var roles []models.Role
// 	config.DB.Model(&user).Association("Roles").Find(&roles)

// 	roleNames := []string{}
// 	permissions := []string{}

// 	for _, role := range roles {
// 		roleNames = append(roleNames, role.Name)

// 		var perms []models.Permission
// 		config.DB.Model(&role).Association("Permissions").Find(&perms)
// 		for _, perm := range perms {
// 			log.Printf("Permission: %s", perm.Name)
// 			permissions = append(permissions, perm.Name)
// 		}
// 	}

// 	log.Printf("Permissions collected: %v", permissions)

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"id":          user.ID,
// 		"Name":        user.Name,
// 		"email":       user.Email,
// 		"roles":       roleNames,
// 		"permissions": permissions,
// 		"exp":         time.Now().Add(time.Hour * 1).Unix(),
// 	})

// 	secretKey := os.Getenv("JWT_SECRET_KEY")
// 	if secretKey == "" {
// 		return "", errors.New("clave secreta no definida en .env")
// 	}

//		return token.SignedString([]byte(secretKey))
//	}
func Authenticate(email, password, moduleKey string) (string, error) {
	var user models.User

	// Verificar si el usuario existe
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return "", errors.New("usuario o contraseña inválidos")
	}

	// Verificar si la cuenta está activa
	if !user.Active {
		return "", errors.New("la cuenta está inactiva")
	}

	// Verificar contraseña
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("usuario o contraseña inválidos")
	}

	// Obtener los roles del usuario
	var roles []models.Role
	config.DB.Model(&user).Association("Roles").Find(&roles)

	// Validar acceso al módulo solicitado
	hasAccess := false
	roleNames := []string{}
	permissions := []string{}

	for _, role := range roles {
		roleNames = append(roleNames, role.Name)

		// Obtener los permisos del rol
		var perms []models.Permission
		config.DB.Model(&role).Association("Permissions").Find(&perms)

		for _, perm := range perms {
			log.Printf("Permission: %s", perm.Name)
			permissions = append(permissions, perm.Name)

			// Validar si el permiso pertenece al módulo solicitado
			var module models.Module
			if err := config.DB.Where("id = ?", perm.ModuleID).First(&module).Error; err == nil {
				if module.ModuleKey == moduleKey {
					hasAccess = true
				}
			}
		}
	}

	if !hasAccess {
		return "", errors.New("no dispone de acceso a este módulo")
	}

	// Generar el token JWT
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

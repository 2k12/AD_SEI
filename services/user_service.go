package services

import (
	"errors"
	"seguridad-api/config"
	"seguridad-api/models"

	"golang.org/x/crypto/bcrypt"
)

// Función para encriptar la contraseña
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// Crear usuario con encriptación
func CreateUser(name, email, password string, active bool) (models.User, error) {
	// Encriptar la contraseña
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return models.User{}, errors.New("error al encriptar la contraseña")
	}

	// Crear el usuario con el role_id
	user := models.User{
		Name:     name,
		Email:    email,
		Password: hashedPassword,
		Active:   active,
	}

	// Insertar el usuario en la base de datos
	result := config.DB.Create(&user)
	if result.Error != nil {
		return models.User{}, result.Error
	}

	return user, nil
}

// Obtener usuarios
//
//	func GetUsers() ([]models.User, error) {
//		var users []models.User
//		result := config.DB.Find(&users)
//		return users, result.Error
//	}
// func GetUsers() ([]models.User, error) {
// 	var users []models.User

// 	// Usamos Preload para cargar roles y permisos asociados a los roles
// 	result := config.DB.Debug().Preload("Role.Permissions").Find(&users)
// 	// result := config.DB.Debug().Preload("Roles.Permissions").Find(&users)

// 	if result.Error != nil {
// 		// Si hay algún error, lo devolvemos
// 		return nil, result.Error
// 	}

//		// Devolvemos los usuarios con las relaciones cargadas
//		return users, nil
//	}
func GetUsers() ([]models.User, error) {
	var users []models.User

	// Usamos Preload para cargar los roles y permisos de los usuarios
	result := config.DB.Debug().Preload("Roles.Permissions").Find(&users)

	if result.Error != nil {
		// Si hay algún error, lo devolvemos
		return nil, result.Error
	}

	// Devolvemos los usuarios con las relaciones cargadas
	return users, nil
}

// Actualizar un usuario existente
func UpdateUser(id string, name, email string, active *bool) (models.User, error) {
	var user models.User

	// Buscar el usuario por ID
	if err := config.DB.First(&user, "id = ?", id).Error; err != nil {
		return models.User{}, errors.New("usuario no encontrado")
	}

	// Solo actualizar los campos que no son nulos o vacíos
	if name != "" {
		user.Name = name
	}
	if email != "" {
		user.Email = email
	}
	if active != nil {
		user.Active = *active
	}
	// if roleID != nil {
	// 	user.RoleID = *roleID // Actualizar el role_id si se pasa
	// }

	// Guardar los cambios en la base de datos
	result := config.DB.Save(&user)
	return user, result.Error
}

func DeleteUser(id string) error {
	var user models.User

	if err := config.DB.First(&user, "id = ?", id).Error; err != nil {
		return errors.New("usuario no encontrado")
	}

	// user.Active = false
	user.Active = !user.Active
	return config.DB.Save(&user).Error
}

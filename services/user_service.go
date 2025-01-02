package services

import (
	"errors"
	"seguridad-api/config"
	"seguridad-api/models"

	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CreateUser(name, email, password string, active bool) (models.User, error) {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return models.User{}, errors.New("error al encriptar la contraseña")
	}

	user := models.User{
		Name:     name,
		Email:    email,
		Password: hashedPassword,
		Active:   active,
	}

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
// func GetUsers() ([]models.User, error) {
// 	var users []models.User

// 	result := config.DB.Debug().Preload("Roles.Permissions").Find(&users)

// 	if result.Error != nil {
// 		return nil, result.Error
// 	}

//		return users, nil
//	}
func GetUsers() ([]models.User, error) {
	var users []models.User

	// Preload Roles, Permissions, and Modules
	result := config.DB.Debug().
		Preload("Roles.Permissions.Module").
		Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

func GetUserPermissions(userID uint) ([]models.Permission, error) {
	var user models.User
	var permissions []models.Permission

	result := config.DB.Debug().Preload("Roles.Permissions").First(&user, userID)

	if result.Error != nil {
		return nil, result.Error
	}

	for _, role := range user.Roles {
		for _, permission := range role.Permissions {
			permissions = append(permissions, permission)
		}
	}

	return permissions, nil
}

func UpdateUser(id string, name, email string, active *bool) (models.User, error) {
	var user models.User

	if err := config.DB.First(&user, "id = ?", id).Error; err != nil {
		return models.User{}, errors.New("usuario no encontrado")
	}

	if name != "" {
		user.Name = name
	}
	if email != "" {
		user.Email = email
	}
	if active != nil {
		user.Active = *active
	}

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

func GetPaginatedUsers(page, pageSize int, filters map[string]interface{}) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := config.DB.Model(&models.User{})

	// Aplicar filtros
	if name, ok := filters["name"]; ok {
		query = query.Where("name ILIKE ?", "%"+name.(string)+"%")
	}
	if email, ok := filters["email"]; ok {
		query = query.Where("email ILIKE ?", "%"+email.(string)+"%")
	}
	if active, ok := filters["active"]; ok {
		query = query.Where("active = ?", active)
	}

	// Contar el total de registros
	query.Count(&total)

	// Aplicar paginación
	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	// Obtener usuarios con relaciones
	err := query.Preload("Roles.Permissions.Module").Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

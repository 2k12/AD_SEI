package services

import (
	"errors"
	"seguridad-api/config"
	"seguridad-api/models"
)

// Crear un rol
func CreateRole(name, description string) (models.Role, error) {
	role := models.Role{Name: name, Description: description}
	result := config.DB.Create(&role)
	return role, result.Error
}

// Obtener todos los roles
func GetRoles() ([]models.Role, error) {
	var roles []models.Role
	result := config.DB.Find(&roles)
	return roles, result.Error
}

// Actualizar un rol existente
func UpdateRole(id int, name string, description string, active bool) (models.Role, error) {
	var role models.Role

	// Buscar el rol por ID
	if err := config.DB.First(&role, id).Error; err != nil {
		return role, errors.New("rol no encontrado")
	}

	// Actualizar los campos
	role.Name = name
	role.Description = description
	role.Active = active // Agregar actualización del estado

	// Guardar los cambios
	if err := config.DB.Save(&role).Error; err != nil {
		return role, errors.New("error al actualizar el rol")
	}

	return role, nil
}

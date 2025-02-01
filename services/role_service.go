package services

import (
	"errors"
	"seguridad-api/config"
	"seguridad-api/models"
)

// Crear un rol
func CreateRole(name, description string, module uint) (models.Role, error) {
	role := models.Role{Name: name, Description: description, IDModule: module}
	result := config.DB.Create(&role)
	return role, result.Error
}

// Obtener todos los roles
func GetRoles() ([]models.Role, error) {
	var roles []models.Role
	result := config.DB.Find(&roles)
	return roles, result.Error
}

// Obtener solo los roles activos
func GetRolesActive() ([]models.Role, error) {
	var roles []models.Role
	result := config.DB.Where("active = ?", true).Find(&roles)
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
func GetPaginatedRoles(page, pageSize int, filters map[string]interface{}) ([]models.Role, int64, error) {
	var roles []models.Role
	var total int64

	query := config.DB.Model(&models.Role{})

	// Aplicar filtros
	if name, ok := filters["name"]; ok {
		query = query.Where("name LIKE ?", "%"+name.(string)+"%")
	}
	if active, ok := filters["active"]; ok {
		query = query.Where("active = ?", active)
	}

	// Contar el total de registros
	query.Count(&total)

	// Aplicar paginación
	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	// Obtener roles
	err := query.Preload("Permissions.Module").Find(&roles).Error // Preload para relaciones si aplica
	if err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}
func UpdateRoleState(id int, active bool) error {
	var role models.Role

	// Buscar el rol por ID
	if err := config.DB.First(&role, id).Error; err != nil {
		return errors.New("rol no encontrado")
	}

	// Actualizar el estado del rol
	role.Active = active

	// Guardar los cambios
	if err := config.DB.Save(&role).Error; err != nil {
		return errors.New("error al actualizar el estado del rol")
	}

	return nil
}

package services

import (
	"seguridad-api/config"
	"seguridad-api/models"
)

func CreatePermission(permission models.Permission) (models.Permission, error) {
	if err := config.DB.Create(&permission).Error; err != nil {
		return models.Permission{}, err
	}
	return permission, nil
}

func GetPermissions() ([]models.Permission, error) {
	var permissions []models.Permission
	if err := config.DB.Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

// GetActivePermissions devuelve los permisos que están activos.
func GetActivePermissions() ([]models.Permission, error) {
	var permissions []models.Permission
	if err := config.DB.Where("active = ?", true).Preload("Module").Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

func GetPaginatedPermissions(page, pageSize int, filters map[string]interface{}) ([]models.Permission, int64, error) {
	var permissions []models.Permission
	var total int64

	query := config.DB.Model(&models.Permission{})

	// Aplicar filtros dinámicos
	if moduleName, ok := filters["moduleName"]; ok && moduleName != "" {
		query = query.Joins("JOIN modules ON modules.id = permissions.module_id").
			Where("LOWER(modules.name) LIKE LOWER(?)", "%"+moduleName.(string)+"%")
	}
	if name, ok := filters["name"]; ok && name != "" {
		query = query.Where("LOWER(permissions.name) LIKE LOWER(?)", "%"+name.(string)+"%")
	}
	if active, ok := filters["active"]; ok {
		query = query.Where("permissions.active = ?", active)
	}

	// Contar el total de registros filtrados
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Aplicar paginación
	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	// Obtener datos
	err := query.Preload("Module").Find(&permissions).Error
	if err != nil {
		return nil, 0, err
	}

	return permissions, total, nil
}

func UpdatePermission(id string, updatedPermission models.Permission) (models.Permission, error) {
	var permission models.Permission
	if err := config.DB.First(&permission, "id = ?", id).Error; err != nil {
		return models.Permission{}, err
	}
	permission.Name = updatedPermission.Name
	permission.Description = updatedPermission.Description
	permission.ModuleID = updatedPermission.ModuleID
	permission.Active = updatedPermission.Active
	if err := config.DB.Save(&permission).Error; err != nil {
		return models.Permission{}, err
	}
	return permission, nil
}

func DeletePermission(id string) error {
	if err := config.DB.Delete(&models.Permission{}, id).Error; err != nil {
		return err
	}
	return nil
}

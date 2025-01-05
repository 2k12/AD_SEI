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

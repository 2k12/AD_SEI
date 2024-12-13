package services

import (
	"seguridad-api/config"
	"seguridad-api/models"
)

func CreatePermission(name, description string) (models.Permission, error) {
	permission := models.Permission{Name: name, Description: description}
	result := config.DB.Create(&permission)
	return permission, result.Error
}

func GetPermissions() ([]models.Permission, error) {
	var permissions []models.Permission
	result := config.DB.Find(&permissions)
	return permissions, result.Error
}

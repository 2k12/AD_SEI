package services

import (
	"seguridad-api/config"
	"seguridad-api/models"
)

func CreateRole(name, description string) (models.Role, error) {
	role := models.Role{Name: name, Description: description}
	result := config.DB.Create(&role)
	return role, result.Error
}

func GetRoles() ([]models.Role, error) {
	var roles []models.Role
	result := config.DB.Find(&roles)
	return roles, result.Error
}

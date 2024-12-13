package services

import (
	"errors"
	"seguridad-api/config"
	"seguridad-api/models"
)

// Crear un nuevo módulo
func CreateModule(name, description string, active bool) error {
	module := models.Module{
		Name:        name,
		Description: description,
		Active:      active,
	}

	result := config.DB.Create(&module)
	return result.Error
}

// Obtener un módulo por su ID
func GetModuleByID(id uint) (*models.Module, error) {
	var module models.Module
	if err := config.DB.First(&module, id).Error; err != nil {
		return nil, err
	}
	return &module, nil
}

// Obtener todos los módulos
func GetAllModules() ([]models.Module, error) {
	var modules []models.Module
	if err := config.DB.Find(&modules).Error; err != nil {
		return nil, err
	}
	return modules, nil
}

// Actualizar un módulo parcialmente (solo los campos enviados)
func UpdateModule(id uint, name, description *string, active *bool) error {
	var module models.Module
	if err := config.DB.First(&module, id).Error; err != nil {
		return errors.New("module not found")
	}

	// Actualizar solo los campos enviados
	if name != nil {
		module.Name = *name
	}
	if description != nil {
		module.Description = *description
	}
	if active != nil {
		module.Active = *active
	}

	result := config.DB.Save(&module)
	return result.Error
}

// "Eliminar" un módulo (cambiar estado de active)
func ToggleModuleActive(id uint) error {
	var module models.Module
	if err := config.DB.First(&module, id).Error; err != nil {
		return errors.New("module not found")
	}

	// Cambiar el estado de active
	module.Active = !module.Active

	result := config.DB.Save(&module)
	return result.Error
}

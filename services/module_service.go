package services

import (
	"errors"
	"seguridad-api/config"
	"seguridad-api/models"
)

func CreateModule(name, description string, active bool) error {
	module := models.Module{
		Name:        name,
		Description: description,
		Active:      active,
	}

	result := config.DB.Create(&module)
	return result.Error
}

func GetModuleByID(id uint) (*models.Module, error) {
	var module models.Module
	if err := config.DB.First(&module, id).Error; err != nil {
		return nil, err
	}
	return &module, nil
}

func GetAllModules() ([]models.Module, error) {
	var modules []models.Module
	if err := config.DB.Find(&modules).Error; err != nil {
		return nil, err
	}
	return modules, nil
}

func UpdateModule(id uint, name, description *string, active *bool) error {
	var module models.Module
	if err := config.DB.First(&module, id).Error; err != nil {
		return errors.New("module not found")
	}

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

func DeleteModule(id uint) error {
	var module models.Module
	if err := config.DB.First(&module, id).Error; err != nil {
		return errors.New("module not found")
	}

	result := config.DB.Delete(&module)
	return result.Error
}

// func ToggleModuleActive(id uint) error {
// 	var module models.Module
// 	if err := config.DB.First(&module, id).Error; err != nil {
// 		return errors.New("module not found")
// 	}

// 	module.Active = !module.Active

// 	result := config.DB.Save(&module)
// 	return result.Error
// }

func GetPaginatedModules(page, pageSize int, filters map[string]interface{}) ([]models.Module, int64, error) {
	var modules []models.Module
	var total int64

	query := config.DB.Model(&models.Module{})

	// Aplicar filtros
	// if email, ok := filters["email"]; ok {
	// 	query = query.Where("email = ?", email)
	// }
	if name, ok := filters["name"]; ok {
		query = query.Where("name LIKE ?", "%"+name.(string)+"%")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	err := query.Find(&modules).Error
	if err != nil {
		return nil, 0, err
	}

	return modules, total, nil
}

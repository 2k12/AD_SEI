package services

import (
	"errors"
	"fmt"
	"seguridad-api/config"
	"seguridad-api/models"
)

// Validar existencia del rol
func RoleExists(roleID uint) (bool, error) {
	var role models.Role
	result := config.DB.First(&role, roleID)
	return result.RowsAffected > 0, result.Error
}

// Validar existencia del permiso
func PermissionExists(permissionID uint) (bool, error) {
	var permission models.Permission
	result := config.DB.First(&permission, permissionID)
	return result.RowsAffected > 0, result.Error
}

func AssignPermissionToRole(roleID, permissionID uint) (models.RolePermission, error) {
	fmt.Printf("Validando existencia del rol: %d\n", roleID)
	exists, err := RoleExists(roleID)
	if err != nil {
		fmt.Println("Error al verificar el rol:", err)
		return models.RolePermission{}, err
	}
	if !exists {
		fmt.Println("Rol no encontrado")
		return models.RolePermission{}, errors.New("rol no encontrado")
	}

	fmt.Printf("Validando existencia del permiso: %d\n", permissionID)
	exists, err = PermissionExists(permissionID)
	if err != nil {
		fmt.Println("Error al verificar el permiso:", err)
		return models.RolePermission{}, err
	}
	if !exists {
		fmt.Println("Permiso no encontrado")
		return models.RolePermission{}, errors.New("permiso no encontrado")
	}

	fmt.Println("Creando la relación rol-permiso")
	rolePermission := models.RolePermission{
		RoleID:       roleID,
		PermissionID: permissionID,
	}
	result := config.DB.Create(&rolePermission)
	if result.Error != nil {
		fmt.Printf("Error al crear relación rol-permiso: %v\n", result.Error)
		return models.RolePermission{}, result.Error
	}

	fmt.Println("Relación rol-permiso creada correctamente")
	return rolePermission, nil
}

// Eliminar un permiso de un rol
func RemovePermissionFromRole(roleID, permissionID uint) error {
	// Validar rol y permiso
	exists, err := RoleExists(roleID)
	if err != nil || !exists {
		return errors.New("rol no encontrado")
	}

	exists, err = PermissionExists(permissionID)
	if err != nil || !exists {
		return errors.New("permiso no encontrado")
	}

	// Eliminar la relación
	result := config.DB.Where("role_id = ? AND permission_id = ?", roleID, permissionID).Delete(&models.RolePermission{})
	return result.Error
}

// Obtener permisos de un rol
func GetPermissionsByRole(roleID uint) ([]models.Permission, error) {
	var permissions []models.Permission
	err := config.DB.Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&permissions).Error
	return permissions, err
}

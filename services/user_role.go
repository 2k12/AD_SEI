package services

import (
	"seguridad-api/config"
	"seguridad-api/models"
)

// Asignar un rol a un usuario
func AssignRoleToUser(userID, roleID uint) (models.UserRole, error) {
	userRole := models.UserRole{
		UserID: userID,
		RoleID: roleID,
	}
	result := config.DB.Create(&userRole) // Inserta el nuevo registro en la tabla user_roles
	return userRole, result.Error
}

// Eliminar un rol de un usuario
func RemoveRoleFromUser(userID, roleID uint) error {
	// Elimina el registro que coincide con el userID y roleID
	result := config.DB.Where("user_id = ? AND role_id = ?", userID, roleID).Delete(&models.UserRole{})
	return result.Error
}

// Obtener roles asignados a un usuario
func GetUserRoles(userID uint) ([]models.Role, error) {
	var roles []models.Role
	err := config.DB.Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Find(&roles).Error
	return roles, err
}

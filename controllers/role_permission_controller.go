package controllers

import (
	"net/http"
	"seguridad-api/config"
	"seguridad-api/models"
	"seguridad-api/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AssignPermission(c *gin.Context) {
	var input struct {
		PermissionID uint `json:"permission_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	roleID, err := strconv.ParseUint(c.Param("role_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de rol inválido"})
		return
	}

	rolePermission, err := services.AssignPermissionToRole(uint(roleID), input.PermissionID)
	if err != nil {
		// Devolver el mensaje de error específico del servicio
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"role_permission": rolePermission})
}

func GetAllPermissions(c *gin.Context) {
	var permissions []models.Permission
	if err := config.DB.Find(&permissions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener permisos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"permissions": permissions})
}

// Eliminar un permiso de un rol
func RemovePermission(c *gin.Context) {
	var input struct {
		PermissionID uint `json:"permission_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roleID, err := strconv.ParseUint(c.Param("role_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de rol inválido"})
		return
	}

	err = services.RemovePermissionFromRole(uint(roleID), input.PermissionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el permiso"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permiso eliminado correctamente"})
}

// Obtener permisos de un rol
func GetRolePermissions(c *gin.Context) {
	roleID, err := strconv.ParseUint(c.Param("role_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de rol inválido"})
		return
	}

	permissions, err := services.GetPermissionsByRole(uint(roleID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los permisos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"permissions": permissions})
}

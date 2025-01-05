package controllers

import (
	"net/http"
	"seguridad-api/config"
	"seguridad-api/models"
	"seguridad-api/services"

	"github.com/gin-gonic/gin"
)

// CreatePermission crea un nuevo permiso
func CreatePermission(c *gin.Context) {
	var input models.Permission
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	permission, err := services.CreatePermission(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el permiso"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"permission": permission})
}

// GetPermissions obtiene todos los permisos
func GetPermissions(c *gin.Context) {
	permissions, err := services.GetPermissions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los permisos"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"permissions": permissions})
}

// UpdatePermission actualiza un permiso existente
func UpdatePermission(c *gin.Context) {
	id := c.Param("id")
	var input models.Permission
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	permission, err := services.UpdatePermission(id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el permiso"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"permission": permission})
}

// DeletePermission elimina un permiso
func DeletePermission(c *gin.Context) {
	id := c.Param("id")
	if err := services.DeletePermission(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el permiso"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Permiso eliminado exitosamente"})
}

// GetPermissionByID obtiene los detalles de un permiso por ID
func GetPermissionByID(c *gin.Context) {
	id := c.Param("id")
	var permission models.Permission
	if err := config.DB.Preload("Module").First(&permission, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Permiso no encontrado"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"permission": permission})
}

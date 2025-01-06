package controllers

import (
	"net/http"
	"seguridad-api/config"
	helpers "seguridad-api/helpers"
	"seguridad-api/models"
	"seguridad-api/services"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// CreatePermission crea un nuevo permiso
func CreatePermission(c *gin.Context) {
	var input struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		ModuleID    uint   `json:"module_id" binding:"required"`
		Active      bool   `json:"active"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	permission, err := services.CreatePermission(models.Permission{
		Name:        input.Name,
		Description: input.Description,
		ModuleID:    input.ModuleID,
		Active:      input.Active,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el permiso"})
		return
	}

	// Obtener el userID del contexto
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No se pudo obtener el ID del usuario desde el contexto"})
		return
	}

	// Si el ID es de tipo float64, conviértelo a uint
	userIDUint := uint(userID.(float64))

	// Registrar evento de auditoría
	event := "INSERT"
	description := "Se creó un permiso con el nombre: " + input.Name
	originService := "SEGURIDAD"
	currentTime := time.Now()
	ecuadorTime := helpers.AdjustToEcuadorTime(currentTime)

	if auditErr := services.RegisterAudit(event, description, userIDUint, originService, ecuadorTime); auditErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Permiso creado, pero no se pudo registrar la auditoría"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"permission": permission})
}

// GetPermissions obtiene todos los permisos con paginación
func GetPermissions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var permissions []models.Permission
	var total int64

	config.DB.Offset(offset).Limit(limit).Find(&permissions)
	config.DB.Model(&models.Permission{}).Count(&total)

	c.JSON(http.StatusOK, gin.H{
		"permissions": permissions,
		"page":        page,
		"limit":       limit,
		"total":       total,
	})
}

// UpdatePermission actualiza un permiso existente
func UpdatePermission(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		ModuleID    uint   `json:"module_id"`
		Active      bool   `json:"active"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	permission, err := services.UpdatePermission(id, models.Permission{
		Name:        input.Name,
		Description: input.Description,
		ModuleID:    input.ModuleID,
		Active:      input.Active,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el permiso"})
		return
	}

	// Obtener el userID del contexto
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No se pudo obtener el ID del usuario desde el contexto"})
		return
	}

	// Si el ID es de tipo float64, conviértelo a uint
	userIDUint := uint(userID.(float64))

	// Registrar evento de auditoría
	event := "UPDATE"
	description := "Se actualizó el permiso con ID: " + id
	originService := "SEGURIDAD"
	currentTime := time.Now()
	ecuadorTime := helpers.AdjustToEcuadorTime(currentTime)

	if auditErr := services.RegisterAudit(event, description, userIDUint, originService, ecuadorTime); auditErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Permiso actualizado, pero no se pudo registrar la auditoría"})
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

	// Obtener el userID del contexto
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No se pudo obtener el ID del usuario desde el contexto"})
		return
	}

	// Si el ID es de tipo float64, conviértelo a uint
	userIDUint := uint(userID.(float64))

	// Registrar evento de auditoría
	event := "DELETE"
	description := "Se eliminó el permiso con ID: " + id
	originService := "SEGURIDAD"
	currentTime := time.Now()
	ecuadorTime := helpers.AdjustToEcuadorTime(currentTime)

	if auditErr := services.RegisterAudit(event, description, userIDUint, originService, ecuadorTime); auditErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Permiso eliminado, pero no se pudo registrar la auditoría"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permiso eliminado exitosamente"})
}

// GetPermissionByID obtiene los detalles de un permiso por ID
func GetPermissionByID(c *gin.Context) {
	id := c.Param("id")
	var permission models.Permission

	if err := config.DB.Preload("Module").Where("id = ?", id).First(&permission).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Permiso no encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"permission": permission})
}

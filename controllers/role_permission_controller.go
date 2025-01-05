package controllers

import (
	"net/http"
	"seguridad-api/config"
	"seguridad-api/models"
	"seguridad-api/services"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// PermissionDataRequest representa la estructura del cuerpo para asignar o eliminar permisos
type PermissionDataRequest struct {
	PermissionID uint `json:"permission_id" example:"1"`
}

// PermissionResponse representa la estructura de un permiso en la respuesta
type PermissionResponse struct {
	ID          uint   `json:"id" example:"1"`
	Name        string `json:"name" example:"Gestion Usuarios"`
	Description string `json:"description" example:"Permiso de Crud Usuarios"`
	ModuleID    uint   `json:"module_id" example:"1"`
	Active      bool   `json:"active" example:"true"`
}

// AssignPermission asigna un permiso a un rol y registra auditoría
// @Summary Asignar permiso a rol
// @Description Asocia un permiso existente a un rol específico
// @Tags Permisos
// @Accept json
// @Produce json
// @Param role_id path int true "ID del rol"
// @Param permissionData body PermissionDataRequest true "Datos del permiso a asignar"
// @Success 200 {object} map[string]interface{} "role_permission"
// @Failure 400 {object} map[string]string "error"
// @Router /roles/{role_id}/permissions [post]
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Obtener el userID del contexto
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No se pudo obtener el ID del usuario desde el contexto"})
		return
	}

	// Ajustar el tiempo al huso horario de Ecuador
	currentTime := time.Now()
	ecuadorTime := currentTime.Add(-5 * time.Hour)

	event := "ASIGNAR"
	description := "Se asignó el permiso con ID " + strconv.Itoa(int(input.PermissionID)) + " al rol con ID " + strconv.Itoa(int(roleID))
	originService := "role_permission_service"

	if err := services.RegisterAudit(event, description, uint(userID.(float64)), originService, ecuadorTime); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Permiso asignado, pero no se pudo registrar la auditoría"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"role_permission": rolePermission})
}

// RemovePermission elimina un permiso de un rol y registra auditoría
// @Summary Eliminar permiso de rol
// @Description Elimina la asociación de un permiso específico con un rol
// @Tags Permisos
// @Accept json
// @Produce json
// @Param role_id path int true "ID del rol"
// @Param permissionData body PermissionDataRequest true "Datos del permiso a eliminar"
// @Success 200 {object} map[string]string "message"
// @Failure 400 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
// @Router /roles/{role_id}/permissions [delete]
func RemovePermission(c *gin.Context) {
	var input PermissionDataRequest
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

	// Registrar auditoría
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	audit := models.Audit{
		Event:         "REMOVE_PERMISSION",
		Description:   "Eliminación del permiso ID: " + strconv.Itoa(int(input.PermissionID)) + " del rol ID: " + strconv.Itoa(int(roleID)),
		UserID:        uint(userID.(float64)),
		OriginService: "ROLE_PERMISSION",
		Date:          time.Now(),
	}

	if err := config.DB.Create(&audit).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Permiso eliminado, pero no se pudo registrar la auditoría"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permiso eliminado correctamente"})
}

// GetAllPermissions obtiene todos los permisos con módulos asociados
// @Summary Obtener todos los permisos
// @Description Lista todos los permisos disponibles, incluyendo los módulos a los que pertenecen
// @Tags Permisos
// @Produce json
// @Success 200 {array} PermissionResponse "Lista de permisos"
// @Failure 500 {object} map[string]string "error"
// @Router /permissions/all [get]
func GetAllPermissions(c *gin.Context) {
	var permissions []models.Permission
	if err := config.DB.Preload("Module").Find(&permissions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener permisos"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"permissions": permissions})
}

// GetRolePermissions obtiene los permisos de un rol
// @Summary Obtener permisos de un rol
// @Description Lista todos los permisos asignados a un rol específico
// @Tags Permisos
// @Produce json
// @Param role_id path int true "ID del rol"
// @Success 200 {array} PermissionResponse "Lista de permisos"
// @Failure 400 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
// @Router /roles/{role_id}/permissions [get]
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

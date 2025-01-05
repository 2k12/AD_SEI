package controllers

import (
	"net/http"
	"strconv"
	"time"

	"seguridad-api/config"
	helpers "seguridad-api/helpers"
	"seguridad-api/models"
	"seguridad-api/services"

	"github.com/gin-gonic/gin"
)

// AssignRoleToUser godoc
// @Summary Asigna un rol a un usuario
// @Description Asigna un rol a un usuario dado su ID
// @Tags Roles
// @Accept json
// @Produce json
// @Param id path int true "ID del usuario"
// @Param body body struct{ RoleID uint `json:"role_id"` } true "JSON con el ID del rol"
// @Success 201 {object} models.UserRole
// @Failure 400 {object} map[string]string "ID del usuario o datos inválidos"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /users/{id}/roles [post]
func AssignRoleToUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil || userID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuario inválido"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El usuario no existe"})
		return
	}

	var payload struct {
		RoleID uint `json:"role_id"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	var role models.Role
	if err := config.DB.First(&role, "id = ?", payload.RoleID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El rol no existe"})
		return
	}

	var existingUserRole models.UserRole
	if err := config.DB.First(&existingUserRole, "user_id = ? AND role_id = ?", userID, payload.RoleID).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El usuario ya tiene asignado este rol"})
		return
	}

	userRole, err := services.AssignRoleToUser(uint(userID), payload.RoleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al asignar el rol: " + err.Error()})
		return
	}

	// Auditoría
	currentTime := time.Now()
	ecuadorTime := helpers.AdjustToEcuadorTime(currentTime)

	authenticatedUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	userIDUint := uint(authenticatedUserID.(float64))
	event := "ASSIGN_ROLE"
	description := "Se asignó el rol " + strconv.Itoa(int(payload.RoleID)) + " al usuario " + strconv.Itoa(userID)
	originService := "SEGURIDAD"

	if auditErr := services.RegisterAudit(event, description, userIDUint, originService, ecuadorTime); auditErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Rol asignado, pero no se pudo registrar la auditoría"})
		return
	}

	c.JSON(http.StatusCreated, userRole)
}

// RemoveRoleFromUser godoc
// @Summary Elimina un rol de un usuario
// @Description Elimina un rol asignado a un usuario dado su ID
// @Tags Roles
// @Accept json
// @Produce json
// @Param id path int true "ID del usuario"
// @Param role_id path int true "ID del rol a eliminar"
// @Success 200 {object} map[string]string "Rol eliminado correctamente"
// @Failure 400 {object} map[string]string "ID de usuario o rol inválido"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /users/{id}/roles/{role_id} [delete]
func RemoveRoleFromUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil || userID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuario inválido"})
		return
	}

	roleID, err := strconv.Atoi(c.Param("role_id"))
	if err != nil || roleID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de rol inválido"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El usuario no existe"})
		return
	}

	var role models.Role
	if err := config.DB.First(&role, "id = ?", roleID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El rol no existe"})
		return
	}

	err = services.RemoveRoleFromUser(uint(userID), uint(roleID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el rol: " + err.Error()})
		return
	}

	// Auditoría
	currentTime := time.Now()
	ecuadorTime := helpers.AdjustToEcuadorTime(currentTime)

	authenticatedUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	userIDUint := uint(authenticatedUserID.(float64))
	event := "REMOVE_ROLE"
	description := "Se eliminó el rol " + strconv.Itoa(roleID) + " del usuario " + strconv.Itoa(userID)
	originService := "SEGURIDAD"

	if auditErr := services.RegisterAudit(event, description, userIDUint, originService, ecuadorTime); auditErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Rol eliminado, pero no se pudo registrar la auditoría"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rol eliminado correctamente"})
}

// GetUserRoles godoc
// @Summary Obtener los roles asignados a un usuario
// @Description Devuelve una lista de roles asignados a un usuario dado su ID
// @Tags Roles
// @Produce json
// @Param id path int true "ID del usuario"
// @Success 200 {array} models.Role "Lista de roles asignados"
// @Failure 400 {object} map[string]string "ID de usuario inválido"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /users/{id}/roles [get]
func GetUserRoles(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil || userID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuario inválido"})
		return
	}

	roles, err := services.GetUserRoles(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener los roles"})
		return
	}

	// Auditoría
	currentTime := time.Now()
	ecuadorTime := helpers.AdjustToEcuadorTime(currentTime)

	authenticatedUserID, exists := c.Get("userID")
	if exists {
		userIDUint := uint(authenticatedUserID.(float64))
		event := "GET_USER_ROLES"
		description := "Se consultaron los roles del usuario " + strconv.Itoa(userID)
		originService := "SEGURIDAD"

		_ = services.RegisterAudit(event, description, userIDUint, originService, ecuadorTime)
	}

	c.JSON(http.StatusOK, roles)
}

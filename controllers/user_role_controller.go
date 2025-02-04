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

const ErrInvalidUserID = "ID de usuario inválido"
const QueryByID = "id = ?"

type AssignRolePayload struct {
	RoleID uint `json:"role_id"`
}

type ErrorResponseD struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

type UserRoleResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	RoleID    uint      `json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RoleResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// AssignRoleToUser godoc
// @Summary Asignar un rol a un usuario
// @Description Asigna un rol a un usuario especificado por su ID.
// @Tags Roles
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID del usuario"
// @Param body body AssignRolePayload true "Estructura JSON con el ID del rol"
// @Success 201 {object} UserRoleResponse "Rol asignado correctamente"
// @Failure 400 {object} ErrorResponseD "ID de usuario o datos inválidos"
// @Failure 500 {object} ErrorResponseD "Error interno del servidor"
// @Router /users/{id}/roles [post]
func AssignRoleToUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil || userID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidUserID})
		return
	}

	var user models.User
	if err := config.DB.First(&user, QueryByID, userID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El usuario no existe"})
		return
	}

	var payload AssignRolePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	var role models.Role
	if err := config.DB.First(&role, QueryByID, payload.RoleID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El rol no existe"})
		return
	}

	// var existingUserRole models.UserRole
	// if err := config.DB.First(&existingUserRole, "user_id = ? AND role_id = ?", userID, payload.RoleID).Error; err == nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "El usuario ya tiene asignado este rol"})
	// 	return
	// }

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
	event := "CREATE"
	description := "Se asignó el rol " + strconv.Itoa(int(payload.RoleID)) + " al usuario " + strconv.Itoa(userID)
	originService := "SEGURIDAD"

	if auditErr := services.RegisterAudit(event, description, userIDUint, originService, ecuadorTime); auditErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Rol asignado, pero no se pudo registrar la auditoría"})
		return
	}

	c.JSON(http.StatusCreated, userRole)
}

// RemoveRoleFromUser godoc
// @Summary Eliminar un rol de un usuario
// @Description Elimina un rol previamente asignado a un usuario dado su ID.
// @Tags Roles
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID del usuario"
// @Param role_id path int true "ID del rol a eliminar"
// @Success 200 {object} SuccessResponse "Mensaje de éxito indicando que el rol fue eliminado"
// @Failure 400 {object} ErrorResponseD "ID de usuario o rol inválido"
// @Failure 500 {object} ErrorResponseD "Error interno del servidor"
// @Router /users/{id}/roles/{role_id} [delete]
func RemoveRoleFromUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil || userID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidUserID})
		return
	}

	roleID, err := strconv.Atoi(c.Param("role_id"))
	if err != nil || roleID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de rol inválido"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, QueryByID, userID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El usuario no existe"})
		return
	}

	var role models.Role
	if err := config.DB.First(&role, QueryByID, roleID).Error; err != nil {
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
	event := "DELETE"
	description := "Se eliminó el rol " + strconv.Itoa(roleID) + " del usuario " + strconv.Itoa(userID)
	originService := "SEGURIDAD"

	if auditErr := services.RegisterAudit(event, description, userIDUint, originService, ecuadorTime); auditErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Rol eliminado, pero no se pudo registrar la auditoría"})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "Rol eliminado correctamente"})
}

// GetUserRoles godoc
// @Summary Obtener roles de un usuario
// @Description Devuelve una lista de los roles asignados a un usuario especificado por su ID.
// @Tags Roles
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID del usuario"
// @Success 200 {array} RoleResponse "Lista de roles asignados"
// @Failure 400 {object} ErrorResponseD "ID de usuario inválido"
// @Failure 500 {object} ErrorResponseD "Error interno del servidor"
// @Router /users/{id}/roles [get]
func GetUserRoles(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil || userID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidUserID})
		return
	}

	roles, err := services.GetUserRoles(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener los roles"})
		return
	}

	// // Auditoría
	// currentTime := time.Now()
	// ecuadorTime := helpers.AdjustToEcuadorTime(currentTime)

	// authenticatedUserID, exists := c.Get("userID")
	// if exists {
	// 	userIDUint := uint(authenticatedUserID.(float64))
	// 	event := "GET_USER_ROLES"
	// 	description := "Se consultaron los roles del usuario " + strconv.Itoa(userID)
	// 	originService := "SEGURIDAD"

	// 	_ = services.RegisterAudit(event, description, userIDUint, originService, ecuadorTime)
	// }

	c.JSON(http.StatusOK, roles)
}

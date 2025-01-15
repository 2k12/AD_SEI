package controllers

import (
	"net/http"
	"strconv"
	"time"

	helpers "seguridad-api/helpers"
	"seguridad-api/services"

	"github.com/gin-gonic/gin"
)

type FastUserPayload struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Active   bool   `json:"active"`
	RoleID   uint   `json:"role_id" binding:"required"`
}

// ChargeFastUsers permite cargar usuarios en masa con asignación de roles
func ChargeFastUsers(c *gin.Context) {
	var inputs []FastUserPayload

	// Vincular la carga JSON
	if err := c.ShouldBindJSON(&inputs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Obtener el ID del usuario autenticado
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}
	userIDUint := uint(userID.(float64))

	// Procesar la carga masiva
	currentTime := time.Now()
	ecuadorTime := helpers.AdjustToEcuadorTime(currentTime)

	for _, input := range inputs {
		// Crear usuario
		user, err := services.CreateUser(input.Name, input.Email, input.Password, input.Active)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el usuario " + input.Email + ": " + err.Error()})
			return
		}

		// Auditoría de creación de usuario
		event := "INSERT"
		description := "Se creó un usuario con el email: " + input.Email
		originService := "SEGURIDAD"

		if auditErr := services.RegisterAudit(event, description, userIDUint, originService, ecuadorTime); auditErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Usuario creado, pero no se pudo registrar la auditoría"})
			return
		}

		// Asignar rol al usuario
		_, err = services.AssignRoleToUser(user.ID, input.RoleID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al asignar el rol al usuario " + input.Email + ": " + err.Error()})
			return
		}

		// Auditoría de asignación de rol
		roleDescription := "Se asignó el rol " + strconv.Itoa(int(input.RoleID)) + " al usuario con email: " + input.Email
		if auditErr := services.RegisterAudit(event, roleDescription, userIDUint, originService, ecuadorTime); auditErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Rol asignado, pero no se pudo registrar la auditoría"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuarios y roles cargados exitosamente"})
}

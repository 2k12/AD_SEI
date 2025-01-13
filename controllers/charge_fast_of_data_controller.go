package controllers

import (
	"net/http"
	"seguridad-api/models"
	"seguridad-api/services"

	"github.com/gin-gonic/gin"
)

// ChargeFastOfData permite cargar múltiples permisos rápidamente
// @Summary Carga rápida de permisos
// @Description Permite la creación masiva de permisos en una sola solicitud. Requiere un Bearer Token.
// @Tags Permisos
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body []models.Permission true "Lista de permisos a crear"
// @Success 200 {object} map[string]interface{} "Lista de permisos creados y errores de auditoría (si los hay)"
// @Failure 400 {object} map[string]string "Error en la validación de datos"
// @Failure 401 {object} map[string]string "Error de autorización"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /permissions/fast [post]
func ChargeFastOfData(c *gin.Context) {
	// Recibir múltiples permisos directamente en el modelo Permission
	var inputs []models.Permission

	// Validar el cuerpo de la solicitud
	if err := c.ShouldBindJSON(&inputs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Obtener el userID del contexto
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No se pudo obtener el ID del usuario desde el contexto"})
		return
	}

	// Convertir el userID a uint
	userIDUint := uint(userID.(float64))

	// Llamar al servicio para procesar la carga masiva
	createdPermissions, auditErrors, err := services.ChargeFastPermissions(inputs, userIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar la carga rápida", "details": err.Error()})
		return
	}

	// Responder con los permisos creados y errores de auditoría (si los hay)
	c.JSON(http.StatusOK, gin.H{
		"permissions": createdPermissions,
		"auditErrors": auditErrors,
	})
}

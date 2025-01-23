package controllers

import (
	"net/http"
	"seguridad-api/models"
	"seguridad-api/services"

	"github.com/gin-gonic/gin"
)

// ChargeFastOfData permite cargar múltiples permisos rápidamente
// @Summary Carga rápida de permisos
// @Description Este endpoint permite la creación masiva de permisos en una sola solicitud. Se espera que la solicitud incluya una lista de permisos a crear. Además, se requiere un token de autenticación Bearer para acceder a este endpoint.
// @Tags Permisos
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body []models.Permission true "Lista de permisos a crear"
// @Success 200 {object} map[string]interface{} "Lista de permisos creados y errores de auditoría (si los hay)"
// @Failure 400 {object} map[string]string "Error en la validación de datos"
// @Failure 401 {object} map[string]string "Error de autorización (no se pudo obtener el ID del usuario)"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /permissions/fast [post]
func ChargeFastOfData(c *gin.Context) {
	var inputs []models.Permission

	// Bind de los datos del cuerpo de la solicitud
	if err := c.ShouldBindJSON(&inputs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificar si el ID del usuario está disponible en el contexto
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No se pudo obtener el ID del usuario desde el contexto"})
		return
	}

	// Convertir el ID del usuario a uint
	userIDUint := uint(userID.(float64))

	// Llamada al servicio para procesar la carga rápida de permisos
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

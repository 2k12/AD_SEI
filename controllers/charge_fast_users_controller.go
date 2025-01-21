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

// ChargeFastUsers permite la creación masiva de usuarios y la asignación de roles
// @Summary Carga rápida de usuarios y asignación de roles
// @Description Este endpoint permite la creación masiva de usuarios, incluyendo la asignación de roles a cada uno. Requiere un token de autenticación Bearer para ser utilizado. La solicitud debe incluir una lista de usuarios a crear junto con su rol correspondiente.
// @Tags Usuarios
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body []FastUserPayload true "Lista de usuarios a crear y asignar roles"
// @Success 200 {object} map[string]string "Mensaje de éxito"
// @Failure 400 {object} map[string]string "Error en la validación de datos"
// @Failure 401 {object} map[string]string "Error de autorización (usuario no autenticado)"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /users/fast [post]
func ChargeFastUsers(c *gin.Context) {
	var inputs []FastUserPayload

	// Bind de los datos del cuerpo de la solicitud
	if err := c.ShouldBindJSON(&inputs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificar si el ID del usuario está disponible en el contexto
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}
	// Convertir el ID del usuario a uint
	userIDUint := uint(userID.(float64))

	// Obtener la hora actual ajustada a la zona horaria de Ecuador
	currentTime := time.Now()
	ecuadorTime := helpers.AdjustToEcuadorTime(currentTime)

	// Procesar cada usuario de la lista de entrada
	for _, input := range inputs {
		// Crear el usuario
		user, err := services.CreateUser(input.Name, input.Email, input.Password, input.Active)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el usuario " + input.Email + ": " + err.Error()})
			return
		}

		// Registrar auditoría de creación de usuario
		event := "INSERT"
		description := "Se creó un usuario con el email: " + input.Email
		originService := "SEGURIDAD"
		if auditErr := services.RegisterAudit(event, description, userIDUint, originService, ecuadorTime); auditErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Usuario creado, pero no se pudo registrar la auditoría"})
			return
		}

		// Asignar el rol al usuario
		_, err = services.AssignRoleToUser(user.ID, input.RoleID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al asignar el rol al usuario " + input.Email + ": " + err.Error()})
			return
		}

		// Registrar auditoría de asignación de rol
		roleDescription := "Se asignó el rol " + strconv.Itoa(int(input.RoleID)) + " al usuario con email: " + input.Email
		if auditErr := services.RegisterAudit(event, roleDescription, userIDUint, originService, ecuadorTime); auditErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Rol asignado, pero no se pudo registrar la auditoría"})
			return
		}
	}

	// Responder con mensaje de éxito
	c.JSON(http.StatusOK, gin.H{"message": "Usuarios y roles cargados exitosamente"})
}

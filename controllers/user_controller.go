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

type RegisterUserInput struct {
	Name     string `json:"name" binding:"required" example:"John Doe"`
	Email    string `json:"email" binding:"required" example:"nuevousuario@gmail.com"`
	Password string `json:"password" binding:"required" example:"12345abcd"`
	Active   bool   `json:"active" binding:"required" example:"true"`
	RolId    uint   `json:"rolId" binding:"required" example:"1"`
}

// CreateUser crea un nuevo usuario
// @Summary Crear usuario
// @Description Crea un nuevo usuario con nombre, email, contraseña y estado activo. Requiere un Bearer Token.
// @Tags Usuarios
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body RegisterUserInput true "Datos del usuario"
// @Success 200 {object} map[string]interface{} "user"
// @Failure 400 {object} map[string]string "error"
// @Failure 401 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
// @Router /users [post]
func CreateUser(c *gin.Context) {
	currentTime := time.Now()

	// Ajustar la hora al huso horario de Ecuador usando el helper
	ecuadorTime := helpers.AdjustToEcuadorTime(currentTime)

	// Estructura de entrada para el usuario
	var input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
		Active   bool   `json:"active"`
		RolId    string `json:"rolId" binding:"required"`
	}

	// Vincular datos de entrada
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convertir RolId de string a uint
	rolID, err := strconv.ParseUint(input.RolId, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "RolId debe ser un número válido"})
		return
	}

	// Crear el usuario
	user, err := services.CreateUser(input.Name, input.Email, input.Password, input.Active)
	if err != nil || user.ID == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el usuario"})
		return
	}

	// Obtener el ID del usuario autenticado desde el contexto
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No se pudo obtener el ID del usuario desde el contexto"})
		return
	}

	userIDFloat, ok := userIDInterface.(float64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "El ID del usuario no tiene el formato esperado"})
		return
	}
	userIDUint := uint(userIDFloat)

	// Registrar la auditoría para la creación del usuario
	event := "INSERT"
	description := "Se creó un usuario con el email: " + input.Email
	originService := "SEGURIDAD"
	if auditErr := services.RegisterAudit(event, description, userIDUint, originService, ecuadorTime); auditErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Usuario creado, pero no se pudo registrar la auditoría"})
		return
	}

	// Validar si el rol existe
	var role models.Role
	if err := config.DB.First(&role, "id = ?", uint(rolID)).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El rol no existe"})
		return
	}

	// Verificar si ya existe la asignación del rol al usuario
	var existingUserRole models.UserRole
	if err := config.DB.First(&existingUserRole, "user_id = ? AND role_id = ?", user.ID, uint(rolID)).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El usuario ya tiene asignado este rol"})
		return
	}

	// Asignar el rol al usuario
	userRole, err := services.AssignRoleToUser(user.ID, uint(rolID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al asignar el rol: " + err.Error()})
		return
	}

	// Registrar la auditoría para la asignación del rol
	descriptionInterm := "Se asignó el rol " + strconv.Itoa(int(rolID)) +
		" al usuario " + strconv.Itoa(int(userIDUint))
	eventInterm := "CREATE"
	originServiceInterm := "SEGURIDAD"
	if auditErr := services.RegisterAudit(eventInterm, descriptionInterm, userIDUint, originServiceInterm, ecuadorTime); auditErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Rol asignado, pero no se pudo registrar la auditoría"})
		return
	}

	// Respuesta final
	c.JSON(http.StatusOK, gin.H{
		"message": "Usuario creado y rol asignado con éxito",
		"user":    user,
		"role":    userRole,
	})
}

// @Summary Obtener usuarios
// @Description Devuelve una lista de usuarios paginada y filtrada.
// @Tags Usuarios
// @Security BearerAuth
// @Produce json
// @Param page query int false "Número de página (por defecto: 1)"
// @Param pageSize query int false "Tamaño de página (por defecto: 10)"
// @Param email query string false "Filtrar por email"
// @Param active query bool false "Filtrar por estado"
// @Success 200 {object} map[string]interface{} "users"
// @Failure 401 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
// @Router /users [get]
func GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	email := c.Query("email")
	active := c.Query("active")

	// Construir filtros
	filters := make(map[string]interface{})
	if email != "" {
		filters["email"] = email
	}
	if active != "" {
		activeBool, _ := strconv.ParseBool(active)
		filters["active"] = activeBool
	}

	users, total, err := services.GetPaginatedUsers(page, pageSize, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los usuarios"})
		return
	}

	// Respuesta con datos paginados
	c.JSON(http.StatusOK, gin.H{
		"users":      users,
		"total":      total,
		"page":       page,
		"pageSize":   pageSize,
		"totalPages": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// GetUserPermissions obtiene los permisos de un usuario específico
// @Summary Obtener permisos de un usuario
// @Description Devuelve la lista de permisos asignados a un usuario específico, dado su ID. Requiere un Bearer Token.
// @Tags Usuarios
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID del usuario"
// @Success 200 {array} map[string]interface{} "permissions"
// @Failure 400 {object} map[string]string "error"
// @Failure 401 {object} map[string]string "error"
// @Failure 404 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
// @Router /users/{id}/permissions [get]
func GetUserPermissions(c *gin.Context) {
	userID := c.Param("id")

	userIDInt, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuario no válido"})
		return
	}

	permissions, err := services.GetUserPermissions(uint(userIDInt))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los permisos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"permissions": permissions})
}

// UpdateUser actualiza un usuario existente
// @Summary Actualizar usuario
// @Description Actualiza los datos de un usuario existente. Requiere un Bearer Token.
// @Tags Usuarios
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID del usuario"
// @Param userData body object true "Datos a actualizar"
// @Success 200 {object} map[string]interface{} "updatedUser"
// @Failure 400 {object} map[string]string "error"
// @Failure 401 {object} map[string]string "error"
// @Router /users/{id} [put]
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	currentTime := time.Now()

	ecuadorTime := helpers.AdjustToEcuadorTime(currentTime)

	var userData struct {
		Name   string `json:"name"`
		Email  string `json:"email"`
		Active *bool  `json:"active"`
	}

	if err := c.ShouldBindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	updatedUser, err := services.UpdateUser(id, userData.Name, userData.Email, userData.Active)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No se pudo obtener el ID del usuario desde el contexto"})
		return
	}

	userIDUint := uint(userID.(float64))

	event := "UPDATE"
	description := "Se actualizó el usuario con ID: " + id
	originService := "SEGURIDAD"
	// date := time.Now()

	if auditErr := services.RegisterAudit(event, description, userIDUint, originService, ecuadorTime); auditErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Usuario actualizado, pero no se pudo registrar la auditoría"})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

// DeleteUser elimina un usuario
// @Summary Cambiar el estado del usuario
// @Description Cambia el estado de un usuario a inactivo. Requiere un Bearer Token.
// @Tags Usuarios
// @Security BearerAuth
// @Produce json
// @Param id path string true "ID del usuario"
// @Success 200 {object} map[string]string "message"
// @Failure 401 {object} map[string]string "error"
// @Failure 404 {object} map[string]string "error"
// @Router /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	authenticatedUserID, exists := c.Get("userID")
	currentTime := time.Now()

	ecuadorTime := helpers.AdjustToEcuadorTime(currentTime)

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	id := c.Param("id")

	if err := services.DeleteUser(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	event := "DELETE"
	description := "Se cambió el estado del usuario con ID: " + id
	originService := "SEGURIDAD"
	// date := time.Now()

	authUserID, ok := authenticatedUserID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener el ID del usuario autenticado"})
		return
	}
	if auditErr := services.RegisterAudit(event, description, authUserID, originService, ecuadorTime); auditErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Usuario eliminado, pero no se pudo registrar la auditoría"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "El estado del Usuario fue cambiado exitosamente"})
}

// GetUsersforDropdown obtiene todos los usuarios para ser utilizados en un dropdown
// @Summary Obtener usuarios para dropdown
// @Description Devuelve una lista de usuarios que pueden ser utilizados en un dropdown o lista de selección.
// @Tags Usuarios
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{} "users"
// @Failure 500 {object} map[string]string "error"
// @Router /users/dropdown [get]
func GetUsersforDropdown(c *gin.Context) {
	users, err := services.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los usuarios"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

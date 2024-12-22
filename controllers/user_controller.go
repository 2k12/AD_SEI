package controllers

import (
	"net/http"
	"seguridad-api/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateUser crea un nuevo usuario
// @Summary Crear usuario
// @Description Crea un nuevo usuario con nombre, email, contraseña y estado activo. Requiere un Bearer Token.
// @Tags Usuarios
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body object true "Datos del usuario"
// @Success 200 {object} map[string]interface{} "user"
// @Failure 400 {object} map[string]string "error"
// @Failure 401 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
// @Router /users [post]
func CreateUser(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
		Active   bool   `json:"active"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := services.CreateUser(input.Name, input.Email, input.Password, input.Active)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el usuario"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// // GetUsers obtiene todos los usuarios
// // @Summary Obtener usuarios
// // @Description Devuelve una lista de todos los usuarios registrados. Requiere un Bearer Token y PIN.
// // @Tags Usuarios
// // @Security BearerAuth
// // @Security XAPI-PIN
// // @Produce json
// // @Success 200 {array} map[string]interface{} "users"
// // @Failure 401 {object} map[string]string "error"
// // @Failure 500 {object} map[string]string "error"
// // @Router /users [get]
// @Summary Obtener usuarios
// @Description Devuelve una lista de usuarios paginada y filtrada.
// @Tags Usuarios
// @Security BearerAuth
// @Produce json
// @Param page query int false "Número de página (por defecto: 1)"
// @Param pageSize query int false "Tamaño de página (por defecto: 10)"
// @Param name query string false "Filtrar por nombre"
// @Param email query string false "Filtrar por email"
// @Param active query bool false "Filtrar por estado activo"
// @Success 200 {object} map[string]interface{} "users"
// @Failure 401 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
// @Router /users [get]

// func GetUsers(c *gin.Context) {
// 	users, err := services.GetUsers()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los usuarios"})
// 		return
// 	}

//		c.JSON(http.StatusOK, gin.H{"users": users})
//	}
func GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	name := c.Query("name")
	email := c.Query("email")
	active := c.Query("active")

	// Construir filtros
	filters := make(map[string]interface{})
	if name != "" {
		filters["name"] = name
	}
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
	userID := c.Param("id") // Obtener el ID del usuario desde los parámetros de la ruta

	// Convertir el ID a uint
	userIDInt, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuario no válido"})
		return
	}

	// Obtener los permisos del usuario
	permissions, err := services.GetUserPermissions(uint(userIDInt))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los permisos"})
		return
	}

	// Devolver la respuesta con los permisos
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

	c.JSON(http.StatusOK, updatedUser)
}

// DeleteUser elimina un usuario
// @Summary Eliminar usuario
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
	id := c.Param("id")

	if err := services.DeleteUser(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "El estado del Usuario fue cambiado exitosamente"})
}

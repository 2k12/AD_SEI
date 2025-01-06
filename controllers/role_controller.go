package controllers

import (
	"net/http"
	helpers "seguridad-api/helpers"
	"seguridad-api/services"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateRoleInput struct {
	Event       string `json:"event" binding:"required" example:"INSERT"`
	Name        string `json:"name" binding:"required" example:"Administrador"`
	Description string `json:"description" example:"Rol con permisos de administración"`
	Active      bool   `json:"active" example:"true"`
}

type CreateRoleResponse struct {
	Message string `json:"message" example:"Rol creado exitosamente"`
}

type ErrorResponseRole struct {
	Error string `json:"error" example:"Error al crear el rol"`
}

// CreateRole crea un nuevo rol
// @Summary Crear rol
// @Description Crea un nuevo rol con nombre, descripción y estado activo. Requiere un Bearer Token.
// @Tags Roles
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body CreateRoleInput true "Datos del rol"
// @Success 200 {object} CreateRoleResponse "Rol registrado exitosamente"
// @Failure 400 {object} ErrorResponseRole "Datos inválidos o formato incorrecto"
// @Failure 500 {object} ErrorResponseRole "Error al registrar el Rol"
// @Router /roles [post]
func CreateRole(c *gin.Context) {

	currentTime := time.Now()
	ecuadorTime := helpers.AdjustToEcuadorTime(currentTime)

	var input struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Active      bool   `json:"active"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role, err := services.CreateRole(input.Name, input.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el rol"})
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

	event := "INSERT"
	description := "Se creó un rol con el nombre: " + input.Name
	originService := "SEGURIDAD"
	// date := time.Now()

	if auditErr := services.RegisterAudit(event, description, userIDUint, originService, ecuadorTime); auditErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Rol creado, pero no se pudo registrar la auditoría"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"role": role})
}

type GetRolesQueryParams struct {
	Page     int    `json:"page" example:"1"`
	PageSize int    `json:"pageSize" example:"10"`
	Name     string `json:"name" example:"Administrador"`
	Active   string `json:"active" example:"true"`
}

type Role struct {
	ID          uint   `json:"id" example:"1"`
	Name        string `json:"name" example:"Administrador"`
	Description string `json:"description" example:"Gestión de usuarios"`
	Active      bool   `json:"active" example:"true"`
}

type GetRolesResponse struct {
	Roles      []Role `json:"roles"`
	Total      int64  `json:"total" example:"15"`
	Page       int    `json:"page" example:"1"`
	PageSize   int    `json:"pageSize" example:"10"`
	TotalPages int64  `json:"totalPages" example:"2"`
}

type ErrorResponseGetRoles struct {
	Error string `json:"error" example:"Error al obtener los roles"`
}

// GetRoles obtiene la lista de roles con paginación y filtros opcionales
// @Summary Obtener roles
// @Description Devuelve una lista paginada de roles, permitiendo filtrar por nombre y estado activo. Los resultados pueden ser paginados utilizando los parámetros `page` y `pageSize`.
// @Tags Roles
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param page query int false "Número de página (por defecto 1)"
// @Param pageSize query int false "Tamaño de página (por defecto 10)"
// @Param name query string false "Filtrar por nombre del rol"
// @Param active query string false "Filtrar por estado activo (true/false)"
// @Success 200 {object} GetRolesResponse "Roles obtenidos exitosamente"
// @Failure 400 {object} ErrorResponseGetRoles "Datos inválidos o formato incorrecto"
// @Failure 401 {object} ErrorResponseGetRoles "No autorizado, se requiere autenticación"
// @Failure 500 {object} ErrorResponseGetRoles "Error interno del servidor al intentar obtener los roles"
// @Router /roles [get]
func GetRoles(c *gin.Context) {
	// Obtener parámetros de consulta para paginación y filtros
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	name := c.Query("name")
	active := c.Query("active")

	// Construir filtros
	filters := make(map[string]interface{})
	if name != "" {
		filters["name"] = name
	}
	if active != "" {
		activeBool, _ := strconv.ParseBool(active)
		filters["active"] = activeBool
	}

	// Llamar al servicio para obtener roles paginados
	roles, total, err := services.GetPaginatedRoles(page, pageSize, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los roles"})
		return
	}

	// Respuesta con datos paginados
	c.JSON(http.StatusOK, gin.H{
		"roles":      roles,
		"total":      total,
		"page":       page,
		"pageSize":   pageSize,
		"totalPages": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

type UpdateRoleInput struct {
	Name        string `json:"name" binding:"required" example:"Administrador"`
	Description string `json:"description" example:"Rol para gestionar usuarios y permisos"`
	Active      bool   `json:"active" example:"true"`
}

type UpdateRoleResponse struct {
	RoleID      uint   `json:"role_id" example:"1"`
	Name        string `json:"name" example:"Administrador"`
	Description string `json:"description" example:"Rol para gestionar usuarios y permisos"`
	Active      bool   `json:"active" example:"true"`
	Status      string `json:"status" example:"updated"`
}

type ErrorResponseUpdateRole struct {
	Error string `json:"error" example:"Error al actualizar el rol"`
}

// UpdateRole actualiza la información de un rol existente
// @Summary Actualizar rol
// @Description Permite actualizar los datos de un rol existente, como el nombre, descripción y estado activo. Se requiere un Bearer Token para la autenticación.
// @Tags Roles
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID del rol a actualizar"
// @Param input body UpdateRoleInput true "Estructura que contiene los datos actualizados del rol"
// @Success 200 {object} UpdateRoleResponse "Rol actualizado exitosamente"
// @Failure 400 {object} ErrorResponseUpdateRole "Datos inválidos o formato incorrecto"
// @Failure 401 {object} ErrorResponseUpdateRole "No autorizado, falta el token de autenticación"
// @Failure 404 {object} ErrorResponseUpdateRole "Rol no encontrado"
// @Failure 500 {object} ErrorResponseUpdateRole "Error interno del servidor al intentar actualizar el rol"
// @Router /roles/{id} [put]
func UpdateRole(c *gin.Context) {
	currentTime := time.Now()
	ecuadorTime := helpers.AdjustToEcuadorTime(currentTime)
	// Obtener el ID del rol desde los parámetros de la URL
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var input struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Active      bool   `json:"active"`
	}

	// Validar y enlazar los datos del cuerpo de la solicitud
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Llamar al servicio para actualizar el rol
	role, err := services.UpdateRole(id, input.Name, input.Description, input.Active)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

	event := "UPDATE"
	description := "Se actualizó el rol con ID: " + c.Param("id")
	originService := "SEGURIDAD"
	// date := time.Now()

	if auditErr := services.RegisterAudit(event, description, userIDUint, originService, ecuadorTime); auditErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Rol actualizado, pero no se pudo registrar la auditoría"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"role": role})
}

type UpdateRoleStateInput struct {
	Active bool `json:"active" example:"true"`
}

type UpdateRoleStateResponse struct {
	Message string `json:"message" example:"Estado del rol actualizado exitosamente"`
}

type ErrorResponseUpdateRoleState struct {
	Error string `json:"error" example:"Error al actualizar el estado del rol"`
}

// UpdateRoleState actualiza el estado de un rol
// @Summary Actualizar estado del rol
// @Description Actualiza el estado de un rol existente (activo/inactivo)
// @Tags Roles
// @Accept json
// @Produce json
// @Param id path int true "ID del rol a actualizar"
// @Param roleState body UpdateRoleStateInput true "Estado del rol a actualizar"
// @Success 200 {object} UpdateRoleStateResponse "Estado del rol actualizado exitosamente"
// @Failure 400 {object} ErrorResponseUpdateRoleState "Datos inválidos o ID incorrecto"
// @Failure 401 {object} ErrorResponseUpdateRoleState "No se pudo obtener el ID del usuario desde el contexto"
// @Failure 500 {object} ErrorResponseUpdateRoleState "Error interno al actualizar el estado del rol o registrar la auditoría"
// @Router /roles/{id}/state [patch]
func UpdateRoleState(c *gin.Context) {
	currentTime := time.Now()
	ecuadorTime := helpers.AdjustToEcuadorTime(currentTime)

	// Obtener el ID del rol desde los parámetros de la URL
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var input struct {
		Active bool `json:"active"` // Solo necesitamos el estado
	}

	// Validar y enlazar los datos del cuerpo de la solicitud
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Llamar al servicio para actualizar solo el estado del rol
	err = services.UpdateRoleState(id, input.Active)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

	event := "UPDATE"
	description := "Se actualizó el estado del rol con ID: " + c.Param("id")
	originService := "SEGURIDAD"
	// date := time.Now()

	if auditErr := services.RegisterAudit(event, description, userIDUint, originService, ecuadorTime); auditErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Estadi del rol actualizado, pero no se pudo registrar la auditoría"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Estado del rol actualizado exitosamente"})
}

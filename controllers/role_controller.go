package controllers

import (
	"net/http"
	helpers "seguridad-api/helpers"
	"seguridad-api/services"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateRole crea un nuevo rol
// @Summary Crear rol
// @Description Crea un nuevo rol con nombre, descripción y estado activo. Requiere un Bearer Token.
// @Tags Roles
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body object true "Datos del rol"
// @Success 200 {object} map[string]interface{} "role"
// @Failure 400 {object} map[string]string "error"
// @Failure 401 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
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

// GetRoles obtiene la lista de roles con paginación y filtros opcionales
// @Summary Obtener roles
// @Description Devuelve una lista paginada de roles, permitiendo filtrar por nombre y estado activo.
// @Tags Roles
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param page query int false "Número de página (por defecto 1)"
// @Param pageSize query int false "Tamaño de página (por defecto 10)"
// @Param name query string false "Filtrar por nombre"
// @Param active query boolean false "Filtrar por estado activo"
// @Success 200 {object} map[string]interface{} "roles"
// @Failure 400 {object} map[string]string "error"
// @Failure 401 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
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

// UpdateRole actualiza la información de un rol existente
// @Summary Actualizar rol
// @Description Actualiza los datos de un rol existente identificándolo por su ID. Requiere un Bearer Token.
// @Tags Roles
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID del rol"
// @Param input body object true "Datos actualizados del rol"
// @Success 200 {object} map[string]interface{} "role"
// @Failure 400 {object} map[string]string "error"
// @Failure 401 {object} map[string]string "error"
// @Failure 404 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
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

// UpdateRoleState actualiza el estado a activo-inactivo de un rol
// @Summary Actualizar estado del rol
// @Description Cambia únicamente el estado activo de un rol identificado por su ID. Requiere un Bearer Token.
// @Tags Roles
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID del rol"
// @Param input body object true "Estado del rol"
// @Success 200 {object} map[string]string "message"
// @Failure 400 {object} map[string]string "error"
// @Failure 401 {object} map[string]string "error"
// @Failure 404 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
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

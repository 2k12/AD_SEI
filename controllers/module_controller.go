package controllers

import (
	"net/http"
	helpers "seguridad-api/helpers"
	"seguridad-api/services"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateModule crea un nuevo módulo
// @Summary Crear módulo
// @Description Crea un nuevo módulo con nombre, descripción y estado activo. Requiere un Bearer Token.
// @Tags Módulos
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body models.Module true "Datos del módulo"
// @Success 200 {object} map[string]string "message"
// @Failure 400 {object} map[string]string "error"
// @Failure 401 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
// @Router /modules [post]
func CreateModule(c *gin.Context) {
	var input struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Active      bool   `json:"active"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := services.CreateModule(input.Name, input.Description, input.Active)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create module"})
		return
	}

	currentTime := helpers.AdjustToEcuadorTime(time.Now())
	userID, exists := c.Get("userID")
	if exists {
		userIDUint := uint(userID.(float64))
		event := "INSERT"
		description := "Se creó un módulo con el nombre: " + input.Name
		_ = services.RegisterAudit(event, description, userIDUint, "SEGURIDAD", currentTime)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Module created successfully"})
}

// GetModules obtiene todos los módulos
// @Summary Obtener módulos
// @Description Devuelve una lista de todos los módulos. Requiere un Bearer Token.
// @Tags Módulos
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.Module "modules"
// @Failure 401 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
// @Router /modules [get]
func GetModules(c *gin.Context) {
	// // Obtén los parámetros de paginación de la consulta
	// page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	// limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	// offset := (page - 1) * limit

	// var modules []models.Module
	// var total int64

	// // Aplica paginación y busca los módulos
	// err := config.DB.Offset(offset).Limit(limit).Find(&modules).Error
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch modules"})
	// 	return
	// }

	// // Cuenta el total de módulos
	// err = config.DB.Model(&models.Module{}).Count(&total).Error
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to count modules"})
	// 	return
	// }

	// // Responde con los datos de los módulos y la información de paginación
	// c.JSON(http.StatusOK, gin.H{
	// 	"modules": modules,
	// 	"page":    page,
	// 	"limit":   limit,
	// 	"total":   total,
	// })
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	name := c.Query("name")

	// Construir filtros
	filters := make(map[string]interface{})
	if name != "" {
		filters["name"] = name
	}

	modules, total, err := services.GetPaginatedModules(page, pageSize, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los modulos"})
		return
	}

	// Respuesta con datos paginados
	c.JSON(http.StatusOK, gin.H{
		"modules":    modules,
		"total":      total,
		"page":       page,
		"pageSize":   pageSize,
		"totalPages": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// GetModule obtiene los detalles de un módulo por ID
// @Summary Obtener módulo por ID
// @Description Devuelve los detalles de un módulo dado su ID. Requiere un Bearer Token.
// @Tags Módulos
// @Security BearerAuth
// @Produce json
// @Param id path string true "ID del módulo"
// @Success 200 {object} models.Module "module"
// @Failure 401 {object} map[string]string "error"
// @Failure 404 {object} map[string]string "error"
// @Router /modules/{id} [get]
func GetModule(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	module, err := services.GetModuleByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Module not found"})
		return
	}

	c.JSON(http.StatusOK, module)
}

// UpdateModule actualiza un módulo existente
// @Summary Actualizar módulo
// @Description Actualiza los datos de un módulo existente. Requiere un Bearer Token.
// @Tags Módulos
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID del módulo"
// @Param input body models.Module true "Datos del módulo"
// @Success 200 {object} map[string]string "message"
// @Failure 400 {object} map[string]string "error"
// @Failure 401 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
// @Router /modules/{id} [put]
func UpdateModule(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var input struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
		Active      *bool   `json:"active"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := services.UpdateModule(uint(id), input.Name, input.Description, input.Active)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	currentTime := helpers.AdjustToEcuadorTime(time.Now())
	userID, exists := c.Get("userID")
	if exists {
		userIDUint := uint(userID.(float64))
		event := "UPDATE"
		description := "Se actualizó el módulo con ID: " + strconv.Itoa(id)
		_ = services.RegisterAudit(event, description, userIDUint, "MÓDULOS", currentTime)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Module updated successfully"})
}

// DeleteModule elimina un módulo
// @Summary Eliminar módulo
// @Description Elimina un módulo existente. Requiere un Bearer Token.
// @Tags Módulos
// @Security BearerAuth
// @Produce json
// @Param id path string true "ID del módulo"
// @Success 200 {object} map[string]string "message"
// @Failure 401 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
// @Router /modules/{id} [delete]
func DeleteModule(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	err := services.DeleteModule(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to delete module"})
		return
	}

	currentTime := helpers.AdjustToEcuadorTime(time.Now())
	userID, exists := c.Get("userID")
	if exists {
		userIDUint := uint(userID.(float64))
		event := "DELETE"
		description := "Se eliminó el módulo con ID: " + strconv.Itoa(id)
		_ = services.RegisterAudit(event, description, userIDUint, "SEGURIDAD", currentTime)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Module deleted successfully"})
}

// func ToggleModuleActive(c *gin.Context) {
// 	id, _ := strconv.Atoi(c.Param("id"))

// 	err := services.ToggleModuleActive(uint(id))
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to change module active state"})
// 		return
// 	}

// 	currentTime := helpers.AdjustToEcuadorTime(time.Now())
// 	userID, exists := c.Get("userID")
// 	if exists {
// 		userIDUint := uint(userID.(float64))
// 		event := "UPDATE"
// 		description := "Se cambió el estado activo del módulo con ID: " + strconv.Itoa(id)
// 		_ = services.RegisterAudit(event, description, userIDUint, "SEGURIDAD", currentTime)
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Module state changed successfully"})
// }

func GetModuleActive(c *gin.Context) {
	// Llamar al servicio para obtener los módulos activos
	modules, err := services.GetActiveModules()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los módulos"})
		return
	}

	// Validar si no hay módulos activos
	if len(modules) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No se encontraron módulos activos"})
		return
	}

	// Responder con los módulos activos
	c.JSON(http.StatusOK, modules)
}

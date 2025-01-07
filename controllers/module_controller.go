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
	modules, err := services.GetAllModules()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch modules"})
		return
	}

	c.JSON(http.StatusOK, modules)
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

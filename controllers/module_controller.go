package controllers

import (
	"net/http"
	"seguridad-api/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Crear módulo
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

	c.JSON(http.StatusOK, gin.H{"message": "Module created successfully"})
}

// Obtener todos los módulos
func GetModules(c *gin.Context) {
	modules, err := services.GetAllModules()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch modules"})
		return
	}

	c.JSON(http.StatusOK, modules)
}

// Obtener un módulo por ID
func GetModule(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	module, err := services.GetModuleByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Module not found"})
		return
	}

	c.JSON(http.StatusOK, module)
}

// Actualizar módulo parcialmente
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

	// Llamar al servicio para actualizar el módulo
	err := services.UpdateModule(uint(id), input.Name, input.Description, input.Active)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Module updated successfully"})
}

// Cambiar el estado de un módulo (activar/desactivar)
func ToggleModuleActive(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	// Llamar al servicio para cambiar el estado de "active"
	err := services.ToggleModuleActive(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to change module active state"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Module state changed successfully"})
}
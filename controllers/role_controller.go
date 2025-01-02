package controllers

import (
	"net/http"
	"seguridad-api/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateRole(c *gin.Context) {
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

	c.JSON(http.StatusOK, gin.H{"role": role})
}

func GetRoles(c *gin.Context) {
	roles, err := services.GetRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los roles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"roles": roles})
}

func UpdateRole(c *gin.Context) {
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

	c.JSON(http.StatusOK, gin.H{"role": role})
}

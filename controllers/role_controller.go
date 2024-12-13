package controllers

import (
	"net/http"
	"seguridad-api/services"

	"github.com/gin-gonic/gin"
)

func CreateRole(c *gin.Context) {
	var input struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Llamar al servicio para crear el rol
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

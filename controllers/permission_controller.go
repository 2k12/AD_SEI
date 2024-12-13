package controllers

import (
	"net/http"
	"seguridad-api/services"

	"github.com/gin-gonic/gin"
)

func CreatePermission(c *gin.Context) {
	var input struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	permission, err := services.CreatePermission(input.Name, input.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el permiso"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"permission": permission})
}

func GetPermissions(c *gin.Context) {
	permissions, err := services.GetPermissions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los permisos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"permissions": permissions})
}

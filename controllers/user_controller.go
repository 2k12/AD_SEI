package controllers

import (
	"net/http"
	"seguridad-api/services"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
		// RoleID   uint   `json:"role_id" binding:"required"` // Recibir role_id
		Active bool `json:"active"`
	}

	// Validar los datos de entrada
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Llamar al servicio para crear el usuario
	user, err := services.CreateUser(input.Name, input.Email, input.Password, input.Active)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el usuario"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func GetUsers(c *gin.Context) {
	users, err := services.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los usuarios"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

// Actualizar usuario
func UpdateUser(c *gin.Context) {
	id := c.Param("id") // Obtener el ID del usuario desde la URL

	var userData struct {
		Name   string `json:"name"`   // Campo opcional
		Email  string `json:"email"`  // Campo opcional
		Active *bool  `json:"active"` // Campo opcional
	}

	// Parsear el JSON del body
	if err := c.ShouldBindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// Llamar al servicio para actualizar el usuario
	updatedUser, err := services.UpdateUser(id, userData.Name, userData.Email, userData.Active)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Responder con el usuario actualizado
	c.JSON(http.StatusOK, updatedUser)
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	if err := services.DeleteUser(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "El estado del Usuario fue cambiado exitosamente"})
}

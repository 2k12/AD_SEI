package controllers

import (
	"net/http"
	"seguridad-api/services"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func RegisterAudit(c *gin.Context) {
	var input struct {
		Event         string `json:"event" binding:"required"`
		Description   string `json:"description" binding:"required"`
		UserID        string `json:"user_id" binding:"required"`
		OriginService string `json:"origin_service" binding:"required"`
		Date          string `json:"date" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convertir UserID de string a uint
	userID, err := strconv.ParseUint(input.UserID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El ID del usuario debe ser un número válido"})
		return
	}

	// Convertir Date de string a time.Time
	// Supón que la fecha llega en formato "2006-01-02T15:04:05Z07:00" o el formato que estés utilizando
	date, err := time.Parse(time.RFC3339, input.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El formato de la fecha es inválido"})
		return
	}

	// Llamar al servicio para registrar la auditoría
	if err := services.RegisterAudit(input.Event, input.Description, uint(userID), input.OriginService, date); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al registrar la auditoría"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Auditoría registrada exitosamente"})
}

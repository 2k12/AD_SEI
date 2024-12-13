package services

import (
	"seguridad-api/config"
	"seguridad-api/models"
	"time"
)

func RegisterAudit(event, description string, userID uint, originService string, date time.Time) error {
	audit := models.Audit{
		Event:         event,
		Description:   description,
		UserID:        userID,
		OriginService: originService,
		Date:          date,
	}

	result := config.DB.Create(&audit)
	return result.Error
}

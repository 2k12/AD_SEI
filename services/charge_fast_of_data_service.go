package services

import (
	"seguridad-api/config"
	"seguridad-api/models"
	"time"
)

func ChargeFastPermissions(inputs []models.Permission, userID uint) ([]models.Permission, []string, error) {
	var createdPermissions []models.Permission
	var auditErrors []string

	for _, input := range inputs {
		input.ID = 0
		input.CreatedAt = time.Now()
		input.UpdatedAt = time.Now()

		if err := config.DB.Create(&input).Error; err != nil {
			return nil, nil, err
		}

		createdPermissions = append(createdPermissions, input)

		event := "INSERT"
		description := "Se creó un permiso con el nombre: " + input.Name
		audit := models.Audit{
			Event:         event,
			Description:   description,
			UserID:        userID,
			OriginService: "SEGURIDAD",
			Date:          time.Now(),
		}

		if err := config.DB.Create(&audit).Error; err != nil {
			auditErrors = append(auditErrors, "Error en auditoría para permiso: "+input.Name)
		}
	}

	return createdPermissions, auditErrors, nil
}

package services

import (
	"seguridad-api/config"
	"seguridad-api/models"
	"time"
)

// ChargeFastPermissions procesa una lista de permisos para una carga masiva
func ChargeFastPermissions(inputs []models.Permission, userID uint) ([]models.Permission, []string, error) {
	var createdPermissions []models.Permission
	var auditErrors []string

	// Procesar cada permiso recibido
	for _, input := range inputs {
		// Reiniciar campos generados automáticamente
		input.ID = 0
		input.CreatedAt = time.Now()
		input.UpdatedAt = time.Now()

		// Guardar el permiso en la base de datos
		if err := config.DB.Create(&input).Error; err != nil {
			return nil, nil, err // Detener el proceso si ocurre un error
		}

		// Agregar a la lista de permisos creados
		createdPermissions = append(createdPermissions, input)

		// Registrar auditoría
		event := "INSERT"
		description := "Se creó un permiso con el nombre: " + input.Name
		audit := models.Audit{
			Event:         event,
			Description:   description,
			UserID:        userID,
			OriginService: "SEGURIDAD",
			Date:          time.Now(), // Usar el campo Date en lugar de EventTime
		}

		// Guardar la auditoría en la base de datos
		if err := config.DB.Create(&audit).Error; err != nil {
			// Registrar errores de auditoría pero continuar con el proceso principal
			auditErrors = append(auditErrors, "Error en auditoría para permiso: "+input.Name)
		}
	}

	return createdPermissions, auditErrors, nil
}

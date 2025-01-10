package services

import (
	"bytes"
	"fmt"
	"seguridad-api/config"
	"seguridad-api/models"
	"seguridad-api/utils"
	"time"
)

func GenerateReport(all bool) (*bytes.Buffer, string, error) {
	var headers []string
	var data [][]string
	var state string

	if all {
		// Reporte de usuarios
		headers = []string{"Nombre", "Correo Electrónico", "Estado", "Clave del Módulo", "F. Creación", "F. Actualización"}
		var users []models.User
		if err := config.DB.Find(&users).Error; err != nil {
			return nil, "", fmt.Errorf("error al consultar los usuarios: %w", err)
		}

		for _, user := range users {

			if user.Active {
				state = "Activo"
			} else {
				state = "Inactivo"
			}

			row := []string{
				user.Name,
				user.Email,
				// fmt.Sprintf("%t", user.Active),
				state,
				user.ModuleKey,
				user.CreatedAt.Format("2006-01-02 15:04:05"),
				user.UpdatedAt.Format("2006-01-02 15:04:05"),
			}
			data = append(data, row)
		}
	}

	title := "Reporte de Usuarios"

	fileBuffer, err := utils.GeneratePDF(title, "", data, headers)
	if err != nil {
		return nil, "", fmt.Errorf("error al generar el PDF: %w", err)
	}

	fileName := "reporte_" + time.Now().Format("20060102_150405") + ".pdf"
	return fileBuffer, fileName, nil
}

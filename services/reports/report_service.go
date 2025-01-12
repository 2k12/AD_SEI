package services

import (
	"bytes"
	"fmt"
	"seguridad-api/config"
	"seguridad-api/models"
	"seguridad-api/utils"
	"time"
)

func GenerateReport(modelName string, filters map[string]interface{}, userName string, format string, option string) (*bytes.Buffer, string, error) {
	var headers []string
	var data [][]string
	var state string

	var query interface{}

	switch modelName {
	case "User":
		if option == "usuariosCompletos" {
			headers = []string{"Nombre", "Roles", "Permisos", "Módulos"}
			var users []models.User
			result := config.DB.Debug().
				Preload("Roles.Permissions.Module").
				Where(filters).
				Find(&users)
			if result.Error != nil {
				return nil, "", fmt.Errorf("error al consultar los datos: %w", result.Error)
			}
			query = users
		} else {
			headers = []string{"Nombre", "Correo Electrónico", "Estado", "Clave del Módulo", "F. Creación", "F. Actualización"}
			query = &[]models.User{}
			dbQuery := config.DB.Model(query)
			for key, value := range filters {
				dbQuery = dbQuery.Where(fmt.Sprintf("%s = ?", key), value)
			}

			if err := dbQuery.Find(query).Error; err != nil {
				return nil, "", fmt.Errorf("error al consultar los datos: %w", err)
			}
		}
	default:
		return nil, "", fmt.Errorf("modelo no soportado")
	}

	var users []models.User
	switch v := query.(type) {
	case []models.User:
		users = v
	case *[]models.User:
		users = *v
	default:
		return nil, "", fmt.Errorf("tipo de 'query' no reconocido")
	}

	for _, user := range users {
		row := []string{}
		if user.Active {
			state = "Activo"
		} else {
			state = "Inactivo"
		}

		if option == "usuariosCompletos" {
			roles := ""
			permissions := ""
			modules := ""
			modulesMap := make(map[string]bool)

			for _, role := range user.Roles {
				roles += role.Name + ", "
				for _, permission := range role.Permissions {
					permissions += permission.Name + ", "

					if !modulesMap[permission.Module.Name] {
						modulesMap[permission.Module.Name] = true
						modules += permission.Module.Name + ", "
					}
				}
			}

			if len(roles) > 0 {
				roles = roles[:len(roles)-2]
			}
			if len(permissions) > 0 {
				permissions = permissions[:len(permissions)-2]
			}
			if len(modules) > 0 {
				modules = modules[:len(modules)-2]
			}

			row = append(row, user.Name, roles, permissions, modules)
		} else {
			row = append(row,
				user.Name,
				user.Email,
				state,
				user.ModuleKey,
				user.CreatedAt.Format("2006-01-02 15:04:05"),
				user.UpdatedAt.Format("2006-01-02 15:04:05"),
			)
		}

		data = append(data, row)
	}

	// Generar el título y los filtros formateados
	title := fmt.Sprintf("Reporte de %s", modelName)
	formattedFilters := formatFilters(filters)

	// Generar el archivo en el formato solicitado
	var fileBuffer *bytes.Buffer
	var fileName string
	var err error

	switch format {
	case "pdf":
		fileBuffer, err = utils.GeneratePDF(title, "Filtros [ "+formattedFilters+" ]", data, headers, userName, option)
		fileName = "reporte_" + time.Now().Format("20060102_150405") + ".pdf"
	case "excel":
		fileBuffer, err = utils.GenerateExcel(title, headers, data, "Filtros [ "+formattedFilters+" ]", userName, option)
		fileName = "reporte_" + time.Now().Format("20060102_150405") + ".xlsx"
	default:
		return nil, "", fmt.Errorf("formato no soportado")
	}

	if err != nil {
		return nil, "", fmt.Errorf("error al generar el archivo: %w", err)
	}

	return fileBuffer, fileName, nil
}

func formatFilters(filters map[string]interface{}) string {
	if len(filters) == 0 {
		return "Ninguno"
	}

	var result string
	for key, value := range filters {
		result += fmt.Sprintf("%s: %v | ", key, value)
	}
	return result[:len(result)-3]
}

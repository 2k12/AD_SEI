package services

import (
	"bytes"
	"fmt"
	"reflect"
	"seguridad-api/config"
	"seguridad-api/models"
	"seguridad-api/utils"
	"time"
)

func GenerateReport(modelName string, filters map[string]interface{}, userName string, format string) (*bytes.Buffer, string, error) {
	var headers []string
	var data [][]string
	var state string

	var query interface{}

	switch modelName {
	case "User":
		headers = []string{"Nombre de Usuario", "Correo Electrónico", "Estado", "Clave del Módulo", "F. Creación", "F. Actualización"}
		query = &[]models.User{}
	case "Role":
		headers = []string{"Nombre del Rol", "Descripción", "Estado", "F. Creación", "F. Actualización"} // Agregada columna "Descripción"
		query = &[]models.Role{}
	case "Module":
		headers = []string{"Nombre del Módulo", "Descripción", "Estado", "F. Creación", "F. Actualización"} // Agregada columna "Descripción"
		query = &[]models.Module{}
	default:
		return nil, "", fmt.Errorf("modelo no soportado")
	}

	dbQuery := config.DB.Model(query)

	// Aplicar filtros específicos según el modelo
	switch modelName {
	case "User":
		for key, value := range filters {
			switch key {
			case "active":
				dbQuery = dbQuery.Where("active = ?", value)
			case "module_key":
				dbQuery = dbQuery.Where("module_key = ?", value)
			default:
				dbQuery = dbQuery.Where(fmt.Sprintf("%s = ?", key), value)
			}
		}
	case "Role", "Module":
		for key, value := range filters {
			switch key {
			case "active":
				dbQuery = dbQuery.Where("active = ?", value)
			case "date_range":
				dateRange, ok := value.(map[string]interface{})
				if ok {
					if start, exists := dateRange["start"]; exists {
						dbQuery = dbQuery.Where("DATE(created_at) >= ?", start)
					}
					if end, exists := dateRange["end"]; exists {
						dbQuery = dbQuery.Where("DATE(created_at) <= ?", end)
					}
				}
			default:
				dbQuery = dbQuery.Where(fmt.Sprintf("%s = ?", key), value)
			}
		}
	}

	// Ejecutar la consulta
	if err := dbQuery.Find(query).Error; err != nil {
		return nil, "", fmt.Errorf("error al consultar los datos: %w", err)
	}

	// Procesar los datos según el modelo
	value := reflect.ValueOf(query).Elem()
	for i := 0; i < value.Len(); i++ {
		row := []string{}

		if modelName == "User" {
			user := value.Index(i).Interface().(models.User)
			state = "Activo"
			if !user.Active {
				state = "Inactivo"
			}
			row = append(row,
				user.Name,
				user.Email,
				state,
				user.ModuleKey,
				user.CreatedAt.Format("2006-01-02 15:04:05"),
				user.UpdatedAt.Format("2006-01-02 15:04:05"),
			)
		} else if modelName == "Role" {
			role := value.Index(i).Interface().(models.Role)
			state = "Activo"
			if !role.Active {
				state = "Inactivo"
			}
			row = append(row,
				role.Name,
				role.Description, // Nueva columna para Descripción
				state,
				role.CreatedAt.Format("2006-01-02 15:04:05"),
				role.UpdatedAt.Format("2006-01-02 15:04:05"),
			)
		} else if modelName == "Module" {
			module := value.Index(i).Interface().(models.Module)
			state = "Activo"
			if !module.Active {
				state = "Inactivo"
			}
			row = append(row,
				module.Name,
				module.Description, // Nueva columna para Descripción
				state,
				module.CreatedAt.Format("2006-01-02 15:04:05"),
				module.UpdatedAt.Format("2006-01-02 15:04:05"),
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
		fileBuffer, err = utils.GeneratePDF(title, "Filtros [ "+formattedFilters+" ]", data, headers, userName)
		fileName = fmt.Sprintf("reporte_%s_%s.pdf", modelName, time.Now().Format("20060102_150405"))
	case "excel":
		fileBuffer, err = utils.GenerateExcel(title, headers, data, "Filtros [ "+formattedFilters+" ]", userName)
		fileName = fmt.Sprintf("reporte_%s_%s.xlsx", modelName, time.Now().Format("20060102_150405"))
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
		return "Ninguno" // Si no hay filtros, devuelve "Ninguno" o algún mensaje apropiado
	}

	var result string
	for key, value := range filters {
		result += fmt.Sprintf("%s: %v | ", key, value)
	}
	return result[:len(result)-3] // Eliminar el último separador " | "
}

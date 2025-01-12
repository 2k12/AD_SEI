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

func GenerateReport(modelName string, filters map[string]interface{}, userName string) (*bytes.Buffer, string, error) {
	var headers []string
	var data [][]string
	var state string

	var query interface{}
	switch modelName {
	case "User":
		headers = []string{"Nombre", "Correo Electrónico", "Estado", "Clave del Módulo", "F. Creación", "F. Actualización"}
		query = &[]models.User{}
	default:
		return nil, "", fmt.Errorf("modelo no soportado")
	}

	dbQuery := config.DB.Model(query)
	for key, value := range filters {
		dbQuery = dbQuery.Where(fmt.Sprintf("%s = ?", key), value)
	}

	if err := dbQuery.Find(query).Error; err != nil {
		return nil, "", fmt.Errorf("error al consultar los datos: %w", err)
	}

	value := reflect.ValueOf(query).Elem()
	for i := 0; i < value.Len(); i++ {
		row := []string{}
		user := value.Index(i).Interface().(models.User)

		if user.Active {
			state = "Activo"
		} else {
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

		data = append(data, row)
	}

	title := "Reporte de Usuarios"

	// Formatear filtros para incluirlos en el PDF
	formattedFilters := formatFilters(filters)

	fileBuffer, err := utils.GeneratePDF(title, "Filtros [ "+formattedFilters+" ]", data, headers, userName)
	if err != nil {
		return nil, "", fmt.Errorf("error al generar el PDF: %w", err)
	}

	fileName := "reporte_" + time.Now().Format("20060102_150405") + ".pdf"
	return fileBuffer, fileName, nil
}

//	func formatFilters(filters map[string]interface{}) string {
//		var result string
//		for key, value := range filters {
//			result += fmt.Sprintf("%s: %v | ", key, value)
//		}
//		return result[:len(result)-3] // Eliminar el último separador " | "
//	}
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

// package services

// import (
// 	"bytes"
// 	"fmt"
// 	"seguridad-api/config"
// 	"seguridad-api/models"
// 	"seguridad-api/utils"
// 	"time"
// )

// func GenerateReport(all bool) (*bytes.Buffer, string, error) {
// 	var headers []string
// 	var data [][]string
// 	var state string

// 	if all {
// 		// Reporte de usuarios
// 		headers = []string{"Nombre", "Correo Electrónico", "Estado", "Clave del Módulo", "F. Creación", "F. Actualización"}
// 		var users []models.User
// 		if err := config.DB.Find(&users).Error; err != nil {
// 			return nil, "", fmt.Errorf("error al consultar los usuarios: %w", err)
// 		}

// 		for _, user := range users {

// 			if user.Active {
// 				state = "Activo"
// 			} else {
// 				state = "Inactivo"
// 			}

// 			row := []string{
// 				user.Name,
// 				user.Email,
// 				// fmt.Sprintf("%t", user.Active),
// 				state,
// 				user.ModuleKey,
// 				user.CreatedAt.Format("2006-01-02 15:04:05"),
// 				user.UpdatedAt.Format("2006-01-02 15:04:05"),
// 			}
// 			data = append(data, row)
// 		}
// 	}

// 	title := "Reporte de Usuarios"

// 	fileBuffer, err := utils.GeneratePDF(title, "Generado por (USUARIO)", data, headers)
// 	if err != nil {
// 		return nil, "", fmt.Errorf("error al generar el PDF: %w", err)
// 	}

// 	fileName := "reporte_" + time.Now().Format("20060102_150405") + ".pdf"
// 	return fileBuffer, fileName, nil
// }

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

func GenerateReport(modelName string, filters map[string]interface{}) (*bytes.Buffer, string, error) {
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

	fileBuffer, err := utils.GeneratePDF(title, "Filtros [ "+formattedFilters+" ]", data, headers)
	if err != nil {
		return nil, "", fmt.Errorf("error al generar el PDF: %w", err)
	}

	fileName := "reporte_" + time.Now().Format("20060102_150405") + ".pdf"
	return fileBuffer, fileName, nil
}

func formatFilters(filters map[string]interface{}) string {
	var result string
	for key, value := range filters {
		result += fmt.Sprintf("%s: %v | ", key, value)
	}
	return result[:len(result)-3] // Eliminar el último separador " | "
}

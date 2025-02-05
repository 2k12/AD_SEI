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

func GenerateReport(modelName string, filters map[string]interface{}, userName string, format string, option string) (*bytes.Buffer, string, error) {
	var headers []string
	var data [][]string
	var query interface{}

	switch modelName {

	case "Permission":
		headers = []string{"Nombre", "Descripción", "Estado", "Módulo", "F. Creación", "F. Actualización"}
		query = &[]models.Permission{}
		dbQuery := config.DB.Model(query).Preload("Module")

		// Aplicar filtros según los parámetros
		for key, value := range filters {
			switch key {
			case "active":
				dbQuery = dbQuery.Where("active = ?", value)
			case "module_id":
				dbQuery = dbQuery.Where("module_id = ?", value)
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

		// Realizar la consulta
		if err := dbQuery.Find(query).Error; err != nil {
			return nil, "", fmt.Errorf("error al consultar los datos: %w", err)
		}

		// Construir las filas para exportar
		rows := reflect.ValueOf(query).Elem()
		for i := 0; i < rows.Len(); i++ {
			permission := rows.Index(i).Interface().(models.Permission)

			// Determinar el estado
			state := "Activo"
			if !permission.Active {
				state = "Inactivo"
			}

			// Verificar el nombre del módulo
			moduleName := "Sin módulo"
			if permission.Module.ID > 0 && permission.Module.Name != "" {
				moduleName = permission.Module.Name
			}

			// Construir la fila
			row := []string{
				permission.Name,
				permission.Description,
				state,
				moduleName,
				permission.CreatedAt.Format("2006-01-02 15:04:05"),
				permission.UpdatedAt.Format("2006-01-02 15:04:05"),
			}
			data = append(data, row)
		}

		// Validación extra para detectar claves foráneas inválidas (opcional)
		if len(data) == 0 {
			fmt.Println("Advertencia: Puede que haya claves foráneas inválidas en 'module_id'")
		}

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
			headers = []string{"Nombre", "Correo Electrónico", "Estado", "F. Creación", "F. Actualización"}
			query = &[]models.User{}
			dbQuery := config.DB.Model(query)
			for key, value := range filters {
				dbQuery = dbQuery.Where(fmt.Sprintf("%s = ?", key), value)
			}
			if err := dbQuery.Find(query).Error; err != nil {
				return nil, "", fmt.Errorf("error al consultar los datos: %w", err)
			}
		}

	case "Role":
		headers = []string{"Nombre del Rol", "Descripción", "Estado", "F. Creación", "F. Actualización"}
		query = &[]models.Role{}
		dbQuery := config.DB.Model(query)
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
		if err := dbQuery.Find(query).Error; err != nil {
			return nil, "", fmt.Errorf("error al consultar los datos: %w", err)
		}

	case "Module":
		headers = []string{"Nombre del Módulo", "Descripción", "Estado", "F. Creación", "F. Actualización"}
		query = &[]models.Module{}
		dbQuery := config.DB.Model(query)
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
		if err := dbQuery.Find(query).Error; err != nil {
			return nil, "", fmt.Errorf("error al consultar los datos: %w", err)
		}

	case "Audit":
		headers = []string{"Evento", "Descripción", "Usuario", "Servicio Origen", "Fecha"}
		var audits []models.Audit
		dbQuery := config.DB.Model(&audits)
		for key, value := range filters {
			switch key {
			case "userId":
				dbQuery = dbQuery.Where("user_id = ?", value)
			case "date_range":
				dateRange, ok := value.(map[string]interface{})
				if ok {
					if start, exists := dateRange["start"]; exists {
						dbQuery = dbQuery.Where("DATE(date) >= ?", start)
					}
					if end, exists := dateRange["end"]; exists {
						dbQuery = dbQuery.Where("DATE(date) <= ?", end)
					}
				}
			}
		}
		if err := dbQuery.Find(&audits).Error; err != nil {
			return nil, "", fmt.Errorf("error al consultar los datos de auditoría: %w", err)
		}
		query = audits

	default:
		return nil, "", fmt.Errorf("modelo no soportado")
	}

	// Procesar los datos
	if query == nil {
		return nil, "", fmt.Errorf("error: la consulta no devolvió resultados o el modelo '%s' no es válido", modelName)
	}

	var rows reflect.Value
	if reflect.TypeOf(query).Kind() == reflect.Ptr {
		rows = reflect.ValueOf(query).Elem()
	} else {
		rows = reflect.ValueOf(query)
	}

	for i := 0; i < rows.Len(); i++ {
		row := []string{}
		switch modelName {
		case "User":
			user := rows.Index(i).Interface().(models.User)
			state := "Activo"
			if !user.Active {
				state = "Inactivo"
			}
			if option == "usuariosCompletos" {
				roles, permissions, modules := formatUserDetails(user)
				row = append(row, user.Name, roles, permissions, modules)
			} else {
				row = append(row,
					user.Name,
					user.Email,
					state,
					// user.ModuleKey,
					user.CreatedAt.Format("2006-01-02 15:04:05"),
					user.UpdatedAt.Format("2006-01-02 15:04:05"),
				)
			}
		case "Role":
			role := rows.Index(i).Interface().(models.Role)
			state := "Activo"
			if !role.Active {
				state = "Inactivo"
			}
			row = append(row,
				role.Name,
				role.Description,
				state,
				role.CreatedAt.Format("2006-01-02 15:04:05"),
				role.UpdatedAt.Format("2006-01-02 15:04:05"),
			)
		case "Module":
			module := rows.Index(i).Interface().(models.Module)
			state := "Activo"
			if !module.Active {
				state = "Inactivo"
			}
			row = append(row,
				module.Name,
				module.Description,
				state,
				module.CreatedAt.Format("2006-01-02 15:04:05"),
				module.UpdatedAt.Format("2006-01-02 15:04:05"),
			)
		case "Audit": // Procesar datos para auditoría
			audit := rows.Index(i).Interface().(models.Audit)
			row = append(row,
				audit.Event,
				audit.Description,
				fmt.Sprintf("%d", audit.UserID),
				audit.OriginService,
				audit.Date.Format("2006-01-02 15:04:05"),
			)
		}
		data = append(data, row)
	}

	// Generar el archivo
	title := fmt.Sprintf("Reporte de %s", modelName)
	formattedFilters := formatFilters(filters)

	var fileBuffer *bytes.Buffer
	var fileName string
	var err error

	switch format {
	case "pdf":
		fileBuffer, err = utils.GeneratePDF(title, "Filtros [ "+formattedFilters+" ]", data, headers, userName, option)
		fileName = fmt.Sprintf("reporte_%s_%s.pdf", modelName, time.Now().Format("20060102_150405"))
	case "excel":
		fileBuffer, err = utils.GenerateExcel(title, headers, data, "Filtros [ "+formattedFilters+" ]", userName, option)
		fileName = fmt.Sprintf("reporte_%s_%s.xlsx", modelName, time.Now().Format("20060102_150405"))
	default:
		return nil, "", fmt.Errorf("formato no soportado")
	}

	if err != nil {
		return nil, "", fmt.Errorf("error al generar el archivo: %w", err)
	}

	return fileBuffer, fileName, nil
}

func formatUserDetails(user models.User) (roles, permissions, modules string) {
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
	return
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

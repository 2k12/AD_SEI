package services

import (
	"fmt"
	"seguridad-api/config"
	"seguridad-api/models"
	"strings"
	"time"
)

// Registrar auditoría
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

// Obtener todas las auditorías
func GetAudit() ([]models.AuditResponse, error) {
	var audits []models.AuditResponse
	result := config.DB.Find(&audits)
	return audits, result.Error
}

func GetPaginatedAudit(page, pageSize int, filters map[string]interface{}) ([]models.AuditResponse, int64, error) {
	var audits []models.AuditResponse
	var total int64

	query := config.DB.Model(&models.Audit{})

	// Aplicar filtros
	if event, ok := filters["event"]; ok {
		query = query.Joins("INNER JOIN users ON users.id = audit.user_id").
			Where("audit.event LIKE ?", "%"+event.(string)+"%").
			Select("audit.id, audit.event, audit.description, users.name AS user, audit.origin_service, audit.date")
	} else {
		query = query.Joins("INNER JOIN users ON users.id = audit.user_id").
			Select("audit.id, audit.event, audit.description, users.name AS user, audit.origin_service, audit.date")
	}

	if userName, ok := filters["userName"]; ok {
		query = query.Where("LOWER(users.name) LIKE LOWER(?)", "%"+userName.(string)+"%")
	}
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Aplicar paginación
	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	// Obtener datos
	err = query.Find(&audits).Error
	if err != nil {
		return nil, 0, err
	}

	return audits, total, nil
}

// Obtener estadísticas de auditoría con filtros dinámicos
// func GetAuditStatistics(event, userID, originService, startDate, endDate string) ([]models.AuditStatisticsResponse, error) {
// 	var stats []models.AuditStatisticsResponse

// 	// Crear la consulta base
// 	query := config.DB.Model(&models.Audit{}).Select("event, origin_service, COUNT(*) as total")

// 	// Aplicar filtros dinámicos con validaciones y logs
// 	if event != "" {
// 		fmt.Println("Aplicando filtro para event:", event)
// 		query = query.Where("event = ?", event)
// 	}
// 	if userID != "" {
// 		fmt.Println("Aplicando filtro para user_id:", userID)
// 		query = query.Where("user_id = ?", userID)
// 	}
// 	if originService != "" {
// 		fmt.Println("Aplicando filtro para origin_service:", originService)
// 		query = query.Where("origin_service = ?", originService)
// 	}
// 	if startDate != "" && endDate != "" {
// 		start, err := time.Parse("2006-01-02", startDate)
// 		if err != nil {
// 			fmt.Println("Error al parsear startDate:", err)
// 			return nil, err
// 		}
// 		end, err := time.Parse("2006-01-02", endDate)
// 		if err != nil {
// 			fmt.Println("Error al parsear endDate:", err)
// 			return nil, err
// 		}
// 		// Ajustar el rango para incluir el final del día
// 		end = end.Add(24 * time.Hour).Add(-time.Nanosecond)
// 		fmt.Printf("Aplicando filtro para rango de fechas: %v - %v\n", start, end)
// 		query = query.Where("date BETWEEN ? AND ?", start, end)
// 	}

// 	// Registrar el SQL generado para depuración
// 	fmt.Println("SQL Generado:", query.Statement.SQL.String())

// 	// Agrupar resultados por evento y servicio
// 	err := query.Group("event, origin_service").Scan(&stats).Error
// 	if err != nil {
// 		fmt.Println("Error al ejecutar consulta:", err)
// 		return nil, err
// 	}

// 	// Verificar resultados obtenidos
// 	if len(stats) == 0 {
// 		fmt.Println("No se encontraron estadísticas para los filtros aplicados.")
// 	} else {
// 		fmt.Printf("Estadísticas obtenidas: %+v\n", stats)
// 	}

// 	return stats, nil
// }

// func GetAuditStatistics(event, username, originService, startDate, endDate string) ([]models.AuditStatisticsResponse, error) {
// 	var stats []models.AuditStatisticsResponse

// 	// Crear la consulta base
// 	query := config.DB.Model(&models.Audit{}).Select("event, origin_service, COUNT(*) as total")

// 	// Aplicar filtros dinámicos con validaciones y logs
// 	if event != "" {
// 		fmt.Println("Aplicando filtro para event:", event)
// 		query = query.Where("event = ?", event)
// 	}
// 	if username != "" {
// 		fmt.Println("Aplicando filtro para username:", username)
// 		query = query.Joins("INNER JOIN users ON users.id = audit.user_id").Where("users.name = ?", username)
// 	}
// 	if originService != "" {
// 		fmt.Println("Aplicando filtro para origin_service:", originService)
// 		query = query.Where("origin_service = ?", originService)
// 	}
// 	if startDate != "" && endDate != "" {
// 		start, err := time.Parse("2006-01-02", startDate)
// 		if err != nil {
// 			fmt.Println("Error al parsear startDate:", err)
// 			return nil, err
// 		}
// 		end, err := time.Parse("2006-01-02", endDate)
// 		if err != nil {
// 			fmt.Println("Error al parsear endDate:", err)
// 			return nil, err
// 		}
// 		// Ajustar el rango para incluir el final del día
// 		end = end.Add(24 * time.Hour).Add(-time.Nanosecond)
// 		fmt.Printf("Aplicando filtro para rango de fechas: %v - %v\n", start, end)
// 		query = query.Where("date BETWEEN ? AND ?", start, end)
// 	}

// 	// Registrar el SQL generado para depuración
// 	fmt.Println("SQL Generado:", query.Statement.SQL.String())

// 	// Agrupar resultados por evento y servicio
// 	err := query.Group("event, origin_service").Scan(&stats).Error
// 	if err != nil {
// 		fmt.Println("Error al ejecutar consulta:", err)
// 		return nil, err
// 	}

// 	// Verificar resultados obtenidos
// 	if len(stats) == 0 {
// 		fmt.Println("No se encontraron estadísticas para los filtros aplicados.")
// 	} else {
// 		fmt.Printf("Estadísticas obtenidas: %+v\n", stats)
// 	}

// 	return stats, nil
// }

//
// ! funcional con fallas
// func GetAuditStatistics(event, startDate, endDate string) ([]models.AuditStatisticsResponse, []models.Audit, error) {
// 	var stats []models.AuditStatisticsResponse
// 	var records []models.Audit

// 	// Crear la consulta base
// 	query := config.DB.Model(&models.Audit{}).Select("event, COUNT(*) as total, MAX(date) as last_date")

// 	// Aplicar filtro por evento si se proporciona
// 	if event != "" {
// 		query = query.Where("event = ?", event)
// 	}

// 	// Aplicar filtro por fechas si se proporcionan
// 	if startDate != "" && endDate != "" {
// 		start, err := time.Parse("2006-01-02", startDate)
// 		if err != nil {
// 			fmt.Println("Error al parsear startDate:", err)
// 			return nil, nil, err
// 		}
// 		end, err := time.Parse("2006-01-02", endDate)
// 		if err != nil {
// 			fmt.Println("Error al parsear endDate:", err)
// 			return nil, nil, err
// 		}

// 		// Ajusta las fechas a la zona horaria y considera toda la jornada
// 		loc, _ := time.LoadLocation("America/Guayaquil")
// 		start = start.In(loc)
// 		end = end.In(loc).Add(24 * time.Hour) // Ajustar hasta el final del día

// 		// Aplica el filtro por rango de fechas
// 		query = query.Where("date BETWEEN ? AND ?", start, end)
// 	}

// 	// Registrar el SQL generado para depuración
// 	fmt.Println("SQL Generado:", query.Statement.SQL.String())

// 	// Agrupar resultados por evento
// 	err := query.Group("event").Scan(&stats).Error
// 	if err != nil {
// 		fmt.Println("Error al ejecutar consulta:", err)
// 		return nil, nil, err
// 	}

// 	// Obtener los registros individuales con los mismos filtros
// 	err = config.DB.Model(&models.Audit{}).Where(query).Find(&records).Error
// 	if err != nil {
// 		fmt.Println("Error al obtener registros:", err)
// 		return nil, nil, err
// 	}

// 	// Verificar resultados obtenidos
// 	if len(stats) == 0 {
// 		fmt.Println("No se encontraron estadísticas para los filtros aplicados.")
// 	} else {
// 		fmt.Printf("Estadísticas obtenidas: %+v\n", stats)
// 	}

//		// Devuelvo ambas respuestas juntas
//		return stats, records, nil
//	}
func GetAuditStatistics(event, module, startDate, endDate string) ([]models.AuditStatisticsResponse, []models.Audit, error) {
	var stats []models.AuditStatisticsResponse
	var records []models.Audit

	// Crear la consulta base
	query := config.DB.Model(&models.Audit{}).
		Select("event, UPPER(origin_service) AS origin_service, COUNT(*) as total, MAX(date) as last_date").
		Group("event, origin_service") // 🔥 Agrupar por evento y módulo

	// Aplicar filtro por evento si se proporciona
	if event != "" {
		query = query.Where("event = ?", event)
	}

	// Aplicar filtro por módulo si se proporciona
	if module != "" {
		query = query.Where("UPPER(origin_service) = ?", strings.ToUpper(module)) //  Filtrar por módulo en mayúsculas
	}

	// Aplicar filtro por fechas si se proporcionan
	if startDate != "" && endDate != "" {
		// Parsear las fechas en formato YYYY-MM-DD
		start, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			fmt.Println("Error al parsear startDate:", err)
			return nil, nil, err
		}
		end, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			fmt.Println("Error al parsear endDate:", err)
			return nil, nil, err
		}

		// Ajustar las fechas a la zona horaria y considera toda la jornada
		loc, _ := time.LoadLocation("America/Guayaquil")
		start = start.In(loc)
		end = end.In(loc).Add(24 * time.Hour).Add(-time.Nanosecond) // Ajusta hasta el final del día

		// Aplica el filtro por rango de fechas
		query = query.Where("date BETWEEN ? AND ?", start, end)
	}

	// Registrar el SQL generado para depuración
	fmt.Println("SQL Generado:", query.Statement.SQL.String())

	// Obtener los datos de estadísticas
	err := query.Scan(&stats).Error
	if err != nil {
		fmt.Println("Error al ejecutar consulta:", err)
		return nil, nil, err
	}

	// Obtener los registros individuales con los mismos filtros
	err = config.DB.Model(&models.Audit{}).Where(query).Find(&records).Error
	if err != nil {
		fmt.Println("Error al obtener registros:", err)
		return nil, nil, err
	}

	// Verificar resultados obtenidos
	if len(stats) == 0 {
		fmt.Println("No se encontraron estadísticas para los filtros aplicados.")
	} else {
		fmt.Printf("Estadísticas obtenidas: %+v\n", stats)
	}

	// Formatear las fechas de la respuesta antes de devolverlas al frontend (sin zona horaria)
	loc, _ := time.LoadLocation("America/Guayaquil")
	for i := range stats {
		stats[i].LastDateFormatted = stats[i].LastDate.In(loc).Format("2006-01-02 15:04:05")
	}

	// Devuelvo ambas respuestas juntas
	return stats, records, nil
}

package services

import (
	"seguridad-api/config"
	"seguridad-api/models"
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

// Obtener auditorías paginadas con filtros
func GetPaginatedAudit(page, pageSize int, filters map[string]interface{}) ([]models.AuditResponse, int64, error) {
	var audits []models.AuditResponse
	var total int64

	query := config.DB.Model(&models.Audit{})

	// Aplicar filtros
	if event, ok := filters["event"]; ok {
		query = query.Joins("INNER JOIN users ON users.id = audit.user_id").
			Where("audit.event LIKE ? COLLATE utf8_general_ci", "%"+event.(string)+"%").
			Select("audit.id, audit.event, audit.description, users.name AS user, audit.origin_service, audit.date")
	} else {
		query = query.Joins("INNER JOIN users ON users.id = audit.user_id").
			Select("audit.id, audit.event, audit.description, users.name AS user, audit.origin_service, audit.date")
	}

	// Contar el total de registros
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
func GetAuditStatistics(event, userID, originService, startDate, endDate string) ([]models.AuditStatisticsResponse, error) {
	var stats []models.AuditStatisticsResponse

	query := config.DB.Model(&models.Audit{}).Select("event, origin_service, COUNT(*) as total")

	// Aplicar filtros dinámicos
	if event != "" {
		query = query.Where("event = ?", event)
	}
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if originService != "" {
		query = query.Where("origin_service = ?", originService)
	}
	if startDate != "" && endDate != "" {
		start, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			return nil, err
		}
		end, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			return nil, err
		}
		query = query.Where("date BETWEEN ? AND ?", start, end)
	}

	// Agrupar resultados por evento y servicio
	err := query.Group("event, origin_service").Scan(&stats).Error
	if err != nil {
		return nil, err
	}

	return stats, nil
}

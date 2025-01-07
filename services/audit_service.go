package services

import (
	"seguridad-api/config"
	"seguridad-api/models"
	"time"
)

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

// Obtener todos las auditorías
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

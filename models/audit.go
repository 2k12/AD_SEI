package models

import (
	"time"
)

type Audit struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Event         string    `gorm:"type:varchar(50);not null" json:"event"`
	Description   string    `gorm:"type:text" json:"description"`
	UserID        uint      `gorm:"foreignKey:UserID;references:ID" json:"user_id"`
	OriginService string    `gorm:"type:varchar(255)" json:"origin_service"`
	Date          time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"date"`
	User          User      `gorm:"foreignKey:UserID;references:ID" json:"user"`
}
type AuditResponse struct {
	ID            uint      `json:"id"`
	Event         string    `json:"event"`
	Description   string    `json:"description"`
	User          string    `json:"user"`
	OriginService string    `json:"origin_service"`
	Date          time.Time `json:"date"`
}
type AuditStatisticsResponse struct {
	Event         string `json:"event"`
	OriginService string `json:"origin_service"`
	Total         int    `json:"total"`
}

func (Audit) TableName() string {
	return "audit"
}

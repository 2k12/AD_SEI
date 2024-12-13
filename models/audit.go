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

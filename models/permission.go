package models

import (
	"time"
)

type Permission struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Description string    `gorm:"type:varchar(255)" json:"description"`
	ModuleID    uint      `gorm:"not null" json:"module_id"`
	Active      bool      `gorm:"default:true" json:"active"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP;autoUpdateTime" json:"updated_at"`
	Module      Module    `gorm:"foreignKey:ModuleID;references:ID" json:"module"`
	Roles       []Role    `gorm:"many2many:role_permissions" json:"roles"`
}

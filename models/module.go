package models

import (
	"time"
)

type Module struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string `gorm:"type:varchar(100);unique;not null" json:"name"`
	Description string `gorm:"type:varchar(255)" json:"description"`
	Active      bool   `gorm:"default:true" json:"active"`
	ModuleKey   string `gorm:"type:varchar(15);not null" json:"module_key"`
	// CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	// UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP;autoUpdateTime" json:"updated_at"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	// Permissions []Permission `gorm:"foreignKey:ModuleID" json:"permissions"`
	Permissions []Permission `gorm:"foreignKey:ModuleID;references:ID" json:"permissions"`
}

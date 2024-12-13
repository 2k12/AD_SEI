package models

import (
	"time"
)

//	type Role struct {
//		ID          uint   `gorm:"primaryKey;autoIncrement"`
//		Name        string `gorm:"unique"`
//		Description string
//		Users       []User       `gorm:"many2many:user_roles"`
//		Permissions []Permission `gorm:"many2many:role_permissions;"`
//	}
type Role struct {
	ID          uint         `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string       `gorm:"type:varchar(100);unique;not null" json:"name"`
	Description string       `gorm:"type:varchar(255)" json:"description"`
	Active      bool         `gorm:"default:true" json:"active"`
	CreatedAt   time.Time    `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time    `gorm:"default:CURRENT_TIMESTAMP;autoUpdateTime" json:"updated_at"`
	Permissions []Permission `gorm:"many2many:role_permissions" json:"permissions"` // Relationship with Permission
}

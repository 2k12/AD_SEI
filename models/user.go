package models

import (
	"time"
)

type User struct {
	ID             uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name           string     `gorm:"type:varchar(100);not null" json:"name"`
	Email          string     `gorm:"type:varchar(150);unique;not null" json:"email"`
	Password       string     `gorm:"type:varchar(255);not null" json:"password"`
	Active         bool       `gorm:"default:true" json:"active"`
	ModuleKey      string     `gorm:"type:varchar(15);not null" json:"module_key"`
	CreatedAt      time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"default:CURRENT_TIMESTAMP;autoUpdateTime" json:"updated_at"`
	FailedAttempts int        `gorm:"default:0" json:"failed_attempts"`
	LockedUntil    *time.Time `json:"locked_until"`
	Roles          []Role     `gorm:"many2many:user_roles" json:"roles"`
}

package models

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	Email     string    `gorm:"type:varchar(150);unique;not null" json:"email"`
	Password  string    `gorm:"type:varchar(255);not null" json:"password"`
	Active    bool      `gorm:"default:true" json:"active"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP;autoUpdateTime" json:"updated_at"`
	Roles     []Role    `gorm:"many2many:user_roles" json:"roles"`
}

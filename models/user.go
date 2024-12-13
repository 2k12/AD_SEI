// package models

// type User struct {
// 	ID       uint `gorm:"primaryKey;autoIncrement"`
// 	Name     string
// 	Email    string `gorm:"unique"`
// 	Password string
// 	Active   bool
// 	Roles    []Role `gorm:"many2many:user_roles"`
// }
// import (
// 	"time"
// )

//	type User struct {
//		ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
//		Name      string    `gorm:"type:varchar(100);not null" json:"name"`
//		Email     string    `gorm:"type:varchar(150);unique;not null" json:"email"`
//		Password  string    `gorm:"type:varchar(255);not null" json:"password"`
//		RoleID    uint      `gorm:"not null" json:"role_id"`
//		Active    bool      `gorm:"default:true" json:"active"`
//		CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
//		UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP;autoUpdateTime" json:"updated_at"`
//		Role      Role      `gorm:"foreignKey:RoleID;references:ID" json:"role"` // Define relationship with Role
//	}
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
	Roles     []Role    `gorm:"many2many:user_roles" json:"roles"` // Relación many-to-many con Roles
}

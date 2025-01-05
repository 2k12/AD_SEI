package models

//import "time"

type RolePermission struct {
	ID           uint `gorm:"primaryKey;autoIncrement" json:"id"`
	RoleID       uint `gorm:"not null" json:"role_id"`
	PermissionID uint `gorm:"not null" json:"permission_id"`
	//CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	//UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP;autoUpdateTime" json:"updated_at"`

	Role       Role       `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE;" json:"role"`
	Permission Permission `gorm:"foreignKey:PermissionID;constraint:OnDelete:CASCADE;" json:"permission"`
}

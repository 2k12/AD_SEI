package models

// "time"

type UserRole struct {
	ID     uint `gorm:"primaryKey;autoIncrement" json:"id"` // Identificador único
	UserID uint `gorm:"not null" json:"user_id"`            // ID del usuario
	RoleID uint `gorm:"not null" json:"role_id"`            // ID del rol
	// CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`                // Fecha de creación
	// UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP;autoUpdateTime" json:"updated_at"` // Última actualización
	// CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	// UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	//User User `gorm:"foreignKey:UserID;references:ID" json:"user"` // Relación con la tabla users
	//Role Role `gorm:"foreignKey:RoleID;references:ID" json:"role"` // Relación con la tabla roles
}

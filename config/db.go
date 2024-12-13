package config

import (
	"log"
	"seguridad-api/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := "root:admin@tcp(127.0.0.1:3306)/bdd_ads?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error conectando a la base de datos: %v", err)
	}

	// Migraciones automáticas
	db.AutoMigrate(&models.User{}, &models.Role{}, &models.Permission{}, &models.Audit{})
	DB = db
}

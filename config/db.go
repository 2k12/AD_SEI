package config

import (
	"log"
	"seguridad-api/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := "root:pan2@tcp(127.0.0.1:3306)/bd_seguridad?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error conectando a la base de datos: %v", err)
	}

	db.AutoMigrate(&models.User{}, &models.Role{}, &models.Permission{}, &models.Audit{}, models.RolePermission{})
	DB = db
}

// package config

// import (
// 	"fmt"
// 	"log"
// 	"os"
// 	"seguridad-api/models"

// 	"gorm.io/driver/mysql"
// 	"gorm.io/gorm"
// )

// var DB *gorm.DB

// func ConnectDB() {
// 	dbUser := os.Getenv("DB_USER")
// 	dbPassword := os.Getenv("DB_PASSWORD")
// 	dbHost := os.Getenv("DB_HOST")
// 	dbPort := os.Getenv("DB_PORT")
// 	dbName := os.Getenv("DB_NAME")

// 	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
// 		dbUser, dbPassword, dbHost, dbPort, dbName)

// 	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		log.Fatalf("Error conectando a la base de datos: %v", err)
// 	}

// 	db.AutoMigrate(&models.User{}, &models.Role{}, &models.Permission{}, &models.Module{}, &models.Audit{})
// 	DB = db
// }

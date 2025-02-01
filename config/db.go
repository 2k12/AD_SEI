// // package config

// // import (
// // 	"fmt"
// // 	"log"
// // 	"os"
// // 	"seguridad-api/models"

// // 	"gorm.io/driver/mysql"
// // 	"gorm.io/gorm"
// // )

// // var DB *gorm.DB

// // func ConnectDB() {
// // 	dbUser := os.Getenv("DB_USER")
// // 	dbPassword := os.Getenv("DB_PASSWORD")
// // 	dbHost := os.Getenv("DB_HOST")
// // 	dbPort := os.Getenv("DB_PORT")
// // 	dbName := os.Getenv("DB_NAME")

// // 	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
// // 		dbUser, dbPassword, dbHost, dbPort, dbName)

// // 	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
// // 	if err != nil {
// // 		log.Fatalf("Error conectando a la base de datos: %v", err)
// // 	}

// // 	db.AutoMigrate(&models.User{}, &models.Role{}, &models.Permission{}, &models.Module{}, &models.Audit{})
// // 	DB = db
// // }
// // package config

// // import (
// // 	"log"
// // 	"seguridad-api/models"

// // 	"gorm.io/driver/mysql"
// // 	"gorm.io/gorm"
// // )

// // var DB *gorm.DB

// // func ConnectDB() {
// // 	dsn := "root:admin@tcp(127.0.0.1:3306)/bdd_ads?charset=utf8mb4&parseTime=True&loc=Local"
// // 	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
// // 	if err != nil {
// // 		log.Fatalf("Error conectando a la base de datos: %v", err)
// // 	}

// //		db.AutoMigrate(&models.User{}, &models.Role{}, &models.Permission{}, &models.Audit{})
// //		DB = db
// //	}
package config

import (
	"fmt"
	"log"
	"os"
	"seguridad-api/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

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

//		db.AutoMigrate(&models.User{}, &models.Role{}, &models.Permission{}, &models.Module{}, &models.Audit{}, models.RolePermission{})
//		// db.AutoMigrate(&models.User{}, &models.Role{}, &models.Permission{}, &models.Module{}, &models.Audit{})
//		DB = db
//	}
func ConnectDB() {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Error conectando a la base de datos: %v", err)
	}

	// Migración
	// db.AutoMigrate(&models.Module{})
	db.AutoMigrate(&models.User{}, &models.Role{}, &models.Permission{}, &models.Module{}, &models.Audit{}, models.RolePermission{})

	// Añade la columna si no existe
	if err := db.Migrator().AddColumn(&models.Module{}, "ModuleKey"); err != nil {
		log.Printf("Error al agregar la columna ModuleKey: %v", err)
	} else {
		log.Println("Columna ModuleKey agregada exitosamente")
	}

	DB = db
}

// @title API SEGURIDAD con Swagger
// @version 1.0
// @description Esta es la documentación de LA API DE SEGURIDAD hecha con Go.
// @termsOfService http://swagger.io/terms/

// @contact.name Pastillo Joan
// @contact.url https://www.utn.edu.ec
// @contact.email jfpastillod@utn.edu.ec

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"log"
	"seguridad-api/config"
	"seguridad-api/routes"

	"github.com/gin-gonic/gin"
	files "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	config.ConnectDB()

	router := gin.Default()

	router.Static("/docs", "./docs")

	routes.SetupRoutes(router)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler, ginSwagger.URL("http://localhost:8080/docs/swagger.json")))

	log.Println("Servidor corriendo en el puerto 8080")
	log.Println(`http://localhost:8080/swagger/index.html`)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}

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

// @securityDefinitions.apikey XAPI-PIN
// @in header
// @name X-API-PIN

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
	// Conexión a la base de datos
	config.ConnectDB()

	// Inicialización del router
	router := gin.Default()

	// Ruta para acceder a la documentación Swagger
	router.Static("/docs", "./docs")

	// Configuración de las rutas de la API
	routes.SetupRoutes(router)

	// Ruta para cargar la documentación Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler, ginSwagger.URL("http://localhost:8080/docs/swagger.json")))

	log.Println("Servidor corriendo en el puerto 8080")
	log.Println(`http://localhost:8080/swagger/index.html`)

	// Iniciar el servidor
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}

// @title API SEGURIDAD con Swagger
// @version 1.0
// @description Esta es la documentación de LA API DE SEGURIDAD hecha con Go.
// @termsOfService http://swagger.io/terms/

// @contact.name Pastillo Joan
// @contact.url https://www.utn.edu.ec
// @contact.email jfpastillod@utn.edu.ec

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host https://seri-api-utn-2024.fly.dev/
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

package main

import (
	"log"
	"seguridad-api/config"
	"seguridad-api/routes"

	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors" // Librería CORS
	files "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	config.ConnectDB()

	router := gin.Default()

	corsHandler := cors.New(cors.Options{
		// AllowedOrigins: []string{"http://localhost:5173"},
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	router.Use(func(c *gin.Context) {
		corsHandler.ServeHTTP(c.Writer, c.Request, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.Next()
		}))
	})
	// router.Static("/docs", "./docs")
	router.Static("/docs", "/app/docs")

	routes.SetupRoutes(router)

	// router.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler, ginSwagger.URL("https://seri-api-utn-2024.fly.dev/docs/swagger.json")))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler, ginSwagger.URL("/docs/swagger.json")))

	log.Println("Servidor corriendo en el puerto 8080")
	log.Println(`http://localhost:8080/swagger/index.html`)

	// if err := router.Run(":8080"); err != nil {
	// 	log.Fatalf("Error al iniciar el servidor: %v", err)
	// }
	if err := router.Run(":" + getPort()); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}

}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Valor por defecto para desarrollo local
	}
	return port
}

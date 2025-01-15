// @title API SEGURIDAD
// @version 1.6
// @description Esta es la documentación de LA API DE SEGURIDAD hecha con Go.
// @termsOfService http://swagger.io/terms/

// @contact.name Pastillo D Joan
// @contact.url https://www.utn.edu.ec
// @contact.email jfpastillod@utn.edu.ec

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host seri-api-utn-2024.fly.dev
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"log"
	"net/http"
	"os"
	"seguridad-api/config"
	"seguridad-api/routes"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
	files "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	config.ConnectDB()

	router := gin.Default()

	corsHandler := cors.New(cors.Options{
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

	router.Static("/docs", "/app/docs")

	routes.SetupRoutes(router)

	swaggerURL := "http://localhost:8080/docs/swagger.json"
	if os.Getenv("SWAGGER_HOST") != "" {
		swaggerURL = "https://" + os.Getenv("SWAGGER_HOST") + "/docs/swagger.json"
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler, ginSwagger.URL(swaggerURL)))

	port := getPort()
	log.Println("Servidor corriendo en el puerto " + port)

	if err := router.Run("0.0.0.0:" + port); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}

// @title API SEGURIDAD
// @version 1.8
// @description Esta es la documentación de LA API DE SEGURIDAD hecha con Go.
// @termsOfService http://swagger.io/terms/

// @contact.name Pastillo D Joan
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

	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
	files "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	config.ConnectDB()

	router := gin.Default()

	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		//AllowedOrigins:   []string{"*"},
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

	swaggerURL := "http://localhost:8080/docs/swagger.json"
	if os.Getenv("SWAGGER_HOST") != "" {
		swaggerURL = "https://" + os.Getenv("SWAGGER_HOST") + "/docs/swagger.json"
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler, ginSwagger.URL(swaggerURL)))

	log.Println("Servidor corriendo en el puerto " + getPort())

	if err := router.Run(":" + getPort()); err != nil {
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

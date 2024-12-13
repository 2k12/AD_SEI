package main

import (
	"log"
	"seguridad-api/config"
	"seguridad-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()

	router := gin.Default()

	routes.SetupRoutes(router)

	log.Println("Servidor corriendo en el puerto 8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}

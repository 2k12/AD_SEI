package main

import (
	"log"
	"seguridad-api/config"
	"seguridad-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Conexión a la base de datos
	config.ConnectDB()

	// Inicialización del servidor
	router := gin.Default()

	// Configuración de rutas
	routes.SetupRoutes(router)

	// Arrancar el servidor
	log.Println("Servidor corriendo en el puerto 8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}

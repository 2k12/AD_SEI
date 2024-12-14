package routes

import (
	"os"
	"seguridad-api/controllers"
	"seguridad-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	// Agrupamos las rutas bajo "/api"
	api := router.Group("/api")
	{
		// Rutas de autenticación
		api.POST("/login", controllers.Login)
		api.POST("/logout", middleware.AuthMiddleware(os.Getenv("JWT_SECRET_KEY")), controllers.Logout)

		// Rutas de usuarios (requiere Bearer Token)
		api.POST("/users", controllers.CreateUser)

		// Usamos middleware de autenticación para el resto de rutas
		api.Use(middleware.AuthMiddleware(os.Getenv("JWT_SECRET_KEY")))
		{
			api.GET("/users", controllers.GetUsers)
			api.GET("/users/:id/permissions", controllers.GetUserPermissions)
			api.PUT("/users/:id", controllers.UpdateUser)
			api.DELETE("/users/:id", controllers.DeleteUser)

			// Rutas de roles y permisos (requiere Bearer Token)
			api.GET("/roles", controllers.GetRoles)
			api.POST("/roles", controllers.CreateRole)

			api.GET("/permissions", controllers.GetPermissions)
			api.POST("/permissions", controllers.CreatePermission)

			// Ruta para auditar
			api.POST("/audit", controllers.RegisterAudit)
		}
	}
}

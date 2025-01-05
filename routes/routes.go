package routes

import (
	"os"
	"seguridad-api/controllers"
	"seguridad-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.POST("/login", controllers.Login)
		api.POST("/logout", middleware.AuthMiddleware(os.Getenv("JWT_SECRET_KEY")), controllers.Logout)

		api.Use(middleware.AuthMiddleware(os.Getenv("JWT_SECRET_KEY")))
		{
			api.GET("/users/:id/permissions", controllers.GetUserPermissions)

			api.POST("/users", controllers.CreateUser)
			api.GET("/users", controllers.GetUsers)
			api.PUT("/users/:id", controllers.UpdateUser)
			api.DELETE("/users/:id", controllers.DeleteUser)

			api.GET("/roles", controllers.GetRoles)
			api.POST("/roles", controllers.CreateRole)

			api.GET("/permissions", controllers.GetPermissions)
			api.POST("/permissions", controllers.CreatePermission)

			api.POST("/audit", controllers.RegisterAudit)

			api.POST("/users/:id/roles", controllers.AssignRoleToUser)              // Asignar rol
			api.DELETE("/users/:id/roles/:role_id", controllers.RemoveRoleFromUser) // Eliminar rol
			api.GET("/users/:id/roles", controllers.GetUserRoles)

			api.POST("/modules", controllers.CreateModule)    // Crear módulo
			api.GET("/modules", controllers.GetModules)       // Obtener todos los módulos
			api.GET("/modules/:id", controllers.GetModule)    // Obtener un módulo por ID
			api.PUT("/modules/:id", controllers.UpdateModule) // Actualizar un módulo

			api.DELETE("/modules/:id", controllers.DeleteModule) // Eliminar un módulo definitivo

			api.PATCH("/modules/:id/toggle-active", controllers.ToggleModuleActive) // Esta ruta cambia estado activo/inactivo

		}
	}
}

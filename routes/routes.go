package routes

import (
	"os"
	"seguridad-api/controllers"
	controllerReport "seguridad-api/controllers/reports"
	"seguridad-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.POST("/login", controllers.Login)
		api.POST("/generate-report", controllerReport.GenerateReport)
		api.POST("/logout", middleware.AuthMiddleware(os.Getenv("JWT_SECRET_KEY")), controllers.Logout)

		api.Use(middleware.AuthMiddleware(os.Getenv("JWT_SECRET_KEY")))
		{
			api.GET("/users/:id/permissions", controllers.GetUserPermissions)

			api.POST("/users", controllers.CreateUser)
			api.GET("/users", controllers.GetUsers)
			api.GET("/users-dropdown", controllers.GetUsersforDropdown)
			api.PUT("/users/:id", controllers.UpdateUser)
			api.DELETE("/users/:id", controllers.DeleteUser)

			api.GET("/roles", controllers.GetRoles)
			api.POST("/roles", controllers.CreateRole)
			api.PUT("/roles/:id", controllers.UpdateRole)
			api.PUT("/roles/:id/state", controllers.UpdateRoleState)

			api.GET("/permissions", controllers.GetPermissions)
			api.POST("/permissions", controllers.CreatePermission)
			api.POST("/permissions/fastCharge", controllers.ChargeFastOfData)
			api.PUT("/permissions/:id", controllers.UpdatePermission)
			api.DELETE("/permissions/:id", controllers.DeletePermission)
			api.GET("/permissions/:id", controllers.GetPermissionByID)

			//api.GET("/modules", controllers.GetModules)

			api.POST("/audit", controllers.RegisterAudit)
			api.GET("/audit", controllers.GetAudit)

			api.POST("/roles/:role_id/permissions", controllers.AssignPermission)
			api.DELETE("/roles/:role_id/permissions", controllers.RemovePermission)
			api.GET("/roles/:role_id/permissions", controllers.GetRolePermissions)
			api.GET("/permissions/all", controllers.GetAllPermissions)

			api.POST("/users/:id/roles", controllers.AssignRoleToUser)
			api.DELETE("/users/:id/roles/:role_id", controllers.RemoveRoleFromUser)
			api.GET("/users/:id/roles", controllers.GetUserRoles)

			api.POST("/modules", controllers.CreateModule)
			api.GET("/modules", controllers.GetModules)
			api.GET("/modules/:id", controllers.GetModule)
			api.PUT("/modules/:id", controllers.UpdateModule)
			api.DELETE("/modules/:id", controllers.DeleteModule)

			// api.PATCH("/modules/:id/toggle-active", controllers.ToggleModuleActive) // Esta ruta cambia estado activo/inactivo

		}
	}
}

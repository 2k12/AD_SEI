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
		api.POST("/request-reset", controllers.RequestPasswordReset)
		api.POST("/reset-password", controllers.ResetPassword)
		api.POST("/login", controllers.Login)
		api.POST("/logout", middleware.AuthMiddleware(os.Getenv("JWT_SECRET_KEY")), controllers.Logout)

		api.Use(middleware.AuthMiddleware(os.Getenv("JWT_SECRET_KEY")))
		{
			api.GET("/users/:id/permissions", controllers.GetUserPermissions)

			api.POST("/users", controllers.CreateUser)
			api.GET("/users", controllers.GetUsers)
			api.GET("/users-dropdown", controllers.GetUsersforDropdown)
			api.PUT("/users/:id", controllers.UpdateUser)
			api.POST("/users/fastCharge", controllers.ChargeFastUsers)
			api.DELETE("/users/:id", controllers.DeleteUser)

			api.GET("/roles", controllers.GetRoles)
			api.GET("/roles-active", controllers.GetRolesActive)
			api.GET("/roles-dropdown", controllers.GetRolesforDropdown)
			api.POST("/roles", controllers.CreateRole)
			api.PUT("/roles/:id", controllers.UpdateRole)
			api.PUT("/roles/:id/state", controllers.UpdateRoleState)

			const PermissionRoute = "/permissions/:id"
			api.GET("/permissions", controllers.GetPermissions)
			api.POST("/permissions", controllers.CreatePermission)
			api.POST("/permissions/fastCharge", controllers.ChargeFastOfData)
			api.PUT("/permissions/:id", controllers.UpdatePermission)
			api.DELETE("/permissions/:id", controllers.DeletePermission)
			api.GET("/permissions/:id", controllers.GetPermissionByID)
			api.GET("/permissions/active", controllers.GetPermissionsForModal)

			//api.GET("/modules", controllers.GetModules)

			api.POST("/audit", controllers.RegisterAudit)
			api.GET("/audit", controllers.GetAudit)
			api.GET("/audit/statistics", controllers.GetAuditoriaEstadisticas)

			const RolePermissionsRoute = "/roles/:role_id/permissions"
			api.POST(RolePermissionsRoute, controllers.AssignPermission)
			api.DELETE(RolePermissionsRoute, controllers.RemovePermission)
			api.GET(RolePermissionsRoute, controllers.GetRolePermissions)
			api.GET("/permissions/all", controllers.GetAllPermissions)
			api.GET("/modules/:id/permissions", controllers.GetPermissionsByModule)

			api.POST("/users/:id/roles", controllers.AssignRoleToUser)
			api.DELETE("/users/:id/roles/:role_id", controllers.RemoveRoleFromUser)
			api.GET("/users/:id/roles", controllers.GetUserRoles)

			const moduleRoute = "/modules/:id"
			api.POST("/modules", controllers.CreateModule)
			api.GET("/modules", controllers.GetModules)
			api.GET(moduleRoute, controllers.GetModule)
			api.GET("/modules/active", controllers.GetModuleActive)

			api.PUT(moduleRoute, controllers.UpdateModule)
			api.DELETE(moduleRoute, controllers.DeleteModule)

			// api.PATCH("/modules/:id/toggle-active", controllers.ToggleModuleActive) // Esta ruta cambia estado activo/inactivo

			api.POST("/generate-report", controllerReport.GenerateReport)

		}
	}
}

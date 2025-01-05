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
			api.PUT("/permissions/:id", controllers.UpdatePermission)
			api.DELETE("/permissions/:id", controllers.DeletePermission)
			api.GET("/permissions/:id", controllers.GetPermissionByID)

			api.GET("/modules", controllers.GetModules)

			api.POST("/audit", controllers.RegisterAudit)

		}
	}
}

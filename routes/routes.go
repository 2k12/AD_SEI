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

		api.POST("/users", controllers.CreateUser)

		api.Use(middleware.AuthMiddleware(os.Getenv("JWT_SECRET_KEY")))
		{
			api.GET("/users/:id/permissions", controllers.GetUserPermissions)

			authenticated := api.Use(middleware.ValidatePIN())
			{
				authenticated.GET("/users", controllers.GetUsers)
				authenticated.PUT("/users/:id", controllers.UpdateUser)
				authenticated.DELETE("/users/:id", controllers.DeleteUser)

				authenticated.GET("/roles", controllers.GetRoles)
				authenticated.POST("/roles", controllers.CreateRole)

				authenticated.GET("/permissions", controllers.GetPermissions)
				authenticated.POST("/permissions", controllers.CreatePermission)

				authenticated.POST("/audit", controllers.RegisterAudit)
			}

		}
	}
}

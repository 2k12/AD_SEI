// package routes

// import (
// 	"seguridad-api/controllers"
// 	"seguridad-api/middleware"

// 	"github.com/gin-gonic/gin"
// )

// func SetupRoutes(router *gin.Engine) {
// 	api := router.Group("/api")

// 	// Rutas de autenticación
// 	api.POST("/login", controllers.Login)
// 	api.POST("/logout", middleware.AuthMiddleware("your_secret_key"), controllers.Logout)

// 	// Rutas de creación de usuario (sin autenticación)
// 	api.POST("/users", controllers.CreateUser)

// 	// Rutas de gestión de usuarios, roles, permisos y auditoría protegidas por middleware
// 	api.Use(middleware.AuthMiddleware("your_secret_key"))
// 	{
// 		api.GET("/users", controllers.GetUsers)
// 		api.GET("/roles", controllers.GetRoles)
// 		api.POST("/roles", controllers.CreateRole)
// 		api.GET("/permissions", controllers.GetPermissions)
// 		api.POST("/permissions", controllers.CreatePermission)

// 		// Ruta para registrar auditoría
// 		api.POST("/audit", controllers.RegisterAudit)
// 	}
// }

// func SetupRoutes(router *gin.Engine) {
// 	api := router.Group("/api")

// 	// Rutas de autenticación (sin middleware)
// 	api.POST("/login", controllers.Login)
// 	api.POST("/logout", middleware.AuthMiddleware("your_secret_key"), controllers.Logout)

// 	// Ruta de creación de usuario (sin middleware de autenticación)
// 	// api.POST("/users", controllers.CreateUser)

// 	// Rutas de gestión de usuarios, roles, permisos y auditoría protegidas por middleware
// 	api.Use(middleware.AuthMiddleware("your_secret_key"))
// 	{
// 		api.GET("/users", controllers.GetUsers)
// 		api.POST("/users", controllers.CreateUser)
// 		api.PUT("/users/:id", controllers.UpdateUser)    // Actualizar usuario
// 		api.DELETE("/users/:id", controllers.DeleteUser) // Cambiar estado de usuario

// 		api.GET("/roles", controllers.GetRoles)
// 		api.POST("/roles", controllers.CreateRole)
// 		// api.POST("/roles", controllers.CreateRole)
// 		// api.PUT("/roles/:id", controllers.UpdateRole)  // Actualizar rol
// 		// api.DELETE("/roles/:id", controllers.DeleteRole) // Cambiar estado de rol

// 		api.GET("/permissions", controllers.GetPermissions)
// 		api.POST("/permissions", controllers.CreatePermission)
// 		// api.PUT("/permissions/:id", controllers.UpdatePermission)    // Actualizar permiso
// 		// api.DELETE("/permissions/:id", controllers.DeletePermission) // Cambiar estado de permiso

//			// Ruta para registrar auditoría
//			api.POST("/audit", controllers.RegisterAudit)
//		}
//	}
package routes

import (
	"os"
	"seguridad-api/controllers"
	"seguridad-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")

	// Rutas de autenticación (sin middleware)
	api.POST("/login", controllers.Login)
	api.POST("/logout", middleware.AuthMiddleware(os.Getenv("JWT_SECRET_KEY")), controllers.Logout)
	api.POST("/users", controllers.CreateUser)

	// Rutas protegidas por middleware
	api.Use(middleware.AuthMiddleware(os.Getenv("JWT_SECRET_KEY")))
	{
		api.GET("/users", controllers.GetUsers)
		api.PUT("/users/:id", controllers.UpdateUser)
		api.DELETE("/users/:id", controllers.DeleteUser)

		api.GET("/roles", controllers.GetRoles)
		api.POST("/roles", controllers.CreateRole)

		api.GET("/permissions", controllers.GetPermissions)
		api.POST("/permissions", controllers.CreatePermission)

		// Ruta para registrar auditoría
		api.POST("/audit", controllers.RegisterAudit)
	}
}

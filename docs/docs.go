// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "Pastillo D Joan",
            "url": "https://www.utn.edu.ec",
            "email": "jfpastillod@utn.edu.ec"
        },
        "license": {
            "name": "MIT",
            "url": "https://opensource.org/licenses/MIT"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/audit": {
            "post": {
                "description": "Registra un evento de auditoría en el sistema",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auditoría"
                ],
                "summary": "Registrar auditoría",
                "parameters": [
                    {
                        "description": "Datos de auditoría a registrar",
                        "name": "auditData",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.RegisterAuditInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Auditoría registrada exitosamente",
                        "schema": {
                            "$ref": "#/definitions/controllers.RegisterAuditResponse"
                        }
                    },
                    "400": {
                        "description": "Datos inválidos o formato incorrecto",
                        "schema": {
                            "$ref": "#/definitions/controllers.ErrorResponseAudit"
                        }
                    },
                    "500": {
                        "description": "Error al registrar la auditoría",
                        "schema": {
                            "$ref": "#/definitions/controllers.ErrorResponseAudit"
                        }
                    }
                }
            }
        },
        "/login": {
            "post": {
                "description": "Autentica un usuario con email y contraseña, devolviendo un token JWT",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Autenticación"
                ],
                "summary": "Iniciar sesión",
                "parameters": [
                    {
                        "description": "Datos de inicio de sesión (email y password)",
                        "name": "loginData",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.LoginData"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "token",
                        "schema": {
                            "$ref": "#/definitions/controllers.TokenResponse"
                        }
                    },
                    "400": {
                        "description": "Datos inválidos",
                        "schema": {
                            "$ref": "#/definitions/controllers.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Credenciales inválidas",
                        "schema": {
                            "$ref": "#/definitions/controllers.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/logout": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Invalida la sesión actual del usuario. Requiere un Bearer Token.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Autenticación"
                ],
                "summary": "Cerrar sesión",
                "responses": {
                    "200": {
                        "description": "message",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "401": {
                        "description": "error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/permissions/all": {
            "get": {
                "description": "Lista todos los permisos disponibles, incluyendo los módulos a los que pertenecen",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Rol_Permisos"
                ],
                "summary": "Obtener todos los permisos",
                "responses": {
                    "200": {
                        "description": "Lista de permisos",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/controllers.PermissionResponse"
                            }
                        }
                    },
                    "500": {
                        "description": "error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/roles": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Devuelve una lista paginada de roles, permitiendo filtrar por nombre y estado activo.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Roles"
                ],
                "summary": "Obtener roles",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Número de página (por defecto 1)",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Tamaño de página (por defecto 10)",
                        "name": "pageSize",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filtrar por nombre",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "boolean",
                        "description": "Filtrar por estado activo",
                        "name": "active",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "roles",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "401": {
                        "description": "error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Crea un nuevo rol con nombre, descripción y estado activo. Requiere un Bearer Token.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Roles"
                ],
                "summary": "Crear rol",
                "parameters": [
                    {
                        "description": "Datos del rol",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "role",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "401": {
                        "description": "error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/roles/{id}": {
            "put": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Actualiza los datos de un rol existente identificándolo por su ID. Requiere un Bearer Token.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Roles"
                ],
                "summary": "Actualizar rol",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID del rol",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Datos actualizados del rol",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "role",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "401": {
                        "description": "error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/roles/{id}/state": {
            "patch": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Cambia únicamente el estado activo de un rol identificado por su ID. Requiere un Bearer Token.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Roles"
                ],
                "summary": "Actualizar estado del rol",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID del rol",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Estado del rol",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "message",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "401": {
                        "description": "error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/roles/{role_id}/permissions": {
            "get": {
                "description": "Lista todos los permisos asignados a un rol específico",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Rol_Permisos"
                ],
                "summary": "Obtener permisos de un rol",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID del rol",
                        "name": "role_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Lista de permisos",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/controllers.PermissionResponse"
                            }
                        }
                    },
                    "400": {
                        "description": "error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            },
            "post": {
                "description": "Asocia un permiso existente a un rol específico",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Rol_Permisos"
                ],
                "summary": "Asignar permiso a rol",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID del rol",
                        "name": "role_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Datos del permiso a asignar",
                        "name": "permissionData",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.PermissionDataRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "role_permission",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            },
            "delete": {
                "description": "Elimina la asociación de un permiso específico con un rol",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Rol_Permisos"
                ],
                "summary": "Eliminar permiso de rol",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID del rol",
                        "name": "role_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Datos del permiso a eliminar",
                        "name": "permissionData",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.PermissionDataRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "message",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/users": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Crea un nuevo usuario con nombre, email, contraseña y estado activo. Requiere un Bearer Token.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Usuarios"
                ],
                "summary": "Crear usuario",
                "parameters": [
                    {
                        "description": "Datos del usuario",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "user",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "401": {
                        "description": "error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/users/{id}": {
            "put": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Actualiza los datos de un usuario existente. Requiere un Bearer Token.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Usuarios"
                ],
                "summary": "Actualizar usuario",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID del usuario",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Datos a actualizar",
                        "name": "userData",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "updatedUser",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "401": {
                        "description": "error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Cambia el estado de un usuario a inactivo. Requiere un Bearer Token.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Usuarios"
                ],
                "summary": "Eliminar usuario",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID del usuario",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "message",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "401": {
                        "description": "error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/users/{id}/permissions": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Devuelve la lista de permisos asignados a un usuario específico, dado su ID. Requiere un Bearer Token.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Usuarios"
                ],
                "summary": "Obtener permisos de un usuario",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID del usuario",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "permissions",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "object",
                                "additionalProperties": true
                            }
                        }
                    },
                    "400": {
                        "description": "error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "401": {
                        "description": "error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "controllers.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "Credenciales inválidas"
                }
            }
        },
        "controllers.ErrorResponseAudit": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "Error al realizar el registro"
                }
            }
        },
        "controllers.LoginData": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string",
                    "example": "user@example.com"
                },
                "password": {
                    "type": "string",
                    "example": "securePassword123"
                }
            }
        },
        "controllers.PermissionDataRequest": {
            "type": "object",
            "properties": {
                "permission_id": {
                    "type": "integer",
                    "example": 1
                }
            }
        },
        "controllers.PermissionResponse": {
            "type": "object",
            "properties": {
                "active": {
                    "type": "boolean",
                    "example": true
                },
                "description": {
                    "type": "string",
                    "example": "Permiso de Crud Usuarios"
                },
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "module_id": {
                    "type": "integer",
                    "example": 1
                },
                "name": {
                    "type": "string",
                    "example": "Gestion Usuarios"
                }
            }
        },
        "controllers.RegisterAuditInput": {
            "type": "object",
            "required": [
                "date",
                "description",
                "event",
                "origin_service",
                "user_id"
            ],
            "properties": {
                "date": {
                    "type": "string",
                    "example": "2024-12-14T15:04:05Z"
                },
                "description": {
                    "type": "string",
                    "example": "Se creó un nuevo usuario con el email user@example.com."
                },
                "event": {
                    "type": "string",
                    "example": "INSERT"
                },
                "origin_service": {
                    "type": "string",
                    "example": "INVENTARIO"
                },
                "user_id": {
                    "type": "string",
                    "example": "123"
                }
            }
        },
        "controllers.RegisterAuditResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "Auditoría registrada exitosamente"
                }
            }
        },
        "controllers.TokenResponse": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string",
                    "example": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "API SEGURIDAD con Swagger",
	Description:      "Esta es la documentación de LA API DE SEGURIDAD hecha con Go.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

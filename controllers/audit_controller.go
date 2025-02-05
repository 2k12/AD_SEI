package controllers

import (
	"fmt"
	"net/http"
	"seguridad-api/config"
	"seguridad-api/models"
	"seguridad-api/services"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type RegisterAuditInput struct {
	Event         string `json:"event" binding:"required" example:"INSERT"`
	Description   string `json:"description" binding:"required" example:"Se creó un nuevo usuario con el email user@example.com."`
	UserID        string `json:"user_id" binding:"required" example:"123"`
	OriginService string `json:"origin_service" binding:"required" example:"INVENTARIO"`
	Date          string `json:"date" binding:"required" example:"2024-12-14T15:04:05Z"`
}

type RegisterAuditResponse struct {
	Message string `json:"message" example:"Auditoría registrada exitosamente"`
}

type ErrorResponseAudit struct {
	Error string `json:"error" example:"Error al realizar el registro"`
}

// RegisterAudit registra un evento de auditoría
// @Summary Registrar auditoría
// @Description Registra un evento de auditoría en el sistema
// @Tags Auditoría
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param auditData body RegisterAuditInput true "Datos de auditoría a registrar"
// @Success 200 {object} RegisterAuditResponse "Auditoría registrada exitosamente"
// @Failure 400 {object} ErrorResponseAudit "Datos inválidos o formato incorrecto"
// @Failure 500 {object} ErrorResponseAudit "Error al registrar la auditoría"
// @Router /audit [post]
func RegisterAudit(c *gin.Context) {
	var input RegisterAuditInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse del ID del usuario
	userID, err := strconv.ParseUint(input.UserID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El ID del usuario debe ser un número válido"})
		return
	}

	// Parse de la fecha recibida y conversión a UTC
	date, err := time.Parse(time.RFC3339, input.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El formato de la fecha es inválido"})
		return
	}

	// Convertir la fecha a UTC
	date = date.UTC()

	// Registrar la auditoría con la fecha en UTC
	if err := services.RegisterAudit(input.Event, input.Description, uint(userID), input.OriginService, date); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al registrar la auditoría"})
		return
	}

	c.JSON(http.StatusOK, RegisterAuditResponse{Message: "Auditoría registrada exitosamente"})
}

func GetAuditStatistics(event, module, startDate, endDate string) ([]models.AuditStatisticsResponse, error) {
	var stats []models.AuditStatisticsResponse

	// Crear la consulta base
	query := config.DB.Model(&models.Audit{}).
		Select("event, UPPER(origin_service) AS origin_service, COUNT(*) as total, MAX(date) as last_date").
		Group("event, origin_service")

	// Aplicar filtros dinámicos con validaciones y logs
	if event != "" {
		fmt.Println("Aplicando filtro para event:", event)
		query = query.Where("event = ?", event)
	}
	if module != "" {
		fmt.Println("Aplicando filtro para origin_service (módulo):", module)
		query = query.Where("UPPER(origin_service) = ?", strings.ToUpper(module))
	}
	if startDate != "" && endDate != "" {
		start, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			fmt.Println("Error al parsear startDate:", err)
			return nil, err
		}
		end, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			fmt.Println("Error al parsear endDate:", err)
			return nil, err
		}
		end = end.Add(24 * time.Hour).Add(-time.Nanosecond)
		fmt.Printf("Aplicando filtro para rango de fechas: %v - %v\n", start, end)
		query = query.Where("date BETWEEN ? AND ?", start, end)
	}

	// Registrar el SQL generado para depuración
	fmt.Println("SQL Generado:", query.Statement.SQL.String())

	// Obtener los resultados
	err := query.Scan(&stats).Error
	if err != nil {
		fmt.Println("Error al ejecutar consulta:", err)
		return nil, err
	}

	// Formatear las fechas antes de devolver la respuesta
	loc, _ := time.LoadLocation("America/Guayaquil")
	for i := range stats {
		stats[i].LastDateFormatted = stats[i].LastDate.In(loc).Format("2006-01-02 15:04:05")
	}

	return stats, nil
}

// GetAudit obtiene auditorías con paginación y filtros
// @Summary Obtener auditorías
// @Description Devuelve una lista de auditorías registradas, con soporte para paginación y filtros (por evento).
// @Tags Auditoría
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param page query int false "Número de página para la paginación (por defecto: 1)"
// @Param pageSize query int false "Número de registros por página (por defecto: 10)"
// @Param event query string false "Filtrar auditorías por tipo de evento"
// @Success 200 {object} map[string]interface{} "audits"
// @Failure 500 {object} ErrorResponseAudit "Error al obtener las auditorías"
// @Router /audit [get]
func GetAudit(c *gin.Context) {
	// Obtener parámetros de consulta para paginación y filtros
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	event := c.Query("event")

	// Construir filtros
	filters := make(map[string]interface{})
	if event != "" {
		filters["event"] = event
	}

	// Llamar al servicio para obtener auditorias paginadas
	audits, total, err := services.GetPaginatedAudit(page, pageSize, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener las auditorias"})
		return
	}

	// Respuesta con datos paginados
	c.JSON(http.StatusOK, gin.H{
		"audits":     audits,
		"total":      total,
		"page":       page,
		"pageSize":   pageSize,
		"totalPages": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// func GetAuditoriaEstadisticas(c *gin.Context) {
// 	// Obtener parámetros de consulta
// 	event := c.Query("event")
// 	// userID := c.Query("user_id")
// 	// originService := c.Query("origin_service")
// 	startDate := c.Query("start_date")
// 	endDate := c.Query("end_date")

// 	// Llamar al servicio
// 	stats, err := services.GetAuditStatistics(event, startDate, endDate)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Devolver resultados
// 	c.JSON(http.StatusOK, stats)
// }

func GetAuditoriaEstadisticas(c *gin.Context) {
	// Obtener parámetros de consulta
	event := c.Query("event")
	module := c.Query("module") // Obtener el módulo de los parámetros
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	// Llamar al servicio
	stats, records, err := services.GetAuditStatistics(event, module, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Construir la respuesta JSON con estadísticas y registros
	response := gin.H{
		"statistics": stats,
		"records":    records,
	}

	// Devolver resultados
	c.JSON(http.StatusOK, response)
}

// func RegisterAudit(c *gin.Context) {
// 	var input RegisterAuditInput

// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	userID, err := strconv.ParseUint(input.UserID, 10, 32)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "El ID del usuario debe ser un número válido"})
// 		return
// 	}

// 	date, err := time.Parse(time.RFC3339, input.Date)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "El formato de la fecha es inválido"})
// 		return
// 	}

// 	if err := services.RegisterAudit(input.Event, input.Description, uint(userID), input.OriginService, date); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al registrar la auditoría"})
// 		return
// 	}

//		c.JSON(http.StatusOK, RegisterAuditResponse{Message: "Auditoría registrada exitosamente"})
//	}

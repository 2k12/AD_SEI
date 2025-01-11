// package controllers

// import (
// 	"net/http"
// 	services "seguridad-api/services/reports"

// 	"github.com/gin-gonic/gin"
// )

// func GenerateReport(c *gin.Context) {
// 	var requestData struct {
// 		All bool `json:"all"`
// 	}

// 	if err := c.ShouldBindJSON(&requestData); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
// 		return
// 	}

// 	fileBuffer, fileName, err := services.GenerateReport(requestData.All)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.DataFromReader(http.StatusOK, int64(fileBuffer.Len()), "application/pdf", fileBuffer, map[string]string{
// 		"Content-Disposition": `attachment; filename="` + fileName + `"`,
// 	})

// }

package controllers

import (
	"net/http"
	services "seguridad-api/services/reports"

	"github.com/gin-gonic/gin"
)

func GenerateReport(c *gin.Context) {
	var requestData struct {
		Filters map[string]interface{} `json:"filters"` // Recibir filtros como un mapa
		Model   string                 `json:"model"`   // Nombre del modelo a consultar
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	fileBuffer, fileName, err := services.GenerateReport(requestData.Model, requestData.Filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.DataFromReader(http.StatusOK, int64(fileBuffer.Len()), "application/pdf", fileBuffer, map[string]string{
		"Content-Disposition": `attachment; filename="` + fileName + `"`})
}

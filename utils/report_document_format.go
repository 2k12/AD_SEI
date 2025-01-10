package utils

import (
	"bytes"
	"fmt"
	"log"

	"github.com/signintech/gopdf"
)

// GeneratePDF genera un documento PDF con título, encabezados y datos en formato de tabla.
func GeneratePDF(title, usernameAndFilters string, data [][]string, headers []string) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	pdf := gopdf.GoPdf{}

	// Configurar el tamaño de la página (A4)
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()

	// Agregar la fuente
	err := pdf.AddTTFFont("arial", "assets/fonts/Roboto-Regular.ttf")
	if err != nil {
		log.Println("Error al agregar fuente:", err)
		return nil, err
	}
	err = pdf.SetFont("arial", "", 8) // Fuente más pequeña
	if err != nil {
		log.Println("Error al establecer la fuente:", err)
		return nil, err
	}

	// Configuración inicial
	marginX := 30.0
	marginY := 30.0
	cellWidth := 80.0  // Ajustar ancho de celda
	cellHeight := 12.0 // Ajustar alto de celda
	maxWidth := gopdf.PageSizeA4.W - 2*marginX
	startX := marginX
	startY := marginY + 50.0 // Espacio para título y encabezado

	// Función para agregar encabezado en cada página
	addHeader := func(pageNum int) {
		pdf.SetX(marginX)
		pdf.SetY(marginY)

		// Título
		pdf.SetFont("arial", "B", 10)
		pdf.Cell(nil, title)
		pdf.Br(12)

		// Usuario y filtros
		pdf.SetFont("arial", "", 8)
		pdf.Cell(nil, usernameAndFilters)
		pdf.Br(15)

		// Numeración de página
		pdf.SetX(maxWidth - 50)
		pdf.Cell(nil, fmt.Sprintf("Página %d", pageNum))
		pdf.Br(10)
	}

	// Agregar encabezado en la primera página
	addHeader(1)

	// Dibujar encabezados de la tabla
	pdf.SetFont("arial", "B", 8)
	currentY := startY
	for i, header := range headers {
		x := startX + float64(i)*cellWidth
		pdf.RectFromUpperLeftWithStyle(x, currentY, cellWidth, cellHeight, "D")
		pdf.SetXY(x+2, currentY+2)
		pdf.Cell(nil, header)
	}
	pdf.Br(cellHeight)

	// Dibujar filas de datos
	pdf.SetFont("arial", "", 8)
	pageNum := 1
	for _, row := range data {
		currentY = pdf.GetY()
		for i, col := range row {
			x := startX + float64(i)*cellWidth
			pdf.RectFromUpperLeftWithStyle(x, currentY, cellWidth, cellHeight, "D")
			pdf.SetXY(x+2, currentY+2)
			pdf.Cell(nil, col)
		}
		pdf.Br(cellHeight)

		// Verificar si se necesita una nueva página
		if currentY+cellHeight > gopdf.PageSizeA4.H-marginY {
			pageNum++
			pdf.AddPage()
			addHeader(pageNum)
			currentY = startY
			pdf.SetY(currentY)
		}
	}

	// Guardar el documento PDF en el buffer
	_, err = pdf.WriteTo(&buf)
	if err != nil {
		log.Println("Error al escribir el archivo PDF:", err)
		return nil, err
	}

	return &buf, nil
}

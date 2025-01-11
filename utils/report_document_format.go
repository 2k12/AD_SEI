package utils

import (
	"bytes"
	"fmt"
	"log"
	helpers "seguridad-api/helpers"
	"time"

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
	err = pdf.SetFont("arial", "", 8)
	if err != nil {
		log.Println("Error al establecer la fuente:", err)
		return nil, err
	}

	// Configuración inicial
	marginX := 30.0
	marginY := 30.0
	cellHeight := 18.0
	lineSpacing := 8.0 // Espaciado entre líneas dentro de una celda
	pageWidth := gopdf.PageSizeA4.W
	pageHeight := gopdf.PageSizeA4.H
	usableWidth := pageWidth - 2*marginX
	columnWidth := usableWidth / float64(len(headers))
	startX := marginX
	startY := marginY + 50.0 // Espacio para título y encabezado
	currentTime := time.Now()

	// Ajustar la hora al huso horario de Ecuador usando el helper
	ecuadorTime := helpers.AdjustToEcuadorTime(currentTime)

	// Función para agregar encabezado en cada página
	addHeader := func(pageNum int) {
		pdf.Br(20)

		// Cargar imagen
		imgPath := "assets/img/security.png"
		err := pdf.Image(imgPath, pageWidth-marginX-50, marginY, &gopdf.Rect{W: 20, H: 20})
		if err != nil {
			log.Println("Error al cargar la imagen:", err)
		}

		// Título
		pdf.SetFont("arial", "B", 20)
		pdf.SetX(marginX) // Establecer la posición horizontal inicial
		pdf.Cell(nil, title)

		// Fecha
		pdf.Br(12)
		pdf.SetFont("arial", "", 8)
		pdf.SetX(marginX)                                                                     // Asegurar margen izquierdo consistente
		pdf.Cell(nil, fmt.Sprintf("Fecha [ %s ]", ecuadorTime.Format("02/01/2006 15:04:05"))) // Formato corregido
		pdf.Br(12)

		pdf.SetFont("arial", "", 8)
		pdf.SetX(marginX) // Asegurar margen izquierdo consistente
		pdf.Cell(nil, usernameAndFilters)
		pdf.Br(15)

	}

	// // Función para agregar encabezado en cada página
	// addHeader := func(pageNum int) {
	// 	pdf.Br(20)

	// 	// Título
	// 	pdf.SetFont("arial", "B", 14)
	// 	pdf.SetX(marginX) // Establecer la posición horizontal inicial
	// 	pdf.Cell(nil, title)
	// 	pdf.Br(12)

	// 	// Usuario y filtros
	// 	pdf.SetFont("arial", "", 8)
	// 	pdf.SetX(marginX) // Asegurar margen izquierdo consistente
	// 	pdf.Cell(nil, usernameAndFilters)
	// 	pdf.Br(12)

	// 	// Hora de Ecuador
	// 	pdf.SetFont("arial", "", 8)
	// 	pdf.SetX(marginX)                                                                            // Asegurar margen izquierdo consistente
	// 	pdf.Cell(nil, fmt.Sprintf("Hora de Ecuador: %s", ecuadorTime.Format("02/01/2006 15:04:05"))) // Formato corregido
	// 	pdf.Br(15)
	// }

	// Función para agregar pie de página
	addFooter := func(pageNum int) {
		pdf.SetFont("arial", "", 8)
		pdf.SetX(marginX)
		pdf.SetY(pageHeight - marginY + 10)
		pdf.Cell(nil, fmt.Sprintf("Página %d", pageNum))
	}

	// Agregar encabezado en la primera página
	addHeader(1)

	// Dibujar encabezados de la tabla sin fondo (solo texto negro)
	pdf.SetFont("arial", "B", 8)
	currentY := startY
	pdf.SetTextColor(0, 0, 0) // Texto negro
	for i, header := range headers {
		x := startX + float64(i)*columnWidth
		pdf.RectFromUpperLeftWithStyle(x, currentY, columnWidth, cellHeight, "D") // Solo dibuja el borde
		pdf.SetXY(x+2, currentY+2)
		pdf.Cell(nil, header)
	}
	pdf.Br(cellHeight)

	// Dibujar filas de datos sin espacio entre filas
	pdf.SetFont("arial", "", 8)
	pageNum := 1
	for _, row := range data {
		currentY = pdf.GetY()

		// Dibujar cada celda de la fila
		for i, col := range row {
			x := startX + float64(i)*columnWidth

			// Dividir el texto en líneas si es necesario (máximo dos líneas)
			lines := wrapText(col, columnWidth-4, &pdf)
			if len(lines) > 2 {
				lines = lines[:2] // Limitar a dos líneas
			}

			// Dibujar celda
			pdf.RectFromUpperLeftWithStyle(x, currentY, columnWidth, cellHeight, "D") // "D" solo dibuja el borde
			for j, line := range lines {
				pdf.SetXY(x+2, currentY+2+float64(j)*lineSpacing)
				pdf.Cell(nil, line)
			}
		}

		currentY += cellHeight
		pdf.SetY(currentY)

		// Verificar si se necesita una nueva página
		if currentY+cellHeight > pageHeight-marginY {
			addFooter(pageNum)
			pageNum++
			pdf.AddPage()
			addHeader(pageNum)
			currentY = startY
			pdf.SetY(currentY)
		}
	}

	// Agregar pie de página en la última página
	addFooter(pageNum)

	// Guardar el documento PDF en el buffer
	_, err = pdf.WriteTo(&buf)
	if err != nil {
		log.Println("Error al escribir el archivo PDF:", err)
		return nil, err
	}

	return &buf, nil
}

// wrapText divide el texto en líneas según el ancho máximo permitido.
func wrapText(text string, maxWidth float64, pdf *gopdf.GoPdf) []string {
	words := []rune(text)
	var lines []string
	var currentLine string
	var currentWidth float64

	for _, word := range words {
		wordWidth, _ := pdf.MeasureTextWidth(string(word))
		if currentWidth+wordWidth > maxWidth {
			lines = append(lines, currentLine)
			currentLine = ""
			currentWidth = 0
		}
		currentLine += string(word)
		currentWidth += wordWidth
	}
	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}

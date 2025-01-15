package utils

import (
	"bytes"
	"fmt"
	"log"
	helpers "seguridad-api/helpers"
	"time"

	"github.com/signintech/gopdf"
	"github.com/xuri/excelize/v2"
)

func GenerateExcel(title string, headers []string, data [][]string, usernameAndFilters string, userName string, option string) (*bytes.Buffer, error) {
	f := excelize.NewFile()

	sheetName := "Reporte"
	if option == "usuariosCompletos" {
		sheetName = "Usuarios Completos"
	}
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}
	f.SetActiveSheet(index)

	titleCell := "A1"
	f.SetCellValue(sheetName, titleCell, title)
	style := &excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 16,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
		},
	}
	styleID, err := f.NewStyle(style)
	if err != nil {
		return nil, fmt.Errorf("error creating style: %w", err)
	}
	f.MergeCell(sheetName, titleCell, fmt.Sprintf("%s1", columnNameFromIndex(len(headers)-1)))
	f.SetCellStyle(sheetName, titleCell, fmt.Sprintf("%s1", columnNameFromIndex(len(headers)-1)), styleID)

	currentTime := time.Now()
	ecuadorTime := helpers.AdjustToEcuadorTime(currentTime)
	f.SetCellValue(sheetName, "A2", fmt.Sprintf("Fecha [ %s ]", ecuadorTime.Format("02/01/2006 15:04:05")))
	f.SetCellValue(sheetName, "A3", fmt.Sprintf("Generado por: [ %s ]", userName))
	f.SetCellValue(sheetName, "A4", usernameAndFilters)

	for i := 0; i < 4; i++ {
		f.SetColWidth(sheetName, columnNameFromIndex(i), columnNameFromIndex(i), 20)
	}

	if option == "usuariosCompletos" {
		headers = []string{"Nombre", "Roles", "Permisos", "Módulos"}
	}

	for i, header := range headers {
		cell := fmt.Sprintf("%s5", columnNameFromIndex(i))
		f.SetCellValue(sheetName, cell, header)
	}

	for rowIndex, row := range data {
		for colIndex, value := range row {
			cell := fmt.Sprintf("%s%d", columnNameFromIndex(colIndex), rowIndex+6)
			f.SetCellValue(sheetName, cell, value)
			if option == "usuariosCompletos" {
				if colIndex == 1 || colIndex == 2 || colIndex == 3 {
					style := &excelize.Style{
						Alignment: &excelize.Alignment{
							WrapText: true,
						},
					}
					styleID, err := f.NewStyle(style)
					if err != nil {
						return nil, fmt.Errorf("error creando estilo de ajuste de texto: %w", err)
					}
					f.SetCellStyle(sheetName, cell, cell, styleID)
				}
			}
		}
	}

	for i := range headers {
		column := columnNameFromIndex(i)
		f.SetColWidth(sheetName, column, column, 40)
	}

	buf := new(bytes.Buffer)
	if err := f.Write(buf); err != nil {
		return nil, err
	}

	return buf, nil
}

func columnNameFromIndex(index int) string {
	name := ""
	for index >= 0 {
		name = string('A'+(index%26)) + name
		index = index/26 - 1
	}
	return name
}

func GeneratePDF(title, usernameAndFilters string, data [][]string, headers []string, userName string, option string) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	pdf := gopdf.GoPdf{}

	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()

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

	marginX := 30.0
	marginY := 30.0
	cellHeight := 18.0
	lineSpacing := 8.0
	pageWidth := gopdf.PageSizeA4.W
	pageHeight := gopdf.PageSizeA4.H
	usableWidth := pageWidth - 2*marginX
	columnWidth := usableWidth / float64(len(headers))
	startX := marginX
	startY := marginY + 50.0
	currentTime := time.Now()

	ecuadorTime := helpers.AdjustToEcuadorTime(currentTime)

	addHeader := func(pageNum int) {
		pdf.Br(20)

		imgPath := "assets/img/security.png"
		err := pdf.Image(imgPath, pageWidth-marginX-50, marginY, &gopdf.Rect{W: 20, H: 20})
		if err != nil {
			log.Println("Error al cargar la imagen:", err)
		}

		pdf.SetFont("arial", "B", 20)
		pdf.SetX(marginX)
		pdf.Cell(nil, title)

		pdf.Br(12)
		pdf.SetFont("arial", "", 8)
		pdf.SetX(marginX)
		pdf.Cell(nil, fmt.Sprintf("Fecha [ %s ]", ecuadorTime.Format("02/01/2006 15:04:05")))
		pdf.Br(12)

		pdf.SetFont("arial", "", 8)
		pdf.SetX(marginX)
		pdf.Cell(nil, fmt.Sprintf("Generado por: [ %s ]", userName))
		pdf.Br(12)

		pdf.SetFont("arial", "", 8)
		pdf.SetX(marginX)
		pdf.Cell(nil, usernameAndFilters)
		pdf.Br(15)
	}

	addFooter := func(pageNum int) {
		pdf.SetFont("arial", "", 8)
		pdf.SetX(marginX)
		pdf.SetY(pageHeight - marginY + 10)
		pdf.Cell(nil, fmt.Sprintf("Página %d", pageNum))
	}

	addHeader(1)

	pdf.SetFont("arial", "B", 8)
	currentY := startY
	pdf.SetTextColor(0, 0, 0)
	for i, header := range headers {
		x := startX + float64(i)*columnWidth
		pdf.RectFromUpperLeftWithStyle(x, currentY, columnWidth, cellHeight, "D")
		pdf.SetXY(x+2, currentY+2)
		pdf.Cell(nil, header)
	}
	pdf.Br(cellHeight)

	pdf.SetFont("arial", "", 8)
	pageNum := 1
	for _, row := range data {
		currentY = pdf.GetY()
		maxLinesInRow := 0

		defaultMaxLines := 3
		if option == "usuariosCompletos" {
			defaultMaxLines = 15
		}

		cellLines := make([]int, len(row))
		for i, col := range row {
			lines := wrapText(col, columnWidth-4, &pdf)
			if len(lines) > defaultMaxLines {
				lines = lines[:defaultMaxLines]
			}
			cellLines[i] = len(lines)

			if len(lines) > maxLinesInRow {
				maxLinesInRow = len(lines)
			}
		}

		if maxLinesInRow < defaultMaxLines {
			maxLinesInRow = defaultMaxLines
		}

		for i, col := range row {
			x := startX + float64(i)*columnWidth
			lines := wrapText(col, columnWidth-4, &pdf)

			if len(lines) > maxLinesInRow {
				lines = lines[:maxLinesInRow]
			}

			pdf.RectFromUpperLeftWithStyle(x, currentY, columnWidth, float64(maxLinesInRow)*lineSpacing+2, "D")
			for j, line := range lines {
				pdf.SetXY(x+2, currentY+2+float64(j)*lineSpacing)
				pdf.Cell(nil, line)
			}
		}

		currentY += float64(maxLinesInRow)*lineSpacing + 2
		pdf.SetY(currentY)

		if currentY+cellHeight > pageHeight-marginY {
			addFooter(pageNum)
			pageNum++
			pdf.AddPage()
			addHeader(pageNum)
			currentY = startY
			pdf.SetY(currentY)
		}
	}

	addFooter(pageNum)

	_, err = pdf.WriteTo(&buf)
	if err != nil {
		log.Println("Error al escribir el archivo PDF:", err)
		return nil, err
	}

	return &buf, nil
}

func wrapText(text string, maxWidth float64, pdf *gopdf.GoPdf) []string {
	words := text
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

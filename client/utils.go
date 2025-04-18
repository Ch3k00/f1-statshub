package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

// getInput obtiene entrada del usuario con un mensaje

// formatLapTime convierte segundos a formato min:seg.miliseg

// printTable crea una tabla con los encabezados dados
func printTable(headers []string, rows [][]string) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	// Convertir encabezados a interface{}
	headerRow := make([]interface{}, len(headers))
	for i, h := range headers {
		headerRow[i] = h
	}
	t.AppendHeader(headerRow)

	// Añadir filas
	for _, row := range rows {
		rowInterface := make([]interface{}, len(row))
		for i, v := range row {
			rowInterface[i] = v
		}
		t.AppendRow(rowInterface)
	}

	t.Render()
}

// parseDriverID convierte el input a número de piloto
func parseDriverID(input string) (int, error) {
	return strconv.Atoi(strings.TrimSpace(input))
}

// clearScreen limpia la pantalla de la terminal
func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

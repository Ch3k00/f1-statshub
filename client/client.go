package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

const APIBaseURL = "http://localhost:8080/api"

func main() {
	for {
		printMenu()
		option := getInput("Seleccione una opción: ")

		switch option {
		case "1":
			listDrivers()
		case "2":
			showDriverDetails()
		case "3":
			listRaces()
		case "4":
			showRaceDetails()
		case "5":
			showSeasonSummary()
		case "6":
			fmt.Println("Fin del programa!")
			return
		default:
			fmt.Println("Opción no válida. Intente nuevamente.")
		}
	}
}

func printMenu() {
	fmt.Println("\nMenu")
	fmt.Println("1. Ver corredores")
	fmt.Println("2. Ver detalle de corredor")
	fmt.Println("3. Ver carreras")
	fmt.Println("4. Ver detalle de carrera")
	fmt.Println("5. Resumen de temporada")
	fmt.Println("6. Salir")
}

func getInput(prompt string) string {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func listDrivers() {
	resp, err := http.Get(APIBaseURL + "/corredor")
	if err != nil {
		fmt.Println("Error al conectar con la API:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error leyendo respuesta:", err)
		return
	}

	var drivers []struct {
		DriverNumber int    `json:"driver_number"`
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
		TeamName     string `json:"team_name"`
		CountryCode  string `json:"country_code"`
	}

	if err := json.Unmarshal(body, &drivers); err != nil {
		fmt.Println("Error decodificando respuesta:", err)
		return
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Nombre", "Apellido", "N Piloto", "Equipo", "País"})

	for i, driver := range drivers {
		t.AppendRow(table.Row{
			i + 1,
			driver.FirstName,
			driver.LastName,
			driver.DriverNumber,
			driver.TeamName,
			driver.CountryCode,
		})
	}

	t.Render()
}

func showDriverDetails() {
	driverID := getInput("Ingrese el número del piloto: ")

	resp, err := http.Get(APIBaseURL + "/corredor/detalle/" + driverID)
	if err != nil {
		fmt.Println("Error al conectar con la API:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("Error:", string(body))
		return
	}

	var result struct {
		DriverID    int    `json:"driver_id"`
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		TeamName    string `json:"team_name"`
		Performance struct {
			Wins         int     `json:"wins"`
			Top3Finishes int     `json:"top_3_finishes"`
			MaxSpeed     float64 `json:"max_speed"`
		} `json:"performance_summary"`
		RaceResults []struct {
			SessionKey       int     `json:"session_key"`
			CircuitShortName string  `json:"circuit_short_name"`
			Position         int     `json:"position"`
			FastestLap       bool    `json:"fastest_lap"`
			MaxSpeed         float64 `json:"max_speed"`
			BestLapDuration  float64 `json:"best_lap_duration"`
		} `json:"race_results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("Error decodificando respuesta:", err)
		return
	}

	// Tabla de resultados de carreras
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Carrera", "Pos Final", "Vuelta rápida", "Velocidad max", "Menor tiempo vuelta"})

	for i, race := range result.RaceResults {
		fastestLap := "No"
		if race.FastestLap {
			fastestLap = "Sí"
		}

		t.AppendRow(table.Row{
			i + 1,
			"GP de " + race.CircuitShortName,
			race.Position,
			fastestLap,
			fmt.Sprintf("%.1f km/h", race.MaxSpeed),
			formatLapTime(race.BestLapDuration),
		})
	}

	t.Render()

	// Resumen del desempeño
	fmt.Println("\nResumen del desempeño del piloto")
	t = table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendRow(table.Row{"Carreras ganadas", result.Performance.Wins})
	t.AppendRow(table.Row{"Veces en el top 3", result.Performance.Top3Finishes})
	t.AppendRow(table.Row{"Velocidad máxima alcanzada", fmt.Sprintf("%.1f km/h", result.Performance.MaxSpeed)})
	t.Render()
}

func formatLapTime(seconds float64) string {
	minutes := int(seconds) / 60
	remainingSeconds := seconds - float64(minutes*60)
	return fmt.Sprintf("%d:%.3f", minutes, remainingSeconds)
}

func showSeasonSummary() {
	resp, err := http.Get(APIBaseURL + "/temporada/resumen")
	if err != nil {
		fmt.Println("Error al conectar con la API:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error leyendo respuesta:", err)
		return
	}

	var summary struct {
		Season            int            `json:"season"`
		Top3Winners       []SummaryEntry `json:"top_3_winners"`
		Top3FastestLaps   []SummaryEntry `json:"top_3_fastest_laps"`
		Top3PolePositions []SummaryEntry `json:"top_3_pole_positions"`
	}

	typeAlias := struct {
		Season            int            `json:"season"`
		Top3Winners       []SummaryEntry `json:"top_3_winners"`
		Top3FastestLaps   []SummaryEntry `json:"top_3_fastest_laps"`
		Top3PolePositions []SummaryEntry `json:"top_3_pole_positions"`
	}{}

	if err := json.Unmarshal(body, &typeAlias); err != nil {
		fmt.Println("Error decodificando respuesta:", err)
		return
	}

	summary = typeAlias

	fmt.Printf("\nTop 3 Pilotos con más Victorias - Temporada %d\n", summary.Season)
	printSummaryTable(summary.Top3Winners, "Victorias")

	fmt.Printf("\nTop 3 Pilotos con más Vueltas Rápidas - Temporada %d\n", summary.Season)
	printSummaryTable(summary.Top3FastestLaps, "Vueltas Rápidas")

	fmt.Printf("\nTop 3 Pilotos con más Pole Positions - Temporada %d\n", summary.Season)
	printSummaryTable(summary.Top3PolePositions, "Pole Positions")
}

type SummaryEntry struct {
	Position    int    `json:"position"`
	Driver      string `json:"driver"`
	Team        string `json:"team"`
	Country     string `json:"country"`
	Wins        int    `json:"wins,omitempty"`
	Poles       int    `json:"poles,omitempty"`
	FastestLaps int    `json:"fastest_laps,omitempty"`
}

func printSummaryTable(data []SummaryEntry, label string) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Posición", "Piloto", "Equipo", "País", label})

	for _, d := range data {
		val := d.Wins
		if label == "Pole Positions" {
			val = d.Poles
		} else if label == "Vueltas Rápidas" {
			val = d.FastestLaps
		}

		t.AppendRow(table.Row{d.Position, d.Driver, d.Team, d.Country, val})
	}

	t.Render()
}

func showRaceDetails() {
	raceID := getInput("Ingrese el número de la carrera (session_key): ")

	resp, err := http.Get(APIBaseURL + "/carrera/detalle/" + raceID)
	if err != nil {
		fmt.Println("Error al conectar con la API:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("Error:", string(body))
		return
	}

	var detail struct {
		RaceID           int    `json:"race_id"`
		CountryName      string `json:"country_name"`
		DateStart        string `json:"date_start"`
		Year             int    `json:"year"`
		CircuitShortName string `json:"circuit_short_name"`
		Results          []struct {
			Position string `json:"position"` // puede ser "1", "2", "Ultimo", etc.
			Driver   string `json:"driver"`
			Team     string `json:"team"`
			Country  string `json:"country"`
		} `json:"results"`
		FastestLap struct {
			Driver    string  `json:"driver"`
			TotalTime float64 `json:"total_time"`
			Sector1   float64 `json:"sector_1"`
			Sector2   float64 `json:"sector_2"`
			Sector3   float64 `json:"sector_3"`
		} `json:"fastest_lap"`
		MaxSpeed struct {
			Driver   string  `json:"driver"`
			SpeedKmh float64 `json:"speed_kmh"`
		} `json:"max_speed"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&detail); err != nil {
		fmt.Println("Error decodificando respuesta:", err)
		return
	}

	// Podio y último lugar
	fmt.Println("\nResultados")
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Posición", "Piloto", "Equipo", "País"})
	for _, r := range detail.Results {
		t.AppendRow(table.Row{r.Position, r.Driver, r.Team, r.Country})
	}
	t.Render()

	// Vuelta más rápida
	fmt.Println("\nVuelta más rápida")
	t = table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Piloto", "Tiempo total", "Sector 1", "Sector 2", "Sector 3"})
	t.AppendRow(table.Row{
		detail.FastestLap.Driver,
		formatLapTime(detail.FastestLap.TotalTime),
		fmt.Sprintf("%.3f", detail.FastestLap.Sector1),
		fmt.Sprintf("%.3f", detail.FastestLap.Sector2),
		fmt.Sprintf("%.3f", detail.FastestLap.Sector3),
	})
	t.Render()

	// Velocidad máxima
	fmt.Println("\nVelocidad máxima alcanzada")
	t = table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Piloto", "Velocidad (km/h)"})
	t.AppendRow(table.Row{detail.MaxSpeed.Driver, fmt.Sprintf("%.1f", detail.MaxSpeed.SpeedKmh)})
	t.Render()
}

func listRaces() {
	fmt.Println("📝 Aquí se listarán las carreras...")
}

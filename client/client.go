package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

const APIBaseURL = "http://localhost:8080/api"

func main() {
	for {
		fmt.Println(`
Menu
1. Ver corredores
2. Ver detalle de corredor
3. Ver carreras
4. Ver detalle de carrera
5. Resumen de temporada
6. Salir`)
		opt := getInput("Seleccione una opcion : ")

		switch opt {
		case "1":
			viewDrivers()
		case "2":
			viewDriverDetail()
		case "3":
			viewRaces()
		case "4":
			viewRaceDetail()
		case "5":
			viewSeasonSummary()
		case "6":
			fmt.Println("Fin del programa !")
			return
		default:
			fmt.Println("Opcion no valida.")
		}
	}
}

func getInput(msg string) string {
	fmt.Print(msg)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func viewDrivers() {
	resp, err := http.Get(APIBaseURL + "/corredor")
	if err != nil {
		fmt.Println("Error al obtener corredores:", err)
		return
	}
	defer resp.Body.Close()

	var drivers []struct {
		DriverNumber int    `json:"driver_number"`
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
		TeamName     string `json:"team_name"`
		CountryCode  string `json:"country_code"`
	}

	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &drivers); err != nil {
		fmt.Println("Error decodificando respuesta:", err)
		return
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Nombre", "Apellido", "N Piloto", "Equipo", "Pais"})
	for i, d := range drivers {
		t.AppendRow(table.Row{i + 1, d.FirstName, d.LastName, d.DriverNumber, d.TeamName, d.CountryCode})
	}
	t.Render()
}

func viewDriverDetail() {
	id := getInput("Ingrese el numero del piloto : ")
	resp, err := http.Get(APIBaseURL + "/corredor/detalle/" + id)
	if err != nil {
		fmt.Println("Error consultando detalle del piloto:", err)
		return
	}
	defer resp.Body.Close()

	var data struct {
		DriverID    int `json:"driver_id"`
		RaceResults []struct {
			Race            string   `json:"race"`
			Position        int      `json:"position"`
			FastestLap      bool     `json:"fastest_lap"`
			MaxSpeed        int      `json:"max_speed"`
			BestLapDuration *float64 `json:"best_lap_duration"`
		} `json:"race_results"`
		PerformanceSummary struct {
			Wins     int `json:"wins"`
			Top3     int `json:"top3"`
			MaxSpeed int `json:"max_speed"`
		} `json:"performance_summary"`
	}

	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println("Error decodificando respuesta:", err)
		return
	}

	podio := 0
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Carrera", "Pos Final", "Vuelta rapida", "Velocidad max", "Menor tiempo vuelta"})
	for i, r := range data.RaceResults {
		fast := "No"
		if r.FastestLap {
			fast = "Si"
		}
		time := "N/A"
		if r.BestLapDuration != nil {
			time = formatLapTime(*r.BestLapDuration)
		}
		if r.Position <= 3 {
			podio++
		}
		t.AppendRow(table.Row{i + 1, r.Race, r.Position, fast, fmt.Sprintf("%d km /h", r.MaxSpeed), time})
	}
	t.Render()

	fmt.Println("\n-----------------------------------------------")
	fmt.Println("| Resumen del desempeno del piloto           |")
	fmt.Println("-----------------------------------------------")
	fmt.Printf("| Carreras ganadas           | %d           |\n", data.PerformanceSummary.Wins)
	fmt.Printf("| Veces en el top 3          | %d           |\n", podio)
	fmt.Printf("| Velocidad maxima alcanzada | %d km /h     |\n", data.PerformanceSummary.MaxSpeed)
	fmt.Println("-----------------------------------------------")
}

func viewRaces() {
	resp, err := http.Get(APIBaseURL + "/carrera")
	if err != nil {
		fmt.Println("Error al obtener carreras:", err)
		return
	}
	defer resp.Body.Close()

	var races []struct {
		SessionKey       int    `json:"session_key"`
		CountryName      string `json:"country_name"`
		DateStart        string `json:"date_start"`
		Year             int    `json:"year"`
		CircuitShortName string `json:"circuit_short_name"`
	}

	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &races)

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "ID carrera", "Pais", "Fecha", "Year", "Circuito"})
	for i, r := range races {
		fecha := strings.Split(r.DateStart, "T")[0]
		fechaFmt := strings.ReplaceAll(fecha, "-", " - ")
		t.AppendRow(table.Row{i + 1, r.SessionKey, r.CountryName, fechaFmt, r.Year, r.CircuitShortName})
	}
	t.Render()
}

func viewRaceDetail() {
	sid := getInput("Ingrese el numero de la carrera : ")
	resp, err := http.Get(APIBaseURL + "/carrera/detalle/" + sid)
	if err != nil {
		fmt.Println("Error al obtener detalle de carrera:", err)
		return
	}
	defer resp.Body.Close()

	var detail struct {
		Results []struct {
			Position string `json:"position"`
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

	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &detail)

	fmt.Println("---------------------------------------------------------------")
	fmt.Println("| Resultados                                                 |")
	fmt.Println("---------------------------------------------------------------")
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Posicion", "Piloto", "Equipo", "Pais"})
	for _, r := range detail.Results {
		t.AppendRow(table.Row{r.Position, r.Driver, r.Team, r.Country})
	}
	t.Render()

	fmt.Println("---------------------------------------------------------------")
	fmt.Println("| Vuelta mas rapida                                          |")
	fmt.Println("---------------------------------------------------------------")
	t = table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Piloto", "Tiempo Total", "Sector 1", "Sector 2", "Sector 3"})
	t.AppendRow(table.Row{detail.FastestLap.Driver, formatLapTime(detail.FastestLap.TotalTime),
		fmt.Sprintf("%.3f", detail.FastestLap.Sector1), fmt.Sprintf("%.3f", detail.FastestLap.Sector2), fmt.Sprintf("%.3f", detail.FastestLap.Sector3)})
	t.Render()

	fmt.Println("---------------------------------------------------------------")
	fmt.Println("| Velocidad maxima alcanzada                                 |")
	fmt.Println("---------------------------------------------------------------")
	t = table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Piloto", "Velocidad ( km / h)"})
	t.AppendRow(table.Row{detail.MaxSpeed.Driver, fmt.Sprintf("%.1f", detail.MaxSpeed.SpeedKmh)})
	t.Render()
}

func viewSeasonSummary() {
	resp, err := http.Get(APIBaseURL + "/temporada/resumen")
	if err != nil {
		fmt.Println("Error al obtener resumen de temporada:", err)
		return
	}
	defer resp.Body.Close()

	var summary struct {
		Season            int            `json:"season"`
		Top3Winners       []SummaryEntry `json:"top_3_winners"`
		Top3FastestLaps   []SummaryEntry `json:"top_3_fastest_laps"`
		Top3PolePositions []SummaryEntry `json:"top_3_pole_positions"`
	}

	json.NewDecoder(resp.Body).Decode(&summary)

	fmt.Printf("\n---------------------------------------------------------------------\n")
	fmt.Printf("| Top 3 Pilotos con mas Victorias - Temporada %d |\n", summary.Season)
	fmt.Printf("---------------------------------------------------------------------\n")
	printSummaryTable(summary.Top3Winners, "Victorias")

	fmt.Printf("\n---------------------------------------------------------------------\n")
	fmt.Printf("| Top 3 Pilotos con mas Vueltas Rapidas - Temporada %d |\n", summary.Season)
	fmt.Printf("---------------------------------------------------------------------\n")
	printSummaryTable(summary.Top3FastestLaps, "Vueltas Rapidas")

	fmt.Printf("\n---------------------------------------------------------------------\n")
	fmt.Printf("| Top 3 Pilotos con mas Pole Positions - Temporada %d |\n", summary.Season)
	fmt.Printf("---------------------------------------------------------------------\n")
	printSummaryTable(summary.Top3PolePositions, "Poles")
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
	t.AppendHeader(table.Row{"Posicion", "Piloto", "Equipo", "Pais", label})
	for _, d := range data {
		val := d.Wins
		if label == "Poles" {
			val = d.Poles
		} else if label == "Vueltas Rapidas" {
			val = d.FastestLaps
		}
		t.AppendRow(table.Row{d.Position, d.Driver, d.Team, d.Country, val})
	}
	t.Render()
}

func formatLapTime(seconds float64) string {
	if seconds == 0 {
		return "N/A"
	}
	minutes := int(seconds) / 60
	remainingSeconds := seconds - float64(minutes*60)
	return fmt.Sprintf("%d:%06.3f", minutes, remainingSeconds)
}

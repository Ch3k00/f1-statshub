package models

type RaceDetail struct {
	RaceID           int               `json:"race_id"`
	CountryName      string            `json:"country_name"`
	DateStart        string            `json:"date_start"`
	Year             int               `json:"year"`
	CircuitShortName string            `json:"circuit_short_name"`
	Results          []RaceResultEntry `json:"results"`
	FastestLap       FastestLapDetail  `json:"fastest_lap"`
	MaxSpeed         MaxSpeedDetail    `json:"max_speed"`
}

type RaceResultEntry struct {
	Position string `json:"position"` // Será "1", "2", ..., "Ultimo"
	Driver   string `json:"driver"`
	Team     string `json:"team"`
	Country  string `json:"country"`
}

type FastestLapDetail struct {
	Driver  string  `json:"driver"`
	Total   float64 `json:"total_time"`
	Sector1 float64 `json:"sector_1"`
	Sector2 float64 `json:"sector_2"`
	Sector3 float64 `json:"sector_3"`
}

type MaxSpeedDetail struct {
	Driver   string  `json:"driver"`
	SpeedKMH float64 `json:"speed_kmh"`
}

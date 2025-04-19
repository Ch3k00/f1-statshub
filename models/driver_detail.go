package models

type DriverDetail struct {
	DriverID           int                `json:"driver_id"`
	PerformanceSummary PerformanceSummary `json:"performance_summary"`
	RaceResults        []RaceResult       `json:"race_results"`
}

type PerformanceSummary struct {
	Wins     int `json:"wins"`
	Top3     int `json:"top_3_finishes"`
	MaxSpeed int `json:"max_speed"`
}

type RaceResult struct {
	SessionKey       int     `json:"session_key"`
	CircuitShortName string  `json:"circuit_short_name"`
	Race             string  `json:"race"`
	Position         int     `json:"position"`
	FastestLap       bool    `json:"fastest_lap"`
	MaxSpeed         int     `json:"max_speed"`
	BestLapDuration  float64 `json:"best_lap_duration"`
}

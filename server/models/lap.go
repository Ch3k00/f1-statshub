package models

type Lap struct {
	DriverNumber    int     `json:"driver_number"`
	SessionKey      int     `json:"session_key"`
	LapNumber       int     `json:"lap_number"`
	LapDuration     float64 `json:"lap_duration"`
	DurationSector1 float64 `json:"duration_sector_1"`
	DurationSector2 float64 `json:"duration_sector_2"`
	DurationSector3 float64 `json:"duration_sector_3"`
	StSpeed         float64 `json:"st_speed"`
	DateStart       string  `json:"date_start"`
}

package models

type Position struct {
	DriverNumber int    `json:"driver_number"`
	SessionKey   int    `json:"session_key"`
	Position     int    `json:"position"`
	Date         string `json:"date"`
}

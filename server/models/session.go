package models

type Session struct {
	SessionKey       int    `json:"session_key"`
	SessionName      string `json:"session_name"`
	SessionType      string `json:"session_type"`
	Location         string `json:"location"`
	CountryName      string `json:"country_name"`
	Year             int    `json:"year"`
	CircuitShortName string `json:"circuit_short_name"`
	DateStart        string `json:"date_start"`
}

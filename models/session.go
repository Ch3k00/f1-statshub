package models

type Session struct {
	SessionKey       int    `json:"session_key"`
	CountryName      string `json:"country_name"`
	DateStart        string `json:"date_start"`
	Year             int    `json:"year"`
	CircuitShortName string `json:"circuit_short_name"`
}

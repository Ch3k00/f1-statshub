package models

type Driver struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	DriverNumber int    `json:"driver_number"`
	TeamName     string `json:"team_name"`
	CountryCode  string `json:"country_code"`
}

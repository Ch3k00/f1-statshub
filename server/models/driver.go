package models

// Driver representa un piloto de Fórmula 1
type Driver struct {
	DriverNumber int    `json:"driver_number"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	NameAcronym  string `json:"name_acronym"`
	TeamName     string `json:"team_name"`
	CountryCode  string `json:"country_code"`
}

// NewDriver crea una nueva instancia de Driver
func NewDriver(number int, firstName, lastName, acronym, team, country string) *Driver {
	return &Driver{
		DriverNumber: number,
		FirstName:    firstName,
		LastName:     lastName,
		NameAcronym:  acronym,
		TeamName:     team,
		CountryCode:  country,
	}
}

package models

type SeasonSummary struct {
	Season            int                `json:"season"`
	Top3Winners       []SeasonResultItem `json:"top_3_winners"`
	Top3FastestLaps   []SeasonResultItem `json:"top_3_fastest_laps"`
	Top3PolePositions []SeasonResultItem `json:"top_3_pole_positions"`
}

type SeasonResultItem struct {
	Position    int    `json:"position"`
	Driver      string `json:"driver"`
	Team        string `json:"team"`
	Country     string `json:"country"`
	Wins        int    `json:"wins,omitempty"`
	Poles       int    `json:"poles,omitempty"`
	FastestLaps int    `json:"fastest_laps,omitempty"`
}

package api

import (
	"encoding/json"
	"f1-statshub/server/models"
	"fmt"
	"io"
	"net/http"
)

const OpenF1BaseURL = "https://api.openf1.org/v1"

func GetDriversBySession(sessionKey int) ([]models.Driver, error) {
	url := fmt.Sprintf("%s/drivers?session_key=%d", OpenF1BaseURL, sessionKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var drivers []models.Driver
	if err := json.Unmarshal(body, &drivers); err != nil {
		return nil, err
	}

	return drivers, nil
}

func GetRaces(year int) ([]models.Session, error) {
	url := fmt.Sprintf("%s/sessions?session_name=Race&year=%d", OpenF1BaseURL, year)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var sessions []models.Session
	if err := json.Unmarshal(body, &sessions); err != nil {
		return nil, err
	}

	return sessions, nil
}

// Agrega esto debajo de las funciones que ya tienes en server/api/openf1.go:

func GetPositions(sessionKey int) ([]models.Position, error) {
	url := fmt.Sprintf("%s/position?session_key=%d", OpenF1BaseURL, sessionKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var positions []models.Position
	if err := json.Unmarshal(body, &positions); err != nil {
		return nil, err
	}

	return positions, nil
}

func GetLaps(sessionKey int) ([]models.Lap, error) {
	url := fmt.Sprintf("%s/laps?session_key=%d", OpenF1BaseURL, sessionKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var laps []models.Lap
	if err := json.Unmarshal(body, &laps); err != nil {
		return nil, err
	}

	return laps, nil
}

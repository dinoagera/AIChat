package client

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/dinoagera/AIChat/config"
)

type ClientWeather struct {
	urlAPI       string
	MasterAPIKey string
	httpClient   *http.Client
	log          *slog.Logger
	searchCity   ClientSearch
}
type Weather struct {
}

func NewClientWeather(cfg *config.Config, log *slog.Logger, clientSearch ClientSearch) *ClientWeather {
	return &ClientWeather{
		urlAPI:       cfg.BaseURLWeather,
		MasterAPIKey: cfg.APIKey,
		log:          log,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
			// Transport: &http.Transport{

			// }, configurations later if will be error from weather api
		},
		searchCity: clientSearch,
	}
}
func (cw *ClientWeather) CurrentWeather(city string) (Weather, error) {
	respSearch, err := cw.searchCity.SearchCordinate(city)
	if err != nil {
		cw.log.Info("failed to get cord by city name", "err", err)
		return Weather{}, err
	}
	url := fmt.Sprintf("%s?lat=%s&lon=%s&key=%s&include=minetly", cw.urlAPI, respSearch.Latitude, respSearch.Longitude, cw.MasterAPIKey)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		cw.log.Info("failed to create req in weather", "err", err)
		return Weather{}, err
	}
	resp, err := cw.httpClient.Do(req)
	if err != nil {
		cw.log.Info("failed to do req in weather", "err", err)
		return Weather{}, err
	}
	if resp.StatusCode != http.StatusOK {
		cw.log.Error("Weather API returned error", "status", resp.StatusCode)
		return Weather{}, err
	}
	var weatherResp Weather
	if err := json.NewDecoder(resp.Body).Decode(&weatherResp); err != nil {
		cw.log.Info("failed to decode", "err", err)
		return Weather{}, err
	}
	return weatherResp, nil
}

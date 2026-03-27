package client

import (
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
		//
	}
	resp, err := cw.httpClient.Do(req)
	if err != nil {
		//
	}
	if resp.StatusCode != http.StatusOK {
		//
	}

}

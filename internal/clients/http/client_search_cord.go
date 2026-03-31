package client

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/dinoagera/AIChat/config"
)

type ClientSearch struct {
	urlAPI     string
	httpClient *http.Client
	log        *slog.Logger
}
type ResponseClientSearch struct {
	Latitude  string `json:"lat"`
	Longitude string `json:"lon"`
}

func NewClientSearch(cfg *config.Config, log *slog.Logger) *ClientSearch {
	return &ClientSearch{
		urlAPI: cfg.BaseURLWeather,
		log:    log,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
			// Transport: &http.Transport{

			// }, configurations later if will be error from weather api
		},
	}
}
func (cs *ClientSearch) SearchCordinate(city string) (ResponseClientSearch, error) {
	url := fmt.Sprintf("%s?city=%s&format=json&limit=1", cs.urlAPI, city)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		cs.log.Info("failed to create new request in search cord", "err", err)
		return ResponseClientSearch{}, err
	}
	req.Header.Set("Accept", "application/json")
	resp, err := cs.httpClient.Do(req)
	if err != nil {
		cs.log.Info("failed to do req in search cord", "err", err)
		return ResponseClientSearch{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		cs.log.Error("CordSearch API returned error", "status", resp.StatusCode)
		return ResponseClientSearch{}, err
	}
	var searchResp ResponseClientSearch
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		cs.log.Info("failed to decode", "err", err)
		return ResponseClientSearch{}, err
	}
	return searchResp, nil
}

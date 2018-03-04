package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type (
	WidgetConfig struct {
		WidgetWeatherApiConfig    WidgetWeatherApiConfig    `json:"widget-weather"`
		WidgetLocationConfig      WidgetLocationConfig      `json:"widget-location"`
		WidgetTorrentStatusConfig WidgetTorrentStatusConfig `json:"widget-torrent-status"`
	}

	WidgetWeatherApiConfig struct {
		ApiToken string `json:"api-token"`
		City     string `json:"city"`
	}

	WidgetLocationConfig struct {
		GoogleMapsToken                  string `json:"google-maps-token"`
		StationaryLocationGpsCoordinates string `json:"stationary-location-gps-coordinates"`
		LocationProviderUrl              string `json:"location-provider-url"`
	}

	WidgetTorrentStatusConfig struct {
		TorrentWebApiUrl string `json:"torrent-web-api-url"`
		Username         string `json:"username"`
		Password         string `json:"password"`
	}
)

func NewConfig() *WidgetConfig {

	workingDir, _ := filepath.Abs(filepath.Dir("."))
	configFile, err := os.Open(workingDir + "/config.json")

	if err != nil {
		errors.New("failed to read configuration file: " + err.Error())
	}
	decoder := json.NewDecoder(configFile)

	widgetConfig := &WidgetConfig{}
	err = decoder.Decode(widgetConfig)
	if err != nil {
		errors.New("Error decoding configuration: " + err.Error())
	}

	return widgetConfig
}

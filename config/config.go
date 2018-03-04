package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type (

	General struct {
       ResourcesDir string `json:"resources-dir"`
	}

	Widget struct {
		General General  `json:"general"`
		WidgetWeatherApiConfig    WidgetWeatherApi    `json:"widget-weather"`
		WidgetLocationConfig      WidgetLocation      `json:"widget-location"`
		WidgetTorrentStatusConfig WidgetTorrentStatus `json:"widget-torrent-status"`
	}

	WidgetWeatherApi struct {
		ApiToken string `json:"api-token"`
		City     string `json:"city"`
	}

	WidgetLocation struct {
		GoogleMapsToken                  string `json:"google-maps-token"`
		StationaryLocationGpsCoordinates string `json:"stationary-location-gps-coordinates"`
		LocationProviderUrl              string `json:"location-provider-url"`
	}

	WidgetTorrentStatus struct {
		TorrentWebApiUrl string `json:"torrent-web-api-url"`
		Username         string `json:"username"`
		Password         string `json:"password"`
	}
)

func NewConfig() *Widget {

	workingDir, _ := filepath.Abs(filepath.Dir("."))
	configFile, err := os.Open(workingDir + "/config.json")

	if err != nil {
		errors.New("failed to read configuration file: " + err.Error())
	}
	decoder := json.NewDecoder(configFile)

	widgetConfig := &Widget{}
	err = decoder.Decode(widgetConfig)
	if err != nil {
		errors.New("Error decoding configuration: " + err.Error())
	}

	return widgetConfig
}

package config

import (
	"encoding/json"
	"errors"
	"flag"
	"os"
)

type (
	AppConfig struct {
		WidgetWeatherApiConfig    WidgetWeatherApi    `json:"widget-weather"`
		WidgetLocationConfig      WidgetLocation      `json:"widget-location"`
		WidgetTorrentStatusConfig WidgetTorrentStatus `json:"widget-torrent-status"`
		WidgetCo2Meter            WidgetCo2Meter      `json:"widget-co2-meter"`
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

	WidgetCo2Meter struct {
		PathToDevice string `json:"path-to-device"`
	}
)

var App *AppConfig
var workingDir = flag.String("working-dir", "", "Path to working directory")

func init() {
	flag.Parse()
}

func InitializeConfiguration() *AppConfig {

	configFile, err := os.Open(*workingDir + "/config.json")

	if err != nil {
		errors.New("failed to read configuration file: " + err.Error())
	}
	decoder := json.NewDecoder(configFile)

	App = &AppConfig{}

	err = decoder.Decode(App)

	if err != nil {
		errors.New("Error decoding configuration: " + err.Error())
	}
	return App
}

func (e *AppConfig) GetResourcesDir() string {
	return *workingDir + "/resources"
}

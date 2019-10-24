package config

import (
	"encoding/json"
	"flag"
	"log"
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
		PathToDevice     string `json:"path-to-device"`
		WarningThreshold int  `json:"warning-threshold"`
	}
)

var App *AppConfig
var WorkingDir = flag.String("working-dir", "", "Path to working directory")

func init() {
	flag.Parse()
}

func InitializeConfiguration() *AppConfig {

	configFile, err := os.Open(*WorkingDir + "/config.json")
	if err != nil {
		log.Fatal("An error occurred while reading configuration file: " + err.Error())
	}
	decoder := json.NewDecoder(configFile)

	App = &AppConfig{}

	err = decoder.Decode(App)
	if err != nil {
		log.Fatal("An error occurred while decoding configuration file: " + err.Error())
	}
	return App
}

func (e *AppConfig) GetResourcesDir() string {
	return *WorkingDir + "/resources"
}

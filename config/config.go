package config

import (
	"encoding/json"
	"errors"
	"os"
)

type (

	//@todo make General config accessible everywhere.
	General struct {
       ResourcesDir string `json:"resources-dir"`
	}

	AppConfig struct {
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

func NewAppConfig() *AppConfig {

	//@todo should come from --working-dir argument
	configFile, err := os.Open( "/home/pi/go/src/github.com/aurisl/ledmatrix/config.json")

	if err != nil {
		errors.New("failed to read configuration file: " + err.Error())
	}
	decoder := json.NewDecoder(configFile)

	widgetConfig := &AppConfig{}
	err = decoder.Decode(widgetConfig)
	if err != nil {
		errors.New("Error decoding configuration: " + err.Error())
	}

	return widgetConfig
}

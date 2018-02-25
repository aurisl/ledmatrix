package main

import (
	"os"
	"errors"
	"encoding/json"
	"path/filepath"
)

type (
	WidgetConfig struct {
		WidgetWeatherApiConfig WidgetWeatherApiConfig `json:"widget-weather" bson:"widget-weather"`
		WidgetLocationConfig   WidgetLocationConfig   `json:"widget-location" bson:"widget-location"`
	}

	WidgetWeatherApiConfig struct {
		ApiToken string `json:"api-token" bson:"api-token"`
		City     string `json:"city" bson:"city"`
	}

	WidgetLocationConfig struct {
		GoogleMapsToken                  string `json:"google-maps-token" bson:"google-maps-token"`
		StationaryLocationGpsCoordinates string `json:"stationary-location-gps-coordinates" bson:"stationary-location-gps-coordinates"`
		LocationProviderUrl              string `json:"location-provider-url" bson:"location-provider-url"`
	}
)

func NewConfig() *WidgetConfig {

	configFile, err := os.Open(filepath.Dir("config.json"))

	if err != nil {
		errors.New("failed to read configuration file: " + err.Error())
	}
	decoder := json.NewDecoder(configFile)

	widgetConfig := &WidgetConfig{}
	err = decoder.Decode(&widgetConfig)
	if err != nil {
		errors.New("Error decoding configuration: " + err.Error())
	}

	return widgetConfig
}

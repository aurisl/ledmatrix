package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"log"
	"net/http"
	"io/ioutil"
	"time"
	"bytes"
	"fmt"
	"io"
	"image"
)

type (
	Weather struct {
		LastUpdated  time.Time    `json:"last_updated" bson:"last_updated"`
		WeatherData  WeatherAPI   `json:"weather_data" bson:"weather_data"`
		WeatherImage WeatherImage `json:"-"`
	}

	WeatherAPI struct {
		WeatherCurrent []WeatherCurrent `json:"weather" bson:"weather"`
		WeatherMain    WeatherMain      `json:"main" bson:"main"`
	}

	WeatherCurrent struct {
		Id          uint16 `json:"id" bson:"id"`
		Main        string `json:"main" bson:"main"`
		Description string `json:"description" bson:"description"`
		Icon        string `json:"icon" bson:"icon"`
	}

	WeatherMain struct {
		Temp     float64 `json:"temp" bson:"temp"`
		Pressure float64 `json:"pressure" bson:"pressure"`
		Humidity uint16  `json:"humidity" bson:"humidity"`
		TempMin  float64 `json:"temp_min" bson:"temp_min"`
		TempMax  float64 `json:"temp_max" bson:"temp_max"`
	}

	WeatherImage struct {
		Img image.Image
		Ico string
	}
)

var (
	weatherFileName = "/weather.json"
	weatherApiEndpoint = "http://api.openweathermap.org/data/2.5/weather"
)

func NewWeather() *Weather {
	return &Weather{}
}

func (w *Weather) ReadWeather(config WidgetWeatherApiConfig) {

	workingDirectory, err := filepath.Abs(filepath.Dir(resourcesDir + weatherFileName))
	if err != nil {
		log.Fatal(err)
	}

	weatherFileLocation := workingDirectory + weatherFileName
	weatherFile, err := os.Open(weatherFileLocation)

	if err != nil {
		w.updateWeatherData(weatherFileLocation, config)
		return
	}

	w.decodeWeatherJsonFile(weatherFile)

	duration := time.Since(w.LastUpdated)
	if duration.Minutes() > 15 {
		w.updateWeatherData(weatherFileLocation, config)
	}
}

func (w *Weather) updateWeatherData(weatherFileLocation string, config WidgetWeatherApiConfig) {
	fmt.Println("Updating weather data...")

	weatherApi := WeatherAPI{}

	apiResponse := readApi(config)

	if apiResponse == nil {
		return
	}

	response := bytes.NewReader(apiResponse)
	decoder := json.NewDecoder(response)

	err := decoder.Decode(&weatherApi)
	if err != nil {
		panic("Weather json decode error: " + err.Error())
		return
	}

	w.LastUpdated = time.Now()
	w.WeatherData = weatherApi

	w.persistCurrentWeatherData(weatherFileLocation)
}

func (w *Weather) decodeWeatherJsonFile(weatherFile io.Reader) {

	decoder := json.NewDecoder(weatherFile)

	err := decoder.Decode(&w)
	if err != nil {
		panic("Weather json decode error: " + err.Error())
		return
	}
}

func readApi(config WidgetWeatherApiConfig) []byte {
	req, _ := http.NewRequest("GET", weatherApiEndpoint, nil)

	q := req.URL.Query()
	q.Add("q", config.City)
	q.Add("appid", config.ApiToken)
	q.Add("units", "metric")
	req.URL.RawQuery = q.Encode()

	resp, err := http.Get(req.URL.String())
	if err != nil {
		fmt.Println("And error accurred when trying access weather API: " + err.Error())
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("And error accurred when trying access weather API: " + err.Error())
		return nil
	}

	resp.Body.Close()

	return body
}

func (w *Weather) SetSelectedImage(img image.Image, ico string) {
	w.WeatherImage = WeatherImage{img, ico}
}

func (w *Weather) persistCurrentWeatherData(weatherFileLocation string) {

	jsonData, _ := json.Marshal(w)

	err := ioutil.WriteFile(weatherFileLocation, []byte(jsonData), 0644)
	fmt.Printf("Writing weather data to %s", weatherFileLocation)
	if err != nil {
		panic("Failed to write weather data: " + err.Error())
		return
	}
}

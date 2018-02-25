package main

import (
	"github.com/mcuadros/go-rpi-rgb-led-matrix"
	"net/url"
)

type (
	Widget struct {
		name   string
		params url.Values
	}
)

var (
	resourcesDir = "resources"
)

func main() {

	widgetConfig := NewConfig()

	config := createMatrixDefaultConfiguration()

	matrix, err := rgbmatrix.NewRGBLedMatrix(config)
	fatal(err)

	toolkit := rgbmatrix.NewToolKit(matrix)
	defer toolkit.Close()

	terminate := make(chan bool)

	widget := &Widget{}

	go StartHttpServer(widget, terminate)

	for {
		switch widget.name {
		case "weather":
			DrawWeatherWidget(toolkit, terminate, widgetConfig.WidgetWeatherApiConfig)
		case "boom":
			DrawExplosion(toolkit, terminate, widget)
		case "gif":
			DrawGif(toolkit, widget, terminate)
		case "fire":
			DrawFire(toolkit, terminate, widget)
		case "location":
			DisplayDistance(toolkit, widget, terminate, widgetConfig.WidgetLocationConfig)
		default:
			DrawWeatherWidget(toolkit, terminate, widgetConfig.WidgetWeatherApiConfig)
		}
	}

}
func createMatrixDefaultConfiguration() *rgbmatrix.HardwareConfig {
	config := &rgbmatrix.DefaultConfig
	config.Rows = 32
	config.Brightness = 50

	return config
}

func fatal(err error) {
	if err != nil {
		panic(err)
	}
}

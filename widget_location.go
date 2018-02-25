package main

import (
	"time"
	"github.com/mcuadros/go-rpi-rgb-led-matrix"
	"net/http"
	"fmt"
	"io/ioutil"
	"bytes"
	"encoding/json"
	"github.com/fogleman/gg"
	"image/color"
	"image"
	"io"
)

type (
	Location struct {
		Rows []LocationElements `json:"rows" bson:"rows"`
	}

	LocationElements struct {
		LocationElement []LocationElement `json:"elements" bson:"elements"`
	}

	LocationElement struct {
		LocationDistance LocationDistance `json:"distance" bson:"distance"`
		Status           string           `json:"status" bson:"status"`
	}

	LocationDistance struct {
		Range string `json:"text" bson:"text"`
	}

	LocationCord struct {
		Lon string `json:"lon" bson:"lon"`
		Lat string `json:"lat" bson:"lat"`
	}
)

var (
	tick                  = 0
	distance              = ""
	googleMapsApiEndpoint = "https://maps.googleapis.com/maps/api/distancematrix/json"
)

type LocationAnimation struct {
	ctx    *gg.Context
	close  chan bool
	config WidgetLocationConfig
}

func DisplayDistance(toolkit *rgbmatrix.ToolKit, widget *Widget, close chan bool, config WidgetLocationConfig) {

	animation := &LocationAnimation{
		ctx:    gg.NewContext(32, 32),
		close:  close,
		config: config,
	}

	toolkit.PlayAnimation(animation)
	widget.name = ""
}

func (animation *LocationAnimation) Next() (image.Image, <-chan time.Time, error) {
	initializeLocationCanvas(animation)

	if err := animation.ctx.LoadFontFace("resources/fonts/PixelOperator.ttf", 16); err != nil {
		panic(err)
	}

	if tick == 0 {
		googleResponse := readDistance(animation.config)

		response := bytes.NewReader(googleResponse)
		decoder := json.NewDecoder(response)

		location := Location{}

		err := decoder.Decode(&location)
		if err != nil {
			fmt.Println("Error decoding current location from google API: " + err.Error())
		}

		if len(location.Rows) == 0 {
			return nil, nil, io.EOF
		}

		if len(location.Rows[0].LocationElement) == 0 {
			return nil, nil, io.EOF
		}

		locationElement := location.Rows[0].LocationElement[0]

		if locationElement.Status == "OK" {
			distance = locationElement.LocationDistance.Range
		}
	}

	if distance != "" {
		animation.ctx.SetColor(color.RGBA{255, 0, 0, 255})
		animation.ctx.DrawString(distance, 1, 18)
	} else {
		animation.ctx.SetColor(color.RGBA{255, 0, 0, 255})
		animation.ctx.DrawString("n/a", 1, 18)
	}

	tick++

	if tick == 10 {
		tick = 0
		return nil, nil, io.EOF
	}

	select {
	case <-animation.close:
		return nil, nil, io.EOF
	default:
		return animation.ctx.Image(), time.After(time.Second * 1), nil
	}

}

func initializeLocationCanvas(animation *LocationAnimation) {

	animation.ctx.SetRGB(0, 0, 0)
	animation.ctx.Clear()
	animation.ctx.SetColor(color.RGBA{255, 255, 255, 255})
	animation.ctx.SetLineWidth(1)
	animation.ctx.Stroke()
}

func readDistance(config WidgetLocationConfig) []byte {
	req, _ := http.NewRequest("GET", googleMapsApiEndpoint, nil)

	q := req.URL.Query()
	q.Add("units", "metric")
	q.Add("origins", config.StationaryLocationGpsCoordinates)

	cCord := readGeo(config)

	response := bytes.NewReader(cCord)
	decoder := json.NewDecoder(response)

	locationCord := LocationCord{}

	err := decoder.Decode(&locationCord)
	if err != nil {
		fmt.Println("Error decoding current location: " + err.Error())
		return nil
	}

	q.Add("destinations", locationCord.Lat+","+locationCord.Lon)
	q.Add("key", config.GoogleMapsToken)
	req.URL.RawQuery = q.Encode()

	resp, err := http.Get(req.URL.String())
	if err != nil {
		fmt.Println("And error accurred when trying access google API: " + err.Error())
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("And error accurred when trying access google API: " + err.Error())
		return nil
	}

	resp.Body.Close()

	return body
}

func readGeo(config WidgetLocationConfig) []byte {
	req, _ := http.NewRequest("GET", config.LocationProviderUrl, nil)

	resp, err := http.Get(req.URL.String())
	if err != nil {
		fmt.Println("And error accurred when trying to get current location: " + err.Error())
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("And error accurred when reading current location response: " + err.Error())
		return nil
	}

	resp.Body.Close()

	return body
}

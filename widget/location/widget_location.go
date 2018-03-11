package location

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fogleman/gg"
	"image"
	"image/color"
	"io"
	"io/ioutil"
	"net/http"
	"time"
	"github.com/aurisl/ledmatrix/config"
	"github.com/aurisl/ledmatrix/draw"
	"github.com/aurisl/ledmatrix/matrix"
)

type (
	Location struct {
		Rows []Elements `json:"rows" bson:"rows"`
	}

	Elements struct {
		LocationElement []Element `json:"elements" bson:"elements"`
	}

	Element struct {
		LocationDistance Distance `json:"distance" bson:"distance"`
		Status           string   `json:"status" bson:"status"`
	}

	Distance struct {
		Range string `json:"text" bson:"text"`
	}

	Coordinates struct {
		Lon string `json:"lon" bson:"lon"`
		Lat string `json:"lat" bson:"lat"`
	}
)

var (
	tick                  = 0
	distance              = ""
	googleMapsApiEndpoint = "https://maps.googleapis.com/maps/api/distancematrix/json"
)

type animation struct {
	ctx    *gg.Context
	config config.WidgetLocation
}

func Draw(toolkit *matrix.LedToolKit,
	config *config.AppConfig) {

	animation := &animation{
		ctx:    gg.NewContext(32, 32),
		config: config.WidgetLocationConfig,
	}

	toolkit.PlayAnimation(animation)
}

func (animation *animation) Next() (image.Image, <-chan time.Time, error) {
	initializeLocationCanvas(animation)

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

	text := distance
	if distance == "" {
		text = "n\a"
	}

	draw.Text(text, 1, 18, animation.ctx, color.RGBA{255, 0, 0, 255})

	tick++

	if tick == 10 {
		tick = 0
		return nil, nil, io.EOF
	}

	return animation.ctx.Image(), time.After(time.Second * 1), nil

}

func initializeLocationCanvas(animation *animation) {

	animation.ctx.SetRGB(0, 0, 0)
	animation.ctx.Clear()
	animation.ctx.SetColor(color.RGBA{255, 255, 255, 255})
	animation.ctx.SetLineWidth(1)
	animation.ctx.Stroke()
}

func readDistance(locationConfig config.WidgetLocation) []byte {
	req, _ := http.NewRequest("GET", googleMapsApiEndpoint, nil)

	q := req.URL.Query()
	q.Add("units", "metric")
	q.Add("origins", locationConfig.StationaryLocationGpsCoordinates)

	cCord := readGeo(locationConfig)

	response := bytes.NewReader(cCord)
	decoder := json.NewDecoder(response)

	locationCord := Coordinates{}

	err := decoder.Decode(&locationCord)
	if err != nil {
		fmt.Println("Error decoding current location: " + err.Error())
		return nil
	}

	q.Add("destinations", locationCord.Lat+","+locationCord.Lon)
	q.Add("key", locationConfig.GoogleMapsToken)
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

func readGeo(locationConfig config.WidgetLocation) []byte {
	req, _ := http.NewRequest("GET", locationConfig.LocationProviderUrl, nil)

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

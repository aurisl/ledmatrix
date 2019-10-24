package location

import (
	"bytes"
	"encoding/json"
	"github.com/aurisl/ledmatrix/config"
	"github.com/aurisl/ledmatrix/draw"
	"github.com/aurisl/ledmatrix/matrix"
	"github.com/fogleman/gg"
	"image"
	"image/color"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
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
}

func Draw(toolkit *matrix.LedToolKit) {

	animation := &animation{
		ctx:    gg.NewContext(32, 32),
	}

	err := toolkit.PlayAnimation(animation)

	if err != nil {
		log.Println("An error occurred while playing location animation: " + err.Error())
	}
}

func (animation *animation) Next() (image.Image, <-chan time.Time, error) {
	initializeLocationCanvas(animation)

	if tick == 0 {
		googleResponse := readDistance(config.App.WidgetLocationConfig)

		response := bytes.NewReader(googleResponse)
		decoder := json.NewDecoder(response)

		location := Location{}

		err := decoder.Decode(&location)
		if err != nil {
			log.Println("Error decoding current location from google API: " + err.Error())
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

	if distance == "" {
		distance = "n\a"
	}

	draw.TextScrolling(distance, 20, animation.ctx, color.RGBA{R: 255, G: 255, A: 255})

	tick++

	if tick == 100 {
		tick = 0
		return nil, nil, io.EOF
	}

	return animation.ctx.Image(), time.After(time.Millisecond * 100), nil

}

func initializeLocationCanvas(animation *animation) {

	animation.ctx.SetRGB(0, 0, 0)
	animation.ctx.Clear()
	animation.ctx.SetColor(color.RGBA{R: 255, G: 255, B: 255, A: 255})
	animation.ctx.SetLineWidth(1)
	animation.ctx.Stroke()
}

func readDistance(locationConfig config.WidgetLocation) []byte {
	req, _ := http.NewRequest("GET", googleMapsApiEndpoint, nil)

	q := req.URL.Query()
	q.Add("units", "metric")
	q.Add("origins", locationConfig.StationaryLocationGpsCoordinates)

	cCord := readGeo(locationConfig)

	stationaryCoordinatesResponse := bytes.NewReader(cCord)
	decoder := json.NewDecoder(stationaryCoordinatesResponse)

	locationCord := Coordinates{}

	err := decoder.Decode(&locationCord)
	if err != nil {
		log.Println("Error decoding current location: " + err.Error())
		return nil
	}

	q.Add("destinations", locationCord.Lat+","+locationCord.Lon)
	q.Add("key", locationConfig.GoogleMapsToken)
	req.URL.RawQuery = q.Encode()

	googleApiResponse, err := http.Get(req.URL.String())
	if err != nil {
		log.Println("An error occurred when trying access google API: " + err.Error())
		return nil
	}

	body, err := ioutil.ReadAll(googleApiResponse.Body)

	if err != nil {
		log.Println("An error occurred when trying access google API: " + err.Error())
		return nil
	}

	err = googleApiResponse.Body.Close()
	if err != nil {
		log.Printf("An error occurred while closing stationaryCoordinatesResponse body '%s'", err.Error())
		return nil
	}

	return body
}

func readGeo(locationConfig config.WidgetLocation) []byte {
	request, _ := http.NewRequest("GET", locationConfig.LocationProviderUrl, nil)

	response, err := http.Get(request.URL.String())
	if err != nil {
		log.Println("An error occurred when trying to get current location: " + err.Error())
		return nil
	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Println("An error occurred when reading current location response: " + err.Error())
		return nil
	}

	err = response.Body.Close()
	if err != nil {
		log.Printf("An error occurred while closing response body '%s'", err.Error())
		return nil
	}

	return body
}

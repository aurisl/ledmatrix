package weather

import (
	"bytes"
	"fmt"
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"github.com/aurisl/ledmatrix/config"
	"github.com/aurisl/ledmatrix/draw"
	"github.com/aurisl/ledmatrix/matrix"
)

var (
	provider    = NewWeather()
	lastUpdate  time.Time
	displayTick      = false
	loopTick= 0
	mainLoop  = time.Second * 200
)

type animation struct {
	ctx           *gg.Context
	close         chan bool
	weatherConfig config.WidgetWeatherApi
}

func Draw(toolkit *matrix.Toolkit, close chan bool, weatherConfig config.WidgetWeatherApi) {

	animation := &animation{
		ctx:           toolkit.Ctx,
		close:         close,
		weatherConfig: weatherConfig,
	}

	toolkit.MatrixToolkit.PlayAnimation(animation)
}

func (animation *animation) Next() (image.Image, <-chan time.Time, error) {

	draw.ClearCanvas(animation.ctx)

	currentDate := time.Now()

	hour, err := strconv.ParseInt(currentDate.Format("15"), 10, 8)
	minute := currentDate.Format("04")
	second, err := strconv.ParseInt(currentDate.Format("05"), 10, 8)

	if err != nil {
		fmt.Println(err)
	}

	if hour >= 23 || hour < 07 {
		animation.ctx.SetRGB(0, 0, 0)
		animation.ctx.Clear()
	} else if minute == "00" && second < 5 {
		mainLoop = time.Microsecond * 20
		draw.GradientFlashing(animation.ctx)
	} else {
		mainLoop = time.Second * 1
		drawTime(animation)
		draw.GradientLine(animation.ctx)
		drawWeatherInformation(animation)
	}

	select {
	case <-animation.close:
		return nil, nil, io.EOF
	default:
		return animation.ctx.Image(), time.After(mainLoop), nil
	}
}

func drawWeatherInformation(animation *animation) {

	weatherData := readWeatherData(animation.weatherConfig)

	if len(weatherData.WeatherCurrent) == 0 {
		return
	}

	if loopTick == 3 {

		loopTick = 0
		icon := weatherData.WeatherCurrent[0].Icon

		if icon == provider.WeatherImage.Ico {
			animation.ctx.DrawImage(provider.WeatherImage.Img, 5, 13)
			return
		}

		workingDirectory, _ := filepath.Abs(filepath.Dir("resources/img"))
		iconPath := workingDirectory + "/resources/img/weather_icons/ " + icon + ".png"
		iconFile, err := os.Open(iconPath)

		var img image.Image

		if err != nil {
			resp, err := http.Get("http://openweathermap.org/img/w/" + icon + ".png")

			if err != nil {
				return
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return
			}

			resp.Body.Close()

			response := bytes.NewReader(body)

			imgOriginal, _, err := image.Decode(response)
			if err != nil {
				panic(err)
			}

			img = resize.Resize(25, 25, imgOriginal, resize.Lanczos3)

			buff := new(bytes.Buffer)
			png.Encode(buff, img)
			ioutil.WriteFile(iconPath, buff.Bytes(), 0777)

		} else {
			img, _, _ = image.Decode(iconFile)
		}

		provider.SetSelectedImage(img, icon)
		animation.ctx.DrawImage(img, 5, 13)

		return
	}

	animation.ctx.SetColor(color.RGBA{255, 255, 0, 255})
	animation.ctx.DrawString(strconv.FormatFloat(weatherData.WeatherMain.Temp, 'f', 0, 64)+"°C", 4, 30)
	loopTick++

}

func drawTime(animation *animation) {

	draw.Text(time.Now().Format("15"), 1, 13, animation.ctx, color.RGBA{255, 0, 0, 255})
	drawTimeSemicolon(animation)
	draw.Text(time.Now().Format("04"), 17, 13, animation.ctx, color.RGBA{255, 0, 0, 255})
}

func drawTimeSemicolon(animation *animation) {
	if displayTick == true {
		draw.Text(":", 14, 12, animation.ctx, color.RGBA{255, 0, 0, 255})
		displayTick = false
		return
	}
	displayTick = true
}

func readWeatherData(weatherConfig config.WidgetWeatherApi) API {
	duration := time.Since(lastUpdate)
	if duration.Minutes() > 15 {
		go provider.ReadWeather(weatherConfig)
		lastUpdate = time.Now()
	}

	return provider.WeatherData
}
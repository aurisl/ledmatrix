package main

import (
	"bytes"
	"fmt"
	"github.com/fogleman/gg"
	"github.com/mcuadros/go-rpi-rgb-led-matrix"
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
)

var (
	WeatherProvider = NewWeather()
	lastUpdate      time.Time
	displayTick     = false
	loopTick        = 0
	radianTick      = 0
	back            = false
	mainLoop        = time.Second * 1
)

type WeatherAnimation struct {
	ctx    *gg.Context
	close  chan bool
	config WidgetWeatherApiConfig
}

func DrawWeatherWidget(toolkit *rgbmatrix.ToolKit, close chan bool, config WidgetWeatherApiConfig) {

	animation := &WeatherAnimation{
		ctx:    gg.NewContext(32, 32),
		close:  close,
		config: config,
	}

	if err := animation.ctx.LoadFontFace("resources/fonts/PixelOperator.ttf", 16); err != nil {
		panic(err)
	}

	toolkit.PlayAnimation(animation)

}

func (animation *WeatherAnimation) Next() (image.Image, <-chan time.Time, error) {

	initializeCanvas(animation)

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
		drawHourBang(animation)

	} else {

		mainLoop = time.Second * 1

		drawTime(animation)
		drawLine(animation)
		drawWeatherInformation(animation)
	}

	select {
	case <-animation.close:
		return nil, nil, io.EOF
	default:
		return animation.ctx.Image(), time.After(mainLoop), nil
	}
}

func initializeCanvas(animation *WeatherAnimation) {

	animation.ctx.SetRGB(0, 0, 0)
	animation.ctx.Clear()
	animation.ctx.SetColor(color.RGBA{255, 255, 255, 255})
	animation.ctx.SetLineWidth(1)
	animation.ctx.Stroke()
}

func drawWeatherInformation(animation *WeatherAnimation) {

	weatherData := readWeatherData(animation.config)

	if len(weatherData.WeatherCurrent) == 0 {
		return
	}

	if loopTick == 3 {

		loopTick = 0

		icon := weatherData.WeatherCurrent[0].Icon

		if icon == WeatherProvider.WeatherImage.Ico {
			animation.ctx.DrawImage(WeatherProvider.WeatherImage.Img, 5, 13)
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

		WeatherProvider.SetSelectedImage(img, icon)
		animation.ctx.DrawImage(img, 5, 13)

		return

	}

	animation.ctx.SetColor(color.RGBA{255, 255, 0, 255})
	animation.ctx.DrawString(strconv.FormatFloat(weatherData.WeatherMain.Temp, 'f', 0, 64)+"°C", 4, 30)
	loopTick++

}

func drawLine(animation *WeatherAnimation) {

	grad := gg.NewLinearGradient(1, 1, 32, 1)
	grad.AddColorStop(0, color.RGBA{0, 255, 0, 255})
	grad.AddColorStop(1, color.RGBA{0, 0, 255, 255})
	grad.AddColorStop(0.5, color.RGBA{255, 0, 0, 255})
	animation.ctx.SetStrokeStyle(grad)

	animation.ctx.DrawLine(1, 17, 32, 17)
	animation.ctx.SetLineWidth(1)
	animation.ctx.Stroke()
}

func drawTime(animation *WeatherAnimation) {

	animation.ctx.SetColor(color.RGBA{255, 0, 0, 255})

	animation.ctx.DrawString(time.Now().Format("15"), 1, 13)
	drawTimeSemicolon(animation)
	animation.ctx.DrawString(time.Now().Format("04"), 17, 13)

}

func drawHourBang(animation *WeatherAnimation) {

	var y1 float64

	if back == false {
		y1 = float64(radianTick + 2)
	} else if back == true {
		y1 = float64(radianTick - 2)
	}

	animation.ctx.SetRGB(0, 0, 0)
	animation.ctx.Clear()

	grad := gg.NewRadialGradient(15, 15, 0, 15, 15, y1)

	grad.AddColorStop(0, color.RGBA{0, 255, 0, 255})
	grad.AddColorStop(1, color.RGBA{0, 0, 255, 255})
	grad.AddColorStop(0.5, color.RGBA{255, 0, 0, 255})

	animation.ctx.SetFillStyle(grad)
	animation.ctx.DrawRectangle(0, 0, 32, 32)
	animation.ctx.Fill()

	if radianTick < 30 && back == false {
		radianTick++

		if radianTick >= 30 {
			back = true
		}

	} else {
		radianTick--

		if radianTick == 5 {
			back = false
		}
	}

}

func drawTimeSemicolon(animation *WeatherAnimation) {
	if displayTick == true {
		animation.ctx.DrawString(":", 14, 12)
		displayTick = false
		return
	}
	displayTick = true
}

func readWeatherData(config WidgetWeatherApiConfig) WeatherAPI {
	duration := time.Since(lastUpdate)
	if duration.Minutes() > 15 {
		go WeatherProvider.ReadWeather(config)
		lastUpdate = time.Now()
	}

	return WeatherProvider.WeatherData
}

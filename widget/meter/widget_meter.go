package meter

import (
	"github.com/aurisl/ledmatrix/config"
	"github.com/aurisl/ledmatrix/draw"
	"github.com/aurisl/ledmatrix/matrix"
	"github.com/fogleman/gg"
	"image"
	"image/color"
	"log"
	"strconv"
	"time"
)

var (
	mainLoop           = time.Second
	CurrentMeasurement = &Measurement{}
)

type animation struct {
	ctx *gg.Context
}

func Draw(toolkit *matrix.LedToolKit) {
	animation := &animation{ctx: toolkit.Ctx}

	err := toolkit.PlayAnimation(animation)
	if err != nil {
		log.Println("An error occurred while player meter animation: " + err.Error())
	}
}

func (animation *animation) Next() (image.Image, <-chan time.Time, error) {
	draw.ClearCanvas(animation.ctx)

	draw.Text(strconv.Itoa(CurrentMeasurement.Co2), 4, 13, animation.ctx, createCo2Color(CurrentMeasurement.Co2))
	draw.GradientLine(animation.ctx)
	animation.ctx.SetColor(color.RGBA{R: 255, G: 255, A: 255})
	animation.ctx.DrawString(strconv.FormatFloat(CurrentMeasurement.Temperature, 'f', 0, 64)+"°C", 4, 30)

	return animation.ctx.Image(), time.After(mainLoop), nil
}

func createCo2Color(co2 int) color.RGBA {
	if co2 <= 800 {
		return color.RGBA{G: 255, A: 255}
	}

	if co2 > 800 && co2 <= 1200 {
		return color.RGBA{R: 255, G: 150, A: 255}
	}

	return color.RGBA{R: 255, A: 255}
}

func Measure() {
	meter := new(Meter)
	err := meter.Open(config.App.WidgetCo2Meter.PathToDevice)
	if err != nil {
		log.Printf("CO2 Meter could not open device file '%s'", config.App.WidgetCo2Meter.PathToDevice)
		return
	}

	for {
		select {
		case <-time.After(time.Second):
			result, err := meter.Read()
			if err != nil {
				log.Printf("CO2 Meter returned error: '%v'", err)
			}
			CurrentMeasurement = result
		}
	}
}

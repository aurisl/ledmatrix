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
	mainLoop    = time.Second
	temperature = 0.0
	co2         = 0.0
	quitCh      = make(chan struct{})
)

type animation struct {
	ctx    *gg.Context
	config config.WidgetLocation
}

func Draw(toolkit *matrix.LedToolKit, config *config.AppConfig) {
	animation := &animation{
		ctx:    gg.NewContext(32, 32),
		config: config.WidgetLocationConfig,
	}

	err := toolkit.PlayAnimation(animation)
	if err != nil {
		log.Fatalf("An error occurred while player meter animation: " + err.Error())
	}

	close(quitCh)
	go measure()
}

func (animation *animation) Next() (image.Image, <-chan time.Time, error) {
	draw.ClearCanvas(animation.ctx)

	draw.Text(strconv.FormatFloat(co2, 'f', 0, 64), 4, 30, animation.ctx, color.RGBA{255, 0, 0, 255})
	draw.GradientLine(animation.ctx)
	animation.ctx.SetColor(color.RGBA{255, 255, 0, 255})
	animation.ctx.DrawString(strconv.FormatFloat(temperature, 'f', 0, 64)+"°C", 4, 30)

	return animation.ctx.Image(), time.After(mainLoop), nil
}

func measure() {
	quitCh = make(chan struct{})
	meter := new(Meter)
	err := meter.Open("/dev/hidraw0")
	if err != nil {
		log.Fatalf("Could not open '/dev/hidraw0'")
		return
	}

	for {
		select {
		case <-quitCh:
			return
		case <-time.After(5 * time.Second):
			result, err := meter.Read()
			if err != nil {
				log.Fatalf("Meter reader returned error: '%v'", err)
			}
			temperature = result.Temperature
			co2 = float64(result.Co2)
		}
	}
}

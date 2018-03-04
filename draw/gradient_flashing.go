package draw

import (
	"github.com/fogleman/gg"
	"image/color"
)

var radianTick = 0
var back = false

func GradientFlashing(ctx *gg.Context) {

	var y1 float64

	if back == false {
		y1 = float64(radianTick + 2)
	} else if back == true {
		y1 = float64(radianTick - 2)
	}

	ctx.SetRGB(0, 0, 0)
	ctx.Clear()

	grad := gg.NewRadialGradient(15, 15, 0, 15, 15, y1)

	grad.AddColorStop(0, color.RGBA{0, 255, 0, 255})
	grad.AddColorStop(1, color.RGBA{0, 0, 255, 255})
	grad.AddColorStop(0.5, color.RGBA{255, 0, 0, 255})

	ctx.SetFillStyle(grad)
	ctx.DrawRectangle(0, 0, 32, 32)
	ctx.Fill()

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

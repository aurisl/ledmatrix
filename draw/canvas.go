package draw

import (
	"github.com/fogleman/gg"
	"image/color"
)

func ClearCanvas(ctx *gg.Context) {

	ctx.SetRGB(0, 0, 0)
	ctx.Clear()
	ctx.SetColor(color.RGBA{255, 255, 255, 255})
}

func RedScreen(ctx *gg.Context, tick int) {
	if tick % 10 == 0 {
		ctx.SetColor(color.RGBA{R: 255, A: 255})
	} else {
		ctx.SetColor(color.RGBA{A: 255})
	}
	ctx.Clear()
}

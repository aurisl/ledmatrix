package draw

import (
	"image/color"
	"github.com/fogleman/gg"
)

func ClearCanvas(ctx *gg.Context) {

	ctx.SetRGB(0, 0, 0)
	ctx.Clear()
	ctx.SetColor(color.RGBA{255, 255, 255, 255})
}

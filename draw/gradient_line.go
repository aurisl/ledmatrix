package draw

import (
	"github.com/fogleman/gg"
	"image/color"
)

func GradientLine(ctx *gg.Context) {

	grad := gg.NewLinearGradient(1, 1, 32, 1)
	grad.AddColorStop(0, color.RGBA{0, 255, 0, 255})
	grad.AddColorStop(1, color.RGBA{0, 0, 255, 255})
	grad.AddColorStop(0.5, color.RGBA{255, 0, 0, 255})
	ctx.SetStrokeStyle(grad)

	ctx.DrawLine(1, 17, 32, 17)
	ctx.SetLineWidth(1)
	ctx.Stroke()
}

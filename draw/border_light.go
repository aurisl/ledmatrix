package draw

import (
	"github.com/fogleman/gg"
	"image/color"
)

type (
	BorderShared struct {
		tick uint8
		step uint8
	}
)

func NewBorderShader() *BorderShared {
	return &BorderShared{0, 1}
}

func (border *BorderShared) SetTick(tick uint8) {
	border.tick = tick
}

func (border *BorderShared) SetStep(step uint8) {
	border.step = step
}

func (border *BorderShared) DrawBorderShader(ctx *gg.Context) {

	if border.tick > 128 {
		border.tick = 0
	}

	switch {
	case border.tick <= 31: //left
		xlCallback := func(i uint8) uint8 { return border.tick + i }
		ylCallback := func(i uint8) uint8 { return 0 }
		border.drawLine(ctx, xlCallback, ylCallback)
	case border.tick > 31 && border.tick <= 63: //down
		xlCallback := func(i uint8) uint8 { return 31 }
		ylCallback := func(i uint8) uint8 { return border.tick - 31 + i }
		border.drawLine(ctx, xlCallback, ylCallback)
	case border.tick > 63 && border.tick <= 95: //right
		xlCallback := func(i uint8) uint8 { return 31 - (border.tick - 64) - i }
		ylCallback := func(i uint8) uint8 { return 31 }
		border.drawLine(ctx, xlCallback, ylCallback)
	case border.tick > 95 && border.tick <= 127: //up
		xlCallback := func(i uint8) uint8 { return 0 }
		ylCallback := func(i uint8) uint8 { return 31 - (border.tick - 96) - i }
		border.drawLine(ctx, xlCallback, ylCallback)
	}

	border.tick = border.tick + border.step
}

func (border *BorderShared) drawLine(ctx *gg.Context, xCallback func(i uint8) uint8, yCallback func(i uint8) uint8) {

	var numberOfSteps uint8 = 10
	var colorShiftStep uint8 = 25
	var lineColor uint8 = 0
	var i uint8 = 1

	for i = 1; i <= numberOfSteps; i++ {
		lineColor = lineColor + colorShiftStep
		ctx.SetColor(color.RGBA{R: lineColor, A: 255})

		x := xCallback(i)
		y := yCallback(i)

		ctx.SetPixel(int(x), int(y))
	}
}

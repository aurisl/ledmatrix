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

func (border *BorderShared) GetTick() uint8 {
	return border.tick
}

func (border *BorderShared) SetStep(step uint8) {
	border.step = step
}

func (border *BorderShared) DrawBorderShader(ctx *gg.Context) {

	if border.tick >= 128 {
		border.tick = 1
	}

	switch {
	case border.tick <= 31: //right
		xlCallback := func(i uint8) uint8 {
			x := border.tick + i
			if x >= 31 {
				return 31
			}
			return x
		}
		ylCallback := func(i uint8) uint8 {
			x := border.tick + i
			if x >= 31 {
				return i - (31 - border.tick)
			}
			return 0
		}
		border.drawLine(ctx, xlCallback, ylCallback)
	case border.tick > 31 && border.tick <= 63: //down
		xlCallback := func(i uint8) uint8 {
			x := border.tick + i
			if x >= 63 {
				return 31 - (border.tick - 63) - i
			}
			return 31
		}
		ylCallback := func(i uint8) uint8 {
			x := border.tick + i
			if x >= 63 {
				return 31
			}
			return border.tick - 31 + i
		}
		border.drawLine(ctx, xlCallback, ylCallback)
	case border.tick > 63 && border.tick <= 95: //left
		xlCallback := func(i uint8) uint8 {
			x := border.tick + i
			if x >= 95 {
				return 0
			}
			return 31 - (border.tick - 64) - i
		}
		ylCallback := func(i uint8) uint8 {
			x := border.tick + i
			if x >= 95 {
				return 63 - (border.tick - 63) - i
			}
			return 31
		}
		border.drawLine(ctx, xlCallback, ylCallback)
	case border.tick > 95 && border.tick <= 127: //up
		xlCallback := func(i uint8) uint8 {
			x := border.tick + i

			if x >= 127 {
				return i - 1 - (126 - border.tick)
			}
			return 0
		}
		ylCallback := func(i uint8) uint8 {
			x := border.tick + i

			if x >= 127 {
				return 0
			}
			return 31 - (border.tick - 96) - i
		}
		border.drawLine(ctx, xlCallback, ylCallback)
	}

	border.tick = border.tick + border.step
}

func (border *BorderShared) drawLine(ctx *gg.Context, xCallback func(i uint8) uint8, yCallback func(i uint8) uint8) {

	var numberOfSteps uint8 = 10
	var colorShiftStep uint8 = 20
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

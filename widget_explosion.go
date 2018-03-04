package main

import (
	"github.com/fogleman/gg"
	"github.com/mcuadros/go-rpi-rgb-led-matrix"
	"image"
	"image/color"
	"io"
	"time"
)

var (
	renderTick    = 0
	numberOfLoops = 0
)

type ExplosionAnimation struct {
	ctx    *gg.Context
	close  chan bool
	widget *Widget
}

func DrawExplosion(toolkit *rgbmatrix.ToolKit, close chan bool, widget *Widget) {

	animation := &ExplosionAnimation{
		ctx:    gg.NewContext(32, 32),
		close:  close,
		widget: widget,
	}

	toolkit.PlayAnimation(animation)
}

func (animation *ExplosionAnimation) Next() (image.Image, <-chan time.Time, error) {

	explosion(animation)

	select {
	case <-animation.close:
		return nil, nil, io.EOF
	default:
		if numberOfLoops == 10 {
			numberOfLoops = 0
			animation.widget.name = ""
			return nil, nil, io.EOF
		}

		return animation.ctx.Image(), time.After(time.Millisecond * 10), nil
	}
}

func explosion(animation *ExplosionAnimation) {

	x := 15
	y := 15

	if renderTick == 13 || numberOfLoops == 0 {
		animation.ctx.SetRGB(0, 0, 0)
		animation.ctx.Clear()
		if renderTick > 0 {
			numberOfLoops++
		}

		renderTick = 0
	}

	animation.ctx.SetColor(color.RGBA{255, 0, 0, 255})

	green := 255 - renderTick*14
	color1 := color.RGBA{R: 255, G: uint8(green), B: 0, A: 255}

	animation.ctx.DrawCircle(16, 16, float64(renderTick)/1.5)

	animation.ctx.SetColor(color1)
	animation.ctx.SetPixel(x+renderTick, y+renderTick)
	animation.ctx.SetPixel(x-renderTick, y-renderTick)

	animation.ctx.SetColor(color1)
	animation.ctx.SetPixel(x, y+renderTick)
	animation.ctx.SetPixel(x, y-renderTick)

	animation.ctx.SetColor(color1)
	animation.ctx.SetPixel(x+renderTick, y)
	animation.ctx.SetPixel(x-renderTick, y)

	animation.ctx.SetColor(color1)
	animation.ctx.SetPixel(x+renderTick, y-renderTick)
	animation.ctx.SetPixel(x-renderTick, y+renderTick)

	animation.ctx.Stroke()

	renderTick++
}

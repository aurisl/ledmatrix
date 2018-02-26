package main

import (
	"errors"
	"github.com/fogleman/gg"
	"github.com/mcuadros/go-rpi-rgb-led-matrix"
	"image"
	"image/color"
	"time"
)

type FireAnimation struct {
	ctx    *gg.Context
	close  chan bool
	widget *Widget
}

func DrawFire(toolkit *rgbmatrix.ToolKit, close chan bool, widget *Widget) {

	animation := &FireAnimation{
		ctx:    gg.NewContext(32, 32),
		close:  close,
		widget: widget,
	}

	toolkit.PlayAnimation(animation)
}

func (animation *FireAnimation) Next() (image.Image, <-chan time.Time, error) {

	animation.ctx.SetColor(color.RGBA{0, 0, 0, 255})
	animation.ctx.Clear()

	YRand := 10
	max := 10
	for x := 1; x <= 32; x++ {

		if x > 15 {
			if x%2 == 0 {
				YRand = 10 + max + 32 - x
			}
		} else {
			if x%2 == 0 {
				YRand = 10 + max - x
			}
		}

		for y := YRand; y < 32; y++ {

			if y > 20 {
				animation.ctx.SetColor(color.RGBA{249, 166, 0, 255})
			} else {
				animation.ctx.SetColor(color.RGBA{255, uint8(10 + y), uint8(10 + y), 255})
			}

			animation.ctx.SetPixel(x, y)
		}

	}

	animation.ctx.Stroke()

	select {
	case <-animation.close:
		return nil, nil, errors.New("EOF")
	default:
		return animation.ctx.Image(), time.After(time.Millisecond * 500), nil
	}
}

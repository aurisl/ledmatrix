package fire

import (
	"github.com/aurisl/ledmatrix/command"
	"github.com/aurisl/ledmatrix/matrix"
	"github.com/fogleman/gg"
	"image"
	"image/color"
	"log"
	"time"
)

type animation struct {
	ctx           *gg.Context
	widgetCommand *command.WidgetCommand
}

func Draw(toolkit *matrix.LedToolKit) {

	animation := &animation{
		ctx: gg.NewContext(32, 32),
	}

	err := toolkit.PlayAnimation(animation)
	if err != nil {
		log.Printf("An error occurred while playing wire animation '%s'", err.Error())
	}
}

func (animation *animation) Next() (image.Image, <-chan time.Time, error) {

	animation.ctx.SetColor(color.RGBA{A: 255})
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
				animation.ctx.SetColor(color.RGBA{R: 249, G: 166, A: 255})
			} else {
				animation.ctx.SetColor(color.RGBA{R: 255, G: uint8(10 + y), B: uint8(10 + y), A: 255})
			}

			animation.ctx.SetPixel(x, y)
		}

	}

	animation.ctx.Stroke()

	return animation.ctx.Image(), time.After(time.Millisecond * 500), nil
}

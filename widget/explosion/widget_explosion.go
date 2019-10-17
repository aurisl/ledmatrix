package explosion

import (
	"github.com/aurisl/ledmatrix/matrix"
	"github.com/fogleman/gg"
	"image"
	"image/color"
	"io"
	"log"
	"time"
)

var (
	renderTick    = 0
	numberOfLoops = 0
)

type animation struct {
	ctx *gg.Context
}

func Draw(toolkit *matrix.LedToolKit) {

	animation := &animation{
		ctx: toolkit.Ctx,
	}

	err := toolkit.PlayAnimation(animation)
	if err != nil {
		log.Printf("An error occurred while playing explosion animation '%s'", err.Error())
	}
}

func (animation *animation) Next() (image.Image, <-chan time.Time, error) {

	explosion(animation)

	if numberOfLoops == 10 {
		numberOfLoops = 0
		return nil, nil, io.EOF
	}

	return animation.ctx.Image(), time.After(time.Millisecond * 10), nil
}

func explosion(animation *animation) {

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

package matrix

import (
	"fmt"
	"image"
	"image/color"
	"os"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/imageutil"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

const windowTitle = "RGB led matrix emulator"

var margin = 10

type Emulator struct {
	PixelPitch int
	Gutter     int
	Width      int
	Height     int

	leds []color.Color
	w    screen.Window
	s    screen.Screen

	isReady bool
}

func NewEmulator(w, h, pixelPitch int, autoInit bool) *Emulator {
	e := &Emulator{
		PixelPitch: pixelPitch,
		Gutter:     int(float64(pixelPitch) * 0.66),
		Width:      w,
		Height:     h,
	}

	if autoInit {
		e.Init()
	}

	return e
}

// Init initialize the emulator, creating a new Window and waiting until is
// painted. If something goes wrong the function panics
func (e *Emulator) Init() {
	e.leds = make([]color.Color, e.Width*e.Height)

	driver.Main(e.mainWindowLoop)
}

func (e *Emulator) mainWindowLoop(s screen.Screen) {
	var err error
	e.s = s
	e.w, err = s.NewWindow(&screen.NewWindowOptions{
		Title: windowTitle,
	})

	if err != nil {
		panic(err)
	}

	defer e.w.Release()

	var sz size.Event
	for {
		evn := e.w.NextEvent()
		switch evn := evn.(type) {
		case paint.Event:
			e.drawContext(sz)
			if e.isReady {
				continue
			}

			e.Apply(make([]color.Color, e.Width*e.Height))
			e.isReady = true
		case size.Event:
			sz = evn

		case error:
			fmt.Fprintln(os.Stderr, e)
		}
	}
}

func (e *Emulator) drawContext(sz size.Event) {
	e.updatePixelPitch(sz.Size())
	for _, r := range imageutil.Border(sz.Bounds(), margin) {
		e.w.Fill(r, color.White, screen.Src)
	}

	e.w.Fill(sz.Bounds().Inset(margin), color.Black, screen.Src)
	e.w.Publish()
}

func (e *Emulator) updatePixelPitch(size image.Point) {
	maxLedSizeInX := (size.X - (margin * 2)) / e.Width
	maxLedSizeInY := (size.Y - (margin * 2)) / e.Height

	maxLedSize := maxLedSizeInY
	if maxLedSizeInX < maxLedSizeInY {
		maxLedSize = maxLedSizeInX
	}

	e.PixelPitch = 2 * (maxLedSize / 3.)
	e.Gutter = maxLedSize / 3
}

func (e *Emulator) Geometry() (width, height int) {
	return e.Width, e.Height
}

func (e *Emulator) Apply(leds []color.Color) error {
	defer func() { e.leds = make([]color.Color, e.Height*e.Width) }()

	for col := 0; col < e.Width; col++ {
		for row := 0; row < e.Height; row++ {
			x := col * (e.PixelPitch + e.Gutter)
			y := row * (e.PixelPitch + e.Gutter)

			x += margin * 2
			y += margin * 2

			color := e.At(col + (row * e.Width))
			led := image.Rect(x, y, x+e.PixelPitch, y+e.PixelPitch)

			e.w.Fill(led, color, screen.Over)
		}
	}

	e.w.Publish()
	return nil
}

func (e *Emulator) Render() error {
	return e.Apply(e.leds)
}

func (e *Emulator) At(position int) color.Color {
	if e.leds[position] == nil {
		return color.Black
	}

	return e.leds[position]
}

func (e *Emulator) Set(position int, c color.Color) {
	e.leds[position] = color.RGBAModel.Convert(c)
}

func (e *Emulator) Close() error {
	return nil
}

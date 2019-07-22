// +build software

package matrix

import (
	"encoding/json"
	"image/color"
	"sync"
)

type Emulator struct {
	Width  int
	Height int

	LEDs []color.Color

	renderCallback func(pixelMap []byte)
}

var colorMap = make(map[int]map[int]map[string]uint32, 32)
var mutex = &sync.Mutex{}

func NewEmulator(w, h int, renderCallback func(pixelMap []byte)) *Emulator {
	emulator := &Emulator{
		Width:          w,
		Height:         h,
		renderCallback: renderCallback,
	}

	emulator.LEDs = make([]color.Color, emulator.Width*emulator.Height)

	return emulator
}

func (e *Emulator) Geometry() (width, height int) {
	return e.Width, e.Height
}

func (e *Emulator) Apply(LEDs []color.Color) error {
	defer func() { e.LEDs = make([]color.Color, e.Height*e.Width) }()

	mutex.Lock()
	for col := 0; col < e.Width; col++ {
		colorMap[col] = make(map[int]map[string]uint32, e.Width)
		for row := 0; row < e.Height; row++ {
			pixelColor := e.At(col + (row * e.Width))

			R, G, B, _ := pixelColor.RGBA()
			colorMap[col][row] = make(map[string]uint32, 3)

			colorMap[col][row]["R"] = R >> 8
			colorMap[col][row]["G"] = G >> 8
			colorMap[col][row]["B"] = B >> 8
		}
	}
	mutex.Unlock()

	pixelMap, _ := json.Marshal(colorMap)
	e.renderCallback(pixelMap)

	return nil
}

func (e *Emulator) Render() error {
	return e.Apply(e.LEDs)
}

func (e *Emulator) At(position int) color.Color {
	if e.LEDs[position] == nil {
		return color.Black
	}

	return e.LEDs[position]
}

func (e *Emulator) Set(position int, c color.Color) {
	e.LEDs[position] = color.RGBAModel.Convert(c)
}

func (e *Emulator) Close() error {
	return nil
}

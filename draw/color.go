package draw

import "image/color"

func BlendingPercentageColor(percentage uint8) color.RGBA {

	var red uint8 = 0
	var green uint8 = 0
	var blue uint8 = 0

	if percentage > 50 {
		red = uint8(1 - float64(2*(percentage-50))/100*255)
		green = 255
	} else {
		red = 255
		green = uint8(float64(2*percentage) / 100 * 255)
	}

	return color.RGBA{R: red, G: green, B: blue, A: 255}
}

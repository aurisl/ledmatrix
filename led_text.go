package main

import (
	"strings"
	"github.com/fogleman/gg"
	"image/color"
)

var characterPixelSize = map[string]int{
	"1": 3,
	"2": 3,
	"3": 3,
	"4": 3,
	"5": 3,
	"6": 3,
	"7": 3,
	"8": 3,
	"9": 3,
	"A": 7,
	"B": 7,
	"C": 7,
	"D": 7,
	"E": 7,
	"F": 7,
	"G": 7,
	"H": 7,
	"L": 7,
	"I": 7,
	"M": 7,
	"N": 7,
	"O": 7,
	"P": 7,
	"Q": 7,
	"R": 7,
	"S": 7,
	"T": 7,
	"U": 7,
	"V": 7,
	"W": 7,
	"X": 7,
	"Y": 7,
	"Z": 7,
	"-": 7,
	"+": 7,
	"=": 3,
	"/": 3,
	"\\": 3,
	":": 1,
	".": 1,
	",": 1,
	";": 1,
	"&": 4,
	"$": 4,
	"!": 1,
	"?": 3,
	"\"": 1,
	"'": 1,
	"*": 1,
	"@": 5,
	"^": 3,
	"~": 3,
	"[": 2,
	"]": 2,
	"{": 2,
	"}": 2,
	"<": 3,
	">": 3,
	"#": 3,
	"_": 3,
	"%": 3,
}

var textPosition = 35.0 //35 is out of bounds

func DrawText(text string, y float64, ctx *gg.Context) bool {

	text = strings.ToUpper(text)

	if err := ctx.LoadFontFace("resources/fonts/PixelOperator.ttf", 16); err != nil {
		panic(err)
	}

	totalPixels := CountWordPixels(text)

	ctx.SetRGB(0, 0, 0)
	ctx.Clear()
	ctx.SetColor(color.RGBA{255, 0, 0, 255})

	ctx.DrawString(text, textPosition, y)

	textPosition--
	if int(textPosition) <= -totalPixels {
		textPosition = 35.0
		return true
	}
	return false
}


func CountWordPixels(text string) int {

	var totalPixels = 0
	for _, char := range text {
		charValue := string(char)

		if size, ok := characterPixelSize[charValue]; ok {
			totalPixels += size
		} else {
			totalPixels += 3
		}
	}
	return totalPixels
}

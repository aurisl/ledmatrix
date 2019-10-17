// +build hardware

package matrix

import (
	"github.com/mcuadros/go-rpi-rgb-led-matrix"
	"github.com/aurisl/ledmatrix/error"
)

var Renderer = matrix.CreateHardwareMatrix()

func CreateHardwareMatrix() Matrix {
	matrixConfig := CreateMatrixDefaultConfiguration()

	matrix, err := rgbmatrix.NewRGBLedMatrix(matrixConfig)
	error.Fatal(err)

	return matrix
}

func CreateMatrixDefaultConfiguration() *rgbmatrix.HardwareConfig {

	matrixConfig := &rgbmatrix.DefaultConfig
	matrixConfig.Rows = 32
	matrixConfig.Brightness = 50

	return matrixConfig
}

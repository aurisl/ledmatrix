package matrix

import (
	"github.com/mcuadros/go-rpi-rgb-led-matrix"
	"github.com/aurisl/ledmatrix/error"
)

func initializeMatrixToolkit() *ToolKit {
	matrixConfig := createMatrixDefaultConfiguration()

	matrix, err := rgbmatrix.NewRGBLedMatrix(matrixConfig)
	error.Fatal(err)

	toolkit := NewToolKit(matrix)

	return toolkit
}

func createMatrixDefaultConfiguration() *rgbmatrix.HardwareConfig {

	matrixConfig := &rgbmatrix.DefaultConfig
	matrixConfig.Rows = 32
	matrixConfig.Brightness = 50

	return matrixConfig
}

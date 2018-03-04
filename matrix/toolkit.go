package matrix

import (
	"github.com/mcuadros/go-rpi-rgb-led-matrix"
	"github.com/fogleman/gg"
	"github.com/aurisl/ledmatrix/error"
)

var matrixToolkit *rgbmatrix.ToolKit

type (
	Toolkit struct {
		MatrixToolkit *rgbmatrix.ToolKit
		Ctx *gg.Context
	}
)

func LoadLedMatrixToolkit() *Toolkit {

	return &Toolkit{
		MatrixToolkit: initializeMatrixToolkit(),
		Ctx:           gg.NewContext(32, 32),
	}
}

func (toolkit *Toolkit) Close() {
	matrixToolkit.Close()
}

func initializeMatrixToolkit() *rgbmatrix.ToolKit  {
	matrixConfig := createMatrixDefaultConfiguration()

	matrix, err := rgbmatrix.NewRGBLedMatrix(matrixConfig)
	error.Fatal(err)

	toolkit := rgbmatrix.NewToolKit(matrix)

	return toolkit

}

func createMatrixDefaultConfiguration() *rgbmatrix.HardwareConfig {
	matrixConfig := &rgbmatrix.DefaultConfig
	matrixConfig.Rows = 32
	matrixConfig.Brightness = 50

	return matrixConfig
}
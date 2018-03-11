package matrix

import (
	"github.com/mcuadros/go-rpi-rgb-led-matrix"
	"github.com/fogleman/gg"
	merror "github.com/aurisl/ledmatrix/error"
	"time"
	"io"
	"image"
	"github.com/aurisl/ledmatrix/command"
)

var matrixToolkit *rgbmatrix.ToolKit

type (
	LedToolKit struct {
		MatrixToolkit *rgbmatrix.ToolKit
		Ctx           *gg.Context
		commandInput  <-chan command.WidgetCommand
		commandOutput chan command.WidgetCommand
	}
)

func NewLedToolkit(commandInput <-chan command.WidgetCommand, commandOutput chan command.WidgetCommand) *LedToolKit {

	return &LedToolKit{
		MatrixToolkit: initializeMatrixToolkit(),
		Ctx:           gg.NewContext(32, 32),
		commandInput:  commandInput,
		commandOutput: commandOutput,
	}
}

func (toolKit *LedToolKit) PlayAnimation(a rgbmatrix.Animation) error {

	var err error
	var i image.Image
	var n <-chan time.Time

	for {
		select {
		case in := <-toolKit.commandInput:
			toolKit.commandOutput <- in
			return nil
		default:

			i, n, err = a.Next()
			if err != nil && err == io.EOF {
				return nil
			}

			if err := toolKit.MatrixToolkit.PlayImageUntil(i, n); err != nil {
				return err
			}
		}

	}

	return err
}

func (toolKit *LedToolKit) Close() {
	matrixToolkit.Close()
}

func initializeMatrixToolkit() *rgbmatrix.ToolKit {
	matrixConfig := createMatrixDefaultConfiguration()

	matrix, err := rgbmatrix.NewRGBLedMatrix(matrixConfig)
	merror.Fatal(err)

	toolkit := rgbmatrix.NewToolKit(matrix)

	return toolkit

}

func createMatrixDefaultConfiguration() *rgbmatrix.HardwareConfig {
	matrixConfig := &rgbmatrix.DefaultConfig
	matrixConfig.Rows = 32
	matrixConfig.Brightness = 50

	return matrixConfig
}
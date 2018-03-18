package matrix

import (
	"github.com/fogleman/gg"
	"time"
	"io"
	"image"
	"github.com/aurisl/ledmatrix/command"
)

var matrixToolkit *ToolKit

type (
	LedToolKit struct {
		MatrixToolkit *ToolKit
		Ctx           *gg.Context
		commandInput  <-chan command.WidgetCommand
		commandOutput chan command.WidgetCommand
	}
)

func NewLedToolkit(commandInput <-chan command.WidgetCommand, commandOutput chan command.WidgetCommand, m Matrix) *LedToolKit {

	toolkit := NewToolKit(m)

	return &LedToolKit{
		MatrixToolkit: toolkit,
		Ctx:           gg.NewContext(32, 32),
		commandInput:  commandInput,
		commandOutput: commandOutput,
	}
}

func (toolKit *LedToolKit) PlayAnimation(a Animation) error {

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
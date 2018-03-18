package main

import (
	"github.com/aurisl/ledmatrix/command"
	"github.com/aurisl/ledmatrix/widget"
	"github.com/aurisl/ledmatrix/matrix"
)

func main() {

	commandInput := make(chan command.WidgetCommand)

	command.Start(commandInput)

	onCanvasRender := command.StartFeed()

	emulatedMatrix := matrix.NewEmulator(32, 32, onCanvasRender)

	widget.Start(commandInput, emulatedMatrix)
}

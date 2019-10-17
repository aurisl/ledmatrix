package main

import (
	"github.com/aurisl/ledmatrix/command"
	"github.com/aurisl/ledmatrix/config"
	"github.com/aurisl/ledmatrix/matrix"
	"github.com/aurisl/ledmatrix/widget"
)

func main() {

	config.InitializeConfiguration()

	commandInput := make(chan command.WidgetCommand)

	command.Start(commandInput)

	//Hardware matrix
	MMatrix := matrix.CreateHardwareMatrix()

	//Software emulated matrix (display in browser)
	//onCanvasRender := command.StartFeed()
	//MMatrix := matrix.NewEmulator(32, 32, onCanvasRender)

	widget.Start(commandInput, MMatrix)
}

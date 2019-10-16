package main

import (
	"github.com/aurisl/ledmatrix/command"
	"github.com/aurisl/ledmatrix/config"
	"github.com/aurisl/ledmatrix/matrix"
	"github.com/aurisl/ledmatrix/widget"
	"github.com/aurisl/ledmatrix/widget/meter"
)

func main() {

	config.InitializeConfiguration()

	commandInput := make(chan command.WidgetCommand)

	command.Start(commandInput)

	MMatrix := matrix.CreateHardwareMatrix()

	//onCanvasRender := command.StartFeed()
	//MMatrix := matrix.NewEmulator(32, 32, onCanvasRender)

	widget.Start(commandInput, MMatrix)

	go meter.Measure()
}

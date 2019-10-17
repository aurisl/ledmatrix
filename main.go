package main

import (
	"github.com/aurisl/ledmatrix/command"
	"github.com/aurisl/ledmatrix/config"
	"github.com/aurisl/ledmatrix/matrix"
	"github.com/aurisl/ledmatrix/widget"
	"log"
)

func main() {

	config.InitializeConfiguration()

	commandInput := make(chan command.WidgetCommand)

	command.Start(commandInput)

	//Hardware matrix
	//MMatrix := matrix.CreateHardwareMatrix()

	//Software emulated matrix (display in browser)
	MMatrix := matrix.NewEmulator(32, 32, command.StartFeed())

	log.Println("Matrix renderer started")
	widget.Start(commandInput, MMatrix)
}

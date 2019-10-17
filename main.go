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

	log.Println("Matrix renderer started")
	widget.Start(commandInput, matrix.Renderer)
}

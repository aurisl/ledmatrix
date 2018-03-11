package main

import (
	"github.com/aurisl/ledmatrix/command"
	"github.com/aurisl/ledmatrix/widget"
)

func main() {

	commandInput := make(chan command.WidgetCommand)

	command.Start(commandInput)
	widget.Start(commandInput)
}

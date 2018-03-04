package main

import (
	"github.com/aurisl/ledmatrix/command"
	"github.com/aurisl/ledmatrix/widget"
)

func main() {

	widgetCommand := &command.WidgetCommand{}

	terminate := make(chan bool)

	command.StartHttpServer(widgetCommand, terminate)
	widget.StartLoader(widgetCommand, terminate)
}

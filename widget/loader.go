package widget

import (
	"github.com/aurisl/ledmatrix/command"
	"github.com/aurisl/ledmatrix/matrix"
	"github.com/aurisl/ledmatrix/widget/explosion"
	"github.com/aurisl/ledmatrix/widget/gif"
	"github.com/aurisl/ledmatrix/widget/location"
	"github.com/aurisl/ledmatrix/widget/meter"
	"github.com/aurisl/ledmatrix/widget/torrent"
	"github.com/aurisl/ledmatrix/widget/weather"
)

func Start(commandInput <-chan command.WidgetCommand, m matrix.Matrix) {

	commandOutput := make(chan command.WidgetCommand, 1)

	ledToolKit := matrix.NewLedToolkit(commandInput, commandOutput, m)
	defer ledToolKit.Close()

	//Start CO2 meter in background and read data from device
	go meter.Measure()

	//The main drawing loop
	for {

		select {
		case widgetCommand := <-commandOutput:

			switch widgetCommand.Name {
			case "weather":
				weather.Draw(ledToolKit)
			case "boom":
				explosion.Draw(ledToolKit)
			case "gif":
				gif.Draw(ledToolKit, widgetCommand.Params)
			case "location":
				location.Draw(ledToolKit)
			case "torrent":
				torrent.Draw(ledToolKit)
			case "meter":
				meter.Draw(ledToolKit)
			default:
				weather.Draw(ledToolKit)
			}
		default:
			weather.Draw(ledToolKit)
		}

	}

}

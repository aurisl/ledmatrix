package widget

import (
	"github.com/aurisl/ledmatrix/command"
	"github.com/aurisl/ledmatrix/config"
	"github.com/aurisl/ledmatrix/matrix"
	"github.com/aurisl/ledmatrix/widget/explosion"
	"github.com/aurisl/ledmatrix/widget/fire"
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

	go meter.Measure()

	for {

		select {
		case widgetCommand := <-commandOutput:

			switch widgetCommand.Name {
			case "weather":
				weather.Draw(ledToolKit, config.App)
			case "boom":
				explosion.Draw(ledToolKit)
			case "gif":
				gif.Draw(ledToolKit, widgetCommand.Params)
			case "fire":
				fire.Draw(ledToolKit)
			case "location":
				location.Draw(ledToolKit, config.App)
			case "torrent":
				torrent.Draw(ledToolKit, config.App)
			case "meter":
				meter.Draw(ledToolKit)
			default:
				weather.Draw(ledToolKit, config.App)
			}
		default:
			weather.Draw(ledToolKit, config.App)
		}

	}

}

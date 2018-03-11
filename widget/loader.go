package widget

import (
	"github.com/aurisl/ledmatrix/command"
	"github.com/aurisl/ledmatrix/widget/explosion"
	"github.com/aurisl/ledmatrix/widget/weather"
	"github.com/aurisl/ledmatrix/widget/image"
	"github.com/aurisl/ledmatrix/widget/fire"
	"github.com/aurisl/ledmatrix/widget/location"
	"github.com/aurisl/ledmatrix/widget/torrent"
	"github.com/aurisl/ledmatrix/matrix"
	"github.com/aurisl/ledmatrix/config"
)

func Start(commandInput <-chan command.WidgetCommand) {

	commandOutput := make(chan command.WidgetCommand, 1)

	ledToolKit := matrix.NewLedToolkit(commandInput, commandOutput)
	defer ledToolKit.Close()

	appConfig := config.NewAppConfig()

	for {

		select {
		case widgetCommand := <-commandOutput:

			switch widgetCommand.Name {
			case "weather":
				weather.Draw(ledToolKit, appConfig)
			case "boom":
				explosion.Draw(ledToolKit)
			case "gif":
				image.Draw(ledToolKit, widgetCommand.Params)
			case "fire":
				fire.Draw(ledToolKit)
			case "location":
				location.Draw(ledToolKit, appConfig)
			case "torrent":
				torrent.Draw(ledToolKit, appConfig)
			default:
				weather.Draw(ledToolKit, appConfig)
			}
		default:
			weather.Draw(ledToolKit, appConfig)
		}

	}

}

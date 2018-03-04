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

func StartLoader(widgetCommand *command.WidgetCommand, terminate chan bool)  {

	ledToolkit := matrix.LoadLedMatrixToolkit()
	defer ledToolkit.Close()

	widgetConfig := config.NewConfig()

	for {

		switch widgetCommand.Name {
		case "weather":
			weather.Draw(ledToolkit, terminate, widgetConfig.WidgetWeatherApiConfig)
		case "boom":
			explosion.Draw(ledToolkit, terminate, widgetCommand)
		case "gif":
			image.Draw(ledToolkit, terminate, widgetCommand)
		case "fire":
			fire.Draw(ledToolkit, terminate, widgetCommand)
		case "location":
			location.Draw(ledToolkit, terminate, widgetCommand, widgetConfig.WidgetLocationConfig)
		case "torrent":
			torrent.Draw(ledToolkit, terminate, widgetCommand, widgetConfig.WidgetTorrentStatusConfig)
		default:
			weather.Draw(ledToolkit, terminate, widgetConfig.WidgetWeatherApiConfig)
		}
	}

}


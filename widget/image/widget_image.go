package image

import (
	"net/http"
	"strconv"
	"time"
	"github.com/aurisl/ledmatrix/command"
	"github.com/aurisl/ledmatrix/matrix"
)

func Draw(toolkit *matrix.Toolkit, close chan bool, WidgetCommand *command.WidgetCommand) {

	imgUrl := WidgetCommand.Params.Get("url")
	response, err := http.Get(imgUrl)
	if err != nil {
		WidgetCommand.Name = ""
		return
	}

	closed, err := toolkit.MatrixToolkit.PlayGIF(response.Body)
	if err != nil {
		WidgetCommand.Name = ""
		return
	}

	animationDuration := WidgetCommand.Params.Get("duration")
	if animationDuration != "" {
		if animationDuration == "-1" {
			select {
			case <-close:
				WidgetCommand.Name = ""
			}
		}

		d, _ := strconv.ParseInt(animationDuration, 10, 32)
		time.Sleep(time.Second * time.Duration(d))
	} else {
		time.Sleep(time.Second * 10)
	}

	WidgetCommand.Name = ""

	closed <- true
}

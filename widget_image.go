package main

import (
	"github.com/mcuadros/go-rpi-rgb-led-matrix"
	"net/http"
	"strconv"
	"time"
)

func DrawGif(toolkit *rgbmatrix.ToolKit, widget *Widget, close chan bool) {

	imgUrl := widget.params.Get("url")
	response, err := http.Get(imgUrl)
	if err != nil {
		widget.name = ""
		return
	}

	closed, err := toolkit.PlayGIF(response.Body)
	if err != nil {
		widget.name = ""
		return
	}

	animationDuration := widget.params.Get("duration")
	if animationDuration != "" {
		if animationDuration == "-1" {
			select {
			case <-close:
				widget.name = ""
			}
		}

		d, _ := strconv.ParseInt(animationDuration, 10, 32)
		time.Sleep(time.Second * time.Duration(d))
	} else {
		time.Sleep(time.Second * 10)
	}

	widget.name = ""

	closed <- true
}

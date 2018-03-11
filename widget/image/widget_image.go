package image

import (
	"net/http"
	"strconv"
	"time"
	"github.com/aurisl/ledmatrix/matrix"
	"net/url"
)

func Draw(toolkit *matrix.LedToolKit, params url.Values) {

	imgUrl := params.Get("url")
	response, err := http.Get(imgUrl)
	if err != nil {
		return
	}

	closed, err := toolkit.MatrixToolkit.PlayGIF(response.Body)
	if err != nil {
		return
	}

	animationDuration := params.Get("duration")
	if animationDuration != "" {
		if animationDuration == "-1" {
			return
		}

		d, _ := strconv.ParseInt(animationDuration, 10, 32)
		time.Sleep(time.Second * time.Duration(d))
	} else {
		time.Sleep(time.Second * 10)
	}

	closed <- true
}

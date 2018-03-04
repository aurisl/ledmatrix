package command

import (
	"fmt"
	"html"
	"net/http"
	"net/url"
	"github.com/aurisl/ledmatrix/error"
)

type (
	WidgetCommand struct {
		Name   string
		Params url.Values
	}
)

func StartHttpServer(widget *WidgetCommand, terminateCurrentRunning chan bool) {

	go func() {
		http.HandleFunc("/exec", func(w http.ResponseWriter, r *http.Request) {
			widgetName := html.EscapeString(r.URL.Query().Get("widget"))
			fmt.Fprintf(w, "Command '%q' received. ", widgetName)
			widget.Name = widgetName
			widget.Params = r.URL.Query()

			terminateCurrentRunning <- true
		})

		err := http.ListenAndServe(":8080", nil)
		error.Fatal(err)
	}()
}


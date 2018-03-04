package main

import (
	"fmt"
	"html"
	"net/http"
)

func StartHttpServer(widget *Widget, terminateCurrentRunning chan bool) {

	http.HandleFunc("/exec", func(w http.ResponseWriter, r *http.Request) {
		widgetName := html.EscapeString(r.URL.Query().Get("widget"))
		fmt.Fprintf(w, "Command '%q' received. ", widgetName)
		widget.name = widgetName
		widget.params = r.URL.Query()

		terminateCurrentRunning <- true
	})

	err := http.ListenAndServe(":8080", nil)
	fatal(err)
}

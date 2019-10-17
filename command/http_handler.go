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

func Start(commandInput chan WidgetCommand) {

	go func() {

		http.HandleFunc("/exec", func(w http.ResponseWriter, r *http.Request) {

			widgetName := html.EscapeString(r.URL.Query().Get("widget"))
			fmt.Fprintf(w, "Executing '%q' widget. ", widgetName)

			commandInput <- WidgetCommand{
				Name:   widgetName,
				Params: r.URL.Query(),
			}

			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", "*")
		})

		fs := http.FileServer(http.Dir("public"))
		http.Handle("/", http.StripPrefix("/", fs))

		err := http.ListenAndServe(":8081", nil)
		error.Fatal(err)
	}()
}

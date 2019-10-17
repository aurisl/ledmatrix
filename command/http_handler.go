package command

import (
	"fmt"
	"github.com/aurisl/ledmatrix/config"
	"github.com/aurisl/ledmatrix/error"
	"html"
	"log"
	"net/http"
	"net/url"
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
			fmt.Fprintf(w, "Executing %q widget. ", widgetName)

			commandInput <- WidgetCommand{
				Name:   widgetName,
				Params: r.URL.Query(),
			}

			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", "*")
		})

		fs := http.FileServer(http.Dir(*config.WorkingDir + "public"))
		http.Handle("/", http.StripPrefix("/", fs))

		log.Printf("API server started at 8081")
		err := http.ListenAndServe(":8081", nil)

		error.Fatal(err)
	}()
}

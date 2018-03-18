package command

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"time"
)

var upgrade = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
	return true
}}

var lastMatrix2 []byte

func feedPixels(w http.ResponseWriter, r *http.Request)  {
	connection, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer connection.Close()

	for {
		connection.WriteMessage(websocket.TextMessage, lastMatrix2)
		time.Sleep(time.Millisecond * 10)
	}
}

func StartFeed() func(pixelMap []byte) {

	onCanvasRender := func(pixelMap []byte) {
		lastMatrix2 = pixelMap
	}

	go func() {
		http.HandleFunc("/pixel", feedPixels)
		log.Fatal(http.ListenAndServe(":8081", nil))
	}()

	return onCanvasRender
}
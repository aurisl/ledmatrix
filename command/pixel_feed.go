package command

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var upgrade = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
	return true
}}

var currentPixelMap []byte
var pushPixelUpdate bool

func feedPixels(w http.ResponseWriter, r *http.Request) {
	connection, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer connection.Close()

	for {
		if pushPixelUpdate == true {
			connection.WriteMessage(websocket.TextMessage, currentPixelMap)
			pushPixelUpdate = false
		}
		time.Sleep(time.Millisecond * 10)
	}
}

func StartFeed() func(pixelMap []byte) {

	onCanvasRender := func(pixelMap []byte) {
		pushPixelUpdate = isUpdated(pixelMap)
		currentPixelMap = pixelMap
	}

	go func() {
		http.HandleFunc("/pixel", feedPixels)
		log.Fatal(http.ListenAndServe(":8082", nil))
	}()

	return onCanvasRender
}

func isUpdated(pixelMap []byte) bool {

	if string(pixelMap) != string(currentPixelMap) {
		return true
	}

	return false
}

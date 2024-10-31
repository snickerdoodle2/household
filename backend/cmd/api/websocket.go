package main

import (
	"net/http"

	"github.com/coder/websocket"
)

func (app *App) upgradeSensorWebsocket(w http.ResponseWriter, r *http.Request) {
	_, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		Subprotocols: []string{"inzynierka"}, // TODO: better name for subprotocol
	})

	if err != nil {
		app.logger.Error("open websocket", "error", err)
		return
	}
}

package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/google/uuid"
)

type connStatus struct {
	mu           sync.Mutex
	auth         bool
	subscribedTo []uuid.UUID
}

func (app *App) upgradeSensorWebsocket(w http.ResponseWriter, r *http.Request) {
	const subprotocol = "inzynierka"
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		Subprotocols: []string{subprotocol}, // TODO: better name for subprotocol
	})

	if err != nil {
		app.logger.Error("open websocket", "error", err)
		return
	}

	defer conn.CloseNow()

	app.logger.Debug("new connection", "subprotocol", conn.Subprotocol())

	if conn.Subprotocol() != subprotocol {
		app.logger.Debug("closing connection", "reason", "wrong subprotocol")
		conn.Close(websocket.StatusPolicyViolation, fmt.Sprintf("client must speak %s subprotocol", subprotocol))
		return
	}

	connStatus := &connStatus{
		auth:         false,
		subscribedTo: make([]uuid.UUID, 0),
	}

	for {
		err = app.handleWebSocketMessage(conn, connStatus)
		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			return
		}
		if err != nil {
			app.logger.Error("unhandled ws error", "error", err)
			return
		}
	}
}

func (app *App) handleWebSocketMessage(conn *websocket.Conn, status *connStatus) error {
	var msg interface{}
	err := wsjson.Read(context.Background(), conn, &msg)
	if err != nil {
		return err
	}
	app.logger.Debug("new message", "msg", msg)

	if !status.auth {
		res := map[string]string{"type": "auth", "message": "NO_AUTH"}
		wsjson.Write(context.Background(), conn, res)
	}
	return nil
}

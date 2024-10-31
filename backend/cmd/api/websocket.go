package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"inzynierka/internal/data"
	"inzynierka/internal/data/validator"
	"net/http"
	"sync"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/google/uuid"
)

type connStatus struct {
	mu           sync.Mutex
	authed       bool
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
		authed:       false,
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

type messageType string

const (
	authMsg     messageType = "auth"
	serverError messageType = "server_error"
)

type websocketMsg struct {
	Type messageType     `json:"type"`
	Data json.RawMessage `json:"data"`
}

func (app *App) handleWebSocketMessage(conn *websocket.Conn, status *connStatus) error {
	var msg websocketMsg

	err := wsjson.Read(context.Background(), conn, &msg)
	if err != nil {
		return err
	}
	app.logger.Debug("new message", "msg", msg)

	if msg.Type == authMsg {
		return app.handleAuthMsg(conn, status, msg.Data)
	}

	if !status.authed {
		res := map[string]string{"type": string(authMsg), "message": "NO_AUTH"}
		return wsjson.Write(context.Background(), conn, res)
	}
	return nil
}

func (app *App) handleAuthMsg(conn *websocket.Conn, status *connStatus, input json.RawMessage) error {
	var msgData struct {
		Token string `json:"token"`
	}

	err := json.Unmarshal(input, &msgData)
	if err != nil {
		return err
	}

	v := validator.New()
	if data.ValidateTokenPlaintext(v, msgData.Token); !v.Valid() {
		return invalidTokenResponse(conn)
	}

	_, err = app.models.Users.GetForToken(msgData.Token)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			return serverErrorResponse(conn)
		default:
			return serverErrorResponse(conn)
		}
	}

	status.authed = true
	return authResponse(conn, "ok")
}

func invalidTokenResponse(conn *websocket.Conn) error {
	return authResponse(conn, "INVALID_TOKEN")
}

func authResponse(conn *websocket.Conn, msg string) error {
	res := map[string]string{"type": string(authMsg), "message": msg}
	return wsjson.Write(context.Background(), conn, res)
}

func serverErrorResponse(conn *websocket.Conn) error {
	res := map[string]string{"type": string(serverError)}
	return wsjson.Write(context.Background(), conn, res)
}

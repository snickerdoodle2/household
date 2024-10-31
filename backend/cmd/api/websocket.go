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
	mu     sync.Mutex
	authed bool
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
		authed: false,
	}

	for {
		err = app.handleWebSocketMessage(conn, connStatus)
		if err != nil {
			switch {
			case websocket.CloseStatus(err) == websocket.StatusNormalClosure:
				return
			default:
				app.logger.Error("unhandled ws error", "error", err)
				return
			}
		}
	}
}

type messageType string

const (
	authMsg      messageType = "auth"
	serverError  messageType = "server_error"
	subscribeMsg messageType = "subscribe"
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

	status.mu.Lock()
	defer status.mu.Unlock()

	if !status.authed {
		return authResponse(conn, "NO_AUTH")
	}

	if msg.Type == subscribeMsg {
		return app.handleSubscribeMsg(conn, msg.Data)
	}

	return nil
}

func (app *App) handleAuthMsg(conn *websocket.Conn, status *connStatus, input json.RawMessage) error {
	var token string

	err := json.Unmarshal(input, &token)
	if err != nil {
		return err
	}

	v := validator.New()
	if data.ValidateTokenPlaintext(v, token); !v.Valid() {
		return invalidTokenResponse(conn)
	}

	_, err = app.models.Users.GetForToken(token)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			return serverErrorResponse(conn)
		default:
			return serverErrorResponse(conn)
		}
	}

	status.mu.Lock()
	defer status.mu.Unlock()
	status.authed = true
	return authResponse(conn, "ok")
}

func (app *App) handleSubscribeMsg(conn *websocket.Conn, input json.RawMessage) error {
	var sensorIDs []json.RawMessage

	err := json.Unmarshal(input, &sensorIDs)
	if err != nil {
		return err
	}

	data := make(map[string]interface{})

	for _, sensorIDraw := range sensorIDs {
		var sensorID uuid.UUID
		err = json.Unmarshal(sensorIDraw, &sensorID)
		if err != nil {
			var tmp string
			err = json.Unmarshal(sensorIDraw, &tmp)
			if err != nil {
				return err
			}
			data[tmp] = sensorErrorMsg("INVALID_UUID")
			continue
		}
		// TODO: handle error
		tmp, _ := handleSensorSubcribtion(sensorID)
		data[sensorID.String()] = tmp
	}

	res := map[string]interface{}{"type": subscribeMsg, "data": data}

	return wsjson.Write(context.Background(), conn, res)
}

func handleSensorSubcribtion(id uuid.UUID) (map[string]interface{}, error) {
	return map[string]interface{}{"foo": "bar"}, nil
}

func sensorErrorMsg(msg string) map[string]string {
	return map[string]string{"status": "error", "message": msg}
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

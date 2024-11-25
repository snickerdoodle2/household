package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"inzynierka/internal/data"
	"inzynierka/internal/data/validator"
	"net/http"
	"reflect"
	"slices"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/google/uuid"
)

type connStatus struct {
	mu     sync.Mutex
	authed bool
	ch     wsChan
}

func (app *App) upgradeSensorWebsocket(w http.ResponseWriter, r *http.Request) {
	const subprotocol = "inzynierka"
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		Subprotocols:   []string{subprotocol}, // TODO: better name for subprotocol
		OriginPatterns: app.config.cors.trustedOrigins,
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
		ch:     make(wsChan, 32),
	}

	go app.sendSensorUpdates(conn, connStatus)

	for {
		err = app.handleWebSocketMessage(conn, connStatus)
		if err != nil {
			connStatus.ch <- wsMsg{
				action: actionClose,
			}

			status := websocket.CloseStatus(err)
			switch status {
			case websocket.StatusNormalClosure, websocket.StatusGoingAway:
				return
			default:
				app.logger.Error("unhandled ws error", "error", err)
				return
			}
		}
	}
}

type wsAction string

const (
	actionClose       wsAction = "CLOSE"
	actionSubscribe   wsAction = "SUBSCRIBE"
	actionUnsubscribe wsAction = "UNSUBSCRIBE"
)

type wsMsg struct {
	action wsAction
	id     uuid.UUID
}

type wsChan chan wsMsg

func (app *App) sendSensorUpdates(conn *websocket.Conn, status *connStatus) {
	// TODO: brokers sending only current value?
	type wsListener struct {
		id    uuid.UUID
		msgCh chan []float64
	}

	for {
		// wait for being authed
		status.mu.Lock()
		if status.authed {
			status.mu.Unlock()
			break
		}
		status.mu.Unlock()
	}

	listeners := make([]wsListener, 0)

	defer (func() {
		for _, tmp := range listeners {
			if listener, ok := app.listeners[tmp.id]; ok {
				listener.Broker.Unsubscribe(tmp.msgCh)
				app.logger.Debug("sendSensorUpdates", "action", "cleanup", "sensorID", tmp.id)
			} else {
				app.logger.Debug("sendSensorUpdates", "action", "cleanup", "note", "listener already closed", "sensorID", tmp.id)
			}
		}
	})()

	defer app.logger.Debug("sendSensorUpdates", "action", "closing")
	channels := make([]reflect.SelectCase, 1)
	channels[0] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(status.ch)}

	for {
		i, msg, ok := reflect.Select(channels)
		if !ok {
			app.logger.Debug("sendSensorUpdates", "error", "read from channel", "channel idx", i)
		}
		if i == 0 {
			action := msg.Interface().(wsMsg)

			switch action.action {
			case actionClose:
				return
			case actionSubscribe:
				app.logger.Debug("sendSensorUpdates", "action", "subscribe", "sensorID", action.id)

				// go to next message if already subscribed
				if slices.IndexFunc(listeners, func(e wsListener) bool { return e.id == action.id }) != -1 {
					app.logger.Debug("sendSensorUpdates", "action", "subscribe", "sensorID", action.id, "error", "already subscribed")
					continue
				}

				listener, ok := app.listeners[action.id] // should be in listeners
				if !ok {
					app.logger.Error("sendSensorUpdates", "action", "subscribe", "sensorID", action.id, "error", "listener not found")
					continue
				}
				msgCh := listener.Broker.Subscribe()
				listeners = append(listeners, wsListener{id: action.id, msgCh: msgCh})
				channels = append(channels, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(msgCh)})
			case actionUnsubscribe:
				app.logger.Debug("sendSensorUpdates", "action", "unsubscribe", "sensorID", action.id)

				idx := slices.IndexFunc(listeners, func(e wsListener) bool { return e.id == action.id })
				if idx == -1 {
					app.logger.Debug("sendSensorUpdates", "action", "unsubscribe", "sensorID", action.id, "error", "not subscribed")
					continue
				}

				app.listeners[listeners[i].id].Broker.Unsubscribe(listeners[i].msgCh)
				listeners = slices.Delete(listeners, i, i+1)
				channels = slices.Delete(channels, i+1, i+2)

			default:
				app.logger.Debug("sendSensorUpdates", "action", action.action, "error", "unhandled")
			}
			continue
		}

		// message fron sensor listener
		values := msg.Interface().([]float64)
		err := sendSensorUpdate(conn, listeners[i-1].id, values[len(values)-1])
		if err != nil {
			app.logger.Error("sendSensorUpdates", "action", "update", "error", err)
		}
	}
}

func sendSensorUpdate(conn *websocket.Conn, id uuid.UUID, value float64) error {
	type Msg struct {
		Type     messageType `json:"type"`
		SensorId uuid.UUID   `json:"sensor_id"`
		Time     time.Time   `json:"time"`
		Value    float64     `json:"value"`
	}

	msg := Msg{
		Type:     measurementMsg,
		SensorId: id,
		Time:     time.Now(),
		Value:    value,
	}

	return wsjson.Write(context.Background(), conn, msg)
}

type messageType string

const (
	authMsg        messageType = "auth"
	serverError    messageType = "server_error"
	subscribeMsg   messageType = "subscribe"
	unsubscribeMsg messageType = "unsubscribe"
	measurementMsg messageType = "measurment"
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
		return app.handleSubscribeMsg(conn, status, msg.Data)
	}

	if msg.Type == unsubscribeMsg {
		return app.handleUnsubscribeMsg(conn, status, msg.Data)
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

func (app *App) handleSubscribeMsg(conn *websocket.Conn, status *connStatus, input json.RawMessage) error {
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

		values, err := app.handleSensorSubcribtion(sensorID)
		if err != nil {
			app.logger.Error("handleSensorSubcribtion: fetch measurements data from db", "error", err)
			data[sensorID.String()] = sensorErrorMsg("SERVER_ERROR")
			continue
		}

		status.ch <- wsMsg{
			action: actionSubscribe,
			id:     sensorID,
		}

		data[sensorID.String()] = values
	}

	res := map[string]interface{}{"type": subscribeMsg, "data": data}

	return wsjson.Write(context.Background(), conn, res)
}

func (app *App) handleSensorSubcribtion(id uuid.UUID) (map[string]interface{}, error) {
	measurements, err := app.models.SensorMeasurements.GetMeasurementsSince(id, 5*time.Hour)
	if err != nil {
		return nil, err
	}

	values := make(map[string]float64)
	for _, measurement := range measurements {
		values[measurement.MeasuredAt.Format(time.RFC3339)] = measurement.MeasuredValue
	}

	return map[string]interface{}{"status": "ok", "values": values}, nil
}

func (app *App) handleUnsubscribeMsg(conn *websocket.Conn, status *connStatus, input json.RawMessage) error {
	var sensorId uuid.UUID
	err := json.Unmarshal(input, &sensorId)
	if err != nil {
		var tmp string
		err = json.Unmarshal(input, &tmp)
		if err != nil {
			return err
		}
		res := map[string]interface{}{"type": unsubscribeMsg, "error": "INVALID_UUID"}
		return wsjson.Write(context.Background(), conn, res)
	}

	status.ch <- wsMsg{
		action: actionUnsubscribe,
		id:     sensorId,
	}

	res := map[string]interface{}{"type": unsubscribeMsg, "data": sensorId}
	return wsjson.Write(context.Background(), conn, res)
}

func sensorErrorMsg(msg string) map[string]interface{} {
	return map[string]interface{}{"status": "error", "message": msg}
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

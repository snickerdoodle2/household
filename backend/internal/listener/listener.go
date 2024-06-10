package listener

import (
	"encoding/json"
	"fmt"
	"inzynierka/internal/broker"
	"inzynierka/internal/data"
	"io"
	"net/http"
	"time"
)

type ListenerT interface {
	Start() error
	GetBroker() *broker.Broker[[]byte]
	GetStopCh() chan struct{}
	GetCurrentValue() ([]byte, error)
}

func New[T any](sensor *data.Sensor) *Listener[T] {
	return &Listener[T]{
		sensor: sensor,
		values: make([]T, 0),
		StopCh: make(chan struct{}, 2),
		Broker: broker.NewBroker[[]byte](),
	}
}

// TODO: Make it work with SensorResult
type Response[T any] struct {
	Values []T    `json:"values"`
	Status string `json:"status"`
}

type Listener[T any] struct {
	sensor *data.Sensor
	values []T
	StopCh chan struct{}
	Broker *broker.Broker[[]byte]
}

func (l *Listener[T]) Start() error {
	var input struct {
		Value T `json:"value"`
	}

	go l.Broker.Start()
	defer l.Broker.Stop()

	delayMultiplier := 1
	for {
		select {
		case <-l.StopCh:
			return nil
		default:
		}
		delay := delayMultiplier * l.sensor.RefreshRate
		time.Sleep(time.Duration(delay) * time.Second)

		res, err := http.Get(fmt.Sprintf("http://%v/value", l.sensor.URI))
		if err != nil {
			fmt.Printf(err.Error())
			msg := Response[T]{
				Status: "OFFLINE",
				Values: make([]T, 0),
			}

			json, err := json.Marshal(msg)
			if err != nil {
				return err
			}
			l.Broker.Publish(json)

			delayMultiplier += 1
			continue
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		err = json.Unmarshal(body, &input)
		if err != nil {
			return err
		}

		l.values = append(l.values, input.Value)
		if len(l.values) > 5 {
			l.values = l.values[1:]
		}

		msg := Response[T]{
			Status: "ONLINE",
			Values: l.values,
		}
		json, err := json.Marshal(msg)
		if err != nil {
			return err
		}
		l.Broker.Publish(json)

		delayMultiplier = 1
	}
}

func (l *Listener[T]) GetBroker() *broker.Broker[[]byte] {
	return l.Broker
}

func (l *Listener[T]) GetStopCh() chan struct{} {
	return l.StopCh
}

func (l *Listener[T]) GetCurrentValue() ([]byte, error) {
	var msg Response[T]
	if l.values == nil {
		msg = Response[T]{
			Status: "OFFLINE",
			Values: make([]T, 0),
		}
	} else {
		msg = Response[T]{
			Status: "ONLINE",
			Values: l.values,
		}
	}

	json, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return json, nil
}

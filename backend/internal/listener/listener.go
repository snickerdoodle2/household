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

func New[T any](sensor *data.Sensor) *Listener[T] {
	return &Listener[T]{
		sensor: sensor,
		values: make([]T, 0),
		StopCh: make(chan struct{}),
		Broker: broker.NewBroker[Response[T]](),
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
	Broker *broker.Broker[Response[T]]
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

		res, err := http.Get(fmt.Sprintf("%v/value", l.sensor.URI))
		if err != nil {
			l.Broker.Publish(Response[T]{
				Status: "OFFLINE",
				Values: make([]T, 0),
			})

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

		l.Broker.Publish(Response[T]{
			Status: "ONLINE",
			Values: l.values,
		})

		delayMultiplier = 1
		delay := delayMultiplier * l.sensor.RefreshRate
		time.Sleep(time.Duration(delay) * time.Second)
	}
}

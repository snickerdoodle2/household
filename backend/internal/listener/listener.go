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

type Response[T data.SensorReturn] struct {
	values []T
	status string
}

type Listener[T data.SensorReturn] struct {
	Sensor *data.Sensor
	Values []T
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

		res, err := http.Get(fmt.Sprintf("%v/value", l.Sensor.URI))
		if err != nil {
			l.Broker.Publish(Response[T]{
				status: "OFFLINE",
				values: make([]T, 0),
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

		l.Values = append(l.Values, input.Value)
		if len(l.Values) > 5 {
			l.Values = l.Values[1:]
		}

		l.Broker.Publish(Response[T]{
			status: "ONLINE",
			values: l.Values,
		})

		delayMultiplier = 1
		delay := delayMultiplier * l.Sensor.RefreshRate
		time.Sleep(time.Duration(delay) * time.Second)
	}
}

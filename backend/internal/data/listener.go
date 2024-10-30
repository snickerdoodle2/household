package data

import (
	"encoding/json"
	"fmt"
	"inzynierka/internal/broker"
	"io"
	"net/http"
	"time"
)

func NewListener[T any](sensor *Sensor, onNewValue func(T)) *Listener[T] {
	return &Listener[T]{
		sensor:     sensor,
		values:     make([]T, 0),
		StopCh:     make(chan struct{}, 2),
		Broker:     broker.NewBroker[[]T](),
		onNewValue: onNewValue,
	}
}

// TODO: Make it work with SensorResult
type Response[T any] struct {
	Values []T    `json:"values"`
	Status string `json:"status"`
}

type Listener[T any] struct {
	sensor     *Sensor
	values     []T
	StopCh     chan struct{}
	Broker     *broker.Broker[[]T]
	onNewValue func(T)
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
			fmt.Print(err.Error())

			l.Broker.Publish(nil)

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

		// TODO: wyrzucic do konfigu
		l.values = append(l.values, input.Value)
		if len(l.values) > 5 {
			l.values = l.values[1:]
		}

		l.onNewValue(input.Value)

		l.Broker.Publish(l.values)

		delayMultiplier = 1
	}
}

func (l *Listener[T]) GetBroker() *broker.Broker[[]T] {
	return l.Broker
}

func (l *Listener[T]) GetStopCh() chan struct{} {
	return l.StopCh
}

func (l *Listener[T]) GetCurrentValue() []T {
	return l.values
}

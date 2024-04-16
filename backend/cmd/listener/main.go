package main

import (
	"fmt"
	"inzynierka/internal/broker"
	"inzynierka/internal/data"
	"inzynierka/internal/listener"
	"time"
)

// Assuming that a decimal sensor is listening on 127.0.0.1:10002
func main() {
	stpCh := make(chan struct{})
	l := listener.Listener[float64]{
		Sensor: &data.Sensor{
			URI:         "http://localhost:10002",
			RefreshRate: 1,
		},
		Values: make([]float64, 0),
		Broker: broker.NewBroker[listener.Response[float64]](),
		StopCh: stpCh,
	}

	clientFunc := func(id int) {
		msgCh := l.Broker.Subscribe()
		for {
			fmt.Printf("Client %d got message: %v\n", id, <-msgCh)
		}
	}

	for i := 0; i < 3; i++ {
		go clientFunc(i)
	}

	fmt.Println("Starting...")

	go l.Start()
	time.Sleep(10 * time.Second)
	l.StopCh <- struct{}{}

	time.Sleep(10 * time.Second)
}

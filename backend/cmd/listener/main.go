package main

import (
	"encoding/json"
	"fmt"
	"inzynierka/internal/data/broker"
	"io"
	"net/http"
	"time"
)

// Assuming that a decimal sensor is listening on 127.0.0.1:10002
func main() {
	b := broker.NewBroker[float32]()
	go b.Start()

	clientFunc := func(id int) {
		msgCh := b.Subscribe()
		for {
			fmt.Printf("Client %d got message: %v\n", id, <-msgCh)
		}
	}

	for i := 0; i < 3; i++ {
		go clientFunc(i)
	}

	go func() {
		for {
			res, err := http.Get("http://127.0.0.1:10002/value")
			if err != nil {
				continue
			}

			body, err := io.ReadAll(res.Body)

			var input struct {
				Value float32 `json:"value"`
			}

			err = json.Unmarshal(body, &input)
			if err != nil {
				continue
			}

			b.Publish(input.Value)

			time.Sleep(5 * time.Second)
		}
	}()

	time.Sleep(time.Minute)
}

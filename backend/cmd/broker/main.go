package main

import (
	"fmt"
	"inzynierka/internal/data/broker"
	"time"
)

func main() {
	b := broker.NewBroker[string]()
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
		for msgId := 0; ; msgId++ {
			b.Publish(fmt.Sprintf("msg#%d", msgId))
			time.Sleep(300 * time.Millisecond)
		}
	}()

	time.Sleep(time.Second)
}

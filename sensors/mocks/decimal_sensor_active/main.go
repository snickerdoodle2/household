package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const (
	defaultAddress = "localhost:8080"
)

type Measurement struct {
	MessageType string  `json:"message-type"`
	SensorType  string  `json:"sensor-type"`
	Value       float64 `json:"value"`
}

func main() {
	interval := flag.Int("interval", 5, "Interval between measurements in seconds")
	address := flag.String("address", defaultAddress, "Target address in format host:port")
	flag.Parse()

	ticker := time.NewTicker(time.Duration(*interval) * time.Second)
	defer ticker.Stop()

	log.Printf("Starting sensor simulator - sending measurements to %s every %d seconds\n", *address, *interval)

	for range ticker.C {
		measurement := Measurement{
			MessageType: "measurement",
			SensorType:  "decimal_sensor",
			Value:       rand.NormFloat64() * 20,
		}

		// Convert measurement to JSON
		jsonData, err := json.Marshal(measurement)
		if err != nil {
			log.Printf("Error marshaling JSON: %v", err)
			continue
		}

		resp, err := http.Post(fmt.Sprintf("http://%s", *address), "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("Error sending request: %v", err)
			continue
		}
		resp.Body.Close()

		log.Printf("Sent measurement: %+v", measurement)
	}
}

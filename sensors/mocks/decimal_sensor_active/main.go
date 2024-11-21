package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"time"
)

const (
	defaultAddress = "localhost:8080"
	endpoint       = "/api/v1/sensor/measurements/active"
)

type Measurement struct {
	MessageType string  `json:"message-type"`
	SensorType  string  `json:"sensor-type"`
	Value       float64 `json:"value"`
}

func main() {
	interval := flag.Int("interval", 5, "Interval between measurements in seconds")
	address := flag.String("address", defaultAddress, "Target address in format host:port")
	srcIP := flag.String("src-ip", "0.0.0.0", "Source IP address to bind to")
	srcPort := flag.Int("src-port", 9002, "Source port to bind to (default 9002)")
	flag.Parse()

	if *srcPort < 0 || *srcPort > 65535 {
		log.Fatalf("Invalid source port: %d. Must be between 0 and 65535.", *srcPort)
	}

	ticker := time.NewTicker(time.Duration(*interval) * time.Second)
	defer ticker.Stop()

	log.Printf("Starting sensor simulator - sending measurements to %s every %d seconds\n", *address, *interval)
	log.Printf("Using source IP: %s and port: %d\n", *srcIP, *srcPort)

	// Configure HTTP client with a custom transport using the source IP and port
	localAddr := &net.TCPAddr{
		IP:   net.ParseIP(*srcIP),
		Port: *srcPort,
	}
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			LocalAddr: localAddr,
			Timeout:   30 * time.Second,
		}).DialContext,
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

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

		// Send POST request
		url := fmt.Sprintf("http://%s%s", *address, endpoint)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("Error creating request: %v", err)
			continue
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Error sending request: %v", err)
			continue
		}
		resp.Body.Close()

		log.Printf("Sent measurement: %+v", measurement)
	}
}
